// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

type PawnCacheEntry struct {
	hash    uint64
	midgame int
	endgame int
}

var pawnCache [8192]PawnCacheEntry

func (e *Evaluator) analyzePawnStructure() {
	hashPawn := e.position.hashPawn
	index := hashPawn % uint64(len(pawnCache))
	entry := &pawnCache[index]

	if entry.hash != hashPawn {
		whiteExtra, blackExtra := e.pawnsScore(White), e.pawnsScore(Black)
		entry.midgame = whiteExtra.midgame - blackExtra.midgame
		entry.endgame = whiteExtra.endgame - blackExtra.endgame
		entry.hash = e.position.hashPawn
	}

	e.midgame += entry.midgame
	e.endgame += entry.endgame
}

// Calculates extra bonus and penalty based on pawn structure. Specifically,
// a bonus is awarded for passed pawns, and penalty applied for isolated and
// doubled pawns.
func (e *Evaluator) pawnsScore(color int) (extra Score) {
	hisPawns := e.position.outposts[pawn(color)]
	herPawns := e.position.outposts[pawn(color^1)]

	pawns := hisPawns
	for pawns != 0 {
		square := pawns.pop()
		column := Col(square)
		//
		// The pawn is passed if a) there are no enemy pawns in the
		// same and adjacent columns; and b) there is no same color
		// pawns in front of us.
		//
		if maskPassed[color][square] & herPawns == 0 && maskInFront[color][square] & hisPawns == 0 {
			square = Flip(color, square)
			extra.midgame += bonusPassedPawn[0][square]
			extra.endgame += bonusPassedPawn[1][square]
		}
		//
		// Check if the pawn is isolated, i.e. has no pawns of the
		// same color on either sides.
		//
		if maskIsolated[column] & hisPawns == 0 {
			extra.midgame += penaltyIsolatedPawn[0][column]
			extra.endgame += penaltyIsolatedPawn[1][column]
		}
	}
	//
	// Penalties for doubled pawns.
	//
	for col := 0; col <= 7; col++ {
		if doubled := (maskFile[col] & hisPawns).count(); doubled > 1 {
			extra.midgame += (doubled - 1) * penaltyDoubledPawn[0][col]
			extra.endgame += (doubled - 1) * penaltyDoubledPawn[1][col]
		}
	}
	return
}

