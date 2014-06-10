// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

type MaterialEntry struct {
	hash   uint64 	// Material hash key.
	flags  uint8    // Evaluation flags based on material balance.
	phase  int 	// Game phase, based on available material.
	score  Score 	// Score adjustment for the given material.
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
		material.phase = e.phase()
		if e.isMaterialEqual() {
			material.score = Score{0, 0}
		} else {
			material.score = e.materialAdjustment()
		}
		if Settings.Trace {
			e.checkpoint(`Material`, material.score)
		}
	}

	return material
}

func (e *Evaluation) isMaterialEqual() bool {
	count := &e.position.count

	return count[Pawn] == count[BlackPawn] &&
	       count[Knight] == count[BlackKnight] &&
	       count[Bishop] == count[BlackBishop] &&
	       count[Rook] == count[BlackRook] &&
	       count[Queen] == count[BlackQueen]
}

// Calculates material score adjustment for the position we are evaluating.
func (e *Evaluation) materialAdjustment() (score Score) {
	count := &e.position.count

	// pawns   := count[Pawn]   - count[BlackPawn]
	// knights := count[Knight] - count[BlackKnight]
	// bishops := count[Bishop] - count[BlackBishop]
	// rooks   := count[Rook]   - count[BlackRook]
	// queens  := count[Queen]  - count[BlackQueen]

	// Bonus for the pair of bishops.
	if count[Bishop] > 1 {
		score.add(bishopPair)
	}
	if count[BlackBishop] > 1 {
		score.subtract(bishopPair)
	}

	return
}

// Calculates game phase based on what pieces are on the board (256 for the
// initial position, 0 for bare kings).
func (e *Evaluation)  phase() int {
	count := &e.position.count

	phase := 12 * (count[Knight] + count[BlackKnight] + count[Bishop] + count[BlackBishop]) +
		 18 * (count[Rook]   + count[BlackRook]) +
		 44 * (count[Queen]  + count[BlackQueen])

	return Min(256, phase)
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
		}
		if bB > 1 {
			material.score.subtract(bishopPair)
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
