// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

//import(`fmt`)

func (p *Position) search(depth, ply int, alpha, beta int) int {
        // Log("\nsearch(depth: %d/%d, color: %s, alpha: %d, beta: %d)\n", depth, ply, C(p.color), alpha, beta)
        p.game.nodes++
        if depth <= 0 && !p.isInCheck(p.color) {
                return p.quiescence(depth, ply, alpha, beta)
        }

        if p.isRepetition() {
                p.game.bestLength[ply] = -1
                return 0
        }

        if ply > MaxPly - 2 {
                p.game.bestLength[ply] = ply
                return p.Evaluate()
        }

	// Checkmate pruning.
	if Checkmate - ply <= alpha {
		return alpha
	} else if -Checkmate + ply >= beta {
		return beta
	}

        // Null move pruning. TODO: skip it if we're following principal variation.
        if !p.isInCheck(p.color) && p.outposts[p.color].count() > 5 && p.Evaluate() >= beta {
                p.color ^= 1
                score := -p.search(depth - 4, ply + 1, -beta, -beta + 1)
                p.color ^= 1
                if score >= beta {
                        return score
                }
        }

        gen := p.StartMoveGen(ply).GenerateMoves().rank()
        movesMade := 0
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        movesMade++
                        score := -position.search(depth - 1, ply + 1, -beta, -alpha)
                        position.TakeBack(move)
                        // Log("Move %d: %s (%d): score: %d, alpha: %d, beta: %d\n", movesMade, C(p.color), depth, score, alpha, beta)

                        if score >= beta {
                                if !p.isInCheck(p.color) && move.capture() == 0 && move != p.killers[0] {
                                        p.killers[1] = p.killers[0]
                                        p.killers[0] = move
                                }
                                return score
                        }
                        if score > alpha {
                                alpha = score
                                p.saveBest(ply, move)
                        }
                }
        }

        if movesMade == 0 { // No moves were available.
                p.game.bestLength[ply] = 0
                if p.isInCheck(p.color) {
                        Lop("Checkmate")
                        return -Checkmate + ply
                } else {
                        Lop("Stalemate")
                        return 0
                }
        }

        // Log("End of search(depth: %d/%d, color: %s, alpha: %d, beta: %d) => %d\n", depth, ply, C(p.color), alpha, beta, alpha)
        return alpha
}
