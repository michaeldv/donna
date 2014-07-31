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
			e.checkpoint(`Imbalance`, material.score)
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

	whitePair, blackPair := 0, 0
	if count[Bishop] > 1 {
		whitePair++
	}
	if count[BlackBishop] > 1 {
		blackPair++
	}

	white := e.imbalance(whitePair, count[Pawn], count[Knight], count[Bishop], count[Rook], count[Queen],
		blackPair, count[BlackPawn], count[BlackKnight], count[BlackBishop], count[BlackRook], count[BlackQueen])
	black := e.imbalance(blackPair, count[BlackPawn], count[BlackKnight], count[BlackBishop], count[BlackRook], count[BlackQueen],
		whitePair, count[Pawn], count[Knight], count[Bishop], count[Rook], count[Queen])

	adjustment := (white - black) / 32
	score.midgame = adjustment
	score.endgame = adjustment

	return
}

// Simplified second-degree polynomial material imbalance by Tord Romstad.
func (e *Evaluation) imbalance(w2, wP, wN, wB, wR, wQ, b2, bP, bN, bB, bR, bQ int) int {
	polynom := func(a, b, c, x int) int {
		return a * (x * x) + (b + c) * x
	}

	return polynom(   0, (   0                                                                                    ),  1852, w2) +
	       polynom(   2, (  39*w2 +                                      37*b2                                    ),  -162, wP) +
	       polynom(  -4, (  35*w2 + 271*wP +                             10*b2 +  62*bP                           ), -1122, wN) +
	       polynom(   0, (   0*w2 + 105*wP +   4*wN +                    57*b2 +  64*bP +  39*bN                  ),  -183, wB) +
	       polynom(-141, ( -27*w2 +  -2*wP +  46*wN + 100*wB +           50*b2 +  40*bP +  23*bN + -22*bB         ),   249, wR) +
	       polynom(   0, (-177*w2 +  25*wP + 129*wN + 142*wB + -137*wR + 98*b2 + 105*bP + -39*bN + 141*bB + 274*bR),  -154, wQ)
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
							for wN := 1; wN <= 2; wN++ {
								count[Knight] = wN
								for bN := 1; bN <= 2; bN++ {
									count[BlackKnight] = bN
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

		material.phase = 12 * (wN + bN + wB + bB) + 18 * (wR + bR) + 44 * (wQ + bQ)

		wPair, bPair := 0, 0
		if wB > 1 {
			wPair++
		}
		if bP > 1 {
			bPair++
		}

		// Cheating with eval global: this should be game.eval.
		white := eval.imbalance(wPair, wP, wN, wB, wR, wQ,  bPair, bP, bN, bB, bR, bQ)
		black := eval.imbalance(bPair, bP, bN, bB, bR, bQ,  wPair, wP, wN, wB, wR, wQ)

		adjustment := (white - black) / 32
		material.score.midgame += adjustment
		material.score.endgame += adjustment
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
