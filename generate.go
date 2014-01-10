// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`sort`)

// All moves.
func (p *Position) Moves(ply int) (moves []*Move) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.color() == p.color {
                        moves = append(moves, p.possibleMoves(i, piece)...)
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

        for targets.isNotEmpty() {
                target := targets.firstSet()
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
                targets.clear(target)
        }
        return
}

// All capture moves for the piece in certain square. This might include
// illegal moves that cause check to the king.
func (p *Position) possibleCaptures(square int, piece Piece) (moves []*Move) {
        targets := p.targets[square]

        for targets.isNotEmpty() {
                target := targets.firstSet()
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
                } else if p.enpassant != 0 && target == p.enpassant {
                        moves = append(moves, NewMove(p, square, target))
                }
                targets.clear(target)
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
                } else if move.promoted != 0 {
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
