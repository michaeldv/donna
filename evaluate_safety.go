// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeSafety() {
	var cover, safety Total
	color := e.position.color
	oppositeBishops := e.oppositeBishops()

	// Any pawn move invalidates king's square in the pawns hash so that we
	// could detect it here.
	whiteKingMoved := e.position.king[White] != e.pawns.king[White]
	blackKingMoved := e.position.king[Black] != e.pawns.king[Black]

	if engine.trace {
		defer func() {
			var his, her Score
			e.checkpoint(`+King`, Total{*his.add(cover.white).add(safety.white), *her.add(cover.black).add(safety.black)})
			e.checkpoint(`-Cover`, cover)
			e.checkpoint(`-Safety`, safety)
		}()
	}

	// If the king has moved then recalculate king/pawn proximity and update
	// cover score and king square in the pawn cache.
	if whiteKingMoved {
		e.pawns.cover[White] = e.kingCover(White)
		e.pawns.cover[White].endgame += e.kingPawnProximity(White)
		e.pawns.king[White] = e.position.king[White]
	}
	if blackKingMoved {
		e.pawns.cover[Black] = e.kingCover(Black)
		e.pawns.cover[Black].endgame += e.kingPawnProximity(Black)
		e.pawns.king[Black] = e.position.king[Black]
	}

	// Fetch king cover score from the pawn cache.
	cover.white.add(e.pawns.cover[White])
	cover.black.add(e.pawns.cover[Black])

	// Compute king's safety for both sides.
	safety.white = e.kingSafety(White)
	if oppositeBishops && e.isKingUnsafe(White) && safety.white.midgame < -onePawn / 10 * 8 {
		safety.white.midgame -= bishopDanger.midgame
	}
	safety.black = e.kingSafety(Black)
	if oppositeBishops && e.isKingUnsafe(Black) && safety.black.midgame < -onePawn / 10 * 8 {
		safety.black.midgame -= bishopDanger.midgame
	}

	// Apply king safety weights, then adjust the final score.
	if color == White {
		cover.white.apply(weightOurKingSafety)
		safety.white.apply(weightOurKingSafety)
		cover.black.apply(weightTheirKingSafety)
		safety.black.apply(weightTheirKingSafety)
	} else {
		cover.white.apply(weightTheirKingSafety)
		safety.white.apply(weightTheirKingSafety)
		cover.black.apply(weightOurKingSafety)
		safety.black.apply(weightOurKingSafety)
	}

	e.score.add(cover.white).add(safety.white).sub(cover.black).sub(safety.black)
}

func (e *Evaluation) kingSafety(color uint8) (score Score) {
	p := e.position

	if e.safety[color].threats > 0 {
		square := int(p.king[color])
		safetyIndex := 0

		// Find squares around the king that are being attacked by the
		// enemy and defended by our king only.
		defended := e.attacks[pawn(color)] | e.attacks[knight(color)] |
		            e.attacks[bishop(color)] | e.attacks[rook(color)] |
		            e.attacks[queen(color)]
		weak := e.attacks[king(color)] & e.attacks[color^1] & ^defended

		// Find possible queen checks on weak squares around the king.
		// We only consider squares where the queen is protected and
		// can't be captured by the king.
		protected := e.attacks[pawn(color^1)] | e.attacks[knight(color^1)] |
		             e.attacks[bishop(color^1)] | e.attacks[rook(color^1)] |
		             e.attacks[king(color^1)]
		checks := weak & e.attacks[queen(color^1)] & protected & ^p.outposts[color^1]
		if checks != 0 {
			safetyIndex += bonusCloseCheck[Queen/2] * checks.count()
		}

		// Find possible rook checks within king's home zone. Unlike
		// queen we must only consider squares where the rook actually
		// gives a check.
		protected = e.attacks[pawn(color^1)] | e.attacks[knight(color^1)] |
		            e.attacks[bishop(color^1)] | e.attacks[queen(color^1)] |
		            e.attacks[king(color^1)]
		checks = weak & e.attacks[rook(color^1)] & protected & ^p.outposts[color^1]
		checks &= rookMagicMoves[square][0]
		if checks != 0 {
			safetyIndex += bonusCloseCheck[Rook/2] * checks.count()
		}

		// Double safety index if the enemy has right to move.
		if p.color == color^1 {
			safetyIndex *= 2
		}

		// Out of all squares available for enemy pieces select the ones
		// that are not under our attack.
		safe := ^(e.attacks[color] | p.outposts[color^1])

		// Are there any safe squares from where enemy Knight could give
		// us a check?
		if checks := knightMoves[square] & safe & e.attacks[knight(color^1)]; checks != 0 {
			safetyIndex += bonusDistanceCheck[Knight/2] * checks.count()
		}

		// Are there any safe squares from where enemy Bishop could give
		// us a check?
		safeBishopMoves := p.bishopMoves(square) & safe
		if checks := safeBishopMoves & e.attacks[bishop(color^1)]; checks != 0 {
			safetyIndex += bonusDistanceCheck[Bishop/2] * checks.count()
		}

		// Are there any safe squares from where enemy Rook could give
		// us a check?
		safeRookMoves := p.rookMoves(square) & safe
		if checks := safeRookMoves & e.attacks[rook(color^1)]; checks != 0 {
			safetyIndex += bonusDistanceCheck[Rook/2] * checks.count()
		}

		// Are there any safe squares from where enemy Queen could give
		// us a check?
		if checks := (safeBishopMoves | safeRookMoves) & e.attacks[queen(color^1)]; checks != 0 {
			safetyIndex += bonusDistanceCheck[Queen/2] * checks.count()
		}

		threatIndex := min(12, e.safety[color].attackers * e.safety[color].threats / 3) + (e.safety[color].attacks + weak.count()) * 2
		safetyIndex = min(63, safetyIndex + threatIndex)

		score.midgame -= kingSafety[safetyIndex]
	}

	return
}

func (e *Evaluation) kingCover(color uint8) (bonus Score) {
	p, square := e.position, int(e.position.king[color])

	// Calculate relative square for the king so we could treat black king
	// as white. Don't bother with the cover if the king is too far.
	flipped := flip(color^1, square)
	if flipped > H3 {
		return
	}

	// If we still have castle rights encourage castle pawns to stay intact
	// by scoring least safe castle.
	bonus.midgame = e.kingCoverBonus(color, square, flipped)
	if p.castles & castleKingside[color] != 0 {
		bonus.midgame = max(bonus.midgame, e.kingCoverBonus(color, homeKing[color] + 2, G1))
	}
	if p.castles & castleQueenside[color] != 0 {
		bonus.midgame = max(bonus.midgame, e.kingCoverBonus(color, homeKing[color] - 2, C1))
	}

	return
}

func (e *Evaluation) kingCoverBonus(color uint8, square, flipped int) (bonus int) {
	r, c := coordinate(flipped)
	from, to := max(0, c - 1), min(7, c + 1)
	bonus = onePawn + onePawn / 3

	// Get friendly pawns adjacent and in front of the king.
	adjacent := maskIsolated[c] & maskRank[row(square)]
	pawns := e.position.outposts[pawn(color)] & (adjacent | maskPassed[color][square])

	// For each of the cover files find the closest friendly pawn. The penalty
	// is carried if the pawn is missing or is too far from the king (more than
	// one rank apart).
	for column := from; column <= to; column++ {
		if cover := (pawns & maskFile[column]); cover != 0 {
			closest := rank(color, cover.closest(color))
			bonus -= penaltyCover[closest - r]
		} else {
			bonus -= coverMissing.midgame
		}
	}

	// Log("penalty[%s] => %+v\n", C(color), penalty)
	return
}

// Calculates endgame penalty to encourage a king stay closer to friendly pawns.
func (e *Evaluation) kingPawnProximity(color uint8) (penalty int) {
	if pawns := e.position.outposts[pawn(color)]; pawns != 0 && pawns & e.attacks[king(color)] == 0 {
		proximity, king := 8, e.position.king[color]

		for pawns != 0 {
			proximity = min(proximity, distance[king][pawns.pop()])
		}
		penalty = -kingByPawn.endgame * (proximity - 1)
	}

	return
}

