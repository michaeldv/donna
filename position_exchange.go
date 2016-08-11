// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

var exchangeScores = []int{
	0, 0, 						// Zero score for non-capture moves.
	valuePawn.midgame, valuePawn.midgame, 		// Pawn/BlackPawn captures.
	valueKnight.midgame, valueKnight.midgame, 	// Knight/BlackKinght captures.
	valueBishop.midgame, valueBishop.midgame, 	// Bishop/BlackBishop captures.
	valueRook.midgame, valueRook.midgame, 		// Rook/BlackRook captures.
	valueQueen.midgame, valueQueen.midgame, 	// Queen/BlackQueen captures.
	valueQueen.midgame * 8, valueQueen.midgame * 8, // King/BlackKing specials.
}

// Static exchange evaluation.
func (p *Position) exchange(move Move) int {
	from, to, piece, capture := move.split()

	score := exchangeScores[capture]
	if promo := move.promo(); promo.some() {
		score += exchangeScores[promo] - exchangeScores[Pawn]
		piece = promo
	}

	board := p.board ^ bit[from]
	return -p.exchangeScore(piece.color()^1, to, -score, exchangeScores[piece], board)
}

// Recursive helper method for the static exchange evaluation.
func (p *Position) exchangeScore(color int, to, score, extra int, board Bitmask) int {
	attackers := p.attackers(color, to, board) & board
	if attackers.empty() {
		return score
	}

	from, best := 0, Checkmate
	for bm := attackers; bm.any(); bm = bm.pop() {
		square := bm.first()
		if index := p.pieces[square]; exchangeScores[index] < best {
			from = square
			best = exchangeScores[index]
		}
	}

	if best != Checkmate {
		board ^= bit[from]
	}

	return max(score, -p.exchangeScore(color^1, to, -(score + extra), best, board))
}
