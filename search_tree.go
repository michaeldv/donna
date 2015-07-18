// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (p *Position) searchTree(alpha, beta, depth int) (score int) {
	ply := ply()

	// Reset principal variation.
	game.pv[ply].size = 0

	// Return if it's time to stop search.
	if ply >= MaxPly || engine.clock.halt {
		return p.Evaluate()
	}

	// Insufficient material and repetition/perpetual check pruning.
	if p.insufficient() || p.repetition() || p.fifty() {
		return 0
	}

	// Checkmate distance pruning.
	if score := abs(ply - Checkmate); score < beta {
		if score <= alpha {
			return alpha
		}
		beta = score
	}

	// Initialize node search conditions.
	isNull := p.isNull()
	inCheck := p.isInCheck(p.color)
	isPrincipal := (beta - alpha > 1)

	// Probe cache.
	cachedMove := Move(0)
	staticScore := UnknownScore
	if cached := p.probeCache(); cached != nil {
		cachedMove = cached.move
		if int(cached.depth) >= depth {
			staticScore = uncache(int(cached.score), ply)
			if (cached.flags == cacheExact && isPrincipal) ||
			   (cached.flags == cacheBeta  && staticScore >= beta) ||
			   (cached.flags == cacheAlpha && staticScore <= alpha) {
				if staticScore >= beta && !inCheck && cachedMove != 0 && cachedMove.isQuiet() {
					game.saveGood(depth, cachedMove)
				}
				return staticScore
			}
		}
	}

	// Quiescence search.
	if !inCheck && depth < 1 {
		return p.searchQuiescence(alpha, beta, 0, inCheck)
	}

	if staticScore == UnknownScore {
		staticScore = p.Evaluate()
	}

	// Razoring and futility margin pruning.
	if !inCheck && !isPrincipal {

		// No razoring if pawns are on 7th rank.
		if cachedMove == Move(0) && depth < 8 && p.outposts[pawn(p.color)] & mask7th[p.color] == 0 {
			razoringMargin := func(depth int) int {
				return 512 + 64 * (depth - 1)
			}

		   	// Special case for razoring at low depths.
			if depth <= 2 && staticScore <= alpha - razoringMargin(5) {
				return p.searchQuiescence(alpha, beta, 0, inCheck)
			}
			
			margin := alpha - razoringMargin(depth)
			if score := p.searchQuiescence(alpha, beta + 1, 0, inCheck); score <= margin {
				return score
			}
		}

		// Futility pruning is only applicable if we don't have winning score
		// yet and there are pieces other than pawns.
		if !isNull && depth < 14 && !isMate(beta) &&
		   (p.outposts[p.color] & ^(p.outposts[king(p.color)] | p.outposts[pawn(p.color)])).any() {
			// Largest conceivable positional gain.
			if gain := staticScore - 256 * depth; gain >= beta {
				return gain
			}
		}

		// Null move pruning.
		if !isNull && depth > 1 && p.outposts[p.color].count() > 5 {
			position := p.makeNullMove()
			game.nodes++
			nullScore := -position.searchTree(-beta, -beta + 1, depth - 1 - 3)
			position.undoNullMove()

			if nullScore >= beta {
				if isMate(nullScore) {
					return beta
				}
				return nullScore
			}
		}
	}

	// Internal iterative deepening.
	if !inCheck && cachedMove == Move(0) && depth > 4 {
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

	bestMove := Move(0)
	moveCount, quietMoveCount := 0, 0
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if !gen.isValid(move) {
			continue
		}

		position := p.makeMove(move)
		moveCount++

		// Search depth extension/reduction.
		newDepth, reduction := depth - 1, 0
		if position.isInCheck(position.color) { // Extend search depth if we're checking.
			newDepth++
		} else if !isPrincipal && !inCheck && depth > 2 && moveCount > 1 && move.isQuiet() && !move.isPawnAdvance() {
			quietMoveCount++
			if quietMoveCount >= 8 {
				reduction++
				if quietMoveCount >= 16 {
					reduction++
					if quietMoveCount >= 24 {
						reduction++
					}
				}
			}
		}

		// Start search with full window.
		if isPrincipal && moveCount == 1 {
			score = -position.searchTree(-beta, -alpha, newDepth)
		} else if reduction > 0 {
			score = -position.searchTree(-alpha - 1, -alpha, max(0, newDepth - reduction))

			// Verify late move reduction and re-run the search if necessary.
			if score > alpha {
				score = -position.searchTree(-alpha - 1, -alpha, newDepth)
			}
		} else {
			score = -position.searchTree(-alpha - 1, -alpha, newDepth)

			// If zero window failed try full window.
			if score > alpha && score < beta {
				score = -position.searchTree(-beta, -alpha, newDepth)
			}
		}
		position.undoLastMove()

		if engine.clock.halt {
			game.nodes += moveCount
			//Log("searchTree at %d (%s): move %s (%d) score %d alpha %d\n", depth, C(p.color), move, moveCount, score, alpha)
			return alpha
		}

		if score > alpha {
			alpha = score
			bestMove = move
			if isPrincipal {
				game.saveBest(ply, move)
			}

			if alpha >= beta {
				break // Stop searching. Happiness is right next to you.
			}
		}
	}

	game.nodes += moveCount

	if moveCount == 0 {
		if inCheck {
			alpha = -Checkmate + ply
		} else {
			alpha = 0
		}
	} else if score >= beta && !inCheck {
		game.saveGood(depth, bestMove)
	}

	score = alpha

	cacheFlags := cacheAlpha
	if score >= beta {
		cacheFlags = cacheBeta
	} else if (isPrincipal && moveCount > 0) {
		cacheFlags = cacheExact
	}
	p.cache(bestMove, score, depth, ply, cacheFlags)

	return
}
