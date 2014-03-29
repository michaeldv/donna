// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

// Search for the node in check.
func (p *Position) searchInCheck(beta, depth int) int {
        p.game.nodes++
        if p.isRepetition() {
                return 0
        }

        bestScore := Ply() - Checkmate
        if bestScore >= beta {
                return beta//bestScore
        }

        gen := p.StartMoveGen(Ply()).GenerateEvasions().quickRank()
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        //Log("%*schck/%s> depth: %d, ply: %d, move: %s\n", Ply()*2, ` `, C(p.color), depth, Ply(), move)
                        inCheck := position.isInCheck(position.color)
                        reducedDepth := depth - 1
                        if inCheck {
                                reducedDepth++
                        }

                        moveScore := 0
                        if reducedDepth == 0 {
                                moveScore = -position.searchQuiescence(-beta, 1 - beta)
                        } else if inCheck {
                                moveScore = -position.searchInCheck(1 - beta, reducedDepth)
                        } else {
                                moveScore = -position.searchWithZeroWindow(1 - beta, reducedDepth)
                        }

                        position.TakeBack(move)
                        if moveScore > bestScore {
                                if moveScore >= beta {
                                        return moveScore
                                }
                                bestScore = moveScore
                        }
                }
        }

        return bestScore
}
