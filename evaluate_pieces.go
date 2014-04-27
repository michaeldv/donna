// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

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
	var white, black [4]Score

	maskSafe := ^e.position.pawnAttacks(Black) 	// Squares not attacked by Black pawns.
	maskEnemy := e.position.outposts[Black]  	// Squares occupied by Black pieces.
	white[0] = e.knights(White, maskSafe, maskEnemy)
	white[1] = e.bishops(White, maskSafe, maskEnemy)
	white[2] = e.rooks(White, maskSafe, maskEnemy)
	white[3] = e.queens(White, maskSafe, maskEnemy)

	maskSafe = ^e.position.pawnAttacks(White) 	// Squares not attacked by White pawns.
	maskEnemy = e.position.outposts[White]  	// Squares occupied by White pieces.
	black[0] = e.knights(Black, maskSafe, maskEnemy)
	black[1] = e.bishops(Black, maskSafe, maskEnemy)
	black[2] = e.rooks(Black, maskSafe, maskEnemy)
	black[3] = e.queens(Black, maskSafe, maskEnemy)

	e.midgame += white[0].midgame + white[1].midgame + white[2].midgame + white[3].midgame -
	             black[0].midgame - black[1].midgame - black[2].midgame - black[3].midgame

     	e.endgame += white[0].endgame + white[1].endgame + white[2].endgame + white[3].endgame -
     	             black[0].endgame - black[1].endgame - black[2].endgame - black[3].endgame
}

func (e *Evaluator) knights(color int, maskSafe, maskEnemy Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[knight(color)]

	for outposts != 0 {
		square := outposts.pop()
		targets := p.targets(square)

		// Attacks.
		attacks := (targets & maskEnemy).count()
		score.midgame += attacks * attackForce.midgame
		score.endgame += attacks * attackForce.endgame

		// Mobility
		mobility := mobilityKnight[(targets & maskSafe).count()]
		score.midgame += mobility.midgame
		score.midgame += mobility.midgame

		// Placement.
		square = Flip(color, square)
		score.midgame += bonusKnight[0][square]
		score.endgame += bonusKnight[1][square]

		// Penalty a knight is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.midgame -= penaltyPawnThreat[2].midgame
			score.endgame -= penaltyPawnThreat[2].endgame
		}

		// Bonus if a knight is behind friendly pawn.
		if RelRow(color, square) < 4 && p.outposts[pawn(color)].isSet(square + eight[color]) {
			score.midgame += behindPawn.midgame
			score.endgame += behindPawn.endgame
		}
	}
	return
}

func (e *Evaluator) bishops(color int, maskSafe, maskEnemy Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[bishop(color)]

	for outposts != 0 {
		square := outposts.pop()
		targets := p.xrayTargets(square)

		// Attacks.
		attacks := (targets & maskEnemy).count()
		score.midgame += attacks * attackForce.midgame
		score.endgame += attacks * attackForce.endgame

		// Mobility
		mobility := mobilityBishop[(targets & maskSafe).count()]
		score.midgame += mobility.midgame
		score.midgame += mobility.midgame

		// Placement.
		square = Flip(color, square)
		score.midgame += bonusBishop[0][square]
		score.endgame += bonusBishop[1][square]

		// Penalty for light/dark square bishop and matching pawns.
		if count := (SameColor(square) & p.outposts[pawn(color)]).count(); count > 0 {
			score.midgame -= bishopPawns.midgame
			score.midgame -= bishopPawns.midgame
		}

		// Penalty a bishop is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.midgame -= penaltyPawnThreat[3].midgame
			score.endgame -= penaltyPawnThreat[3].endgame
		}

		// Bonus if a bishop is behind friendly pawn.
		if RelRow(color, square) < 4 && p.outposts[pawn(color)].isSet(square + eight[color]) {
			score.midgame += behindPawn.midgame
			score.endgame += behindPawn.endgame
		}
	}

	// Bonus for the pair of bishops.
	if bishops := p.count[bishop(color)]; bishops >= 2 {
		e.midgame += bishopPair.midgame
		e.endgame += bishopPair.endgame
	}
	return
}


func (e *Evaluator) rooks(color int, maskSafe, maskEnemy Bitmask) (score Score) {
	p := e.position
	hisPawns := p.outposts[pawn(color)]
	herPawns := p.outposts[pawn(color^1)]
	outposts := p.outposts[rook(color)]

	// Bonus if rook is on 7th rank and enemy's king trapped on 8th.
	if count := (outposts & mask7th[color]).count(); count > 0 && p.outposts[king(color^1)] & mask8th[color] != 0 {
		score.midgame += count * rookOn7th.midgame
		score.endgame += count * rookOn7th.endgame
	}
	for outposts != 0 {
		square := outposts.pop()
		targets := p.xrayTargets(square)

		// Attacks.
		attacks := (targets & maskEnemy).count()
		score.midgame += attacks * attackForce.midgame
		score.endgame += attacks * attackForce.endgame

		// Mobility
		mobility := mobilityRook[(targets & maskSafe).count()]
		score.midgame += mobility.midgame
		score.midgame += mobility.midgame

		// Placement.
		square = Flip(color, square)
		score.midgame += bonusRook[0][square]
		score.endgame += bonusRook[1][square]

		// Penalty a rook is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.midgame -= penaltyPawnThreat[4].midgame
			score.endgame -= penaltyPawnThreat[4].endgame
		}

		// Bonus if rook is attacking enemy's pawns.
		if count := (targets & p.outposts[pawn(color^1)]).count(); count > 0 {
			score.midgame += count * rookOnPawn.midgame
			score.endgame += count * rookOnPawn.endgame
		}

		// Bonuses if rook is on open or semi-open file.
		column := Col(square)
		if hisPawns & maskFile[column] == 0 {
			if herPawns & maskFile[column] == 0 {
				score.midgame += rookOnOpen.midgame
				score.endgame += rookOnOpen.endgame
			} else {
				score.midgame += rookOnSemiOpen.midgame
				score.endgame += rookOnSemiOpen.endgame
			}
		}
	}
	return
}

func (e *Evaluator) queens(color int, maskSafe, maskEnemy Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[queen(color)]

	// Bonus if queen is on 7th rank and enemy's king trapped on 8th.
	if count := (outposts & mask7th[color]).count(); count > 0 && p.outposts[king(color^1)] & mask8th[color] != 0 {
		score.midgame += count * queenOn7th.midgame
		score.endgame += count * queenOn7th.endgame
	}
	for outposts != 0 {
		square := outposts.pop()
		targets := p.targets(square)

		// Attacks.
		attacks := (targets & maskEnemy).count()
		score.midgame += attacks * attackForce.midgame
		score.endgame += attacks * attackForce.endgame

		// Mobility
		mobility := mobilityQueen[Max(15, (targets & maskSafe).count())]
		score.midgame += mobility.midgame
		score.midgame += mobility.midgame

		// Placement.
		square = Flip(color, square)
		score.midgame += bonusQueen[0][square]
		score.endgame += bonusQueen[1][square]

		// Penalty if queen is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.midgame -= penaltyPawnThreat[5].midgame
			score.endgame -= penaltyPawnThreat[5].endgame
		}

		// Bonus if queen is out and attacking enemy's pawns.
		if count := (targets & p.outposts[pawn(color^1)]).count(); count > 0 && RelRow(color, square) > 3 {
			score.midgame += count * queenOnPawn.midgame
			score.endgame += count * queenOnPawn.endgame
		}
	}
	return
}
