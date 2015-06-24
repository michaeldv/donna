// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Root node search. Basic principle is expressed by Boob's Law: you always find
// something in the last place you look.
func (p *Position) search(alpha, beta, depth int) (score int) {
	inCheck := p.isInCheck(p.color)

	// Root move generator makes sure all generated moves are valid. The
	// best move found so far is always the first one we search.
	gen := NewRootGen(p, depth)
	if depth == 1 {
		gen.generateRootMoves()
	} else if depth == 11 { // Skip moves that failed all iterations so far.
		gen.cleanupRootMoves(depth)
	} else {
		gen.reset()
	}

	moveCount, bestMove, rootAlpha := 0, Move(0), alpha
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		position := p.makeMove(move)
		moveCount++
		if engine.uci {
			engine.uciMove(move, moveCount, depth)
		}

		// Search depth extension.
		newDepth := depth - 1
		if position.isInCheck(p.color^1) { // Give check.
			newDepth++
		}

		if moveCount == 1 {
			game.deepening = true
			score = -position.searchTree(-beta, -alpha, newDepth)
		} else {
			game.deepening = false
			score = -position.searchTree(-alpha - 1, -alpha, newDepth)
			if score > alpha { // && score < beta {
				score = -position.searchTree(-beta, -alpha, newDepth)
			}
		}
		position.undoLastMove()

		if engine.clock.halt {
			//Log("searchRoot: bestMove %s pv[0][0] %s alpha %d\n", bestMove, game.pv[0][0], alpha)
			game.nodes += moveCount
			if engine.uci { // Report alpha as score since we're returning alpha.
				engine.uciScore(depth, alpha, alpha, beta)
			}
			return alpha
		}

		if moveCount == 1 || score > alpha {
			bestMove = move
			game.saveBest(0, move)
			gen.scoreMove(depth, score)
			gen.rearrangeRootMoves()

			if moveCount > 1 {
				game.volatility++
			}

			alpha = max(score, alpha)
			if alpha >= beta {
				break // Tap out.
			}
			p.cache(bestMove, score, depth, ply(), cacheBeta)
		} else {
			gen.scoreMove(depth, -depth)
		}
	}


	if moveCount == 0 {
		score = 0 // <-- Stalemate.
		if inCheck {
			score = -Checkmate
		}
		if engine.uci {
			engine.uciScore(depth, score, alpha, beta)
		}
		return
	}

	game.nodes += moveCount
	if score >= beta && !inCheck {
		game.saveGood(depth, bestMove)
	}
	score = alpha

	p.cacheDelta(bestMove, score, depth, ply(), rootAlpha, beta)

	if engine.uci {
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

	gen := NewGen(p, depth).generateAllMoves()
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if !gen.isValid(move) {
			continue
		}
		position := p.makeMove(move)
		total += position.Perft(depth - 1)
		position.undoLastMove()
	}
	return
}
