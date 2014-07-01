// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

// Root node search.
func (p *Position) searchRoot(alpha, beta, depth int) (bestMove Move, bestScore int) {
	var score, reducedDepth int

	gen := NewRootGen(p, depth)
	if gen.onlyMove() {
		p.game.saveBest(Ply(), gen.list[0].move)
		return gen.list[0].move, p.Evaluate()
	}

	bestMove = gen.list[0].move
	bestScore = alpha

	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		position := p.MakeMove(move)
		//Log("%*sroot/%s> depth: %d, ply: %d, move: %s\n", Ply()*2, ` `, C(p.color), depth, Ply(), move)
		inCheck := position.isInCheck(position.color)
		if inCheck {
			reducedDepth = depth
		} else {
			reducedDepth = depth - 1
		}

		if bestScore != -Checkmate && reducedDepth > 0 {
			if inCheck {
				score = -position.searchInCheck(-alpha, reducedDepth)
			} else {
				score = -position.searchWithZeroWindow(-alpha, reducedDepth)
			}
			if score > alpha {
				score = -position.searchPrincipal(-Checkmate, -alpha, reducedDepth)
			}
		} else {
			score = -position.searchPrincipal(alpha, beta, reducedDepth)
		}

		position.TakeBack(move)
		if score > bestScore {
			bestScore = score
			position.game.saveBest(Ply(), move)
			if bestScore > alpha {
				alpha = bestScore
				bestMove = move
			}
		}
	} // next move.

	// fmt.Printf("depth: %d, node: %d\nbestline %v\nkillers %v\n", depth, node, p.game.pv, p.game.killers)

	return
}

// Helps with testing root search by initializing move genarator at given depth and
// bypassing iterative deepening altogether.
func (p *Position) search(depth int) Move {
	NewGen(p, 0).generateAllMoves().validOnly(p)
	move, _ := p.searchRoot(-Checkmate, Checkmate, depth)

	return move
}
