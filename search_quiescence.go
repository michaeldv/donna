// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Quiescence search.
func (p *Position) searchQuiescence(alpha, beta, depth, iteration int) (score int) {
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

	isPrincipal := (beta - alpha > 1)

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

	inCheck := p.isInCheck(p.color)
	if !inCheck {
		staticScore = p.Evaluate()
		if staticScore >= beta {
			p.cache(Move(0), staticScore, 0, ply, cacheBeta)
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
	}
	gen.quickRank()

	cacheFlags := cacheAlpha
	moveCount, bestMove, king := 0, Move(0), int(p.king[p.color^1])
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if !inCheck && !move.piece().isKing() {
			// Prune useless captures that are not checks.
			useless := !isPrincipal && !move.isPromo() && staticScore + pieceValue[move.capture()] + 72 < alpha
			if p.targetsFor(move.to(), move.piece()).off(king) && (useless || p.exchange(move) < 0) {
				continue
			}
		}

		if !gen.isValid(move) {
			continue
		}

		position := p.makeMove(move)
		moveCount++
		score = -position.searchQuiescence(-beta, -alpha, depth, iteration+1)
		position.undoLastMove()

		if score > alpha {
			alpha = score
			bestMove = move
			// if isPrincipal {
			// 	game.saveBest(ply, bestMove)
			// }
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

	if !inCheck && iteration < 1 {
		gen = NewGen(p, ply).generateChecks().quickRank()
		for move := gen.NextMove(); move != 0; move = gen.NextMove() {
			if p.exchange(move) < 0 || !gen.isValid(move) {
				continue
			}

			position := p.makeMove(move)
			moveCount++
			score = -position.searchQuiescence(-beta, -alpha, depth, iteration+1)
			position.undoLastMove()

			if score > alpha {
				alpha = score
				bestMove = move
				// if isPrincipal {
				// 	game.saveBest(ply, bestMove)
				// }
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
	}

	game.qnodes += moveCount

	score = alpha
	if inCheck && moveCount == 0 {
		score = -Checkmate + ply
	}
	p.cache(bestMove, score, depth, ply, cacheFlags)

	return
}
