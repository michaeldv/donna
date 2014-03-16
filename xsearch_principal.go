// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`fmt`)

// Search principal variation.
func (p *Position) xSearchPrincipal(alpha, beta, depth int) int {
        if depth == 0 {
                return p.xSearchQuiescence(alpha, beta, 1)
        }

        if Ply() > maxDepth {
                score := p.Evaluate()
                fmt.Printf("%*sprin/%s> DEPTH: %d, PLY: %d, SCORE: %d\n", depth*2, ` `, C(p.color), depth, Ply(), score)
                return score
        }

        if p.isRepetition() {
                return 0
        }

        bestScore := Ply() - Checkmate
        if bestScore >= beta {
                return bestScore
        }

        gen := p.StartMoveGen(Ply())
        if !p.isInCheck(p.color) {
                gen.GenerateMoves()
        } else {
                gen.GenerateEvasions()
        }
        gen.rank()

        firstMove := true
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        fmt.Printf("%*sprin/%s> depth: %d, ply: %d, move: %s\n", depth*2, ` `, C(p.color), depth, Ply(), move)
                        inCheck := p.isInCheck(p.color)
                        reducedDepth := depth - 1
                        if inCheck {
                                reducedDepth++
                        }

                        moveScore := 0
                        if firstMove {
                                moveScore = -position.xSearchPrincipal(-beta, -alpha, reducedDepth)
                        } else {
                                if reducedDepth == 0 {
                                        moveScore = -position.xSearchQuiescence(-alpha - 1, -alpha, reducedDepth)
                                } else if inCheck {
                                        moveScore = -position.xSearchInCheck(-alpha, reducedDepth)
                                } else {
                                        moveScore = -position.xSearchWithZeroWindow(-alpha, reducedDepth)
                                }
                                if moveScore > alpha {
                                        moveScore = -position.xSearchPrincipal(-beta, -alpha, reducedDepth)
                                }
                        }

                        firstMove = false
                        position.TakeBack(move)

                        if moveScore > bestScore {
                                if moveScore > alpha {
                                        if moveScore >= beta {
                                                return moveScore
                                        }
                                        alpha = moveScore
                                }
                                bestScore = moveScore
                        }
                }
        } // next move.

        if firstMove { // Checkmate
                if p.isInCheck(p.color) {
                        return bestScore
                } else { // Stalemate
                        return 0
                }
        }

        return bestScore
}
