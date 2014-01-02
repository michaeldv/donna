// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`sort`)

// All moves.
func (p *Position) Moves(ply int) (moves []*Move) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.color() == p.color {
                        moves = append(moves, p.PossibleMoves(i, piece)...)
                }
        }
        moves = p.Reorder(moves, best[0][ply])
        Log("%d candidates for %s: %v\n", len(moves), C(p.color), moves)
        return
}

func (p *Position) Captures() (moves []*Move) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.color() == p.color {
                        moves = append(moves, p.PossibleCaptures(i, piece)...)
                }
        }
        sort.Sort(byScore{moves})
        Log("%d capture candidates for %s: %v\n", len(moves), C(p.color), moves)

        return
}

// All moves for the piece in certain square. This might include illegal
// moves that cause check to the king.
func (p *Position) PossibleMoves(square int, piece Piece) (moves []*Move) {
        targets := p.targets[square]

        for targets.IsNotEmpty() {
                target := targets.FirstSet()
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
                targets.Clear(target)
        }
        return
}

// All capture moves for the piece in certain square. This might include
// illegal moves that cause check to the king.
func (p *Position) PossibleCaptures(square int, piece Piece) (moves []*Move) {
        targets := p.targets[square]

        for targets.IsNotEmpty() {
                target := targets.FirstSet()
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
                targets.Clear(target)
        }
        return
}

func (p *Position) Reorder(moves []*Move, bestMove *Move) []*Move {
        var principal, captures, promotions, remaining []*Move

        for _, move := range moves {
                if bestMove != nil && len(principal) == 0 && move.is(bestMove) {
                        principal = append(principal, move)
                } else if move.captured != 0 {
                        captures = append(captures, move)
                } else if move.promoted != 0 {
                        promotions = append(promotions, move)
                } else {
                        remaining = append(remaining, move)
                }
        }
        sort.Sort(byScore{captures})
        sort.Sort(byScore{remaining})
        return append(append(append(append(principal, captures...), promotions...), remaining...))
}

// Sorting moves by their relative score based on piece/square for regular moves
// or least valuaeable attacker/most valueable victim for captures.
type byScore struct {
        moves     []*Move
}
func (her byScore) Len() int           { return len(her.moves)}
func (her byScore) Swap(i, j int)      { her.moves[i], her.moves[j] = her.moves[j], her.moves[i] }
func (her byScore) Less(i, j int) bool { return her.moves[i].score > her.moves[j].score }
