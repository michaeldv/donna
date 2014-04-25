// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

var mobilityKnight = [9]Score{
	{-15, -12}, {-10, -8}, {-4, -4}, {1, 0}, {6, 3}, {11, 7}, {15, 10}, {17, 12}, {18, 12},
}

var mobilityBishop = [16]Score{
	{-10, -12}, {-4, -6}, { 2,  0}, { 8,  6}, {14, 12}, {20, 18}, {25, 23}, {28, 26},
	{ 31,  28}, {32, 30}, {33, 31}, {34, 32}, {35, 33}, {36, 34}, {36, 34}, {36, 34},
}

var mobilityRook = [16]Score{
	{-8, -15}, {-5, -7}, {-3,  0}, { 0,  6}, { 2, 13}, { 5, 20}, { 7, 27}, { 9, 34},
	{ 10, 41}, {11, 46}, {12, 49}, {13, 51}, {14, 52}, {14, 52}, {15, 53}, {15, 53},
}

var mobilityQueen = [32]Score{
	{-5, -8}, {-4, -6}, {-2, -3}, {-1, -1}, {0,  1}, {1,  4}, {2,  6}, {3,  9},
	{ 5, 11}, { 6, 13}, { 7, 15}, { 7, 15}, {8, 16}, {9, 16}, {9, 16}, {9, 16},
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
	outposts := e.position.outposts[knight(color)]
	for outposts != 0 {
		square := outposts.pop()
		targets := e.position.targets(square)

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
	}
	return
}

func (e *Evaluator) bishops(color int, maskSafe, maskEnemy Bitmask) (score Score) {
	outposts := e.position.outposts[bishop(color)]
	for outposts != 0 {
		square := outposts.pop()
		targets := e.position.targets(square)

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
	}
	//
	// Bonus for the pair of bishops.
	//
	if bishops := e.position.count[bishop(color)]; bishops >= 2 {
		e.midgame += bishopPair.midgame
		e.endgame += bishopPair.endgame
	}
	return
}


func (e *Evaluator) rooks(color int, maskSafe, maskEnemy Bitmask) (score Score) {
	hisPawns := e.position.outposts[pawn(color)]
	herPawns := e.position.outposts[pawn(color^1)]
	outposts := e.position.outposts[rook(color)]
	//
	// Bonus if rooks are on 7th rank.
	//
	if count := (outposts & mask7th[color]).count(); count > 0 {
		score.midgame += count * rookOn7th.midgame
		score.endgame += count * rookOn7th.endgame
	}
	for outposts != 0 {
		square := outposts.pop()
		targets := e.position.targets(square)

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
		//
		// Bonuses if rooks are on open or semi-open files.
		//
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
	outposts := e.position.outposts[queen(color)]
	for outposts != 0 {
		square := outposts.pop()
		targets := e.position.targets(square)

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
	}
	return
}
