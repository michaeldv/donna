// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`sort`)

type MoveList struct {
        position  *Position
        moves     [256]Move
        ply       int
        head      int
        tail      int
}

var moveList [MaxPly]MoveList

func (p *Position) StartMoveGen(ply int) (ml *MoveList) {
        ml = &moveList[ply]
        ml.position = p
        ml.moves = [256]Move{}
        ml.ply = ply
        ml.head, ml.tail = 0, 0
        return
}

func (ml *MoveList) NextMove() (move *Move) {
        if ml.head == ml.tail {
                return nil
        }
        move = &ml.moves[ml.head]
        ml.head++
        return
}

func (ml *MoveList) GenerateMoves() *MoveList {
        for square, piece := range ml.position.pieces {
                if piece != 0 && piece.color() == ml.position.color {
                        ml.possibleMoves(square, piece)
                }
        }
        return ml
}

func (ml *MoveList) possibleMoves(square int, piece Piece) *MoveList {
        targets := ml.position.targets[square]

        for target := targets.firstPop(); target >= 0; target = targets.firstPop() {
                if !ml.position.isPawnPromotion(piece, target) {
                        ml.moves[ml.tail].add(ml.position, square, target)
                        ml.tail++
                } else {
                        for _,name := range([]int{ QUEEN, ROOK, BISHOP, KNIGHT }) {
                                ml.moves[ml.tail].add(ml.position, square, target).promote(name)
                                ml.tail++
                        }
                }
        }
        return ml
}

func (ml *MoveList) GenerateCaptures() *MoveList {
        for square, piece := range ml.position.pieces {
                if piece != 0 && piece.color() == ml.position.color {
                        ml.possibleCaptures(square, piece)
                }
        }
        return ml
}

func (ml *MoveList) possibleCaptures(square int, piece Piece) *MoveList {
        targets := ml.position.targets[square]

        for target := targets.firstPop(); target >= 0; target = targets.firstPop() {
                capture := ml.position.pieces[target]
                if capture != 0 {
                        if !ml.position.isPawnPromotion(piece, target) {
                                ml.moves[ml.tail].add(ml.position, square, target)
                                ml.tail++
                        } else {
                                for _,name := range([]int{ QUEEN, ROOK, BISHOP, KNIGHT }) {
                                        ml.moves[ml.tail].add(ml.position, square, target).promote(name)
                                        ml.tail++
                                }
                        }
                } else if ml.position.flags.enpassant != 0 && target == ml.position.flags.enpassant {
                        ml.moves[ml.tail].add(ml.position, square, target)
                        ml.tail++
                }
        }
        return ml
}

// All moves.
func (p *Position) Moves(ply int) (moves []*Move) {
        for square, piece := range p.pieces {
                if piece != 0 && piece.color() == p.color {
                        moves = append(moves, p.possibleMoves(square, piece)...)
                }
        }
        moves = p.reorderMoves(moves, p.game.bestLine[0][ply], p.game.killers[ply])
        Log("%d candidates for %s: %v\n", len(moves), C(p.color), moves)
        return
}

func (p *Position) Captures(ply int) (moves []*Move) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.color() == p.color {
                        moves = append(moves, p.possibleCaptures(i, piece)...)
                }
        }
        if bestMove := p.game.bestLine[0][ply]; bestMove != nil && bestMove.captured != 0 {
                moves = p.reorderCaptures(moves, bestMove)
        } else {
                sort.Sort(byScore{moves})
        }

        Log("%d capture candidates for %s: %v\n", len(moves), C(p.color), moves)
        return
}

// All moves for the piece in certain square. This might include illegal
// moves that cause check to the king.
func (p *Position) possibleMoves(square int, piece Piece) (moves []*Move) {
        targets := p.targets[square]

        for target := targets.firstPop(); target >= 0; target = targets.firstPop() {
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
func (p *Position) possibleCaptures(square int, piece Piece) (moves []*Move) {
        targets := p.targets[square]

        for target := targets.firstPop(); target >= 0; target = targets.firstPop() {
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

func (p *Position) reorderMoves(moves []*Move, bestMove *Move, goodMove [2]*Move) []*Move {
        var principal, killers, captures, promotions, remaining []*Move

        for _, move := range moves {
                if len(principal) == 0 && bestMove != nil && move.is(bestMove) {
                        principal = append(principal, move)
                } else if move.captured != 0 {
                        captures = append(captures, move)
                } else if move.flags & isPromotion != 0 {
                        promotions = append(promotions, move)
                } else if (goodMove[0] != nil && move.is(goodMove[0])) || (goodMove[1] != nil && move.is(goodMove[1])) {
                        killers = append(killers, move)
                } else {
                        remaining = append(remaining, move)
                }
        }
        if len(killers) > 1 && killers[0] == goodMove[1] {
                killers[0], killers[1] = killers[1], killers[0]
        }

        sort.Sort(byScore{captures})
        sort.Sort(byScore{remaining})
        return append(append(append(append(append(principal, captures...), promotions...), killers...), remaining...))
}

func (p *Position) reorderCaptures(moves []*Move, bestMove *Move) []*Move {
        var principal, remaining []*Move

        for _, move := range moves {
                if len(principal) == 0 && move.is(bestMove) {
                        principal = append(principal, move)
                } else {
                        remaining = append(remaining, move)
                }
        }
        sort.Sort(byScore{remaining})
        return append(principal, remaining...)
}

// Sorting moves by their relative score based on piece/square for regular moves
// or least valuaeable attacker/most valueable victim for captures.
type byScore struct {
        moves     []*Move
}
func (her byScore) Len() int           { return len(her.moves)}
func (her byScore) Swap(i, j int)      { her.moves[i], her.moves[j] = her.moves[j], her.moves[i] }
func (her byScore) Less(i, j int) bool { return her.moves[i].score > her.moves[j].score }
