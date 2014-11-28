// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

func (p *Position) searchTree(alpha, beta, depth int) (score int) {
	ply := ply()

	// Reset principal variation.
	game.pv[ply] = game.pv[ply][:0]

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
		beta = score
		if score <= alpha {
			return alpha
		}
	}

	// Initialize node search conditions.
	isNull := p.isNull()
	inCheck := p.isInCheck(p.color)
	isPrincipal := (beta - alpha > 1)

	// Probe cache.
	cachedMove := Move(0)
	cacheFlags := uint8(cacheAlpha)
	if cached := p.probeCache(); cached != nil {
		cachedMove = cached.move
		if cached.depth >= depth {
			score := cached.score
			if score > Checkmate - MaxPly && score <= Checkmate {
				score -= ply
			} else if score >= -Checkmate && score < -Checkmate + MaxPly {
				score += ply
			}
			if (cached.flags == cacheExact && isPrincipal) ||
			   (cached.flags == cacheBeta  && score >= beta) ||
			   (cached.flags == cacheAlpha && score <= alpha) {
				if score >= beta && !inCheck && cachedMove != 0 && cachedMove.isQuiet() {
					game.saveGood(depth, cachedMove)
				}
				return score
			}
		}
	}

	// Quiescence search.
	if !inCheck && depth < 1 {
		return p.searchQuiescence(alpha, beta, depth)
	}

	// Razoring and futility margin pruning.
	if !inCheck && !isPrincipal {
		staticScore := p.Evaluate()

		// No razoring if pawns are on 7th rank.
		if cachedMove == Move(0) && depth < 8 && p.outposts[pawn(p.color)] & mask7th[p.color] == 0 {
			razoringMargin := func(depth int) int {
				return 512 + 64 * (depth - 1)
			}

		   	// Special case for razoring at low depths.
			if depth <= 2 && staticScore <= alpha - razoringMargin(5) {
				return p.searchQuiescence(alpha, beta, 0)
			}
			
			margin := alpha - razoringMargin(depth)
			if score := p.searchQuiescence(alpha, beta + 1, 0); score <= margin {
				return score
			}
		}

		// Futility pruning is only applicable if we don't have winning score
		// yet and there are pieces other than pawns.
		if !isNull && depth < 14 && abs(beta) < Checkmate - MaxPly &&
		   p.outposts[p.color] & ^(p.outposts[king(p.color)] | p.outposts[pawn(p.color)]) != 0 {
			futilityMargin := func (depth int) int {
				return 256 * depth
			}

			// Largest conceivable positional gain.
			if gain := staticScore - futilityMargin(depth); gain >= beta {
				return gain
			}
		}
	}

	// Null move pruning.
	if !inCheck && !isNull && depth > 1 && p.outposts[p.color].count() > 5 {
		position := p.makeNullMove()
		game.nodes++
		nullScore := -position.searchTree(-beta, -beta + 1, depth - 1 - 3)
		position.undoNullMove()

		if nullScore >= beta {
			if abs(nullScore) >= Checkmate - MaxPly {
				return beta
			}
			return nullScore
		}
	}

	// Internal iterative deepening.
	if cachedMove == 0 && depth > 4 {
		p.searchTree(alpha, beta, depth - 4)
		if len(game.pv[ply]) > 0 {
			cachedMove = game.pv[ply][0]
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
		newDepth := depth - 1

		// Search depth extension.
		giveCheck := position.isInCheck(position.color)
		if giveCheck {
			newDepth++
		}

		// Late move reduction.
		lateMoveReduction := false
		if depth >= 3 && !isPrincipal && !inCheck && !giveCheck && move.isQuiet() {
			quietMoveCount++
			if newDepth > 0 && quietMoveCount >= 8 {
				newDepth--
				lateMoveReduction = true
				if quietMoveCount >= 16 {
					newDepth--
					if quietMoveCount >= 24 {
						newDepth--
					}
				}
			}
		}

		// Start search with full window.
		if moveCount == 1 {
			score = -position.searchTree(-beta, -alpha, newDepth)
		} else if lateMoveReduction {
			score = -position.searchTree(-alpha - 1, -alpha, newDepth)

			// Verify late move reduction and re-run the search if necessary.
			if score > alpha {
				score = -position.searchTree(-alpha - 1, -alpha, newDepth + 1)
			}
		} else {
			if newDepth < 2 {
				score = -position.searchQuiescence(-alpha - 1, -alpha, 0)
			} else {
				score = -position.searchTree(-alpha - 1, -alpha, newDepth)
			}

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
			cacheFlags = cacheExact
			game.saveBest(ply, move)

			if alpha >= beta {
				cacheFlags = cacheBeta
				break
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
	p.cache(bestMove, score, depth, cacheFlags)

	return
}
