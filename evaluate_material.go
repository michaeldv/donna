// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

const (
	whiteKingSafety    = 0x01  // Should we worry about white king's safety?
	blackKingSafety    = 0x02  // Ditto for the black king.
	materialDraw       = 0x04  // King vs. King (with minor)
	knownEndgame       = 0x08  // Where we calculate exact score.
	lesserKnownEndgame = 0x10  // Where we set score markdown value.
	oppositeBishops    = 0x20  // Sides have one bishop on opposite color squares.
)

type Function func(*Evaluation) int

type MaterialEntry struct {
	hash      uint64 	// Material hash key.
	flags     uint8    	// Evaluation flags based on material balance.
	phase     int 		// Game phase, based on available material.
	score     Score 	// Score adjustment for the given material.
	endgame   Function 	// Function to analyze an endgame position.
}

var materialCache [8192]MaterialEntry

func (e *Evaluation) analyzeMaterial() {
	e.material = e.fetchMaterial()
	e.score.add(e.material.score)
}

func (e *Evaluation) fetchMaterial() *MaterialEntry {
	key := e.position.hashMaterial

	// Just like with pawns cache use faster 32-bit indexing.
	index := uint32(key) % uint32(len(materialCache))
	material := &materialCache[index]

	// Bypass material cache if evaluation tracing is enabled.
	if material.hash != key || Settings.Trace {
		material.hash = key
		material.phase = e.materialPhase()
		material.flags, material.endgame = e.materialFlagsAndFunction(key)
		material.score = e.materialScore()

		if Settings.Trace {
			e.checkpoint(`Material`, material.score)
		}
	}

	return material
}

// Set up evaluation flags based on the material balance.
func (e *Evaluation) materialFlagsAndFunction(key uint64) (flags uint8, endgame Function) {
	count := &e.position.count

	// Calculate material balances for both sides to simplify comparisons.
	whiteForce := count[Pawn]      + (count[Knight]      + count[Bishop]     ) * 10 + count[Rook]      * 100 + count[Queen]      * 1000
	blackForce := count[BlackPawn] + (count[BlackKnight] + count[BlackBishop]) * 10 + count[BlackRook] * 100 + count[BlackQueen] * 1000

	noPawns := (count[Pawn] + count[BlackPawn] == 0)
	bareKing := (whiteForce * blackForce == 0) // Bare king (white, black or both).

	// Set king safety flags if the opposing side has a queen and at least one piece.
	if whiteForce >= 1010 {
		flags |= blackKingSafety
	}
	if blackForce >= 1010 {
		flags |= whiteKingSafety
	}

	// Insufficient material endgames that don't require further evaluation:
	// 1) Two bare kings.
	if whiteForce + blackForce == 0 {
		flags |= materialDraw

	// 2) No pawns and king with a minor.
	} else if noPawns && whiteForce <= 10 && blackForce <= 10 {
		flags |= materialDraw

	// 3) No pawns and king with two knights.
	} else if whiteForce + blackForce == 20 && count[Knight] + count[BlackKnight] == 2 {
		flags |= materialDraw

	// Known endgame: king and a pawn vs. bare king.
	} else if key == 0x5355F900C2A82DC7 || key == 0x9D39247E33776D41 {
		flags |= knownEndgame
		endgame = (*Evaluation).kingAndPawnVsBareKing

	// Known endgame: king with a knight and a bishop vs. bare king.
	} else if key == 0xE6F0FBA55BF280F1 || key == 0x29D8066E0A562122 {
		flags |= knownEndgame
		endgame = (*Evaluation).knightAndBishopVsBareKing

	// Known endgame: king with some winning material vs. bare king.
	} else if bareKing && Abs(whiteForce - blackForce) > 100 {
		flags |= knownEndgame
		endgame = (*Evaluation).winAgainstBareKing

	// Lesser known endgame: king and two or more pawns vs. bare king.
	} else if bareKing && whiteForce + blackForce <= 8 {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).kingAndPawnsVsBareKing

	// Lesser known endgame: queen vs. rook with pawn(s)
	} else if (blackForce == 1000 && whiteForce - count[Pawn] == 100 && count[Pawn] > 0) ||
		  (whiteForce == 1000 && blackForce - count[BlackPawn] == 100 && count[BlackPawn] > 0) {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).queenVsRookAndPawns

	// Lesser known endgame: king and pawn vs. king and pawn.
	} else if key == 0xCE6CDD7EF1DF4086 {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).kingAndPawnVsKingAndPawn

	// Lesser known endgame: bishop and pawn vs. bare king.
	} else if key == 0x70E2F7DBDBFDE978 || key == 0xE2A24E8FD880E6EE {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).bishopAndPawnVsBareKing

	// Lesser known endgame: rook and pawn vs. rook.
	} else if key == 0x29F14397EB52ECA8 || key == 0xE79D9EE91A8DAC2E {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).rookAndPawnVsRook
	}

	// Do we have opposite-colored bishops?
	if count[Bishop] * count[BlackBishop] == 1 && flags & (materialDraw | knownEndgame) == 0 {
		bishops := e.position.outposts[Bishop] | e.position.outposts[BlackBishop]
		if bishops & maskDark != 0 && bishops & ^maskDark != 0 {
			flags |= oppositeBishops
		}
	}

	return
}

// Calculates game phase based on what pieces are on the board (256 for the
// initial position, 0 for bare kings).
func (e *Evaluation) materialPhase() int {
	count := &e.position.count

	phase := 12 * (count[Knight] + count[BlackKnight] + count[Bishop] + count[BlackBishop]) +
		 18 * (count[Rook]   + count[BlackRook]) +
		 44 * (count[Queen]  + count[BlackQueen])

	return Min(256, phase)
}

// Calculates material score adjustment for the position we are evaluating.
func (e *Evaluation) materialScore() (score Score) {
	count := &e.position.count

	// Bonus for the pair of bishops.
	if count[Bishop] > 1 {
		score.add(bishopPair)
		if count[Pawn] > 5 {
			score.subtract(bishopPairPawn.times(count[Pawn] - 5))
		}
	}
	if count[BlackBishop] > 1 {
		score.subtract(bishopPair)
		if count[BlackPawn] > 5 {
			score.add(bishopPairPawn.times(count[BlackPawn] - 5))
		}
	}

	return
}

// Pre-populates material cache with the most common middle game material
// balances, namely zero or one queen, one or two rooks/bishops/knights, and
// four to eight pawns. Total number of pre-populated entries is
// (2*2) * (2*2) * (2*2) * (2*2) * (5*5) = 6400.
func (g *Game) warmUpMaterialCache() {
	var key uint64
	var index uint32
	var count [14]int
	var material *MaterialEntry

	for wQ := 0; wQ <= 1; wQ++ {
		count[Queen] = wQ
		for bQ := 0; bQ <= 1; bQ++ {
			count[BlackQueen] = bQ
			for wR := 1; wR <=2; wR++ {
				count[Rook] = wR
				for bR := 1; bR <= 2; bR++ {
					count[BlackRook] = bR
					for wB := 1; wB <= 2; wB++ {
						count[Bishop] = wB
						for bB := 1; bB <= 2; bB++ {
							count[BlackBishop] = bB
							for wK := 1; wK <= 2; wK++ {
								count[Knight] = wK
								for bK := 1; bK <= 2; bK++ {
									count[BlackKnight] = bK
									for wP := 4; wP <= 8; wP++ {
										count[Pawn] = wP
										for bP := 4; bP <= 8; bP++ {
											count[BlackPawn] = bP
		// Compute material hash key for the current material balance.
		key = 0
		for piece := Pawn; piece <= BlackQueen; piece++ {
			for i := 0; i < count[piece]; i++ {
				key ^= Piece(piece).polyglot(i)
			}
		}

		// Compute index and populate material cache entry.
		index = uint32(key) % uint32(len(materialCache))
		material = &materialCache[index]
		material.hash = key

		material.phase = 12 * (wK + bK + wB + bB) + 18 * (wR + bR) + 44 * (wQ + bQ)

		// Bonus for the pair of bishops.
		if wB > 1 {
			material.score.add(bishopPair)
			if wP > 5 {
				material.score.subtract(bishopPairPawn.times(wP - 5))
			}
		}
		if bB > 1 {
			material.score.subtract(bishopPair)
			if bP > 5 {
				material.score.add(bishopPairPawn.times(bP - 5))
			}
		}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
