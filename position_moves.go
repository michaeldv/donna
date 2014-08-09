// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`regexp`
)

// Returns true if *non-evasion* move is valid, i.e. it is possible to make
// the move in current position without violating chess rules. If the king is
// in check the generator is expected to generate valid evasions where extra
// validation is not needed.
func (p *Position) isValid(move Move, pins Bitmask) bool {
	color := move.color() // TODO: make color part of move split.
	from, to, piece, capture := move.split()

	// For rare en-passant pawn captures we validate the move by actually
	// making it, and then taking it back.
	if p.enpassant != 0 && to == p.enpassant && capture.isPawn() {
		if position := p.MakeMove(move); position != nil {
			position.UndoLastMove()
			return true
		}
		return false
	}

	// King's move is valid when a) the move is a castle or b) the destination
	// square is not being attacked by the opponent.
	if piece.isKing() {
		return (move & isCastle != 0) || !p.isAttacked(to, color^1)
	}

	// For all other peices the move is valid when it doesn't cause a
	// check. For pinned sliders this includes moves along the pinning
	// file, rank, or diagonal.
	return pins == 0 || pins.isClear(from) || IsBetween(from, to, p.king[color])
}

// Returns a bitmask of all pinned pieces preventing a check for the king on
// given square. The color of the pieces match the color of the king.
func (p *Position) pinnedMask(square int) (mask Bitmask) {
	color := p.pieces[square].color()
	enemy := color ^ 1
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
	if Abs(square - target) == 16 && p.causesEnpassant(target) {
		return NewEnpassant(p, square, target)
	}

	return NewMove(p, square, target)
}

func (p *Position) pawnPromotion(square, target int) (Move, Move, Move, Move) {
	return NewMove(p, square, target).promote(Queen),
	       NewMove(p, square, target).promote(Rook),
	       NewMove(p, square, target).promote(Bishop),
	       NewMove(p, square, target).promote(Knight)
}

// Returns true if a pawn jump causes en-passant. This is done by checking whether
// the enemy pawns occupy squares ajacent to the target square.
func (p *Position) causesEnpassant(target int) bool {
	pawns := p.outposts[pawn(p.color^1)] // Opposite color pawns.

	return maskIsolated[Col(target)] & maskRank[Row(target)] & pawns != 0
}

func (p *Position) NewMoveFromString(e2e4 string) (move Move) {
	re := regexp.MustCompile(`([KkQqRrBbNn]?)([a-h])([1-8])-?([a-h])([1-8])([QqRrBbNn]?)`)
	arr := re.FindStringSubmatch(e2e4)

	if len(arr) > 0 {
		name := arr[1]
		from := Square(int(arr[3][0]-'1'), int(arr[2][0]-'a'))
		to := Square(int(arr[5][0]-'1'), int(arr[4][0]-'a'))
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
		if (p.pieces[from] != piece) || (p.targets(from)&bit[to] == 0) {
			move = 0 // Invalid move.
		} else {
			move = NewMove(p, from, to)
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
		from := p.king[p.color]
		to := G1
		if e2e4 == `0-0-0` {
			to = C1
		}
		if p.color == Black {
			to += 56
		}
		move = NewCastle(p, from, to)
		if !move.isCastle() {
			move = 0
		}
	}
	return
}
