// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

// Root node search.
func (p *Position) searchRoot(depth int) (bestMove Move, bestScore int) {
        gen := NewRootGen(p, depth)
        if gen.onlyMove() {
                p.game.saveBest(Ply(), gen.list[0].move)
                return gen.list[0].move, p.Evaluate()
        }

        alpha := -Checkmate
        bestMove, bestScore = gen.list[0].move, -Checkmate

        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                position := p.MakeMove(move)
                //Log("%*sroot/%s> depth: %d, ply: %d, move: %s\n", Ply()*2, ` `, C(p.color), depth, Ply(), move)
                inCheck := position.isInCheck(position.color)
                reducedDepth := depth - 1
                if inCheck {
                        reducedDepth++
                }

                moveScore := 0
                if bestScore != -Checkmate && reducedDepth > 0 {
                        if inCheck {
                                moveScore = -position.searchInCheck(-alpha, reducedDepth)
                        } else {
                                moveScore = -position.searchWithZeroWindow(-alpha, reducedDepth)
                        }
                        if moveScore > alpha {
                                moveScore = -position.searchPrincipal(-Checkmate, -alpha, reducedDepth)
                        }
                } else {
                        moveScore = -position.searchPrincipal(-Checkmate, Checkmate, reducedDepth)
                }

                position.TakeBack(move)
                if moveScore > bestScore {
                        bestScore = moveScore
                        position.game.saveBest(Ply(), move)
                        if bestScore > alpha {
                                alpha = bestScore
                                bestMove = move
                                // if alpha > 32000 { // <-- Not in puzzle solving mode.
                                //         break
                                // }
                        }
                }
        } // next move.

        // fmt.Printf("depth: %d, node: %d\nbestline %v\nkillers %v\n", depth, node, p.game.bestLine, p.game.killers)

        return
}

// Helps with testing root search by initializing move genarator at given depth and
// bypassing iterative deepening altogether.
func (p *Position) search(depth int) Move {
        NewGen(p, 0).generateAllMoves().validOnly(p)
        move, _ := p.searchRoot(depth)

        return move
}
