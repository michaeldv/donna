// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

type Score struct {
	midgame  int
	endgame  int
}

type Evaluator struct {
        stage     int
        midgame   int
        endgame   int
        position  *Position
}

func (p *Position) Evaluate() (score int) {
        evaluator := &Evaluator{ 0, 0, 0, p }

        evaluator.analyzeMaterial()
        evaluator.analyzeCoordination()
        evaluator.analyzePawnStructure()
        // evaluator.analyzePassedPawns()
        // evaluator.analyzeKingSafety()

        score = (evaluator.midgame * p.stage + evaluator.endgame * (256 - p.stage)) / 256
        return
}

func (e *Evaluator) analyzeMaterial() {
        color, opposite := e.position.color, e.position.color^1

        for _,piece := range []int{ PAWN, KNIGHT, BISHOP, ROOK, QUEEN } {
                count := e.position.count[Piece(piece|color)] - e.position.count[Piece(piece|opposite)]
                midgame, endgame := Piece(piece).value()
                e.midgame += midgame * count
                e.endgame += endgame * count
        }
}

func (e *Evaluator) analyzeCoordination() {
        var moves, attacks [2]int
        var bonus [2]Score

        for square, piece := range e.position.pieces {
                if piece == 0 {
                        continue
                }
                color := piece.color()

                // Mobility: how many moves are available to squares not attacked by
                // the opponent?
                targets := e.position.targets[square]
                moves[color] += targets.Intersect(e.position.attacks[color^1]).Count()

                // Agressivness: how many opponent's pieces are being attacked?
                targets = e.position.targets[square]
                attacks[color] += targets.Intersect(e.position.board[color^1]).Count()

                // Calculate bonus or penalty for a piece being at the given square.
                midgame, endgame := piece.bonus(flip[color][square])
                bonus[color].midgame += midgame
                bonus[color].endgame += endgame
        }

        color, opposite := e.position.color, e.position.color^1
        e.midgame += bonus[color].midgame - bonus[opposite].midgame
        e.endgame += bonus[color].endgame - bonus[opposite].endgame

        mobility := moves[color] - moves[opposite]
        e.midgame += mobility * movesAvailable.midgame
        e.endgame += mobility * movesAvailable.endgame

        aggression := attacks[color] - attacks[opposite]
        e.midgame += aggression * attackForce.midgame
        e.endgame += aggression * attackForce.endgame

        if bishops := e.position.count[Bishop(color)]; bishops >= 2 {
                e.midgame += bishopPair.midgame
                e.endgame += bishopPair.endgame
        }
        if bishops := e.position.count[Bishop(opposite)]; bishops >= 2 {
                e.midgame -= bishopPair.midgame
                e.endgame -= bishopPair.endgame
        }
}

func (e *Evaluator) analyzePawnStructure() {
        var bonus, penalty [2]Score
        pawn := [2]Piece{ Pawn(WHITE), Pawn(BLACK) }

        for color := WHITE; color <= BLACK; color++ {
                var doubled [8]int // Number of doubled pawns in each column.

                pawns := e.position.outposts[pawn[color]]
                for pawns.IsNotEmpty() {
                        square := pawns.FirstSet()
                        column := Col(square)
                        //
                        // Count doubled pawns in the column as they carry a penalty.
                        //
                        doubled[column] = (maskFile[column] & e.position.outposts[pawn[color]]).Count()
                        //
                        // The pawn is passed if a) there are no enemy pawns in the
                        // same and adjacent columns; and b) there is no same color
                        // pawns in front of us.
                        //
                        if (maskPassed[color][square] & e.position.outposts[pawn[color^1]]).IsEmpty() &&
                           (maskInFront[color][square] & e.position.outposts[pawn[color]]).IsEmpty() {
                                   bonus[color].midgame += bonusPassedPawn[0][flip[color][square]]
                                   bonus[color].endgame += bonusPassedPawn[1][flip[color][square]]
                        }
                        pawns.Clear(square)
                }
                //
                // Penalties for doubled pawns.
                //
                for i := 0;  i < len(doubled); i++ {
                        if doubled[i] > 0 {
                                penalty[color].midgame += (doubled[i] - 1) * penaltyDoubledPawn[0][i]
                                penalty[color].endgame += (doubled[i] - 1) * penaltyDoubledPawn[1][i]
                        }
                }
        }

        color, opposite := e.position.color, e.position.color^1
        e.midgame += bonus[color].midgame   - bonus[opposite].midgame
        e.endgame += bonus[color].endgame   - bonus[opposite].endgame
        e.midgame += penalty[color].midgame - penalty[opposite].midgame
        e.endgame += penalty[color].endgame - penalty[opposite].endgame
}

func (e *Evaluator) analyzePassedPawns() {
}

func (e *Evaluator) analyzeKingSafety() {
}
