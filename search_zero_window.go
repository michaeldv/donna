// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

var razoringMargin = [4]int{0, 256, 256+256, 256+256+256}
var futilityMargin = [4]int{0, 512, 512+128, 512+128+128}

// Search with zero window.
func (p *Position) searchWithZeroWindow(beta, depth int) int {
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
	cachedMove := Move(0)
	if cached := p.probeCache(); cached != nil {
		cachedMove = cached.move
		if cached.depth >= depth {
			score := cached.score
			if score > Checkmate - MaxPly && score <= Checkmate {
				score -= ply
			} else if score >= -Checkmate && score < -Checkmate + MaxPly {
				score += ply
			}

			if (cached.flags == cacheExact) ||
				(cached.flags == cacheBeta && score >= beta) ||
				(cached.flags == cacheAlpha && score <= beta) {
				return score
			}
		}
	}

	score := p.Evaluate()

	// Razoring and futility pruning. TODO: disable or tune-up in puzzle solving mode.
	if depth < len(razoringMargin) {
		if margin := beta - razoringMargin[depth]; score < margin && beta < Checkmate - MaxPly && cachedMove == Move(0) {
			if p.outposts[pawn(p.color)]&mask7th[p.color] == 0 { // No pawns on 7th.
				razorScore := p.searchQuiescence(margin-1, margin)
				if razorScore < margin {
					return razorScore
				}
			}
		}

		if margin := score - futilityMargin[depth]; margin >= beta && beta > -Checkmate + MaxPly {
			if p.outposts[p.color] & ^p.outposts[king(p.color)] & ^p.outposts[pawn(p.color)] != 0 {
				return margin
			}
		}
	}

	// Null move pruning.
	if depth > 1 && score >= beta && p.outposts[p.color].count() > 5 /*&& beta > -31000*/ {
		reduction := 3 + depth/4
		if score-100 > beta {
			reduction++
		}

		position := p.MakeNullMove()
		if depth <= reduction {
			score = -position.searchQuiescence(-beta, 1-beta)
		} else {
			score = -position.searchWithZeroWindow(1-beta, depth-reduction)
		}
		position.TakeBackNullMove()

		if score >= beta {
			return score
		}
	}

	moveCount := 0
	gen := NewGen(p, ply).generateMoves().rank(cachedMove)
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if position := p.MakeMove(move); position != nil {
			//Log("%*szero/%s> depth: %d, ply: %d, move: %s\n", ply*2, ` `, C(p.color), depth, ply, move)
			inCheck, giveCheck := position.isInCheck(position.color), position.isInCheck(position.color^1)

			reducedDepth := depth
			if !inCheck && !giveCheck && move&(isCapture|isPromo) == 0 && depth >= 3 && moveCount >= 8 {
				reducedDepth = depth - 2 // Late move reduction. TODO: disable or tune-up in puzzle solving mode.
				if reducedDepth > 0 && moveCount >= 16 {
					reducedDepth--
					if reducedDepth > 0 && moveCount >= 32 {
						reducedDepth--
					}
				}
			} else if !inCheck {
				reducedDepth = depth - 1
			}

			moveScore := 0
			if reducedDepth == 0 {
				moveScore = -position.searchQuiescence(-beta, 1-beta)
			} else if inCheck {
				moveScore = -position.searchInCheck(1-beta, reducedDepth)
			} else {
				moveScore = -position.searchWithZeroWindow(1-beta, reducedDepth)

				// Verify late move reduction.
				if reducedDepth < depth-1 && moveScore >= beta {
					moveScore = -position.searchWithZeroWindow(1-beta, depth-1)
				}
			}

			position.TakeBack(move)
			moveCount++

			if moveScore > bestScore {
				if moveScore >= beta {
					p.cache(move, moveScore, depth, cacheBeta)
					p.game.saveGood(depth, move)
					return moveScore
				}
				bestScore = moveScore
			}
		}
	} // next move.

	if moveCount == 0 {
		return 0
	}

	return bestScore
}
