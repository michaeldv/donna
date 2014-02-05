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

type MoveWithScore struct {
        move   Move
        score  int
}

type MoveGen struct {
        p     *Position
        list  [256]MoveWithScore
        head  int
        tail  int
        step  int
        ply   int
}

var moveList [MaxPly]MoveGen

func (p *Position) StartMoveGen(ply int) (gen *MoveGen) {
        gen = &moveList[ply]
        gen.p = p
        gen.list = [256]MoveWithScore{}
        gen.head, gen.tail = 0, 0
        gen.ply = ply
        return
}

func (gen *MoveGen) NextMove() (move Move) {
        if gen.head == gen.tail {
                return 0
        }
        move = gen.list[gen.head].move
        gen.head++
        return
}


func (gen *MoveGen) GenerateQuiets() *MoveGen {
        return gen
}

// Return a list of generated moves by continuously calling the next move
// until the list is empty.
func (gen *MoveGen) allMoves() (moves []Move) {
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		moves = append(moves, move)
	}
	return
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

func (p *Position) pawnMovesMask(color int) (mask Bitmask) {
        if color == White {
                mask = (p.outposts[Pawn(White)] << 8)
        } else {
                mask = (p.outposts[Pawn(Black)] >> 8)
        }
        mask &= ^p.board[2]
        return
}

func (p *Position) pawnJumpsMask(color int) (mask Bitmask) {
        if color == White {
                mask = maskRank[3] & (p.outposts[Pawn(White)] << 16)
        } else {
                mask = maskRank[4] & (p.outposts[Pawn(Black)] >> 16)
        }
        mask &= ^p.board[2]
        return
}
