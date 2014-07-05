// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

// Search for the node in check.
func (p *Position) searchInCheck(beta, depth int) int {
	p.game.nodes++
	if p.isRepetition() {
		return 0
	}

	ply := Ply()
	bestScore := ply - Checkmate
	if bestScore >= beta {
		return bestScore //beta
	}

	// Probe cache.
	if cached := p.probeCache(); cached != nil {
		if cached.depth >= depth {
			score := cached.score
			if score > Checkmate-MaxPly && score <= Checkmate {
				score -= ply
			} else if score >= -Checkmate && score < -Checkmate+MaxPly {
				score += ply
			}

			if (cached.flags == cacheExact) ||
				(cached.flags == cacheBeta && score >= beta) ||
				(cached.flags == cacheAlpha && score <= beta) {
				return score
			}
		}
	}

	p.game.pvsize[ply] = ply
	gen := NewGen(p, ply).generateEvasions().quickRank()
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if position := p.MakeMove(move); position != nil {
			//Log("%*schck/%s> depth: %d, ply: %d, move: %s\n", ply*2, ` `, C(p.color), depth, ply, move)
			inCheck := position.isInCheck(position.color)
			reducedDepth := depth - 1
			if inCheck {
				reducedDepth++
			}

			moveScore := 0
			if reducedDepth == 0 {
				moveScore = -position.searchQuiescence(-beta, 1-beta)
			} else if inCheck {
				moveScore = -position.searchInCheck(1-beta, reducedDepth)
			} else {
				moveScore = -position.searchWithZeroWindow(1-beta, reducedDepth)
			}

			position.TakeBack(move)
			if moveScore > bestScore {
				p.game.saveBest(ply, move)
				if moveScore >= beta {
					p.cache(move, moveScore, depth, cacheBeta)
					return moveScore
				}
				bestScore = moveScore
			}
		}
	}

	return bestScore
}
