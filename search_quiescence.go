// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Quiescence search.
func (p *Position) searchQuiescence(alpha, beta, depth int) int {
	return p.searchQuiescenceWithFlag(alpha, beta, depth, false)
}

func (p *Position) searchQuiescenceWithFlag(alpha, beta, depth int, capturesOnly bool) (score int) {
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
	moveCount, bestMove := 0, Move(0)
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if (!inCheck && p.exchange(move) < 0) || !gen.isValid(move) {
			continue
		}

		// Check if the move is an useless capture.
		useless := !inCheck && !isPrincipal && !move.isPromo() && staticScore + pieceValue[move.capture()] + 72 < alpha

		position := p.makeMove(move)
		moveCount++

		// Prune useless captures -- but make sure it's not a capture move
		// that checks.
		if useless && !position.isInCheck(position.color) {
			position.undoLastMove()
			continue
		}

		score = -position.searchQuiescenceWithFlag(-beta, -alpha, depth, true)
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

	if !inCheck && !capturesOnly {
		gen = NewMoveGen(p).generateChecks().quickRank()
		for move := gen.NextMove(); move != 0; move = gen.NextMove() {
			if p.exchange(move) < 0 || !gen.isValid(move) {
				continue
			}

			position := p.makeMove(move)
			moveCount++
			score = -position.searchQuiescenceWithFlag(-beta, -alpha, depth, true)
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
	}

	game.qnodes += moveCount

	score = alpha
	if inCheck && moveCount == 0 {
		score = -Checkmate + ply
	}
	p.cache(bestMove, score, depth, ply, cacheFlags)

	return
}
