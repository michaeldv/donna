// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeSafety() {
	var cover, safety Total
	oppositeBishops := e.oppositeBishops()

	if engine.trace {
		defer func() {
			var his, her Score
			e.checkpoint(`+King`, Total{*his.add(cover.white).add(safety.white), *her.add(cover.black).add(safety.black)})
			e.checkpoint(`-Cover`, cover)
			e.checkpoint(`-Safety`, safety)
		}()
	}

	// If any of the pawns or a king have moved then recalculate cover score.
	if e.position.king[White] != e.pawns.king[White] {
		e.pawns.cover[White] = e.kingCover(White)
		e.pawns.king[White] = e.position.king[White]
	}
	if e.position.king[Black] != e.pawns.king[Black] {
		e.pawns.cover[Black] = e.kingCover(Black)
		e.pawns.king[Black] = e.position.king[Black]
	}

	// Fetch king cover score from the pawn cache.
	cover.white.add(e.pawns.cover[White])
	cover.black.add(e.pawns.cover[Black])

	// Compute king's safety for both sides.
	if e.safety[White].threats > 0 {
		safety.white = e.kingSafety(White)
	}
	if e.safety[Black].threats > 0 {
		safety.black = e.kingSafety(Black)
	}

	// Less safe with opposite bishops.
	if oppositeBishops && e.isKingUnsafe(White) && safety.white.midgame < -onePawn / 10 * 8 {
		safety.white.midgame -= bishopDanger.midgame
	}
	if oppositeBishops && e.isKingUnsafe(Black) && safety.black.midgame < -onePawn / 10 * 8 {
		safety.black.midgame -= bishopDanger.midgame
	}

	e.score.add(cover.white).add(safety.white).sub(cover.black).sub(safety.black)
}

func (e *Evaluation) kingSafety(color uint8) (score Score) {
	p, rival := e.position, color^1
	safetyIndex, square := 0, int(p.king[color])

	// Find squares around the king that are being attacked by the
	// enemy and defended by our king only.
	defended := e.attacks[pawn(color)] | e.attacks[knight(color)] |
	            e.attacks[bishop(color)] | e.attacks[rook(color)] |
	            e.attacks[queen(color)]
	weak := e.attacks[king(color)] & e.attacks[rival] & ^defended

	// Find possible queen checks on weak squares around the king.
	// We only consider squares where the queen is protected and
	// can't be captured by the king.
	protected := e.attacks[pawn(rival)] | e.attacks[knight(rival)] |
		     e.attacks[bishop(rival)] | e.attacks[rook(rival)] |
		     e.attacks[king(rival)]
	checks := weak & e.attacks[queen(rival)] & protected & ^p.outposts[rival]
	if checks.any() {
		safetyIndex += queenCheck * checks.count()
	}

	// Out of all squares available for enemy pieces select the ones
	// that are not under our attack.
	safe := ^(e.attacks[color] | p.outposts[rival])

	// Are there any safe squares from where enemy Knight could give
	// us a check?
	if checks := knightMoves[square] & safe & e.attacks[knight(rival)]; checks.any() {
		safetyIndex += checks.count()
	}

	// Are there any safe squares from where enemy Bishop could give us a check?
	safeBishopMoves := p.bishopMoves(square) & safe
	if checks := safeBishopMoves & e.attacks[bishop(rival)]; checks.any() {
		safetyIndex += checks.count()
	}

	// Are there any safe squares from where enemy Rook could give us a check?
	safeRookMoves := p.rookMoves(square) & safe
	if checks := safeRookMoves & e.attacks[rook(rival)]; checks.any() {
		safetyIndex += checks.count()
	}

	// Are there any safe squares from where enemy Queen could give us a check?
	if checks := (safeBishopMoves | safeRookMoves) & e.attacks[queen(rival)]; checks.any() {
		safetyIndex += queenCheck / 2 * checks.count()
	}

	threatIndex := min(16, e.safety[color].attackers * e.safety[color].threats / 2) +
			(e.safety[color].attacks + weak.count()) * 3 +
			rank(color, square) - e.pawns.cover[color].midgame / 16
	safetyIndex = min(63, max(0, safetyIndex + threatIndex))

	score.midgame -= kingSafety[safetyIndex]

	return
}

func (e *Evaluation) kingCover(color uint8) (bonus Score) {
	p, square := e.position, int(e.position.king[color])

	// Don't bother with the cover if the king is too far out.
	if rank(color, square) <= A3H3 {
		// If we still have castle rights encourage castle pawns to stay intact
		// by scoring least safe castle.
		bonus.midgame = e.kingCoverBonus(color, square)
		if p.castles & castleKingside[color] != 0 {
			bonus.midgame = max(bonus.midgame, e.kingCoverBonus(color, homeKing[color] + 2))
		}
		if p.castles & castleQueenside[color] != 0 {
			bonus.midgame = max(bonus.midgame, e.kingCoverBonus(color, homeKing[color] - 2))
		}
	}

	bonus.endgame = e.kingPawnProximity(color, square)

	return
}

func (e *Evaluation) kingCoverBonus(color uint8, square int) (bonus int) {
	bonus = onePawn + onePawn / 3

	// Get pawns adjacent to and in front of the king.
	row, col := coordinate(square)
	area := maskRank[row] | maskPassed[color][square]
	cover := e.position.outposts[pawn(color)] & area
	storm := e.position.outposts[pawn(color^1)] & area

	// For each of the cover files find the closest friendly pawn. The penalty
	// is carried if the pawn is missing or is too far from the king.
	from, to := max(B1, col) - 1, min(G1, col) + 1
	for c := from; c <= to; c++ {

		// Friendly pawns protecting the kings.
		closest := 0
		if pawns := (cover & maskFile[c]); pawns.any() {
			closest = rank(color, pawns.closest(color))
		}
		bonus -= penaltyCover[closest]

		// Enemy pawns facing the king.
		if pawns := (storm & maskFile[c]); pawns.any() {
			farthest := rank(color, pawns.farthest(color^1))
			if closest == 0 { // No opposing friendly pawn.
				bonus -= penaltyStorm[farthest]
			} else if farthest == closest + 1 {
				bonus -= penaltyStormBlocked[farthest]
			} else {
				bonus -= penaltyStormUnblocked[farthest]
			}
		}
	}

	return
}

// Calculates endgame penalty to encourage a king stay closer to friendly pawns.
func (e *Evaluation) kingPawnProximity(color uint8, square int) (penalty int) {
	if pawns := e.position.outposts[pawn(color)]; pawns.any() && (pawns & e.attacks[king(color)]).empty() {
		proximity := 8

		for pawns.any() {
			proximity = min(proximity, distance[square][pawns.pop()])
		}

		penalty = -kingByPawn.endgame * (proximity - 1)
	}

	return
}

