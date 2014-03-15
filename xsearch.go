// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`fmt`)

const maxDepth = 4

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
                        fmt.Printf("depth: %d, move: %s\n", depth, move)
                        if position := p.MakeMove(move); position != nil {
                                inCheck := p.isInCheck(p.color)
                                reducedDepth := depth - 1
                                if inCheck {
                                        reducedDepth++
                                }

                                moveScore := 0
                                if bestScore != -Checkmate && reducedDepth > 0 {
                                        if inCheck {
                                                moveScore = -p.xSearchInCheck(-alpha, reducedDepth)
                                        } else {
                                                moveScore = -p.xSearchWithZeroWindow(-alpha, reducedDepth)
                                        }
                                        if moveScore > alpha {
                                                moveScore = -p.xSearchPrincipal(-Checkmate, -alpha, reducedDepth)
                                        }
                                } else {
                                        moveScore = -p.xSearchPrincipal(-Checkmate, Checkmate, reducedDepth)
                                }

                                position.TakeBack(move)
                                if moveScore > bestScore {
                                        bestScore = moveScore
                                        if bestScore > alpha {
                                                alpha = bestScore
                                                fmt.Printf("make first => depth: %d, move: %s\n", depth, move)
                                                gen.makeFirst()
                                                if alpha > 32000 {
                                                        break
                                                }
                                        }
                                        //>> printBestLine(depth, bestScore, gen.list[0])
                                }
                        } // if position
                } // next move.
                //>> printBestLine(depth, bestScore, gen.list[0])
                //>> prevScore = bestScore

                if bestScore < -32000 || bestScore > 32000 {
                    break // from next depth loop.
            }

            gen.head = 0
        } // next depth.

        return gen.list[0].move
}
