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

	threats.white = e.threats(White)
	threats.black = e.threats(Black)
	score.add(threats.white).sub(threats.black).apply(weightThreats)
	e.score.add(score)

	if e.material.turf != 0 && e.material.flags & (whiteKingSafety | blackKingSafety) != 0 {
		center.white = e.center(White)
		center.black = e.center(Black)
		score.clear().add(center.white).sub(center.black).apply(weightCenter)
		e.score.add(score)
	}
}

func (e *Evaluation) threats(our uint8) (score Score) {
	p, their := e.position, our^1

	// Get our protected and non-hanging pawns.
	pawns := p.outposts[pawn(our)] & (e.attacks[our] | ^e.attacks[their])

	// Find enemy pieces attacked by our protected/non-hanging pawns.
	pieces := p.outposts[their] ^ p.outposts[king(their)] 	// All pieces except king.
	majors := pieces ^ p.outposts[pawn(their)]		// All pieces except king and pawns.
	targets := majors & p.pawnTargets(our, pawns)

	// Bonus for each enemy piece attacked by our pawn.
	for targets.any() {
		piece := p.pieces[targets.pop()]
		score.add(bonusPawnThreat[piece.id()])
	}

	// Find enemy pieces that might be our likely targets: major pieces
	// attacked by our pawns and all attacked pieces not defended by pawns.
	defended := majors & e.attacks[pawn(their)]
	undefended := pieces & ^e.attacks[pawn(their)] & e.attacks[our]

	if likely := defended | undefended; likely.any() {
		// Bonus for enemy pieces attacked by knights and bishops.
		targets = likely & (e.attacks[knight(our)] | e.attacks[bishop(our)])
		for targets.any() {
			piece := p.pieces[targets.pop()]
			score.add(bonusMinorThreat[piece.id()])
		}

		// Bonus for enemy pieces attacked by rooks.
		targets = (undefended | p.outposts[queen(their)]) & e.attacks[rook(our)]
		for targets.any() {
			piece := p.pieces[targets.pop()]
			score.add(bonusRookThreat[piece.id()])
		}

		// Bonus for enemy pieces attacked by the king.
		if targets = undefended & e.attacks[king(our)]; targets.any() {
			if count := targets.count(); count == 1 {
				score.add(Score{2, 30})
			} else if count > 1 {
				score.add(Score{2, 30}.times(2))
			}
		}

		// Extra bonus when attacking enemy pieces that are hanging.
		if targets = undefended & ^e.attacks[their]; targets.any() {
			if count := targets.count(); count > 0 {
				score.add(hangingAttack.times(count))
			}
		}
	}

	return score
}

func (e *Evaluation) center(our uint8) (score Score) {
	p, their := e.position, our^1

	turf := p.outposts[pawn(our)]
	safe := homeTurf[our] & ^turf & ^e.attacks[pawn(their)] & (e.attacks[our] | ^e.attacks[their])

	if our == White {
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
