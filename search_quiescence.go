// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Quiescence search.
func (p *Position) searchQuiescence(alpha, beta, depth int) int {
	return p.searchQuiescenceWithFlag(alpha, beta, depth, false)
}

func (p *Position) searchQuiescenceWithFlag(alpha, beta, depth int, capturesOnly bool) (score int) {
	ply := Ply()

	// Reset principal variation.
	game.pv[ply] = game.pv[ply][:0]

	// Return if it's time to stop search.
	if ply >= MaxPly || engine.clock.halt {
		return p.Evaluate()
	}

	// Insufficient material and repetition/perpetual check pruning.
	if p.insufficient() || p.repetition() {
		return 0
	}


	// Probe cache.
	isPrincipal := (beta - alpha > 1)
	cacheFlags := uint8(cacheAlpha)
	if cached := p.probeCache(); cached != nil {
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
				return score
			}
		}
	}

	inCheck := p.isInCheck(p.color)
	staticScore := alpha
	if !inCheck {
		staticScore = p.Evaluate()
		alpha = max(alpha, staticScore)
		if alpha >= beta {
			return beta
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

	moveCount, bestMove := 0, Move(0)
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if !gen.isValid(move) || (!inCheck && p.exchange(move) < 0) {
			continue
		}

		// Check if the move is an useless capture.
		useless := !inCheck && !isPrincipal && !move.isPromo() && staticScore + pieceValue[move.capture()] + 64 < alpha

		position := p.MakeMove(move)
		moveCount++

		// Prune useless captures -- but make sure it's not a capture move
		// that checks.
		if useless && !position.isInCheck(position.color) {
			position.UndoLastMove()
			continue
		}

		score = -position.searchQuiescenceWithFlag(-beta, -alpha, depth, true)
		position.UndoLastMove()

		if score > alpha {
			alpha = score
			bestMove = move
			cacheFlags = cacheExact

			if alpha >= beta {
				cacheFlags = cacheBeta
				break
			}
			game.saveBest(ply, move)
		}
	}

	if !inCheck && !capturesOnly {
		gen = NewGen(p, Ply()).generateChecks().quickRank()
		for move := gen.NextMove(); move != 0; move = gen.NextMove() {
			if !gen.isValid(move) || p.exchange(move) < 0 {
				continue
			}

			position := p.MakeMove(move)
			moveCount++
			score = -position.searchQuiescenceWithFlag(-beta, -alpha, depth, false)
			position.UndoLastMove()

			if engine.clock.halt {
				game.qnodes += moveCount
				//Log("searchQui at %d (%s): move %s (%d) score %d alpha %d\n", depth, C(p.color), move, moveCount, score, alpha)
				return alpha
			}

			if score > alpha {
				alpha = score
				bestMove = move
				cacheFlags = cacheExact

				if alpha >= beta {
					cacheFlags = cacheBeta
					break
				}
				game.saveBest(ply, move)
			}
		}
	}

	game.qnodes += moveCount

	score = alpha
	if inCheck && moveCount == 0 {
		score = -Checkmate + ply
	}
	p.cache(bestMove, score, depth, cacheFlags)

	return
}
