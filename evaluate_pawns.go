// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type PawnEntry struct {
	id       uint64 	// Pawn hash key.
	score    Score 		// Static score for the given pawn structure.
	king     [2]uint8 	// King square for both sides.
	cover    [2]Score 	// King cover penalties for both sides.
	passers  [2]Bitmask 	// Passed pawn bitmasks for both sides.
}

type PawnCache [8192*2]PawnEntry

func (e *Evaluation) analyzePawns() {
	key := e.position.pawnId

	// Since pawn hash is fairly small we can use much faster 32-bit index.
	index := uint32(key) % uint32(len(game.pawnCache))
	e.pawns = &game.pawnCache[index]

	// Bypass pawns cache if evaluation tracing is enabled.
	if e.pawns.id != key || engine.trace {
		white, black := e.pawnStructure(White), e.pawnStructure(Black)
		white.apply(weights[1]); black.apply(weights[1]) // <-- Pawn structure weight.
		e.pawns.score.clear().add(white).subtract(black)
		e.pawns.id = key

		// Force full king shelter evaluation since any legit king square
		// will be viewed as if the king has moved.
		e.pawns.king[White], e.pawns.king[Black] = 0xFF, 0xFF

		if engine.trace {
			e.checkpoint(`Pawns`, Total{white, black})
		}
	}

	e.score.add(e.pawns.score)
}

func (e *Evaluation) analyzePassers() {
	var white, black Score

	if engine.trace {
		defer func() {
			e.checkpoint(`Passers`, Total{white, black})
		}()
	}

	white, black = e.pawnPassers(White), e.pawnPassers(Black)
	white.apply(weights[2]); black.apply(weights[2]) // <-- Passed pawns weight.
	e.score.add(white).subtract(black)
}

// Calculates extra bonus and penalty based on pawn structure. Specifically,
// a bonus is awarded for passed pawns, and penalty applied for isolated and
// doubled pawns.
func (e *Evaluation) pawnStructure(color uint8) (score Score) {
	rival := color ^ 1
	hisPawns := e.position.outposts[pawn(color)]
	herPawns := e.position.outposts[pawn(rival)]
	e.pawns.passers[color] = 0

	// Encourage center pawn moves in the opening.
	pawns := hisPawns

	for pawns != 0 {
		square := pawns.pop()
		row, col := coordinate(square)

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
		supported := (maskIsolated[col] & (maskRank[row] | maskRank[row].pushed(rival)) & hisPawns != 0)
		if supported {
			flipped := flip(color, square)
			score.add(Score{bonusSupportedPawn[flipped], bonusSupportedPawn[flipped]})
		}

		// The pawn is passed if a) there are no enemy pawns in the same
		// and adjacent columns; and b) there are no same color pawns in
		// front of us.
		passed := (maskPassed[color][square] & herPawns == 0 && !doubled)
		if passed {
			e.pawns.passers[color] |= bit[square]
		}

		// Penalty if the pawn is backward.
		backward := false
		if (!passed && !supported && !isolated) {

			// Backward pawn should not be attacking enemy pawns.
			if pawnMoves[color][square] & herPawns == 0 {

				// Backward pawn should not have friendly pawns behind.
				if maskPassed[rival][square] & maskIsolated[col] & hisPawns == 0 {

					// Backward pawn should face enemy pawns on the next two ranks
					// preventing its advance.
					enemy := pawnMoves[color][square].pushed(color)
					if (enemy | enemy.pushed(color)) & herPawns != 0 {
						backward = true
						if !exposed {
							score.subtract(penaltyBackwardPawn[col])
						} else {
							score.subtract(penaltyWeakBackwardPawn[col])
						}
					}
				}
			}
		}

		// Bonus if the pawn has good chance to become a passed pawn.
		if exposed && supported && !passed && !backward {
			his := maskPassed[rival][square + eight[color]] & maskIsolated[col] & hisPawns
			her := maskPassed[color][square] & maskIsolated[col] & herPawns
			if his.count() >= her.count() {
				score.add(bonusSemiPassedPawn[rank(color, square)])
			}
		}

		//\\ Encourage center pawn moves.
		//\\ if maskCenter.on(square) {
		//\\ 	score.midgame += bonusPawn[0][flip(color, square)] / 2
		//\\ }
	}

	return
}

func (e *Evaluation) pawnPassers(color uint8) (score Score) {
	p := e.position
	rival := color ^ 1

	// If opposing side has no pieces other than pawns then need to check if passers are unstoppable.
	chase := (p.outposts[rival] ^ p.outposts[pawn(rival)] ^ p.outposts[king(rival)]).empty()

	pawns := e.pawns.passers[color]
	for pawns.any() {
		square := pawns.pop()
		rank := rank(color, square)
		bonus := bonusPassedPawn[rank]

		if rank > A2H2 {
			extra := extraPassedPawn[rank]
			nextSquare := square + eight[color]

			// Adjust endgame bonus based on how close the kings are from the
			// step forward square.
			bonus.endgame += (distance[p.king[rival]][nextSquare] * 5 - distance[p.king[color]][nextSquare] * 2) * extra

			// Check if the pawn can step forward.
			if p.board.off(nextSquare) {
				boost := 0

				// Assume all squares in front of the pawn are under attack.
				attacked := maskInFront[color][square]
				protected := attacked & e.attacks[color]

				// Boost the bonus if squares in front of the pawn are protected.
				if protected == attacked {
					boost += 6 // All squares.
				} else if protected.on(nextSquare) {
					boost += 4 // Next square only.
				}

				// Check who is attacking the squares in front of the pawn including
				// queen and rook x-ray attacks from behind.
				enemy := maskInFront[rival][square] & (p.outposts[queen(rival)] | p.outposts[rook(rival)])
				if enemy == 0 || enemy & p.rookMoves(square) == 0 {

					// Since nobody attacks the pawn from behind adjust the attacked
					// bitmask to only include squares attacked or occupied by the enemy.
					attacked &= (e.attacks[rival] | p.outposts[rival])
				}

				// Boost the bonus if passed pawn is free to advance to the 8th rank
				// or at least safely step forward.
				if attacked == 0 {
					boost += 15 // Remaining squares are not under attack.
				} else if attacked.off(nextSquare) {
					boost += 9  // Next square is not under attack.
				}

				if boost > 0 {
					bonus.adjust(extra * boost)
				}
			}
		}

		// Before chasing the unstoppable make sure own pieces are not blocking the passer.
		if chase && (p.outposts[color] & maskInFront[color][square]).empty() {
			// Pick square rule bitmask for the pawn. If defending king has the right
			// to move then pick extended square mask.
			mask := Bitmask(0)
			if p.color == color {
				mask = maskSquare[color][square]
			} else {
				mask = maskSquareEx[color][square]
			}
			if (mask & p.outposts[king(rival)]).empty() {
				bonus.endgame += unstoppablePawn
			}
		}

		score.add(bonus)
	}

	return
}

