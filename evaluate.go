// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
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

func (p *Position) Evaluate() int {
        evaluator := &Evaluator{ 0, rightToMove.midgame, rightToMove.endgame, p }
        evaluator.analyzeMaterial()
        evaluator.analyzeCoordination()
        evaluator.analyzePawnStructure()
        evaluator.analyzeRooks()
        evaluator.analyzeKingShield()
        // evaluator.analyzeKingSafety()

        if p.color == Black {
                evaluator.midgame = -evaluator.midgame
                evaluator.endgame = -evaluator.endgame
        }
        return p.score(evaluator.midgame, evaluator.endgame)
}

func (e *Evaluator) analyzeMaterial() {
        for _,piece := range []int{ Pawn, Knight, Bishop, Rook, Queen } {
                count := e.position.count[piece] - e.position.count[piece|Black]
                midgame, endgame := Piece(piece).value()
                e.midgame += midgame * count
                e.endgame += endgame * count
        }
}

func (e *Evaluator) analyzeCoordination() {
        var moves, attacks [2]int
        var bonus [2]Score

        notAttacked := [2]Bitmask{ ^e.position.attacks(White), ^e.position.attacks(Black) }
        board := e.position.board
        for board != 0 {
                square := board.pop()
                piece := e.position.pieces[square]
                color := piece.color()
                targets := e.position.targets(square)

                // Mobility: how many moves are available to squares not attacked by
                // the opponent?
                moves[color] += (targets & notAttacked[color^1]).count()

                // Agressivness: how many opponent's pieces are being attacked?
                attacks[color] += (targets & e.position.outposts[color^1]).count()

                // Calculate bonus or penalty for a piece being at the given square.
                midgame, endgame := piece.bonus(flip[color][square])
                bonus[color].midgame += midgame
                bonus[color].endgame += endgame
        }

        e.midgame += bonus[White].midgame - bonus[Black].midgame
        e.endgame += bonus[White].endgame - bonus[Black].endgame

        mobility := moves[White] - moves[Black]
        e.midgame += mobility * movesAvailable.midgame
        e.endgame += mobility * movesAvailable.endgame

        aggression := attacks[White] - attacks[Black]
        e.midgame += aggression * attackForce.midgame
        e.endgame += aggression * attackForce.endgame

        if bishops := e.position.count[Bishop]; bishops >= 2 {
                e.midgame += bishopPair.midgame
                e.endgame += bishopPair.endgame
        }
        if bishops := e.position.count[BlackBishop]; bishops >= 2 {
                e.midgame -= bishopPair.midgame
                e.endgame -= bishopPair.endgame
        }
}

func (e *Evaluator) analyzePawnStructure() {
        whiteBonus, whitePenalty := e.pawnsScore(White)
        blackBonus, blackPenalty := e.pawnsScore(Black)
        e.midgame += whiteBonus.midgame + whitePenalty.midgame - blackBonus.midgame - blackPenalty.midgame
        e.endgame += whiteBonus.endgame + whitePenalty.endgame - blackBonus.endgame - blackPenalty.endgame
}

func (e *Evaluator) pawnsScore(color int) (bonus, penalty Score){
        hisPawns := e.position.outposts[pawn(color)]
        herPawns := e.position.outposts[pawn(color^1)]

        pawns := hisPawns
        for pawns != 0 {
                square := pawns.pop()
                column := Col(square)
                //
                // The pawn is passed if a) there are no enemy pawns in the
                // same and adjacent columns; and b) there is no same color
                // pawns in front of us.
                //
                if maskPassed[color][square] & herPawns == 0 &&
                   maskInFront[color][square] & hisPawns == 0 {
                           bonus.midgame += bonusPassedPawn[0][flip[color][square]]
                           bonus.endgame += bonusPassedPawn[1][flip[color][square]]
                }
                //
                // Check if the pawn is isolated, i.e. has no pawns of the
                // same color on either sides.
                //
                if maskIsolated[column] & hisPawns == 0 {
                        penalty.midgame += penaltyIsolatedPawn[0][column]
                        penalty.endgame += penaltyIsolatedPawn[1][column]
                }
        }
        //
        // Penalties for doubled pawns.
        //
        for col := 0;  col <= 7; col++ {
                if doubled := (maskFile[col] & hisPawns).count(); doubled > 1 {
                        penalty.midgame += (doubled - 1) * penaltyDoubledPawn[0][col]
                        penalty.endgame += (doubled - 1) * penaltyDoubledPawn[1][col]
                }
        }
        return
}

func (e *Evaluator) analyzeRooks() {
        white := e.rooksScore(White)
        black := e.rooksScore(Black)
        e.midgame += white.midgame - black.midgame
        e.endgame += white.endgame - black.endgame
}

func (e *Evaluator) rooksScore(color int) (bonus Score) {
        rooks := e.position.outposts[rook(color)]
        if rooks == 0 {
                return bonus
        }
        //
        // Bonus if rooks are on 7th rank.
        //
        if count := (rooks & mask7th[color]).count(); count > 0 {
                bonus.midgame += count * rookOn7th.midgame
                bonus.endgame += count * rookOn7th.endgame
        }
        //
        // Bonuses if rooks are on open or semi-open files.
        //
        hisPawns := e.position.outposts[pawn(color)]
        herPawns := e.position.outposts[pawn(color^1)]
        for rooks != 0 {
                square := rooks.pop()
                column := Col(square)
                if hisPawns & maskFile[column] == 0 {
                        if herPawns & maskFile[column] == 0 {
                                bonus.midgame += rookOnOpen.midgame
                                bonus.endgame += rookOnOpen.endgame
                        } else {
                                bonus.midgame += rookOnSemiOpen.midgame
                                bonus.endgame += rookOnSemiOpen.endgame
                        }
                }
        }
        return
}

func (e *Evaluator) analyzeKingShield() {
        // No endgame bonus or penalty.
        e.midgame += e.kingShieldScore(White) - e.kingShieldScore(Black)
}

func (e *Evaluator) kingShieldScore(color int) (penalty int) {
        kings, pawns := e.position.outposts[king(color)], e.position.outposts[pawn(color)]
        //
        // Pass if a) the king is missing, b) the king is on the initial square
        // or c) the opposite side doesn't have a queen with one major piece.
        //
        if kings == 0 || kings == bit[homeKing[color]] || !e.strongEnough(color^1) {
                return
        }
        //
        // Calculate relative square for the king so we could treat black king
        // as white. Don't bother with the shield if the king is too far.
        //
        square := flip[color^1][kings.first()]
        if square > H3 {
                return
        }
        row, col := Coordinate(square)
        from, to := Max(0, col - 1), Min(7, col + 1)
        //
        // For each of the shield columns find the closest same color pawn. The
        // penalty is carried if the pawn is missing or is too far from the king
        // (more than one row apart).
        //
        for column := from; column <= to; column++ {
                if shield := (pawns & maskFile[column]); shield != 0 {
                        closest := flip[color^1][shield.first()] // Make it relative.
                        if distance := Abs(Row(closest) - row); distance > 1 {
                                penalty += distance * -shieldDistance.midgame
                        }
                } else {
                        penalty += -shieldMissing.midgame
                }
        }
        // Log("penalty[%s] => %d\n", C(color), penalty)
        return
}

func (e *Evaluator) strongEnough(color int) bool {
        return e.position.count[queen(color)] > 0 &&
               (e.position.count[rook(color)] > 0 || e.position.count[bishop(color)] > 0 || e.position.count[knight(color)] > 0)
}
