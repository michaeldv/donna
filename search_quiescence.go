// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

// Quiescence search.
func (p *Position) searchQuiescence(alpha, beta int) int {
        p.game.bestLength[Ply()] = Ply()
        return p.quiescence(alpha, beta, false)
}

func (p *Position) quiescence(alpha, beta int, capturesOnly bool) int {
        p.game.qnodes++
        if p.isRepetition() {
                return 0
        }

        bestScore := p.Evaluate()
        if Ply() > MaxDepth {
                return bestScore
        }

        if bestScore > alpha {
                if bestScore >= beta {
                        return bestScore//beta
                }
                alpha = bestScore
        }

        gen := NewGen(p, Ply()).GenerateCaptures().quickRank()
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        //Log("%*squie/%s> ply: %d, move: %s\n", Ply()*2, ` `, C(p.color), Ply(), move)
                        moveScore := 0
                        if position.isInCheck(position.color) {
                                moveScore = -position.quiescenceInCheck(-beta, -alpha)
                        } else {
                                moveScore = -position.quiescence(-beta, -alpha, true)
                        }

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
        }

        if capturesOnly {
                return bestScore
        }

        gen = NewGen(p, Ply()).GenerateChecks().quickRank()
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        //Log("%*squix/%s> ply: %d, move: %s\n", Ply()*2, ` `, C(p.color), Ply(), move)
                        moveScore := -position.quiescenceInCheck(-beta, -alpha)

                        position.TakeBack(move)
                        if moveScore > bestScore {
                                p.game.saveBest(Ply(), move)
                                if moveScore > alpha {
                                        if moveScore >= beta {
                                                return moveScore
                                        }
                                        alpha = moveScore
                                }
                                bestScore = moveScore
                        }
                }
        }

        return bestScore
}

// Quiescence search (in check).
func (p *Position) quiescenceInCheck(alpha, beta int) int {
        if p.isRepetition() {
                return 0
        }

        bestScore := Ply() - Checkmate
        if bestScore >= beta {
                return bestScore//beta
        }

        gen := NewGen(p, Ply()).GenerateEvasions().quickRank()
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        //Log("%*squic/%s> ply: %d, move: %s\n", Ply()*2, ` `, C(p.color), Ply(), move)
                        moveScore := 0
                        if position.isInCheck(position.color) {
                                moveScore = -position.quiescenceInCheck(-beta, -alpha)
                        } else {
                                moveScore = -position.quiescence(-beta, -alpha, true)
                        }

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
        }

        return bestScore
}
