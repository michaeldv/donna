// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

type Score struct {
	midgame int
	endgame int
}

type Evaluator struct {
	stage    int
	midgame  int
	endgame  int
	position *Position
}

func (p *Position) Evaluate() int {
	evaluator := &Evaluator{0, rightToMove.midgame, rightToMove.endgame, p}
	evaluator.analyzeMaterial()
	evaluator.analyzeCoordination()
	evaluator.analyzePawnStructure()
	evaluator.analyzeRooks()
	evaluator.analyzeKingShield()
	// evaluator.analyzeKingSafety()

	if p.color == Black {
		evaluator.midgame = -evaluator.midgame
		evaluator.endgame = -evaluator.endgame
	}
	return p.score(evaluator.midgame, evaluator.endgame)
}

func (e *Evaluator) analyzeMaterial() {
	counters := &e.position.count
	for _, piece := range []Piece{Pawn, Knight, Bishop, Rook, Queen} {
		count := counters[piece] - counters[piece|Black]
		midgame, endgame := piece.value()
		e.midgame += midgame * count
		e.endgame += endgame * count
	}
}

func (e *Evaluator) analyzeCoordination() {
	var moves, attacks [2]int
	var bonus [2]Score

	p := e.position
	notAttacked := [2]Bitmask{^p.attacks(White), ^p.attacks(Black)}
	board := p.board
	for board != 0 {
		square := board.pop()
		piece := p.pieces[square]
		color := piece.color()
		targets := p.targets(square)

		// Mobility: how many moves are available to squares not attacked by
		// the opponent?
		moves[color] += (targets & notAttacked[color^1]).count()

		// Agressivness: how many opponent's pieces are being attacked?
		attacks[color] += (targets & p.outposts[color^1]).count()

		// Calculate bonus or penalty for a piece being at the given square.
		midgame, endgame := piece.bonus(flip[color][square])
		bonus[color].midgame += midgame
		bonus[color].endgame += endgame
	}

	e.midgame += bonus[White].midgame - bonus[Black].midgame
	e.endgame += bonus[White].endgame - bonus[Black].endgame

	mobility := moves[White] - moves[Black]
	e.midgame += mobility * movesAvailable.midgame
	e.endgame += mobility * movesAvailable.endgame

	aggression := attacks[White] - attacks[Black]
	e.midgame += aggression * attackForce.midgame
	e.endgame += aggression * attackForce.endgame

	if bishops := p.count[Bishop]; bishops >= 2 {
		e.midgame += bishopPair.midgame
		e.endgame += bishopPair.endgame
	}
	if bishops := p.count[BlackBishop]; bishops >= 2 {
		e.midgame -= bishopPair.midgame
		e.endgame -= bishopPair.endgame
	}
}

func (e *Evaluator) analyzeRooks() {
	white := e.rooksScore(White)
	black := e.rooksScore(Black)
	e.midgame += white.midgame - black.midgame
	e.endgame += white.endgame - black.endgame
}

func (e *Evaluator) rooksScore(color int) (bonus Score) {
	p := e.position
	rooks := p.outposts[rook(color)]
	if rooks == 0 {
		return bonus
	}
	//
	// Bonus if rooks are on 7th rank.
	//
	if count := (rooks & mask7th[color]).count(); count > 0 {
		bonus.midgame += count * rookOn7th.midgame
		bonus.endgame += count * rookOn7th.endgame
	}
	//
	// Bonuses if rooks are on open or semi-open files.
	//
	hisPawns := p.outposts[pawn(color)]
	herPawns := p.outposts[pawn(color^1)]
	for rooks != 0 {
		square := rooks.pop()
		column := Col(square)
		if hisPawns&maskFile[column] == 0 {
			if herPawns&maskFile[column] == 0 {
				bonus.midgame += rookOnOpen.midgame
				bonus.endgame += rookOnOpen.endgame
			} else {
				bonus.midgame += rookOnSemiOpen.midgame
				bonus.endgame += rookOnSemiOpen.endgame
			}
		}
	}
	return
}

func (e *Evaluator) analyzeKingShield() {
	// No endgame bonus or penalty.
	e.midgame += e.kingShieldScore(White) - e.kingShieldScore(Black)
}

func (e *Evaluator) kingShieldScore(color int) (penalty int) {
	p := e.position
	kings, pawns := p.outposts[king(color)], p.outposts[pawn(color)]
	//
	// Pass if a) the king is missing, b) the king is on the initial square
	// or c) the opposite side doesn't have a queen with one major piece.
	//
	if kings == 0 || kings == bit[homeKing[color]] || !e.strongEnough(color^1) {
		return
	}
	//
	// Calculate relative square for the king so we could treat black king
	// as white. Don't bother with the shield if the king is too far.
	//
	square := flip[color^1][p.king[color]]
	if square > H3 {
		return
	}
	row, col := Coordinate(square)
	from, to := Max(0, col-1), Min(7, col+1)
	//
	// For each of the shield columns find the closest same color pawn. The
	// penalty is carried if the pawn is missing or is too far from the king
	// (more than one row apart).
	//
	for column := from; column <= to; column++ {
		if shield := (pawns & maskFile[column]); shield != 0 {
			closest := flip[color^1][shield.first()] // Make it relative.
			if distance := Abs(Row(closest) - row); distance > 1 {
				penalty += distance * -shieldDistance.midgame
			}
		} else {
			penalty += -shieldMissing.midgame
		}
	}
	// Log("penalty[%s] => %d\n", C(color), penalty)
	return
}

func (e *Evaluator) strongEnough(color int) bool {
	p := e.position
	return p.count[queen(color)] > 0 &&
		(p.count[rook(color)] > 0 || p.count[bishop(color)] > 0 || p.count[knight(color)] > 0)
}
