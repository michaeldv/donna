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
        endgame   int
        midgame   int
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

        for square, piece := range e.position.pieces {
                if piece == 0 {
                        continue
                }
                color := piece.Color()

                // Mobility and agressivness.
                targets := e.position.targets[square]
                moves[color] += targets.Count()
                attacks[color] += targets.Intersect(e.position.board[color^1]).Count()

                // Piece/square adjustments.
                if color == WHITE {
                        square = flip[square]
                }
                switch piece.Kind() {
                case PAWN:
                        e.midgame += bonusPawn[square]
                        e.endgame += bonusPawn[square]
                case KNIGHT:
                        e.midgame += bonusKnight[square]
                        e.endgame += bonusKnight[square]
                case BISHOP:
                        e.midgame += bonusBishop[square]
                        e.endgame += bonusBishop[square]
                // case ROOK:
                //         bonus = bonusRook[square]
                // case QUEEN:
                //         bonus = bonusQueen[square]
                case KING:
                        e.midgame += bonusKing[square]
                        e.endgame += bonusKingEndgame[square]
                }
        }
        mobility := moves[e.position.color] - moves[e.position.color^1]
        if mobility != 0 {
                e.midgame += movesAvailable.midgame * mobility / Abs(mobility)
                e.endgame += movesAvailable.endgame * mobility / Abs(mobility)
        }
        aggression := attacks[e.position.color] - attacks[e.position.color^1]
        if aggression != 0 {
                e.midgame += attackForce.midgame * aggression / Abs(aggression)
                e.endgame += attackForce.endgame * aggression / Abs(aggression)
        }
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
