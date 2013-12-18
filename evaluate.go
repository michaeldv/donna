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

        evaluator.determineGameStage()
        evaluator.analyzeMaterial()
        evaluator.analyzeCoordination()
        // evaluator.analyzePawnStructure(p)
        // evaluator.analyzePassedPawns(p)
        // evaluator.analyzeKingSafety(p)

        score = (evaluator.midgame * evaluator.stage + evaluator.endgame * (256 - evaluator.stage)) / 256
        return
}

// Determine game stage by counting how many pieces are present on the board.
func (e *Evaluator) determineGameStage() {
        e.stage  =  2 * (e.position.count[Pawn(WHITE)]   + e.position.count[Pawn(BLACK)])
        e.stage +=  6 * (e.position.count[Knight(WHITE)] + e.position.count[Knight(BLACK)])
        e.stage += 12 * (e.position.count[Bishop(WHITE)] + e.position.count[Bishop(BLACK)])
        e.stage += 16 * (e.position.count[Rook(WHITE)]   + e.position.count[Rook(BLACK)])
        e.stage += 44 * (e.position.count[Queen(WHITE)]  + e.position.count[Queen(BLACK)])
}

func (e *Evaluator) analyzeMaterial() {
        color, opposite := e.position.color, e.position.color^1

        count := e.position.count[Pawn(color)] - e.position.count[Pawn(opposite)]
        e.endgame += valuePawn.endgame * count
        e.midgame += valuePawn.midgame * count

        count = e.position.count[Knight(color)] - e.position.count[Knight(opposite)]
        e.endgame += valueKnight.endgame * count
        e.midgame += valueKnight.midgame * count

        count = e.position.count[Bishop(color)] - e.position.count[Bishop(opposite)]
        e.endgame += valueBishop.endgame * count
        e.midgame += valueBishop.midgame * count

        count = e.position.count[Rook(color)] - e.position.count[Rook(opposite)]
        e.endgame += valueRook.endgame * count
        e.midgame += valueRook.midgame * count

        count = e.position.count[Queen(color)] - e.position.count[Queen(opposite)]
        e.endgame += valueQueen.endgame * count
        e.midgame += valueQueen.midgame * count
}

func (e *Evaluator) analyzeCoordination() {
        var moves, attacks [2]int
        var bonus [2]Score

        for square, piece := range e.position.pieces {
                if piece == 0 {
                        continue
                }
                color := piece.Color()

                // Mobility: how many moves are available to squares not attacked by
                // the opponent?
                targets := e.position.targets[square]
                moves[color] += targets.Intersect(e.position.attacks[color^1]).Count()

                // Agressivness: how many opponent's pieces are being attacked?
                targets = e.position.targets[square]
                attacks[color] += targets.Intersect(e.position.board[color^1]).Count()

                // Piece/square adjustments.
                if color == WHITE {
                        square = flip[square]
                }
                switch piece.Kind() {
                case PAWN:
                        bonus[color].midgame += bonusPawn[square]
                        bonus[color].endgame += bonusPawn[square]
                case KNIGHT:
                        bonus[color].midgame += bonusKnight[square]
                        bonus[color].endgame += bonusKnight[square]
                case BISHOP:
                        bonus[color].midgame += bonusBishop[square]
                        bonus[color].endgame += bonusBishop[square]
                // case ROOK:
                //         bonus = bonusRook[square]
                // case QUEEN:
                //         bonus = bonusQueen[square]
                case KING:
                        bonus[color].midgame += bonusKing[square]
                        bonus[color].endgame += bonusKingEndgame[square]
                }
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
}

func (e *Evaluator) analyzePawnStructure() {
        // for color := WHITE; color <= BLACK; color++ {
        //     outposts = p->outposts(Pawn(color))
        //
        //     for outposts.IsNotEmpty() {
        //             square := outposts.FirstSet()
        //             outposts.Clear(target)
        //     }
        // }
}

func (e *Evaluator) analyzePassedPawns() {
}

func (e *Evaluator) analyzeKingSafety() {
}
