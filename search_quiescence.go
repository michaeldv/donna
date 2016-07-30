// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
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
	isNull := p.isNull()
	isPrincipal := (beta - alpha > 1)
	if isPrincipal {
		game.pv[ply].size = 0 // Reset principal variation.
	}

	// Use fixed depth for caching.
	newDepth := let(inCheck || depth >= 0, 0, -1)

	// Probe cache.
	cachedMove := Move(0)
	cached := p.probeCache()
	if cached != nil {
		cachedMove = cached.move
		if int(cached.depth) >= newDepth {
			cachedScore := uncache(int(cached.score), ply)
			if !isPrincipal &&
			   ((cached.flags == cacheBeta  && cachedScore >= beta) ||
			   (cached.flags == cacheAlpha && cachedScore <= alpha)) {
				return cachedScore
			}
		}
	}

	if inCheck {
		p.score = Unknown
	} else {
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
		if p.score >= beta {
			return p.score
		}
		if isPrincipal {
			alpha = max(alpha, p.score)
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

	bestAlpha := alpha
	bestScore := let(p.score != Unknown, p.score, matedIn(ply))
	bestMove, moveCount := Move(0), 0
	for move := gen.nextMove(); move.some(); move = gen.nextMove() {
		capture := move.capture()
		if (!inCheck && capture.some() && p.exchange(move) < 0) || !move.isValid(p, gen.pins) {
			continue
		}

		position := p.makeMove(move)
		moveCount++; game.qnodes++
		giveCheck := position.isInCheck(position.color)

		// Prune useless captures -- but make sure it's not a capture move that checks.
		if !inCheck && !giveCheck && !isPrincipal && capture != 0 && !move.isPromo() && p.score + pieceValue[capture.id()] + 72 < alpha {
			position.undoLastMove()
			continue
		}
		score = -position.searchQuiescence(-beta, -alpha, depth - 1, giveCheck)
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
					p.cache(move, score, newDepth, ply, cacheBeta)
					return score
				}
			}
		}
	}

	score = let(inCheck && moveCount == 0, matedIn(ply), bestScore)

	cacheFlags := cacheAlpha
	if isPrincipal && score > bestAlpha {
		cacheFlags = cacheExact
	}
	p.cache(bestMove, score, newDepth, ply, cacheFlags)

	return score
}
