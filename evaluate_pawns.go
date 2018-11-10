// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

type PawnEntry struct {
	id       uint64 	// Pawn hash key.
	score    Score 		// Static score for the given pawn structure.
	king     [2]int 	// King square for both sides.
	cover    [2]Score 	// King cover penalties for both sides.
	passers  [2]Bitmask 	// Passed pawn bitmasks for both sides.
}

type PawnCache [8192*2]PawnEntry

func (e *Evaluation) analyzePawns() {
	key := e.position.pid

	// Since pawn hash is fairly small we can use much faster 32-bit index.
	index := uint32(key) % uint32(len(game.pawnCache))
	e.pawns = &game.pawnCache[index]

	// Bypass pawns cache if evaluation tracing is enabled.
	if e.pawns.id != key || engine.traceʔ {
		white, black := e.pawnStructure(White), e.pawnStructure(Black)
		e.pawns.score.clear().add(white).sub(black).apply(weightPawnStructure)
		e.pawns.id = key

		// Force full king shelter evaluation since any legit king square
		// will be viewed as if the king has moved.
		e.pawns.king[White], e.pawns.king[Black] = 0xFF, 0xFF

		if engine.traceʔ {
			e.checkpoint(`Pawns`, Total{white, black})
		}
	}

	e.score.add(e.pawns.score)
}

func (e *Evaluation) analyzePassers() {
	var white, black, score Score

	if engine.traceʔ {
		defer func() {
			e.checkpoint(`Passers`, Total{white, black})
		}()
	}

	white, black = e.pawnPassers(White), e.pawnPassers(Black)
	score.add(white).sub(black).apply(weightPassedPawns)
	e.score.add(score)
}

// Calculates extra bonus and penalty based on pawn structure. Specifically,
// a bonus is awarded for passed pawns, and penalty applied for isolated and
// doubled pawns.
func (e *Evaluation) pawnStructure(our int) (score Score) {
	their := our^1
	ourPawns := e.position.outposts[pawn(our)]
	theirPawns := e.position.outposts[pawn(their)]
	e.pawns.passers[our] = 0

	for bm := ourPawns; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		row, col := coordinate(square)

		isolatedʔ := (maskIsolated[col] & ourPawns).noneʔ()
		exposedʔ := (maskInFront[our][square] & theirPawns).noneʔ()
		doubledʔ := (maskInFront[our][square] & ourPawns).anyʔ()
		supportedʔ := (maskIsolated[col] & (maskRank[row] | maskRank[row].up(their)) & ourPawns).anyʔ()

		// The pawn is passed if a) there are no enemy pawns in the same
		// and adjacent columns; and b) there are no same our pawns in
		// front of us.
		passedʔ := !doubledʔ && (maskPassed[our][square] & theirPawns).noneʔ()
		if passedʔ {
			e.pawns.passers[our] |= bit(square)
		}

		// Penalty if the pawn is isolated, i.e. has no friendly pawns
		// on adjacent files. The penalty goes up if isolated pawn is
		// exposed on semi-open file.
		if isolatedʔ {
			if !exposedʔ {
				score.sub(penaltyIsolatedPawn[col])
			} else {
				score.sub(penaltyWeakIsolatedPawn[col])
			}
		} else if !supportedʔ {
			score.sub(pawnAlone) // Small penalty if the pawn is not supported by a fiendly pawn.
		}

		// Penalty if the pawn is doubled, i.e. there is another friendly
		// pawn in front of us.
		if doubledʔ {
			score.sub(penaltyDoubledPawn[col])
		}

		// Penalty if the pawn is backward.
		backwardʔ := false
		if (!passedʔ && !supportedʔ && !isolatedʔ) {

			// Backward pawn should not be attacking enemy pawns.
			if (pawnAttacks[our][square] & theirPawns).noneʔ() {

				// Backward pawn should not have friendly pawns behind.
				if (maskPassed[their][square] & maskIsolated[col] & ourPawns).noneʔ() {

					// Backward pawn should face enemy pawns on the next two ranks
					// preventing its advance.
					enemy := pawnAttacks[our][square].up(our)
					if ((enemy | enemy.up(our)) & theirPawns).anyʔ() {
						backwardʔ = true
						if !exposedʔ {
							score.sub(penaltyBackwardPawn[col])
						} else {
							score.sub(penaltyWeakBackwardPawn[col])
						}
					}
				}
			}
		}

		// Bonus if the pawn has good chance to become a passed pawn.
		if exposedʔ && !isolatedʔ && !passedʔ && !backwardʔ {
			his := maskPassed[their][square + up[our]] & maskIsolated[col] & ourPawns
			her := maskPassed[our][square] & maskIsolated[col] & theirPawns
			if his.count() >= her.count() {
				score.add(bonusSemiPassedPawn[rank(our, square)])
			}
		}
	}

	return score
}

func (e *Evaluation) pawnPassers(our int) (score Score) {
	p, their := e.position, our^1

	// If opposing side has no pieces other than pawns then need to check if passers are unstoppable.
	chase := (p.outposts[their] ^ p.outposts[pawn(their)] ^ p.outposts[king(their)]).noneʔ()

	for bm := e.pawns.passers[our]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		rank := rank(our, square)
		bonus := bonusPassedPawn[rank]

		if rank > A2H2 {
			extra := extraPassedPawn[rank]
			nextSquare := square + up[our]

			// Adjust endgame bonus based on how close the kings are from the
			// step forward square.
			bonus.endgame += (distance[p.king[their]][nextSquare] * 5 - distance[p.king[our]][nextSquare] * 2) * extra

			// Check if the pawn can step forward.
			if p.board.offʔ(nextSquare) {
				boost := 0

				// Assume all squares in front of the pawn are under attack.
				attacked := maskInFront[our][square]
				defended := attacked & e.attacks[our]

				// Boost the bonus if squares in front of the pawn are defended.
				if defended == attacked {
					boost += 6 // All squares.
				} else if defended.onʔ(nextSquare) {
					boost += 4 // Next square only.
				}

				// Check who is attacking the squares in front of the pawn including
				// queen and rook x-ray attacks from behind.
				enemy := maskInFront[their][square] & (p.outposts[queen(their)] | p.outposts[rook(their)])
				if enemy.noneʔ() || (enemy & p.rookMoves(square)).noneʔ() {

					// Since nobody attacks the pawn from behind adjust the attacked
					// bitmask to only include squares attacked or occupied by the enemy.
					attacked &= (e.attacks[their] | p.outposts[their])
				}

				// Boost the bonus if passed pawn is free to advance to the 8th rank
				// or at least safely step forward.
				if attacked.noneʔ() {
					boost += 15 // Remaining squares are not under attack.
				} else if attacked.offʔ(nextSquare) {
					boost += 9  // Next square is not under attack.
				}

				if boost > 0 {
					bonus.adjust(extra * boost)
				}
			}
		}

		// Before chasing the unstoppable make sure own pieces are not blocking the passer.
		if chase && (p.outposts[our] & maskInFront[our][square]).noneʔ() {
			// Pick square rule bitmask for the pawn. If defending king has the right
			// to move then pick extended square mask.
			bits := maskNone
			if p.color == our {
				bits = maskSquare[our][square]
			} else {
				bits = maskSquareEx[our][square]
			}
			if (bits & p.outposts[king(their)]).noneʔ() {
				bonus.endgame += unstoppablePawn
			}
		}

		score.add(bonus)
	}

	return score
}

