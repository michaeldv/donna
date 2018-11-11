// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (e *Evaluation) analyzeSafety() {
	var score Score
	var cover, safety Total

	if engine.traceʔ {
		defer func() {
			var our, their Score
			e.checkpoint(`+King`, Total{*our.add(cover.white).add(safety.white), *their.add(cover.black).add(safety.black)})
			e.checkpoint(`-Cover`, cover)
			e.checkpoint(`-Safety`, safety)
		}()
	}

	// If any of the pawns or a king have moved then recalculate cover score.
	if e.position.white.home != e.pawns.king[White] {
		e.pawns.cover[White] = e.kingCover(White)
		e.pawns.king[White] = e.position.white.home
	}
	if e.position.black.home != e.pawns.king[Black] {
		e.pawns.cover[Black] = e.kingCover(Black)
		e.pawns.king[Black] = e.position.black.home
	}

	// Fetch king cover score from the pawn cache.
	cover.white.add(e.pawns.cover[White])
	cover.black.add(e.pawns.cover[Black])

	// Calculate king's safety for both sides.
	if e.safety[White].threats > 0 {
		safety.white = e.kingSafety(White)
	}
	if e.safety[Black].threats > 0 {
		safety.black = e.kingSafety(Black)
	}

	// Calculate total king safety and pawn cover score.
	score.add(safety.white).sub(safety.black).apply(weightSafety)
	score.add(cover.white).sub(cover.black)
	e.score.add(score)
}

func (e *Evaluation) kingSafety(our int) (score Score) {
	p, their := e.position, our^1
	safetyIndex, checkers, square := 0, 0, p.pick(our).home

	// Find squares around the king that are being attacked by the
	// enemy and defended by our king only.
	defended := e.attacks[pawn(our)] | e.attacks[knight(our)] |
	            e.attacks[bishop(our)] | e.attacks[rook(our)] |
	            e.attacks[queen(our)]
	weak := e.attacks[king(our)] & e.attacks[their] & ^defended

	// Find possible queen checks on weak squares around the king.
	// We only consider squares where the queen is protected and
	// can't be captured by the king.
	protected := e.attacks[pawn(their)] | e.attacks[knight(their)] |
		     e.attacks[bishop(their)] | e.attacks[rook(their)] |
		     e.attacks[king(their)]
	checks := weak & e.attacks[queen(their)] & protected & ^p.outposts[their]
	if checks.anyʔ() {
		checkers++
		safetyIndex += queenCheck * checks.count()
	}

	// Out of all squares available for enemy pieces select the ones
	// that are not under our attack.
	safe := ^(e.attacks[our] | p.outposts[their])

	// Are there any safe squares from where enemy Knight could give
	// us a check?
	if checks := knightMoves[square] & safe & e.attacks[knight(their)]; checks.anyʔ() {
		checkers++
		safetyIndex += checks.count()
	}

	// Are there any safe squares from where enemy Bishop could give us a check?
	safeBishopMoves := p.bishopMoves(square) & safe
	if checks := safeBishopMoves & e.attacks[bishop(their)]; checks.anyʔ() {
		checkers++
		safetyIndex += checks.count()
	}

	// Are there any safe squares from where enemy Rook could give us a check?
	safeRookMoves := p.rookMoves(square) & safe
	if checks := safeRookMoves & e.attacks[rook(their)]; checks.anyʔ() {
		checkers++
		safetyIndex += checks.count()
	}

	// Are there any safe squares from where enemy Queen could give us a check?
	if checks := (safeBishopMoves | safeRookMoves) & e.attacks[queen(their)]; checks.anyʔ() {
		checkers++
		safetyIndex += queenCheck / 2 * checks.count()
	}

	threatIndex := min(16, e.safety[our].attackers * e.safety[our].threats / 2) +
			(e.safety[our].attacks + weak.count()) * 3 +
			rank(our, square) - e.pawns.cover[our].midgame / 16
	safetyIndex = min(63, max(0, safetyIndex + threatIndex))

	score.midgame -= kingSafety[safetyIndex]

	if checkers > 0 {
		score.add(rightToMove)
		if checkers > 1 {
			score.add(rightToMove)
		}
	}

	return score
}

func (e *Evaluation) kingCover(our int) (score Score) {
	p, square := e.position, e.position.pick(our).home

	// Don't bother with the cover if the king is too far out.
	if rank(our, square) <= A3H3 {
		// If we still have castle rights encourage castle pawns to stay intact
		// by scoring least safe castle.
		score.midgame = e.kingCoverBonus(our, square)
		if p.castles & castleKingside[our] != 0 {
			score.midgame = max(score.midgame, e.kingCoverBonus(our, homeKing[our] + 2))
		}
		if p.castles & castleQueenside[our] != 0 {
			score.midgame = max(score.midgame, e.kingCoverBonus(our, homeKing[our] - 2))
		}
	}

	score.endgame = e.kingPawnProximity(our, square)

	return score
}

func (e *Evaluation) kingCoverBonus(our int, square int) (bonus int) {
	bonus = onePawn + onePawn / 3

	// Get pawns adjacent to and in front of the king.
	row, col := coordinate(square)
	area := maskRank[row] | maskPassed[our][square]
	cover := e.position.outposts[pawn(our)] & area
	storm := e.position.outposts[pawn(our^1)] & area

	// For each of the cover files find the closest friendly pawn. The penalty
	// is carried if the pawn is missing or is too far from the king.
	from, to := max(B1, col) - 1, min(G1, col) + 1
	for c := from; c <= to; c++ {

		// Friendly pawns protecting the kings.
		closest := 0
		if pawns := (cover & maskFile[c]); pawns.anyʔ() {
			closest = rank(our, pawns.closest(our))
		}
		bonus -= penaltyCover[closest]

		// Enemy pawns facing the king.
		if pawns := (storm & maskFile[c]); pawns.anyʔ() {
			farthest := rank(our, pawns.farthest(our^1))
			if closest == 0 { // No opposing friendly pawn.
				bonus -= penaltyStorm[farthest]
			} else if farthest == closest + 1 {
				bonus -= penaltyStormBlocked[farthest]
			} else {
				bonus -= penaltyStormUnblocked[farthest]
			}
		}
	}

	return bonus
}

// Calculates endgame penalty to encourage a king stay closer to friendly pawns.
func (e *Evaluation) kingPawnProximity(our int, square int) (penalty int) {
	pawns := e.position.outposts[pawn(our)]

	if pawns.anyʔ() && (pawns & e.attacks[king(our)]).noneʔ() {
		proximity := 8

		for bm := pawns; bm.anyʔ(); bm = bm.pop() {
			proximity = min(proximity, distance[square][bm.first()])
		}

		penalty = -kingByPawn.endgame * (proximity - 1)
	}

	return penalty
}

