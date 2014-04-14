// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

var exchangeScores = []int{
	valuePawn.midgame, 	// 2->0
	valueKnight.midgame, 	// 4->1
	valueBishop.midgame, 	// 6->2
	valueRook.midgame, 	// 8->3
	valueQueen.midgame, 	//10->4
}

func (p *Position) exchangeScore(to, color, score, extra int, board Bitmask) int {
	attackers := p.attackers(to, color, board) & board
	if attackers == 0 {
		return score
	}

	from, best := 0, Checkmate
	for attackers != 0 {
		square := attackers.pop()
		index := p.pieces[square].kind() / 2 - 1
		if exchangeScores[index] < best {
			from = square
			best = exchangeScores[index]
		}
	}

	if best != Checkmate {
		board ^= bit[from]
	}

	return Max(score, -p.exchangeScore(to, color^1, -(score + extra), best, board))
}

func (p *Position) exchange(move Move) int {
	from, to, piece, capture := move.split()
	color := move.piece().color()

	score := exchangeScores[capture]
	if promo := move.promo(); promo != 0 {
		score += exchangeScores[promo] - exchangeScores[0] // Pawn
		piece = promo
	}

	board := p.board ^ bit[from]
	return -p.exchangeScore(to, color^1, -score, exchangeScores[piece.kind() / 2 - 1], board)
}
