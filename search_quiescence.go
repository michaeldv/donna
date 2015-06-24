// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Quiescence search.
func (p *Position) searchQuiescence(alpha, beta, iteration int, inCheck bool) (score int) {
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

	// If you pick up a starving dog and make him prosperous, he will not
	// bite you. This is the principal difference between a dog and a man.
        // ―- Mark Twain
	isPrincipal := (beta - alpha > 1)

	// Use fixed depth for caching.
	depth := 0
	if !inCheck && iteration > 0 {
		depth--
	}

	// Probe cache.
	staticScore := alpha
	if cached := p.probeCache(); cached != nil {
		if int(cached.depth) >= depth {
			staticScore = uncache(int(cached.score), ply)
			if (cached.flags == cacheExact && isPrincipal) ||
			   (cached.flags == cacheBeta  && staticScore >= beta) ||
			   (cached.flags == cacheAlpha && staticScore <= alpha) {
				return staticScore
			}
		}
	}

	if !inCheck {
		staticScore = p.Evaluate()
		if staticScore >= beta {
			p.cache(Move(0), staticScore, depth, ply, cacheBeta)
			return staticScore
		}
		if isPrincipal {
			alpha = max(alpha, staticScore)
		}
	}

	// Generate check evasions or captures.
	gen := NewGen(p, ply)
	if inCheck {
		gen.generateEvasions()
	} else {
		gen.generateCaptures()
		if iteration < 1 {
			gen.generateChecks()
		}
	}
	gen.quickRank()

	cacheFlags := cacheAlpha
	moveCount, bestMove := 0, Move(0)
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		capture := move.capture()
		if (!inCheck && capture != 0 && p.exchange(move) < 0) || !gen.isValid(move) {
			continue
		}

		position := p.makeMove(move)
		moveCount++
		giveCheck := position.isInCheck(position.color)

		// Prune useless captures -- but make sure it's not a capture move that checks.
		if !inCheck && !giveCheck && !isPrincipal && capture != 0 && !move.isPromo() && staticScore + pieceValue[capture] + 72 < alpha {
			position.undoLastMove()
			continue
		}
		score = -position.searchQuiescence(-beta, -alpha, iteration + 1, giveCheck)
		position.undoLastMove()

		if score > alpha {
			alpha = score
			bestMove = move
			if alpha >= beta {
				p.cache(bestMove, score, depth, ply, cacheBeta)
				game.qnodes += moveCount
				return
			}
			cacheFlags = cacheExact
		}
		if engine.clock.halt {
			game.qnodes += moveCount
			return alpha
		}
	}

	game.qnodes += moveCount

	score = alpha
	if inCheck && moveCount == 0 {
		score = -Checkmate + ply
	}
	p.cache(bestMove, score, depth, ply, cacheFlags)

	return
}
