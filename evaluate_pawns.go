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

func (e *Evaluator) analyzePawns() {
	hashPawn := e.position.hashPawn
	index := hashPawn % uint64(len(pawnCache))
	entry := &pawnCache[index]

	if entry.hash != hashPawn {
		white, black := e.pawns(White), e.pawns(Black)
		entry.midgame = white.midgame - black.midgame
		entry.endgame = white.endgame - black.endgame
		entry.hash = e.position.hashPawn
	}

	e.midgame += entry.midgame
	e.endgame += entry.endgame
}

// Calculates extra bonus and penalty based on pawn structure. Specifically,
// a bonus is awarded for passed pawns, and penalty applied for isolated and
// doubled pawns.
func (e *Evaluator) pawns(color int) (score Score) {
	hisPawns := e.position.outposts[pawn(color)]
	herPawns := e.position.outposts[pawn(color^1)]

	pawns := hisPawns
	for pawns != 0 {
		square := pawns.pop()
		column := Col(square)

		// The pawn is passed if a) there are no enemy pawns in the same
		// and adjacent columns; and b) there are no same color pawns in
		// front of us.
		if maskPassed[color][square] & herPawns == 0 && maskInFront[color][square] & hisPawns == 0 {
			square := Flip(color, square)
			score.midgame += bonusPassedPawn[0][square]
			score.endgame += bonusPassedPawn[1][square]
		}

		// Check if the pawn is isolated, i.e. has no pawns of the same
		// color on either sides.
		if maskIsolated[column] & hisPawns == 0 {
			score.midgame += penaltyIsolatedPawn[column].midgame
			score.endgame += penaltyIsolatedPawn[column].midgame
		}

		// Bonus for pawn's position on the board.
		square = Flip(color, square)
		score.midgame += bonusPawn[0][square]
		score.endgame += bonusPawn[1][square]
	}

	// Penalty for doubled pawns.
	for col := 0; col <= 7; col++ {
		if doubled := (maskFile[col] & hisPawns).count(); doubled > 1 {
			score.midgame += (doubled - 1) * penaltyDoubledPawn[col].midgame
			score.endgame += (doubled - 1) * penaltyDoubledPawn[col].endgame
		}
	}

	// Penalty for blocked pawns.
	blocked := (Push(color, hisPawns) & (e.position.outposts[White] | e.position.outposts[Black])).count()
	score.midgame -= blocked * pawnBlocked.midgame
	score.endgame -= blocked * pawnBlocked.endgame

	return
}

