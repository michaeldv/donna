// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

func (p *Position) search(depth, ply int, alpha, beta int) int {
        Log("\nsearch(depth: %d/%d, color: %s, alpha: %d, beta: %d)\n", depth, ply, C(p.color), alpha, beta)
        if depth <= 0 && !p.inCheck {
                return p.quiescence(depth, ply, alpha, beta)
        }

        if ply > 14 {
                bestlen[ply] = ply
                return p.Evaluate()
        }

	// Checkmate pruning.
	if CHECKMATE - ply <= alpha {
		return alpha
	} else if -CHECKMATE + ply >= beta {
		return beta
	}

        // Null move pruning. TODO: skip it if we're following principal variation.
        if !p.inCheck && p.board[p.color].Count() > 5 && p.Evaluate() >= beta {
                p.color ^= 1
                score := -p.search(depth - 4, ply + 1, -beta, -beta + 1)
                p.color ^= 1
                if score >= beta {
                        return score
                }
        }

        moves := p.Moves(ply)
        nodes := p.game.nodes
        for i, move := range moves {
                if position := p.MakeMove(move); !position.isCheck(p.color) {
                        p.game.nodes++
                        score := -position.search(depth - 1, ply + 1, -beta, -alpha)
                        Log("Move %d/%d: %s (%d): score: %d, alpha: %d, beta: %d\n", i+1, len(moves), C(p.color), depth, score, alpha, beta)

                        if score >= beta {
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
                        return 0
                }
        }

        Log("End of search(depth: %d/%d, color: %s, alpha: %d, beta: %d) => %d\n", depth, ply, C(p.color), alpha, beta, alpha)
        return alpha
}

func (p *Position) quiescence(depth, ply int, alpha, beta int) int {
        Log("\nquiescence(depth: %d/%d, color: %s, alpha: %d, beta: %d)\n", depth, ply, C(p.color), alpha, beta)

        if ply > 14 {
                return p.Evaluate()
        }

	// Checkmate pruning.
	if CHECKMATE - ply <= alpha {
		return alpha
	} else if -CHECKMATE + ply >= beta {
		return beta
	}

        if p.inCheck {
                return p.quiescenceInCheck(depth, ply, alpha, beta)
        }
        return p.quiescenceStayPat(depth, ply, alpha, beta)
}

func (p *Position) quiescenceInCheck(depth, ply int, alpha, beta int) int {
        score, bestScore := 0, -CHECKMATE
        quietAlpha, quietBeta := alpha, beta

        moves := p.Moves(ply) // TODO: check evasions only.
        qnodes := p.game.qnodes
        for i, move := range moves {
                if position := p.MakeMove(move); !position.isCheck(p.color) {
                        Log("%d out of %d: evasion %s for %s\n", i, len(moves), move, C(move.piece.color()))
                        p.game.qnodes++

                        score = -position.quiescence(depth - 1, ply + 1, -quietBeta, -quietAlpha)
                        if alpha + 1 != beta && score > quietAlpha && quietAlpha + 1 == quietBeta {
                                score = -position.quiescence(depth - 1, ply + 1, -beta, -quietAlpha)
                        }

                        if score >= beta {
                                return score
                        }
                        if score > bestScore {
                                bestScore = score
                                if score > quietAlpha {
                                        quietAlpha = score
                                        p.saveBest(ply, move)
                                }
                        }
                        quietBeta = quietAlpha + 1
                }
        }

        if qnodes == p.game.qnodes {
                bestlen[ply] = ply
                return -CHECKMATE + ply
        }

        Log("End of quiescenceInCheck(depth: %d/%d, color: %s, alpha: %d, beta: %d) => %d\n", depth, ply, C(p.color), alpha, beta, alpha)
        return bestScore
}

func (p *Position) quiescenceStayPat(depth, ply int, alpha, beta int) int {
        score := p.Evaluate()
        if score >= beta {
                return score
        }

        bestScore, quietAlpha, quietBeta := score, alpha, beta
        if score > alpha {
                bestlen[ply] = ply
                quietAlpha = score
        }

        moves := p.Captures() // TODO: followed by quiet checks.
        for i, move := range moves {
                if position := p.MakeMove(move); !position.isCheck(p.color) {
                        Log("%d out of %d: capture %s for %s\n", i, len(moves), move, C(move.piece.color()))
                        p.game.qnodes++

                        score = -position.quiescence(depth - 1, ply + 1, -quietBeta, -quietAlpha)
                        if quietAlpha + 1 != beta && score > quietAlpha && quietAlpha + 1 == quietBeta {
                                score = -position.quiescence(depth - 1, ply + 1, -beta, -quietAlpha)
                        }

                        if score >= beta {
                                return score
                        }
                        if score > bestScore {
                                bestScore = score
                                if score > quietAlpha {
                                        quietAlpha = score
                                        p.saveBest(ply, move)
                                }
                        }
                        quietBeta = quietAlpha + 1
                }
        }

        Log("End of quiescenceStayPat(depth: %d/%d, color: %s, alpha: %d, beta: %d) => %d\n", depth, ply, C(p.color), alpha, beta, alpha)
        return bestScore
}
