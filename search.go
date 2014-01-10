// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

func (p *Position) search(depth, ply int, alpha, beta int) int {
        Log("\nsearch(depth: %d/%d, color: %s, alpha: %d, beta: %d)\n", depth, ply, C(p.color), alpha, beta)
        p.game.nodes++
        if depth <= 0 && !p.inCheck {
                return p.quiescence(depth, ply, alpha, beta)
        }

        if p.isRepetition() {
                return 0
        }

        if ply > MaxPly - 2 {
                bestlen[ply] = ply
                return p.Evaluate()
        }

	// Checkmate pruning.
	if Checkmate - ply <= alpha {
		return alpha
	} else if -Checkmate + ply >= beta {
		return beta
	}

        // Null move pruning. TODO: skip it if we're following principal variation.
        if !p.inCheck && p.board[p.color].count() > 5 && p.Evaluate() >= beta {
                p.color ^= 1
                score := -p.search(depth - 4, ply + 1, -beta, -beta + 1)
                p.color ^= 1
                if score >= beta {
                        return score
                }
        }

        moves, movesMade := p.Moves(ply), 0
        for i, move := range moves {
                if position := p.MakeMove(move); !position.isCheck(p.color) {
                        movesMade++
                        score := -position.search(depth - 1, ply + 1, -beta, -alpha)
                        Log("Move %d/%d: %s (%d): score: %d, alpha: %d, beta: %d\n", i+1, len(moves), C(p.color), depth, score, alpha, beta)

                        if score >= beta {
                                if !p.inCheck && move.captured == 0 && (killer[ply][0] == nil || !move.is(killer[ply][0])) {
                                        killer[ply][1] = killer[ply][0]
                                        killer[ply][0] = move
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
                if p.inCheck {
                        Lop("Checkmate")
                        return -Checkmate + ply
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
        p.game.qnodes++

        if p.isRepetition() {
                return 0
        }

        if ply > MaxPly - 2 {
                return p.Evaluate()
        }

	// Checkmate pruning.
	if Checkmate - ply <= alpha {
		return alpha
	} else if -Checkmate + ply >= beta {
		return beta
	}

        if p.inCheck {
                return p.quiescenceInCheck(depth, ply, alpha, beta)
        }
        return p.quiescenceStayPat(depth, ply, alpha, beta)
}

func (p *Position) quiescenceInCheck(depth, ply int, alpha, beta int) int {
        score, bestScore := 0, -Checkmate
        quietAlpha, quietBeta := alpha, beta

        moves, movesMade := p.Moves(ply), 0 // TODO: check evasions only.
        for i, move := range moves {
                if position := p.MakeMove(move); !position.isCheck(p.color) {
                        Log("%d out of %d: evasion %s for %s\n", i, len(moves), move, C(move.piece.color()))

                        movesMade++
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

        if movesMade == 0 {
                bestlen[ply] = ply
                return -Checkmate + ply
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

        moves := p.Captures(ply) // TODO: followed by quiet checks.
        for i, move := range moves {
                if position := p.MakeMove(move); !position.isCheck(p.color) {
                        Log("%d out of %d: capture %s for %s\n", i, len(moves), move, C(move.piece.color()))

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
