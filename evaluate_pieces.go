// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluator) analyzePieces() {
	var white, black [4]Score

	maskSafe := ^e.position.pawnAttacks(Black) // Squares not attacked by Black pawns.
	white[0] = e.knights(White, maskSafe)
	white[1] = e.bishops(White, maskSafe)
	white[2] = e.rooks(White, maskSafe)
	white[3] = e.queens(White, maskSafe)

	maskSafe = ^e.position.pawnAttacks(White) // Squares not attacked by White pawns.
	black[0] = e.knights(Black, maskSafe)
	black[1] = e.bishops(Black, maskSafe)
	black[2] = e.rooks(Black, maskSafe)
	black[3] = e.queens(Black, maskSafe)


	e.midgame += white[0].midgame + white[1].midgame + white[2].midgame + white[3].midgame -
	             black[0].midgame - black[1].midgame - black[2].midgame - black[3].midgame

     	e.endgame += white[0].endgame + white[1].endgame + white[2].endgame + white[3].endgame -
     	             black[0].endgame - black[1].endgame - black[2].endgame - black[3].endgame
}

func (e *Evaluator) knights(color int, maskSafe Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[knight(color)]

	for outposts != 0 {
		square := outposts.pop()
		targets := p.targets(square)

		// Bonus for knight's mobility.
		score.add(mobilityKnight[(targets & maskSafe).count()])

		// Penalty if knight is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Knight/2])
		}

		// Bonus if knight is behind friendly pawn.
		if RelRow(color, square) < 4 && p.outposts[pawn(color)].isSet(square + eight[color]) {
			score.add(behindPawn)
		}

		// Track if knight attacks squares around enemy's king.
		if targets & p.kingAttacks(color^1) != 0 {
			e.attacks[color]++
			e.threats[color] += bonusKingThreat[Knight/2]
		}

		// Bonus for knight's board position.
		flip := Flip(color, square)
		score.midgame += bonusKnight[0][flip]
		score.endgame += bonusKnight[1][flip]

		// Extra bonus if knight is in the center. Increase the extra
		// bonus if the knight is supported by a pawn and can't be
		// exchanged.
		if extra := extraKnight[flip]; extra > 0 {
			if p.pawnAttacks(color).isSet(square) {
				if p.count[knight(color^1)] == 0 {
					extra *= 2 // No knights to exchange.
				}
				extra += extra / 2 // Supported by a pawn.
			}
			score.increment(extra)
		}
	}
	return
}

func (e *Evaluator) bishops(color int, maskSafe Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[bishop(color)]

	for outposts != 0 {
		square := outposts.pop()
		targets := p.xrayTargets(square)

		// Bonus for bishop's mobility
		score.add(mobilityBishop[(targets & maskSafe).count()])

		// Penalty for light/dark square bishop and matching pawns.
		if count := (Same(square) & p.outposts[pawn(color)]).count(); count > 0 {
			score.subtract(bishopPawns)
		}

		// Penalty if bishop is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Bishop/2])
		}

		// Bonus if bishop is behind friendly pawn.
		if RelRow(color, square) < 4 && p.outposts[pawn(color)].isSet(square + eight[color]) {
			score.add(behindPawn)
		}

		// Middle game penalty for boxed bishop.
		if color == White {
			if (square == C1 && p.pieces[D2].isPawn() && p.pieces[D3] != 0) ||
			   (square == F1 && p.pieces[E2].isPawn() && p.pieces[E3] != 0) {
				score.midgame -= bishopBoxed.midgame
			}
		} else {
			if (square == C8 && p.pieces[D7].isPawn() && p.pieces[D6] != 0) ||
			   (square == F8 && p.pieces[E7].isPawn() && p.pieces[E6] != 0) {
				score.midgame -= bishopBoxed.midgame
			}
		}

		// Track if bishop attacks squares around enemy's king.
		if targets & p.kingAttacks(color^1) != 0 {
			e.attacks[color]++
			e.threats[color] += bonusKingThreat[Bishop/2]
		}

		// Bonus for bishop's board position.
		flip := Flip(color, square)
		score.midgame += bonusBishop[0][flip]
		score.endgame += bonusBishop[1][flip]

		// Extra bonus if bishop is in the center. Increase the extra
		// bonus if the bishop is supported by a pawn and can't be
		// exchanged.
		if extra := extraBishop[flip]; extra > 0 {
			if p.pawnAttacks(color).isSet(square) {
				if p.count[bishop(color^1)] == 0 {
					extra *= 2 // No bishops to exchange.
				}
				extra += extra / 2 // Supported by a pawn.
			}
			score.increment(extra)
		}
	}

	// Bonus for the pair of bishops.
	if bishops := p.count[bishop(color)]; bishops >= 2 {
		score.add(bishopPair)
	}
	return
}


func (e *Evaluator) rooks(color int, maskSafe Bitmask) (score Score) {
	p := e.position
	hisPawns := p.outposts[pawn(color)]
	herPawns := p.outposts[pawn(color^1)]
	outposts := p.outposts[rook(color)]

	// Bonus if rook is on 7th rank and enemy's king trapped on 8th.
	if count := (outposts & mask7th[color]).count(); count > 0 && p.outposts[king(color^1)] & mask8th[color] != 0 {
		score.add(rookOn7th.multiply(count))
	}
	for outposts != 0 {
		square := outposts.pop()
		targets := p.xrayTargets(square)

		// Bonus for rook's mobility
		score.add(mobilityRook[(targets & maskSafe).count()])

		// Penalty if rook is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Rook/2])
		}

		// Bonus if rook is attacking enemy's pawns.
		if count := (targets & p.outposts[pawn(color^1)]).count(); count > 0 {
			score.add(rookOnPawn.multiply(count))
		}

		// Bonuses if rook is on open or semi-open file.
		column := Col(square)
		if hisPawns & maskFile[column] == 0 {
			if herPawns & maskFile[column] == 0 {
				score.add(rookOnOpen)
			} else {
				score.add(rookOnSemiOpen)
			}
		}

		// Middle game penalty if a rook is boxed. Extra penalty if castle
		// rights have been lost.
		if bit[square] & rookBoxA[color] != 0 && p.outposts[king(color)] & castleQueen[color] != 0 {
			score.midgame -= rookBoxed.midgame
			if p.castles & castleQueenside[color] == 0 {
				score.midgame -= rookBoxed.midgame
			}
		}
		if bit[square] & rookBoxH[color] != 0 && p.outposts[king(color)] & castleKing[color] != 0 {
			score.midgame -= rookBoxed.midgame
			if p.castles & castleKingside[color] == 0 {
				score.midgame -= rookBoxed.midgame
			}
		}

		// Track if rook attacks squares around enemy's king.
		if targets & p.kingAttacks(color^1) != 0 {
			e.attacks[color]++
			e.threats[color] += bonusKingThreat[Rook/2]
		}

		// Bonus for rook's board position.
		flip := Flip(color, square)
		score.midgame += bonusRook[0][flip]
		score.endgame += bonusRook[1][flip]
	}
	return
}

func (e *Evaluator) queens(color int, maskSafe Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[queen(color)]

	// Bonus if queen is on 7th rank and enemy's king trapped on 8th.
	if count := (outposts & mask7th[color]).count(); count > 0 && p.outposts[king(color^1)] & mask8th[color] != 0 {
		score.add(queenOn7th.multiply(count))
	}
	for outposts != 0 {
		square := outposts.pop()
		targets := p.targets(square)

		// Bonus for queen's mobility
		score.add(mobilityQueen[Min(15, (targets & maskSafe).count())])

		// Penalty if queen is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Queen/2])
		}

		// Bonus if queen is out and attacking enemy's pawns.
		if count := (targets & p.outposts[pawn(color^1)]).count(); count > 0 && RelRow(color, square) > 3 {
			score.add(queenOnPawn.multiply(count))
		}

		// Track if queen attacks squares around enemy's king.
		if targets & p.kingAttacks(color^1) != 0 {
			e.attacks[color]++
			e.threats[color] += bonusKingThreat[Queen/2]
		}

		// Bonus for queen's board position.
		flip := Flip(color, square)
		score.midgame += bonusQueen[0][flip]
		score.endgame += bonusQueen[1][flip]
	}
	return
}
