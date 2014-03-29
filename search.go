// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

// Root node search.
func (p *Position) searchRoot(depth int) (bestMove Move, bestScore int) {
        gen := p.rootMoveGen(depth)

        if gen.onlyMove() {
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
                        position.saveBest(Ply(), move)
                        if bestScore > alpha {
                                alpha = bestScore
                                bestMove = move
                                // if alpha > 32000 { // <-- Not in puzzle solving mode.
                                //         break
                                // }
                        }
                }
        } // next move.

        // if bestScore < -32000 || bestScore > 32000 { // <-- Not in puzzle solving mode.
        //         break // from next depth loop.
        // }

        return
}

// Initializes move generator for the initial step of iterative deepening (depth == 1)
// and returns existing generator for subsequent iterations (depth > 1).
func (p *Position) rootMoveGen(depth int) (gen *MoveGen) {
        if depth > 1 {
                return p.UseMoveGen(0).rank()
        }

        gen = p.StartMoveGen(0)
        if p.isInCheck(p.color) {
                gen.GenerateEvasions()
        } else {
                gen.GenerateMoves()
        }

        // Get rid of invalid moves so that we don't do it on each iteration.
        gen.validOnly(p).quickRank() // No best move/killers yet.

        // fmt.Printf("depth: %d, node: %d, killers %s/%s\n", depth, node, p.killers[0], p.killers[1])
        // for i := gen.head; i < gen.tail; i++ {
        //         fmt.Printf("\t%d: %s (%d)\n", i, gen.list[i].move, gen.list[i].score)
        // }

        return
}

// Helps with testing root search by initializing move genarator at given depth and
// bypassing iterative deepening altogether.
func (p *Position) search(depth int) Move {
        gen := p.StartMoveGen(0)
        if p.isInCheck(p.color) {
                gen.GenerateEvasions()
        } else {
                gen.GenerateMoves()
        }

        // Get rid of invalid moves just like when starting iterative deepening.
        gen.validOnly(p)

        move, _ := p.searchRoot(depth)

        return move
}
