// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

// Root node search.
func (p *Position) search(alpha, beta, depth int) (score int) {
	ply := 0

	inCheck := p.isInCheck(p.color)
	cacheFlags := uint8(cacheAlpha)

	// Root move generator makes sure all generated moves are valid.
	gen := NewRootGen(p, depth)
	if depth == 1 {
		gen.generateRootMoves()
	} else {
		gen.rearrangeRootMoves()
	}

	moveCount, bestMove := 0, Move(0)
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		position := p.MakeMove(move)
		moveCount++

		// Search depth extension.
		newDepth := depth - 1
		if position.isInCheck(p.color^1) { // Give check.
			newDepth++
		}

		if moveCount == 1 {
			score = -position.searchTree(-beta, -alpha, newDepth)
		} else {
			score = -position.searchTree(-alpha - 1, -alpha, newDepth)
			if score > alpha && score < beta {
				score = -position.searchTree(-beta, -alpha, newDepth)
			}
		}
		position.UndoLastMove()

		if p.game.clock.halt {
			p.game.nodes += moveCount
			//Log("searchRoot: bestMove %s pv[0][0] %s alpha %d\n", bestMove, p.game.pv[0][0], alpha)
			return alpha
		}

		if moveCount == 1 {
			bestMove = move
			p.game.pv[ply] = p.game.pv[ply][:0]
			p.game.saveBest(ply, move)
		}

		if score > alpha {
			alpha = score
			bestMove = move
			cacheFlags = cacheExact
			p.game.saveBest(ply, move)

			if alpha >= beta {
				cacheFlags = cacheBeta
				break
			}
		}
	}

	p.game.nodes += moveCount

	if moveCount == 0 {
		if inCheck {
			alpha = -Checkmate + ply
		} else {
			alpha = 0
		}
	} else if score >= beta && !inCheck {
		p.game.saveGood(depth, bestMove)
	}

	score = alpha
	p.cache(bestMove, score, depth, cacheFlags)

	return
}

// Testing helper method to test root search.
func (p *Position) solve(depth int) Move {
	if depth != 1 {
		NewRootGen(p, 1).generateRootMoves()
	}
	p.search(-Checkmate, Checkmate, depth)
	return p.game.pv[0][0]
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
		position := p.MakeMove(move)
		total += position.Perft(depth - 1)
		position.UndoLastMove()
	}
	return
}
