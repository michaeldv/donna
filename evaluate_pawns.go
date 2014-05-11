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
	hisPawns := e.position.outposts[pawn(color)]
	herPawns := e.position.outposts[pawn(color^1)]

	pawns := hisPawns
	for pawns != 0 {
		square := pawns.pop()
		row, col := Coordinate(square)

		// Penalty if the pawn is isolated, i.e. has no friendly pawns
		// on adjacent files. The penalty goes up if isolated pawn is
		// exposed on semi-open file.
		isolated := (maskIsolated[col] & hisPawns == 0)
		exposed := (maskInFront[color][square] & herPawns == 0)
		if isolated {
			if !exposed {
				score.subtract(penaltyIsolatedPawn[col])
			} else {
				score.subtract(penaltyWeakIsolatedPawn[col])
			}
		}

		// Penalty if the pawn is doubled, i.e. there is another friendly
		// pawn in front of us. The penalty goes up if doubled pawns are
		// isolated.
		doubled := (maskInFront[color][square] & hisPawns != 0)
		if doubled {
			score.subtract(penaltyDoubledPawn[col])
		}

		// Bonus if the pawn is supported by friendly pawn(s) on the same
		// or previous ranks.
		supported := (maskIsolated[col] & (maskRank[row] | maskRank[row].pushed(color^1)) & hisPawns != 0)
		if supported {
			flip := Flip(color, square)
			score.add(Score{bonusSupportedPawn[flip], bonusSupportedPawn[flip]})
		}

		// The pawn is passed if a) there are no enemy pawns in the same
		// and adjacent columns; and b) there are no same color pawns in
		// front of us.
		passed := (maskPassed[color][square] & herPawns == 0 && !doubled)
		if passed {
			flip := Flip(color, square)
			score.midgame += bonusPassedPawn[0][flip]
			score.endgame += bonusPassedPawn[1][flip]
		}

		// Penalty if the pawn is backward.
		if (!passed && !supported && !isolated) {

			// Backward pawn should not be attacking enemy pawns.
			if pawnMoves[color][square] & herPawns == 0 {

				// Backward pawn should not have friendly pawns behind.
				if maskPassed[color^1][square] & maskIsolated[col] & hisPawns == 0 {

					// Backward pawn should face enemy pawns on the next two ranks
					// preventing its advance.
					enemy := pawnMoves[color][square].pushed(color)
					if (enemy | enemy.pushed(color)) & herPawns != 0 {
						if !exposed {
							score.subtract(penaltyBackwardPawn[col])
						} else {
							score.subtract(penaltyWeakBackwardPawn[col])
						}
					}
				}
			}
		}

		// TODO: Bonus if the pawn has good chance to become a passed pawn.
	}

	// Penalty for blocked pawns.
	blocked := (hisPawns.pushed(color) & e.position.board).count()
	score.subtract(pawnBlocked.times(blocked))

	return
}

