// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`sort`)

// All moves.
func (p *Position) Moves() (moves []*Move) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.Color() == p.color {
                        moves = append(moves, p.PossibleMoves(i, piece)...)
                }
        }
        if len(moves) > 1 {
		moves = p.Reorder(moves)
        }
        Log("%d candidates for %s: %v\n", len(moves), C(p.color), moves)

        return
}

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
                        move := NewMove(square, target, piece, capture)
                        if !p.isInvalidCastle(move) {
                                moves = append(moves, move)
                        }
                } else {
                        for _,name := range([]int{ QUEEN, ROOK, BISHOP, KNIGHT }) {
                                candidate := NewMove(square, target, piece, capture).Promote(name)
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

func (p *Position) Reorder(moves []*Move) []*Move {
        var captures, promotions, remaining []*Move

        for _, move := range moves {
                if move.Captured != 0 {
                        captures = append(captures, move)
                } else if move.Promoted != 0 {
                        promotions = append(promotions, move)
                } else {
                        remaining = append(remaining, move)
                }
        }
        sort.Sort(byScore{remaining, p})
        return append(append(append(captures, promotions...), remaining...))
}

// Sorting moves by their relative score based on piece/square.
type byScore struct {
        moves     []*Move
        position  *Position
}
func (her byScore) Len() int           { return len(her.moves)}
func (her byScore) Swap(i, j int)      { her.moves[i], her.moves[j] = her.moves[j], her.moves[i] }
func (her byScore) Less(i, j int) bool { return her.moves[i].score(her.position) > her.moves[j].score(her.position) }
