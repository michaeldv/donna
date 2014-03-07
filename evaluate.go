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

func (p *Position) Evaluate() int {
        evaluator := &Evaluator{ 0, 0, 0, p }
        evaluator.analyzeMaterial()
        evaluator.analyzeCoordination()
        evaluator.analyzePawnStructure()
        evaluator.analyzeRooks()
        evaluator.analyzeKingShield()
        // evaluator.analyzeKingSafety()

        // Right to move: positive bonus for white, and negative for black.
        evaluator.midgame += rightToMove.midgame * (1 - 2 * p.color)
        evaluator.endgame += rightToMove.endgame * (1 - 2 * p.color)

        return p.score(evaluator.midgame, evaluator.endgame)
}

func (e *Evaluator) analyzeMaterial() {
        color, opposite := e.position.color, e.position.color^1

        for _,piece := range []int{ Pawn, Knight, Bishop, Rook, Queen } {
                count := e.position.count[Piece(piece|color)] - e.position.count[Piece(piece|opposite)]
                midgame, endgame := Piece(piece).value()
                e.midgame += midgame * count
                e.endgame += endgame * count
        }
}

func (e *Evaluator) analyzeCoordination() {
        var moves, attacks [2]int
        var bonus [2]Score

        notAttacked := [2]Bitmask{^e.position.attacks(White), ^e.position.attacks(Black)}
        for square, piece := range e.position.pieces {
                if piece == 0 {
                        continue
                }
                color := piece.color()

                // Mobility: how many moves are available to squares not attacked by
                // the opponent?
                moves[color] += (e.position.targets(square) & notAttacked[color^1]).count()

                // Agressivness: how many opponent's pieces are being attacked?
                attacks[color] += (e.position.targets(square) & e.position.outposts[color^1]).count()

                // Calculate bonus or penalty for a piece being at the given square.
                midgame, endgame := piece.bonus(flip[color][square])
                bonus[color].midgame += midgame
                bonus[color].endgame += endgame
        }

        e.adjust(bonus)

        color, opposite := e.position.color, e.position.color^1
        mobility := moves[color] - moves[opposite]
        e.midgame += mobility * movesAvailable.midgame
        e.endgame += mobility * movesAvailable.endgame

        aggression := attacks[color] - attacks[opposite]
        e.midgame += aggression * attackForce.midgame
        e.endgame += aggression * attackForce.endgame

        if bishops := e.position.count[bishop(color)]; bishops >= 2 {
                e.midgame += bishopPair.midgame
                e.endgame += bishopPair.endgame
        }
        if bishops := e.position.count[bishop(opposite)]; bishops >= 2 {
                e.midgame -= bishopPair.midgame
                e.endgame -= bishopPair.endgame
        }
}

func (e *Evaluator) analyzePawnStructure() {
        whiteBonus, whitePenalty := e.pawnsScore(White)
        blackBonus, blackPenalty := e.pawnsScore(Black)
        if e.position.color == White {
                e.midgame += whiteBonus.midgame + whitePenalty.midgame - blackBonus.midgame - blackPenalty.midgame
                e.endgame += whiteBonus.endgame + whitePenalty.endgame - blackBonus.endgame - blackPenalty.endgame
        } else {
                e.midgame += blackBonus.midgame + blackPenalty.midgame - whiteBonus.midgame - whitePenalty.midgame
                e.endgame += blackBonus.endgame + blackPenalty.endgame - whiteBonus.endgame - whitePenalty.endgame
        }
}

func (e *Evaluator) pawnsScore(color int) (bonus, penalty Score){
        var doubled [8]int // Number of doubled pawns in each column.
        hisPawns := e.position.outposts[pawn(color)]
        herPawns := e.position.outposts[pawn(color^1)]

        pawns := hisPawns
        for pawns != 0 {
                square := pawns.pop()
                column := Col(square)
                //
                // count doubled pawns in the column as they carry a penalty.
                //
                doubled[column] = (maskFile[column] & hisPawns).count()
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
        for i := 0;  i < len(doubled); i++ {
                if doubled[i] > 0 {
                        penalty.midgame += (doubled[i] - 1) * penaltyDoubledPawn[0][i]
                        penalty.endgame += (doubled[i] - 1) * penaltyDoubledPawn[1][i]
                }
        }
        return
}

func (e *Evaluator) analyzeRooks() {
        white := e.rooksScore(White)
        black := e.rooksScore(Black)
        if e.position.color == White {
                e.midgame += white.midgame - black.midgame
                e.endgame += white.endgame - black.endgame
        } else {
                e.midgame += black.midgame - white.midgame
                e.endgame += black.endgame - white.endgame
        }
}

func (e *Evaluator) rooksScore(color int) (bonus Score) {
        rooks := e.position.outposts[rook(color)]
        if rooks == 0 {
                return bonus
        }
        //
        // Bonus if rooks are on 7th rank.
        //
        seventh := [2]Bitmask{ maskRank[6], maskRank[1] }
        if count := (rooks & seventh[color]).count(); count > 0 {
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
        white := e.kingShieldScore(White)
        black := e.kingShieldScore(Black)

        // No endgame bonus or penalty.
        if e.position.color == White {
                e.midgame += white - black
        } else {
                e.midgame += black - white
        }
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

func (e *Evaluator) adjust(bonus [2]Score) *Evaluator {
        color, opposite := e.position.color, e.position.color^1

        e.midgame += bonus[color].midgame - bonus[opposite].midgame
        e.endgame += bonus[color].endgame - bonus[opposite].endgame

        return e
}

func (e *Evaluator) strongEnough(color int) bool {
        return e.position.count[queen(color)] > 0 &&
               (e.position.count[rook(color)] > 0 || e.position.count[bishop(color)] > 0 || e.position.count[knight(color)] > 0)
}
