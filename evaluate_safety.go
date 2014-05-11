// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluator) analyzeSafety() {
	var white, black Score

	// White king's threats.
	square := e.position.king[White] ^ A8
	if e.attacks[Black] > 0 {
		white.midgame -= Max(3, e.attacks[Black]) * e.threats[Black]
		white.endgame -= bonusKing[1][square]
	}

	// Black king's threats.
	square = e.position.king[Black]
	if e.attacks[White] > 0 {
		black.midgame -= Max(3, e.attacks[White]) * e.threats[White]
		black.endgame -= bonusKing[1][square]
	}

	e.score.add(white).subtract(black)

	// No endgame bonus or penalty for king shield.
	e.score.midgame += e.kingShieldScore(White) - e.kingShieldScore(Black)
}

func (e *Evaluator) kingShieldScore(color int) (penalty int) {
	p := e.position
	kings, pawns := p.outposts[king(color)], p.outposts[pawn(color)]

	// Pass if a) the king is missing, b) the king is on the initial square
	// or c) the opposite side doesn't have a queen with one major piece.
	if kings == 0 || kings == bit[homeKing[color]] || !e.strongEnough(color^1) {
		return
	}

	// Calculate relative square for the king so we could treat black king
	// as white. Don't bother with the shield if the king is too far.
	square := Flip(color^1, p.king[color])
	if square > H3 {
		return
	}
	row, col := Coordinate(square)
	from, to := Max(0, col-1), Min(7, col+1)

	// For each of the shield columns find the closest same color pawn. The
	// penalty is carried if the pawn is missing or is too far from the king
	// (more than one row apart).
	for column := from; column <= to; column++ {
		if shield := (pawns & maskFile[column]); shield != 0 {
			closest := Flip(color^1, shield.first()) // Make it relative.
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

