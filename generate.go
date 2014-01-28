// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import() //(`sort`)

const (
        stepPrincipal = iota
        stepCaptures
        stepPromotions
        stepKillers
        stepRemaining
)

type MoveEx struct {
        move   Move
        score  int
}

type MoveList struct {
        position  *Position
        moves     [256]MoveEx
        ply       int
        head      int
        tail      int
        step      int
}

var moveList [MaxPly]MoveList

func (p *Position) StartMoveGen(ply int) (ml *MoveList) {
        ml = &moveList[ply]
        ml.position = p
        ml.moves = [256]MoveEx{}
        ml.ply = ply
        ml.head, ml.tail = 0, 0
        return
}

func (ml *MoveList) NextMove() (move Move) {
        if ml.head == ml.tail {
                return 0
        }
        move = ml.moves[ml.head].move
        ml.head++
        return
}

func (ml *MoveList) GenerateMoves() *MoveList {
        color := ml.position.color
        ml.pawnMoves(color)
        ml.pieceMoves(color)
        ml.kingMoves(color)
        return ml
}

func (ml *MoveList) pawnMoves(color int) *MoveList {
        pawns := ml.position.outposts[Pawn(color)]

        for pawns != 0 {
                square := pawns.pop()
                targets := ml.position.targets[square]
                for targets != 0 {
                        target := targets.pop()
                        if target > H1 && target < A8 {
                                ml.moves[ml.tail].move = ml.position.pawnMove(square, target)
                                ml.tail++
                        } else { // Promotion.
                                m1, m2, m3, m4 := ml.position.pawnPromotion(square, target)
                                ml.moves[ml.tail].move = m1
                                ml.tail++
                                ml.moves[ml.tail].move = m2
                                ml.tail++
                                ml.moves[ml.tail].move = m3
                                ml.tail++
                                ml.moves[ml.tail].move = m4
                                ml.tail++
                        }
                }
        }
        return ml
}

func (ml *MoveList) pieceMoves(color int) *MoveList {
	for _, kind := range [4]int{ KNIGHT, BISHOP, ROOK, QUEEN } {
	        outposts := ml.position.outposts[Piece(kind|color)]
	        for outposts != 0 {
	                square := outposts.pop()
	                targets := ml.position.targets[square]
	                for targets != 0 {
	                        ml.moves[ml.tail].move = NewMove(ml.position, square, targets.pop())
	                        ml.tail++
	                }
	        }
	}
        return ml
}

func (ml *MoveList) kingMoves(color int) *MoveList {
        king := ml.position.outposts[King(color)]
        if king != 0 {
                square := king.pop()
                targets := ml.position.targets[square]
                for targets != 0 {
                        target := targets.pop()
                        if square == homeKing[color] && Abs(square - target) == 2 {
                                ml.moves[ml.tail].move = NewCastle(ml.position, square, target)
                        } else {
                                ml.moves[ml.tail].move = NewMove(ml.position, square, target)
                        }
                        ml.tail++
                }
        }
        return ml
}


func (ml *MoveList) GenerateCaptures() *MoveList {
        color := ml.position.color
        ml.pawnCaptures(color)
        ml.pieceCaptures(color)
        return ml
}

// Generates all pseudo-legal pawn captures and Queen promotions.
func (ml *MoveList) pawnCaptures(color int) *MoveList {
        pawns := ml.position.outposts[Pawn(color)]

        for pawns != 0 {
                square := pawns.pop()
                //
                // First check capture targets on rows 2-7 (no promotions).
                //
                targets := ml.position.targets[square] & ml.position.board[color^1] & 0x00FFFFFFFFFFFF00
                for targets != 0 {
                        ml.moves[ml.tail].move = NewMove(ml.position, square, targets.pop())
                        ml.tail++
                }
                //
                // Now check promo rows. The might include capture targets as well
                // as empty promo square in front of the pawn.
                //
                if RelRow(square, color) == 6 {
                        //
                        // Select maskRank[7] for white and maskRank[0] for black.
                        //
                        targets  = ml.position.targets[square] & maskRank[7 - 7 * color]
                        targets |= ml.position.board[2] & Bit(square + eight[color])

                        for targets != 0 {
                                ml.moves[ml.tail].move = NewMove(ml.position, square, targets.pop()).promote(QUEEN)
                                ml.tail++
                        }
                }
        }
        return ml
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (ml *MoveList) pieceCaptures(color int) *MoveList {
	for _, kind := range [5]int{ KNIGHT, BISHOP, ROOK, QUEEN, KING } {
	        outposts := ml.position.outposts[Piece(kind|color)]
	        for outposts != 0 {
	                square := outposts.pop()
	                targets := ml.position.targets[square] & ml.position.board[color^1]
	                for targets != 0 {
	                        ml.moves[ml.tail].move = NewMove(ml.position, square, targets.pop())
	                        ml.tail++
	                }
	        }
	}
	return ml
}


func (ml *MoveList) GenerateEvasions() *MoveList {
        color := ml.position.color
        enemy := ml.position.color^1
        square := ml.position.outposts[King(color)].first()
        pawn, knight, bishop, rook, queen := Pawn(enemy), Knight(enemy), Bishop(enemy), Rook(enemy), Queen(enemy)
        //
        // Find out what pieces are checking the king. Usually it's a single
        // piece but double check is also a possibility.
        //
        checkers := maskPawn[enemy][square] & ml.position.outposts[pawn]
        checkers |= ml.position.Targets(square, knight) & ml.position.outposts[knight]
        checkers |= ml.position.Targets(square, bishop) & (ml.position.outposts[bishop] | ml.position.outposts[queen])
        checkers |= ml.position.Targets(square, rook) & (ml.position.outposts[rook] | ml.position.outposts[queen])
        //
        // Generate possible king retreats first, i.e. moves to squares not
        // occupied by friendly pieces and not attacked by the opponent.
        //
        retreats := kingMoves[square] & ^ml.position.board[color] & ^ml.position.attacks[enemy]
        //
        // If the attacking piece is bishop, rook, or queen then exclude the
        // square behind the king using avasion mask. Note that knight's
        // evasion mask is full board so we only check if the attacking piece
        // is not a pawn.
        //
        attackSquare := checkers.pop()
        if ml.position.pieces[attackSquare] != pawn {
                retreats &= maskEvade[square][attackSquare]
        }
        //
        // If checkers mask is not empty then we've got double check and
        // retreat is the only option.
        //
        if checkers != 0 {
                attackSquare = checkers.first()
                if ml.position.pieces[attackSquare] != pawn {
                        retreats &= maskEvade[square][attackSquare]
                }
                for retreats != 0 {
                        ml.moves[ml.tail].move = NewMove(ml.position, square, retreats.pop())
                        ml.tail++
                }
                return ml
        }
        //
        // Generate king retreats.
        //
        for retreats != 0 {
                ml.moves[ml.tail].move = NewMove(ml.position, square, retreats.pop())
                ml.tail++
        }
        //
        // Pawn captures.
        //
        pawns := maskPawn[color][attackSquare] & ml.position.outposts[Pawn(color)]
        for pawns != 0 {
                ml.moves[ml.tail].move = NewMove(ml.position, square, pawns.pop())
                ml.tail++
        }
        //
        // Rare case when the check could be avoided by en-passant capture.
        // For example: Ke4, c5, e5 vs. Ke8, d7. Black's d7-d5+ could be
        // evaded by c5xd6 or e5xd6 en-passant captures.
        //
        if enpassant := attackSquare + eight[color]; ml.position.flags.enpassant == enpassant {
                pawns := maskPawn[color][enpassant] & ml.position.outposts[Pawn(color)]
                for pawns != 0 {
                        ml.moves[ml.tail].move = NewEnpassant(ml.position, square, pawns.pop())
                        ml.tail++
                }
        }
        //
        // See if the check could be blocked.
        //
        block := maskBlock[square][attackSquare]
        //
        // Handle one square pawn pushes: promote to Queen when reaching last rank.
        //
        pawns = (ml.position.outposts[Pawn(color)] >> uint(eight[color])) & ^(ml.position.board[2]) & block
        for pawns != 0 {
                targetSquare := square + eight[color]
                ml.moves[ml.tail].move = NewMove(ml.position, targetSquare, pawns.pop())
                if targetSquare >= A8 || targetSquare <= H1 {
                        ml.moves[ml.tail].move.promote(QUEEN)
                }
                ml.tail++
        }

        return ml
}

// All moves.
func (p *Position) Moves(ply int) (moves []Move) {
        for square, piece := range p.pieces {
                if piece != 0 && piece.color() == p.color {
                        moves = append(moves, p.possibleMoves(square, piece)...)
                }
        }
        moves = p.reorderMoves(moves, p.game.bestLine[0][ply], p.game.killers[ply])
        Log("%d candidates for %s: %v\n", len(moves), C(p.color), moves)
        return
}

func (p *Position) Captures(ply int) (moves []Move) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.color() == p.color {
                        moves = append(moves, p.possibleCaptures(i, piece)...)
                }
        }
        if bestMove := p.game.bestLine[0][ply]; bestMove != 0 && bestMove.capture() != 0 {
                moves = p.reorderCaptures(moves, bestMove)
        } else {
                //sort.Sort(byScore{moves})
        }

        Log("%d capture candidates for %s: %v\n", len(moves), C(p.color), moves)
        return
}

// All moves for the piece in certain square. This might include illegal
// moves that cause check to the king.
func (p *Position) possibleMoves(square int, piece Piece) (moves []Move) {
        targets := p.targets[square]

        for targets != 0 {
                target := targets.pop()
                //
                // For regular moves each target square represents one possible
                // move. For pawn promotion, however, we have to generate four
                // possible moves, one for each promoted piece.
                //
                if !p.isPawnPromotion(piece, target) {
                        moves = append(moves, NewMove(p, square, target))
                } else {
                        for _,name := range([]int{ QUEEN, ROOK, BISHOP, KNIGHT }) {
                                candidate := NewMove(p, square, target).promote(name)
                                moves = append(moves, candidate)
                        }
                }
        }
        return
}

// All capture moves for the piece in certain square. This might include
// illegal moves that cause check to the king.
func (p *Position) possibleCaptures(square int, piece Piece) (moves []Move) {
        targets := p.targets[square]

        for targets != 0 {
                target := targets.pop()
                capture := p.pieces[target]
                if capture != 0 {
                        if !p.isPawnPromotion(piece, target) {
                                moves = append(moves, NewMove(p, square, target))
                        } else {
                                for _,name := range([]int{ QUEEN, ROOK, BISHOP, KNIGHT }) {
                                        candidate := NewMove(p, square, target).promote(name)
                                        moves = append(moves, candidate)
                                }
                        }
                } else if p.flags.enpassant != 0 && target == p.flags.enpassant {
                        moves = append(moves, NewMove(p, square, target))
                }
        }
        return
}

func (p *Position) reorderMoves(moves []Move, bestMove Move, goodMove [2]Move) []Move {
        var principal, killers, captures, promotions, remaining []Move

        for _, move := range moves {
                if len(principal) == 0 && bestMove != 0 && move == bestMove {
                        principal = append(principal, move)
                } else if move.capture() != 0 {
                        captures = append(captures, move)
                } else if move.promo() != 0 {
                        promotions = append(promotions, move)
                } else if (goodMove[0] != 0 && move == goodMove[0]) || (goodMove[1] != 0 && move == goodMove[1]) {
                        killers = append(killers, move)
                } else {
                        remaining = append(remaining, move)
                }
        }
        if len(killers) > 1 && killers[0] == goodMove[1] {
                killers[0], killers[1] = killers[1], killers[0]
        }

        //sort.Sort(byScore{captures})
        //sort.Sort(byScore{remaining})
        return append(append(append(append(append(principal, captures...), promotions...), killers...), remaining...))
}

func (p *Position) reorderCaptures(moves []Move, bestMove Move) []Move {
        var principal, remaining []Move

        for _, move := range moves {
                if len(principal) == 0 && move == bestMove {
                        principal = append(principal, move)
                } else {
                        remaining = append(remaining, move)
                }
        }
        //sort.Sort(byScore{remaining})
        return append(principal, remaining...)
}

// Sorting moves by their relative score based on piece/square for regular moves
// or least valuaeable attacker/most valueable victim for captures.
// type byScore struct {
//         moves []Move
// }
// func (her byScore) Len() int           { return len(her.moves)}
// func (her byScore) Swap(i, j int)      { her.moves[i], her.moves[j] = her.moves[j], her.moves[i] }
// func (her byScore) Less(i, j int) bool { return her.moves[i].score > her.moves[j].score }

func (p *Position) pawnMove(square, target int) Move {
        if RelRow(square, p.color) == 1 && RelRow(target, p.color) == 3 {
                if p.isEnpassant(target, p.color) {
                        return NewEnpassant(p, square, target)
                } else {
                        return NewPawnJump(p, square, target)
                }
        }

        return NewMove(p, square, target)
}

func (p *Position) pawnPromotion(square, target int) (Move, Move, Move, Move) {
        return NewMove(p, square, target).promote(QUEEN),
               NewMove(p, square, target).promote(ROOK),
               NewMove(p, square, target).promote(BISHOP),
               NewMove(p, square, target).promote(KNIGHT)
}

func (p *Position) isEnpassant(target, color int) bool {
        pawns := p.outposts[Pawn(color^1)] // Opposite color pawns.
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
