// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

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

        gen := p.StartMoveGen(ply).GenerateEvasions()
        movesMade := 0
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        movesMade++
                        Log("%d: evasion %s for %s\n", movesMade, move, C(move.color()))

                        score = -position.quiescence(depth - 1, ply + 1, -quietBeta, -quietAlpha)
                        if alpha + 1 != beta && score > quietAlpha && quietAlpha + 1 == quietBeta {
                                score = -position.quiescence(depth - 1, ply + 1, -beta, -quietAlpha)
                        }
                        position.TakeBack(move)

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
                p.game.bestLength[ply] = ply
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
                p.game.bestLength[ply] = ply
                quietAlpha = score
        }

        gen := p.StartMoveGen(ply).GenerateCaptures() // TODO: followed by quiet checks.
        movesMade := 0
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        movesMade++
                        Log("%d: capture %s for %s\n", movesMade, move, C(move.color()))

                        score = -position.quiescence(depth - 1, ply + 1, -quietBeta, -quietAlpha)
                        if quietAlpha + 1 != beta && score > quietAlpha && quietAlpha + 1 == quietBeta {
                                score = -position.quiescence(depth - 1, ply + 1, -beta, -quietAlpha)
                        }
                        position.TakeBack(move)

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
