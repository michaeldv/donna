// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

func (p *Position) alphaBeta(depth, ply int, alpha, beta int) int {
        Log("\nalphaBeta(depth: %d/%d, color: %s, alpha: %d, beta: %d)\n", depth, ply, C(p.color), alpha, beta)
        if depth == 0 {
                // if !p.inCheck {
                //         return p.quietAlphaBeta(depth, ply, alpha, beta)
                // }
                return p.Evaluate()
        }

	// Checkmate pruning.
	if CHECKMATE - ply <= alpha {
		return alpha
	} else if -CHECKMATE + ply >= beta {
		return beta
	}

        moves := p.Moves()
        nodes := p.game.nodes
        for i, move := range moves {
		position := p.MakeMove(move)
                if !position.isCheck(p.color) {
                        p.game.nodes++
                        score := -position.alphaBeta(depth - 1, ply + 1, -beta, -alpha)
                        Log("Move %d/%d: %s (%d): score: %d, alpha: %d, beta: %d\n", i+1, len(moves), C(p.color), depth, score, alpha, beta)
                        if score >= beta {
                                Log("\n  Done at depth %d after move %d out of %d for %s\n", depth, i+1, len(moves), C(p.color))
                                Log("  Searched %v\n", moves[:i+1])
                                Log("  Skipping %v\n", moves[i+1:])
                                Log("  Picking %v\n\n", move)
                                return score
                        }
                        if score > alpha {
                                alpha = score
                                p.saveBest(ply, move)
                        }
                }
        }

        if nodes == p.game.nodes { // No moves were available.
                if p.inCheck {
                        Lop("Checkmate")
                        return -CHECKMATE + ply
                } else {
                        Lop("Stalemate")
                        alpha = 0.0
                }
        }

        Log("End of AlphaBeta(depth: %d/%d, color: %s, alpha: %d, beta: %d) => %d\n", depth, ply, C(p.color), alpha, beta, alpha)
	return alpha
}

func (p *Position) quietAlphaBeta(depth, ply int, alpha, beta int) int {
        Log("\nquietAlphaBeta(depth: %d/%d, color: %s, alpha: %d, beta: %d)\n", depth, ply, C(p.color), alpha, beta)

	score := p.Evaluate()
	if score >= beta || ply > 14 {
		return score
        }
	if score > alpha {
                bestlen[ply] = ply
		alpha = score
        }

        moves := p.Captures()
        for i, move := range moves {
                if position := p.MakeMove(move); !position.isCheck(p.color) {
                        Log("Making capture %s for %s\n", move, C(move.Piece.Color()))
                        p.game.nodes++
                        score = -position.MakeMove(move).quietAlphaBeta(depth - 1, ply + 1, -beta, -alpha)
                        Log("Capture %d/%d: %s (%d): score: %d, alpha: %d, beta: %d\n", i+1, len(moves), C(p.color), depth, score, alpha, beta)
                        if score >= beta {
                                Log("\n  Done at depth %d after move %d out of %d for %s\n", depth, i+1, len(moves), C(p.color))
                                Log("  Searched %v\n", moves[:i+1])
                                Log("  Skipping %v\n", moves[i+1:])
                                Log("  Picking %v\n\n", move)
                                return score
                        }
                        if score > alpha {
                                alpha = score
                                p.saveBest(ply, move)
                        }
                }
        }

        Log("End of quietAlphaBeta(depth: %d/%d, color: %s, alpha: %d, beta: %d) => %d\n", depth, ply, C(p.color), alpha, beta, alpha)
	return alpha
}
