// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Quiescence search.
func (p *Position) searchQuiescence(alpha, beta, depth int, inCheck bool) (score int) {
	ply := ply()

	// Return if it's time to stop search.
	if ply >= MaxPly || engine.clock.halt {
		return p.Evaluate()
	}

	// Insufficient material and repetition/perpetual check pruning.
	if p.fifty() || p.insufficient() || p.repetition() {
		return 0
	}

	// Checkmate distance pruning.
	alpha, beta = mateDistance(alpha, beta, ply)
	if alpha >= beta {
		return alpha
	}

	// If you pick up a starving dog and make him prosperous, he will not
	// bite you. This is the principal difference between a dog and a man.
        // ―- Mark Twain
	isPrincipal := (beta - alpha > 1)
	if isPrincipal {
		game.pv[ply].size = 0 // Reset principal variation.
	}

	// Use fixed depth for caching.
	newDepth := let(inCheck || depth >= 0, 0, -1)

	// Probe cache.
	cachedMove := Move(0)
	staticScore := matedIn(ply)
	if cached := p.probeCache(); cached != nil {
		cachedMove = cached.move
		if int(cached.depth) >= newDepth {
			staticScore = uncache(int(cached.score), ply)
			if !isPrincipal &&
			   ((cached.flags == cacheBeta  && staticScore >= beta) ||
			   (cached.flags == cacheAlpha && staticScore <= alpha)) {
				return staticScore
			}
		}
	}

	if !inCheck {
		staticScore = p.Evaluate()
		if staticScore >= beta {
			p.cache(Move(0), staticScore, newDepth, ply, cacheBeta)
			return staticScore
		}
		if isPrincipal {
			alpha = max(alpha, staticScore)
		}
	}

	// Generate check evasions or captures.
	gen := NewGen(p, ply)
	if inCheck {
		gen.generateEvasions().quickRank()
	} else {
		gen.generateCaptures()
		if depth >= 0 {
			gen.generateChecks()
		}
		gen.rank(cachedMove)
	}

	bestScore := staticScore
	moveCount, bestMove, bestAlpha := 0, Move(0), alpha
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		capture := move.capture()
		if (!inCheck && capture != 0 /*&& bestScore > Checkmate - 2 * MaxPly*/ && p.exchange(move) < 0) || !move.isValid(p, gen.pins) {
			continue
		}

		position := p.makeMove(move)
		moveCount++; game.qnodes += moveCount
		giveCheck := position.isInCheck(position.color)

		// Prune useless captures -- but make sure it's not a capture move that checks.
		if !inCheck && !giveCheck && !isPrincipal && capture != 0 && !move.isPromo() && staticScore + pieceValue[capture] + 72 < alpha {
			position.undoLastMove()
			continue
		}
		score = -position.searchQuiescence(-beta, -alpha, depth - 1, giveCheck)
		position.undoLastMove()

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
					p.cache(move, staticScore, depth, ply, cacheBeta)
					return score
				}
			}
		}

		if engine.clock.halt {
			return alpha
		}
	}

	score = let(inCheck && moveCount == 0, matedIn(ply), bestScore)

	cacheFlags := cacheAlpha
	if isPrincipal && score > bestAlpha {
		cacheFlags = cacheExact
	}
	p.cache(bestMove, staticScore, newDepth, ply, cacheFlags)

	return
}
