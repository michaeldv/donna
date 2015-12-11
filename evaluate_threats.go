// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeThreats() {
	var score Score
	var threats, center Total

	if engine.trace {
		defer func() {
			e.checkpoint(`Threats`, threats)
			if e.material.turf != 0 && e.material.flags & (whiteKingSafety | blackKingSafety) != 0 {
				e.checkpoint(`Center`, center)
			}
		}()
	}

	threats.white = e.threats(White, e.attacks[White], e.attacks[Black])
	threats.black = e.threats(Black, e.attacks[Black], e.attacks[White])
	e.score.add(threats.white).sub(threats.black)

	if e.material.turf != 0 && e.material.flags & (whiteKingSafety | blackKingSafety) != 0 {
		center.white = e.center(White, e.attacks[White], e.attacks[Black], e.attacks[pawn(Black)])
		center.black = e.center(Black, e.attacks[Black], e.attacks[White], e.attacks[pawn(White)])
		score.add(center.white).sub(center.black).apply(weightCenter)
		e.score.add(score)
	}
}

func (e *Evaluation) threats(color uint8, hisAttacks, herAttacks Bitmask) (score Score) {
	p := e.position
	rival := color^1

	// Find enemy pieces under attack excluding king and pawns.
	weak := p.outposts[rival] & ^(p.outposts[king(rival)] | p.outposts[pawn(rival)])
	weak &= hisAttacks

	if weak.any() {

		// Threat bonus for enemy pieces attacked by our pawns.
		targets := weak & e.attacks[pawn(color)]
		for targets.any() {
			piece := p.pieces[targets.pop()]
			score.add(bonusPawnThreat[piece.kind()/2])
		}

		// Threat bonus for enemy pieces attacked by knights and bishops.
		targets = weak & (e.attacks[knight(color)] | e.attacks[bishop(color)])
		for targets.any() {
			piece := p.pieces[targets.pop()]
			score.add(bonusMinorThreat[piece.kind()/2])
		}

		// Threat bonus for enemy pieces attacked by rooks.
		targets = weak & e.attacks[rook(color)]
		for targets.any() {
			piece := p.pieces[targets.pop()]
			score.add(bonusRookThreat[piece.kind()/2])
		}

		// Threat bonus for enemy pieces attacked by the king.
		targets = weak & e.attacks[king(color)]
		if count := targets.count(); count == 1 {
			score.add(Score{1, 29})
		} else if count > 1 {
			score.add(Score{1, 29}.times(2))
		}

		// Extra bonus when attacking enemy pieces that are hanging.
		if hanging := (weak & ^herAttacks).count(); hanging > 0 {
			score.add(hangingAttack.times(hanging))
		}
	}

	return score
}

func (e *Evaluation) center(color uint8, hisAttacks, herAttacks, herPawnAttacks Bitmask) (score Score) {
	turf := e.position.outposts[pawn(color)]
	safe := homeTurf[color] & ^turf & ^herPawnAttacks & (hisAttacks | ^herAttacks)

	if color == White {
		turf |= turf >> 8   // A4..H4 -> A3..H3
		turf |= turf >> 16  // A4..H4 | A3..H3 -> A2..H2 | A1..H1
		turf &= safe 	    // Keep safe squares only.
		safe <<= 32 	    // Move up to black's half of the board.
	} else {
		turf |= turf << 8   // A5..H5 -> A6..H6
		turf |= turf << 16  // A5..H5 | A6..H6 -> A7..H7 | A8..H8
		turf &= safe 	    // Keep safe squares only.
		safe >>= 32 	    // Move down to white's half of the board.
	}

	score.midgame = (safe | turf).count() * e.material.turf / 3

	return score
}
