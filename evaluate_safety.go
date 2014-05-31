// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeSafety() {
	var white, black [2]Score

	if Settings.Trace {
		defer func() {
			var his, her Score
			e.checkpoint(`+King`, Total{*his.add(white[0]).add(white[1]), *her.add(black[0]).add(black[1])})
			e.checkpoint(`-Cover`, Total{white[0], black[0]})
			e.checkpoint(`-Safety`, Total{white[1], black[1]})
		}()
	}

	if e.strongEnough(Black) {
		white[0] = e.kingCover(White)
		white[1] = e.kingSafety(White)
		e.score.add(white[0]).add(white[1])
	}
	if e.strongEnough(White) {
		black[0] = e.kingCover(Black)
		black[1] = e.kingSafety(Black)
		e.score.subtract(black[0]).subtract(black[1])
	}
}

func (e *Evaluation) kingSafety(color int) (score Score) {
	p := e.position

	if e.safety[color].threats > 0 {
		square := p.king[color]
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

		threatIndex := Min(12, e.safety[color].attackers * e.safety[color].threats / 3) + (e.safety[color].attacks + weak.count()) * 2
		safetyIndex = Min(63, safetyIndex + threatIndex)

		score.midgame -= kingSafety[safetyIndex]
		score.endgame -= bonusKing[1][Flip(color, square)]
	}
	return
}

func (e *Evaluation) kingCover(color int) (penalty Score) {
	p := e.position
	kings, pawns := p.outposts[king(color)], p.outposts[pawn(color)]

	// Pass if a) the king is missing, b) the king is on the initial square
	// or c) the opposite side doesn't have a queen with one major piece.
	if kings == 0 || kings == bit[homeKing[color]] || !e.strongEnough(color^1) {
		return
	}

	// Calculate relative square for the king so we could treat black king
	// as white. Don't bother with the cover if the king is too far.
	square := Flip(color^1, p.king[color])
	if square > H3 {
		return
	}
	row, col := Coordinate(square)
	from, to := Max(0, col-1), Min(7, col+1)

	// For each of the cover columns find the closest same color pawn. The
	// penalty is carried if the pawn is missing or is too far from the king
	// (more than one row apart).
	for column := from; column <= to; column++ {
		if cover := (pawns & maskFile[column]); cover != 0 {
			closest := Flip(color^1, cover.first()) // Make it relative.
			if distance := Abs(Row(closest) - row); distance > 1 {
				penalty.midgame += distance * -coverDistance.midgame
			}
		} else {
			penalty.midgame += -coverMissing.midgame
		}
	}
	// Log("penalty[%s] => %d\n", C(color), penalty)
	return
}

