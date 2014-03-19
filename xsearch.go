// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

const maxDepth = 7

// Root node search.
func (p *Position) xSearch(requestedDepth int) Move {
        gen := p.StartMoveGen(0)
        if !p.isInCheck(p.color) {
                gen.GenerateMoves()
        } else {
                gen.GenerateEvasions()
        }

        if gen.theOnlyMove() {
                return gen.list[0].move
        }
        gen.rank()

        rootNode = node
        bestMoveIdx := 0
        for depth := 1; depth <= Min(maxDepth, requestedDepth); depth++ {
                alpha, bestScore := -Checkmate, -Checkmate

                for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                        if position := p.MakeMove(move); position != nil {
                                //Log("%*sroot/%s> depth: %d, ply: %d, move: %s\n", Ply()*2, ` `, C(p.color), depth, Ply(), move)
                                inCheck := position.isInCheck(position.color)
                                reducedDepth := depth - 1
                                if inCheck {
                                        reducedDepth++
                                }

                                moveScore := 0
                                if bestScore != -Checkmate && reducedDepth > 0 {
                                        if inCheck {
                                                moveScore = -position.xSearchInCheck(-alpha, reducedDepth)
                                        } else {
                                                moveScore = -position.xSearchWithZeroWindow(-alpha, reducedDepth)
                                        }
                                        if moveScore > alpha {
                                                moveScore = -position.xSearchPrincipal(-Checkmate, -alpha, reducedDepth)
                                        }
                                } else {
                                        moveScore = -position.xSearchPrincipal(-Checkmate, Checkmate, reducedDepth)
                                }

                                position = position.TakeBack(move)
                                if moveScore > bestScore {
                                        bestScore = moveScore
                                        position.saveBest(Ply(), move)
                                        if bestScore > alpha {
                                                alpha = bestScore
                                                bestMoveIdx = gen.head - 1
                                                // if alpha > 32000 { // Not in puzzle solving mode.
                                                //         break
                                                // }
                                        }
                                }
                        } // if position
                } // next move.

                Log("=> %d) %d %s => %v\n", depth, bestScore, gen.list[0].move, p.game.bestLine[0][0 : p.game.bestLength[0]])

                // if bestScore < -32000 || bestScore > 32000 { // Not in puzzle solving mode.
                //         break // from next depth loop.
                // }

            gen.head = 0
        } // next depth.

        return gen.list[bestMoveIdx].move
}
