// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`fmt`)

const maxDepth = 5

// Root node search.
func (p *Position) xSearch() Move {
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

        //>> prevScore := p.Evaluate()

        rootNode = node
        for depth := 1; depth <= maxDepth; depth++ {
                alpha, bestScore := -Checkmate, -Checkmate

                for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                        if position := p.MakeMove(move); position != nil {
                                fmt.Printf("%*sroot/%s> depth: %d, ply: %d, move: %s\n", Ply()*2, ` `, C(p.color), depth, Ply(), move)
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

                                position.TakeBack(move)
                                if moveScore > bestScore {
                                        bestScore = moveScore
                                        if bestScore > alpha {
                                                alpha = bestScore
                                                fmt.Printf("make first => depth: %d, move: %s\n", depth, move)
                                                gen.makeFirst()
                                                // if alpha > 32000 { // Not in puzzle solving mode.
                                                //         break
                                                // }
                                        }
                                        fmt.Printf("-> %d) %d %s\n", depth, bestScore, gen.list[0].move)
                                }
                        } // if position
                } // next move.
                fmt.Printf("=> %d) %d %s\n", depth, bestScore, gen.list[0].move)
                //>> prevScore = bestScore

                // if bestScore < -32000 || bestScore > 32000 { // Not in puzzle solving mode.
                //         break // from next depth loop.
                // }

            gen.head = 0
        } // next depth.

        return gen.list[0].move
}
