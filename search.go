// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

// Root node search. Basic principle is expressed by Boob's Law: you always find
// something in the last place you look.
func (p *Position) search(alpha, beta, depth int) (score int) {
	ply, inCheckʔ := ply(), p.inCheckʔ(p.color)

	// Root move generator makes sure all generated moves are valid. The
	// best move found so far is always the first one we search.
	gen := NewRootGen(p, depth)
	if depth == 1 {
		gen.generateRootMoves()
	} else {
		gen.reset()
	}

	bestAlpha, bestScore := alpha, alpha
	bestMove, moveCount := Move(0), 0
	for move := gen.nextMove(); move.someʔ(); move = gen.nextMove() {
		position := p.makeMove(move)
		moveCount++; game.nodes++
		if engine.uciʔ {
			engine.uciMove(move, moveCount, depth)
		}

		// Reduce search depth if we're not checking.
		giveCheck := position.inCheckʔ(position.color)
		newDepth := let(giveCheck && p.exchange(move) >= 0, depth, depth - 1)

		// Start search with full window.
		game.deepeningʔ = (moveCount == 1)
		if moveCount == 1 {
			score = -position.searchTree(-beta, -alpha, newDepth)
		} else {
			reduction := 0
			if !inCheckʔ && !giveCheck && depth > 2 && move.quietʔ() && !move.killerʔ(ply) && !move.pawnAdvanceʔ() {
				reduction = lateMoveReductions[(moveCount-1) & 63][depth & 63]
				if game.history[move.piece()][move.to()] < 0 {
					reduction++
				}
			}

			score = -position.searchTree(-alpha - 1, -alpha, max(0, newDepth - reduction))

			// Verify late move reduction and re-run the search if necessary.
			if reduction > 0 && score > alpha {
				score = -position.searchTree(-alpha - 1, -alpha, newDepth)
			}

			// If zero window fails then try full window.
			if score > alpha {
				score = -position.searchTree(-beta, -alpha, newDepth)
			}
		}
		position.undoLastMove()

		// Don't touch anything if the time has elapsed and we need to abort th search.
		if engine.clock.haltʔ {
			return alpha
		}

		if moveCount == 1 || score > alpha {
			bestMove = move
			game.saveBest(0, move)
			gen.scoreMove(depth, score).rearrangeRootMoves()
			if moveCount > 1 {
				game.volatility++
			}
		} else {
			gen.scoreMove(depth, -depth)
		}

		if score > bestScore {
			bestScore = score
			if score > alpha {
				game.saveBest(ply, move)
				if score < beta {
					alpha = score
					bestMove = move
				} else {
					p.cache(move, score, depth, ply, cacheBeta)
					if !inCheckʔ && alpha > bestAlpha {
						game.saveGood(depth, bestMove).updatePoor(depth, bestMove, gen.reset())
					}
					return score
				}
			}
		}
	}


	if moveCount == 0 {
		score = let(inCheckʔ, -Checkmate, 0) // Mate if in check, stalemate otherwise.
		if engine.uciʔ {
			engine.uciScore(depth, score, alpha, beta)
		}
		return score
	}
	score = bestScore

	if !inCheckʔ && alpha > bestAlpha {
		game.saveGood(depth, bestMove).updatePoor(depth, bestMove, gen.reset())
	}

	cacheFlags := cacheAlpha
	if score >= beta {
		cacheFlags = cacheBeta
	} else if bestMove.someʔ() {
		cacheFlags = cacheExact
	}
	p.cache(bestMove, score, depth, ply, cacheFlags)
	if engine.uciʔ {
		engine.uciScore(depth, score, alpha, beta)
	}

	return
}

// Testing helper method to test root search.
func (p *Position) solve(depth int) Move {
	if depth != 1 {
		NewRootGen(p, 1).generateRootMoves()
	}
	p.search(-Checkmate, Checkmate, depth)
	return game.pv[0].moves[0]
}

func (p *Position) Perft(depth int) (total int64) {
	if depth == 0 {
		return 1
	}

	gen := NewGen(p, depth).generateAllMoves().validOnly()
	if depth == 1 {
		return int64(len(gen.list))
	}

	for move := gen.nextMove(); move != 0; move = gen.nextMove() {
		position := p.makeMove(move)
		total += position.Perft(depth - 1)
		position.undoLastMove()
	}

	return
}
