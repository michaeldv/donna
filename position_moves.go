// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna
import(`regexp`)

func (p *Position) NewMove(from, to int) Move {
        piece, capture := p.pieces[from], p.pieces[to]

        if p.flags.enpassant != 0 && to == p.flags.enpassant {
                capture = pawn(piece.color()^1)
        }

        return Move(from | (to << 8) | (int(piece) << 16) | (int(capture) << 20))
}

func (p *Position) NewCastle(from, to int) Move {
        return Move(from | (to << 8) | (int(p.pieces[from]) << 16) | isCastle)
}

func (p *Position) NewEnpassant(from, to int) Move {
        return Move(from | (to << 8) | (int(p.pieces[from]) << 16) | isEnpassant)
}

func (p *Position) NewPawnJump(from, to int) Move {
        return Move(from | (to << 8) | (int(p.pieces[from]) << 16) | isPawnJump)
}

// Returns true if *non-evasion* move is valid, i.e. it is possible to make
// the move in current position without violating chess rules. If the king is
// in check the generator is expected to generate valid evasions where extra
// validation is not needed.
func (p *Position) isValid(move Move) bool {
        color := move.color() // TODO: make color part of move split.
        from, to, piece, capture := move.split()
        square := p.outposts[king(color)].first()
        pinned := p.pinnedMask(square)
        //
        // For rare en-passant pawn captures we validate the move by actually
        // making it, and then taking it back.
        //
        if p.flags.enpassant != 0 && to == p.flags.enpassant && capture.isPawn() {
                if position := p.MakeMove(move); position != nil {
                        position.TakeBack(move)
                        return true
                }
                return false
        }
        //
        // King's move is valid when the destination square is not being
        // attacked by the opponent or when the move is a castle.
        //
        if piece.isKing() {
                return !p.isAttacked(to, color^1) || (move & isCastle != 0)
        }
        //
        // For all other peices the move is valid when it doesn't cause a
        // check. For pinned sliders this includes moves along the pinning
        // file, rank, or diagonal.
        //
        return pinned == 0 || pinned.isClear(from) || IsBetween(from, to, square)
}

// Returns a bitmask of all pinned pieces preventing a check for the king on
// given square. The color of the pieces match the color of the king.
func (p *Position) pinnedMask(square int) (mask Bitmask) {
        color := p.pieces[square].color()
        enemy := color^1
        attackers := (p.outposts[bishop(enemy)] | p.outposts[queen(enemy)]) & bishopMagicMoves[square][0]
        attackers |= (p.outposts[rook(enemy)] | p.outposts[queen(enemy)]) & rookMagicMoves[square][0]

        for attackers != 0 {
                attackSquare := attackers.pop()
                blockers := maskBlock[square][attackSquare] & ^bit[attackSquare] & p.board

                if blockers.count() == 1 {
                        mask |= blockers & p.outposts[color] // Only friendly pieces are pinned.
                }
        }
        return
}

func (p *Position) pawnMove(square, target int) Move {
        if RelRow(square, p.color) == 1 && RelRow(target, p.color) == 3 {
                if p.causesEnpassant(target) {
                        return p.NewEnpassant(square, target)
                } else {
                        return p.NewPawnJump(square, target)
                }
        }

        return p.NewMove(square, target)
}

func (p *Position) pawnPromotion(square, target int) (Move, Move, Move, Move) {
        return p.NewMove(square, target).promote(Queen),
               p.NewMove(square, target).promote(Rook),
               p.NewMove(square, target).promote(Bishop),
               p.NewMove(square, target).promote(Knight)
}

func (p *Position) causesEnpassant(target int) bool {
        pawns := p.outposts[pawn(p.color^1)] // Opposite color pawns.
        switch col := Col(target); col {
        case 0:
                return pawns.isSet(target + 1)
        case 7:
                return pawns.isSet(target - 1)
        default:
                return pawns.isSet(target + 1) || pawns.isSet(target - 1)
        }
        return false
}

func (p *Position) pawnMovesMask(color int) (mask Bitmask) {
        if color == White {
                mask = (p.outposts[Pawn] << 8)
        } else {
                mask = (p.outposts[BlackPawn] >> 8)
        }
        mask &= ^p.board
        return
}

func (p *Position) pawnJumpsMask(color int) (mask Bitmask) {
        if color == White {
                mask = maskRank[3] & (p.outposts[Pawn] << 16)
        } else {
                mask = maskRank[4] & (p.outposts[BlackPawn] >> 16)
        }
        mask &= ^p.board
        return
}

func (p *Position) NewMoveFromString(e2e4 string) (move Move) {
	re := regexp.MustCompile(`([KkQqRrBbNn]?)([a-h])([1-8])-?([a-h])([1-8])([QqRrBbNn]?)`)
	arr := re.FindStringSubmatch(e2e4)

	if len(arr) > 0 {
		name  := arr[1]
		from  := Square(int(arr[3][0]-'1'), int(arr[2][0]-'a'))
		to    := Square(int(arr[5][0]-'1'), int(arr[4][0]-'a'))
		promo := arr[6]

		var piece Piece
		switch name {
		case `K`, `k`:
			piece = king(p.color)
		case `Q`, `q`:
			piece = queen(p.color)
		case `R`, `r`:
			piece = rook(p.color)
		case `B`, `b`:
			piece = bishop(p.color)
		case `N`, `n`:
			piece = knight(p.color)
		default:
			piece = p.pieces[from] // <-- Makes piece character optional.
		}
                if (p.pieces[from] != piece) || (p.targets(from) & bit[to] == 0) {
                        move = 0 // Invalid move.
                } else {
                        move = p.NewMove(from, to)
                        if len(promo) > 0 {
                                switch promo {
                                case `Q`, `q`:
                                        move = move.promote(Queen)
                                case `R`, `r`:
                                        move = move.promote(Rook)
                                case `B`, `b`:
                                        move = move.promote(Bishop)
                                case `N`, `n`:
                                        move = move.promote(Knight)
                                default:
                                        move = 0
                                }
                        }
                }
	} else if e2e4 == `0-0` || e2e4 == `0-0-0` {
                from := p.outposts[king(p.color)].first()
                to := G1
                if e2e4 == `0-0-0` {
                        to = C1
                }
                if p.color == Black {
                        to += 56
                }
                move = p.NewCastle(from, to)
                if !move.isCastle() {
                        move = 0
                }
	}
	return
}
