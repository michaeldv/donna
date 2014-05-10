// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

type PawnCacheEntry struct {
	hash   uint64
	score  Score
}

var pawnCache [8192]PawnCacheEntry

func (e *Evaluator) analyzePawns() {
	hashPawn := e.position.hashPawn
	index := hashPawn % uint64(len(pawnCache))
	entry := &pawnCache[index]

	if entry.hash != hashPawn {
		white, black := e.pawns(White), e.pawns(Black)
		entry.score.clear().add(white).subtract(black)
		entry.hash = e.position.hashPawn
	}

	e.score.add(entry.score)
}

// Calculates extra bonus and penalty based on pawn structure. Specifically,
// a bonus is awarded for passed pawns, and penalty applied for isolated and
// doubled pawns.
func (e *Evaluator) pawns(color int) (score Score) {
	var passed, isolated [8]bool

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
			flip := Flip(color, square)
			score.midgame += bonusPassedPawn[0][flip]
			score.endgame += bonusPassedPawn[1][flip]
			passed[column] = true
		}

		// Check if the pawn is isolated, i.e. has no pawns of the same
		// color on either sides.
		if maskIsolated[column] & hisPawns == 0 {
			score.add(penaltyIsolatedPawn[column])
			isolated[column] = true
		}
	}

	// Penalty for doubled pawns.
	for col := 0; col <= 7; col++ {
		if doubled := (maskFile[col] & hisPawns).count(); doubled > 1 {
			penalty := penaltyDoubledPawn[col]

			// Increate the penalty if doubled pawns are isolated
			// but not passed.
			if isolated[col] && !passed[col] {
				penalty = penalty.times(2)
			}
			score.add(penalty.times(doubled - 1))
		}
	}

	// Penalty for blocked pawns.
	blocked := (hisPawns.pushed(color) & e.position.board).count()
	score.subtract(pawnBlocked.times(blocked))

	return
}

