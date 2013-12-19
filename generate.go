// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

// All moves.
func (p *Position) Moves() (positions []*Position) {
        var moves []*Move
        for i, piece := range p.pieces {
                if piece != 0 && piece.Color() == p.color {
                        moves = append(moves, p.PossibleMoves(i, piece)...)
                }
        }

        positions = p.Reorder(moves)
        Log("%d candidates for %s: %v\n", len(moves), C(p.color), moves)

        return
}

// TODO: refactor to return positions.
func (p *Position) Captures() (moves []*Move) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.Color() == p.color {
                        moves = append(moves, p.PossibleCaptures(i, piece)...)
                }
        }
        Log("%d capture candidates for %s: %v\n", len(moves), C(p.color), moves)

        return
}

// All moves for the piece in certain square. This might include illegal
// moves that cause check to the king.
func (p *Position) PossibleMoves(square int, piece Piece) (moves []*Move) {
        targets := p.targets[square]

        for targets.IsNotEmpty() {
                target := targets.FirstSet()
                capture := p.pieces[target]
                //
                // For regular moves each target square represents one possible
                // move. For pawn promotion, however, we have to generate four
                // possible moves, one for each promoted piece.
                //
                if !p.isPawnPromotion(piece, target) {
                        moves = append(moves, NewMove(square, target, piece, capture))
                } else {
                        for _,name := range([]int{ QUEEN, ROOK, BISHOP, KNIGHT }) {
                                candidate := NewMove(square, target, piece, capture).Promote(name)
                                moves = append(moves, candidate)
                        }
                }
                targets.Clear(target)
        }
        if castle := p.tryCastle(); castle != nil {
                moves = append(moves, castle)
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
                if capture != 0  {
                        if !p.isPawnPromotion(piece, target) {
                                moves = append(moves, NewMove(square, target, piece, capture))
                        } else {
                                for _,name := range([]int{ QUEEN, ROOK, BISHOP, KNIGHT }) {
                                        candidate := NewMove(square, target, piece, capture).Promote(name)
                                        moves = append(moves, candidate)
                                }
                        }
                }
                targets.Clear(target)
        }
        return
}

func (p *Position) Reorder(moves []*Move) []*Position {
        var checks, promotions, captures, remaining []*Position

        for _, move := range moves {
                position := p.MakeMove(move)
                if position.inCheck {
                        checks = append(checks, position)
                } else if move.Promoted != 0 {
                        promotions = append(promotions, position)
                } else if move.Captured != 0 {
                        captures = append(captures, position)
                } else {
                        remaining = append(remaining, position)
                }
        }

        return append(append(append(captures, promotions...), checks...), remaining...)
}

