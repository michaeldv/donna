// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (p *Position) searchTree(alpha, beta, depth int) (score int) {
	ply := ply()

	// Return if it's time to stop search.
	if ply >= MaxPly || engine.clock.halt {
		return p.Evaluate()
	}

	// Reset principal variation.
	game.pv[ply].size = 0

	// Insufficient material and repetition/perpetual check pruning.
	if p.fifty() || p.insufficient() || p.repetition() {
		return 0
	}

	// Checkmate distance pruning.
	alpha, beta = mateDistance(alpha, beta, ply)
	if alpha >= beta {
		return alpha
	}

	// Initialize node search conditions.
	isNull := p.isNull()
	inCheck := p.isInCheck(p.color)
	isPrincipal := (beta - alpha > 1)

	// Probe cache.
	cached, cachedMove := p.probeCache(), Move(0)
	if cached != nil {
		cachedMove = cached.move
		if !isPrincipal && cached.depth() >= depth {
			bounds, score := cached.bounds(), cached.score(ply)
			if (bounds == cacheBeta && score >= beta) || (bounds == cacheAlpha && score <= alpha) {
				if score >= beta && !inCheck && cachedMove.some() {
					game.saveGood(depth, cachedMove)
				}
				return score
			}
		}
	}

	if !inCheck {
		if depth < 1 {
			return p.searchQuiescence(alpha, beta, 0, inCheck)
		}
		if cached != nil {
			if p.score == Unknown {
				p.score = p.Evaluate()
			}
		} else {
			if isNull {
				p.score = rightToMove.midgame * 2 - tree[node-1].score
			} else {
				p.score = p.Evaluate()
			}
		}
	}

	// Razoring and futility margin pruning.
	if !inCheck && !isPrincipal {

		// No razoring if pawns are on 7th rank.
		if cachedMove.null() && depth < 3 && p.outposts[pawn(p.color)] & mask7th[p.color] == 0 {
			razoringMargin := func(depth int) int {
				return 96 + 64 * (depth - 1)
			}

		   	// Special case for razoring at low depths.
			if p.score <= alpha - razoringMargin(5) {
				return p.searchQuiescence(alpha, beta, 0, inCheck)
			}

			margin := alpha - razoringMargin(depth)
			if score := p.searchQuiescence(margin, margin + 1, 0, inCheck); score <= margin {
				return score
			}
		}

		// Futility pruning is only applicable if we don't have winning score
		// yet and there are pieces other than pawns.
		if !isNull && depth < 14 && !isMate(beta) &&
		   (p.outposts[p.color] & ^(p.outposts[king(p.color)] | p.outposts[pawn(p.color)])).any() {
			// Largest conceivable positional gain.
			if gain := p.score - 256 * depth; gain >= beta {
				return gain
			}
		}

		// Null move pruning.
		if !isNull && depth > 1 && p.outposts[p.color].count() > 5 {
			position := p.makeNullMove()
			game.nodes++
			nullScore := -position.searchTree(-beta, -beta + 1, depth - 1 - 3)
			position.undoLastMove()

			if nullScore >= beta {
				if isMate(nullScore) {
					return beta
				}
				return nullScore
			}
		}
	}

	// Internal iterative deepening.
	if !inCheck && cachedMove.null() && depth > 4 {
		newDepth := depth / 2
		if isPrincipal {
			newDepth = depth - 2
		}
		p.searchTree(alpha, beta, newDepth)
		if cached := p.probeCache(); cached != nil {
			cachedMove = cached.move
		}
	}

	gen := NewGen(p, ply)
	if inCheck {
		gen.generateEvasions().quickRank()
	} else {
		gen.generateMoves().rank(cachedMove)
	}

	bestScore := alpha
	bestMove, moveCount := Move(0), 0
	for move := gen.nextMove(); move.some(); move = gen.nextMove() {
		if !move.valid(p, gen.pins) {
			continue
		}

		position := p.makeMove(move)
		moveCount++; game.nodes++

		// Reduce search depth if we're not checking.
		giveCheck := position.isInCheck(position.color)
		newDepth := let(giveCheck && p.exchange(move) >= 0, depth, depth - 1)

		// Start search with full window.
		if isPrincipal && moveCount == 1 {
			score = -position.searchTree(-beta, -alpha, newDepth)
		} else {
			reduction := 0
			if !inCheck && !giveCheck && depth > 2 && move.isQuiet() && !move.isKiller(ply) && !move.isPawnAdvance() {
				reduction = lateMoveReductions[min(63, moveCount-1)][min(63, depth)]
				if isPrincipal {
					reduction /= 2
				} else if game.history[move.piece()][move.to()] < 0 {
					reduction++
				}
			}

			score = -position.searchTree(-alpha - 1, -alpha, max(0, newDepth - reduction))

			// Verify late move reduction and re-run the search if necessary.
			if reduction > 0 && score > alpha {
				score = -position.searchTree(-alpha - 1, -alpha, newDepth)
			}

			// If zero window failed try full window.
			if isPrincipal && score > alpha && score < beta {
				score = -position.searchTree(-beta, -alpha, newDepth)
			}
		}
		position.undoLastMove()

		// Don't touch anything if the time has elapsed and we need to abort th search.
		if engine.clock.halt {
			return alpha
		}

		if score > bestScore {
			bestScore = score
			if score > alpha {
				if isPrincipal {
					game.saveBest(ply, move)
				}
				if isPrincipal && score < beta {
					alpha = score
					bestMove = move
				} else {
					p.cache(move, score, depth, ply, cacheBeta)
					return score
				}
			}
		}
	}

	if moveCount == 0 {
		score = let(inCheck, matedIn(ply), 0)
	} else {
		score = bestScore
		if !inCheck {
			game.saveGood(depth, bestMove)
		}
	}

	cacheFlags := cacheAlpha
	if score >= beta {
		cacheFlags = cacheBeta
	} else if isPrincipal && bestMove.some() {
		cacheFlags = cacheExact
	}
	p.cache(bestMove, score, depth, ply, cacheFlags)

	return score
}
