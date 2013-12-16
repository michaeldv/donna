// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

type Score struct {
	midgame  int
	endgame  int
}

type Brain struct {
        player   *Player
        color    int
        stage    int
        endgame  int
        midgame  int
}

func NewBrain(player *Player) *Brain {
        brain := new(Brain)

        brain.player = player
        brain.color = player.Color

        return brain
}

func (b *Brain) Evaluate(p *Position) (score int) {
        b.endgame, b.midgame = 0, 0

        b.determineGameStage(p)
        b.analyzeMaterial(p)
        b.analyzeCoordination(p)
        // b.analyzePawnStructure(p)
        // b.analyzePassedPawns(p)
        // b.analyzeKingSafety(p)

        score = (b.midgame * b.stage + b.endgame * (256 - b.stage)) / 256
        return
}

// Determine game stage by counting how many pieces are present on the board.
func (b *Brain) determineGameStage(p *Position) {
        b.stage  =  2 * (p.count[Pawn(WHITE)]   + p.count[Pawn(BLACK)])
        b.stage +=  6 * (p.count[Knight(WHITE)] + p.count[Knight(BLACK)])
        b.stage += 12 * (p.count[Bishop(WHITE)] + p.count[Bishop(BLACK)])
        b.stage += 16 * (p.count[Rook(WHITE)]   + p.count[Rook(BLACK)])
        b.stage += 44 * (p.count[Queen(WHITE)]  + p.count[Queen(BLACK)])
}

func (b *Brain) analyzeMaterial(p *Position) {
        color, opposite := b.color, b.color^1

        count := p.count[Pawn(color)] - p.count[Pawn(opposite)]
        b.endgame += valuePawn.endgame * count
        b.midgame += valuePawn.midgame * count

        count = p.count[Knight(color)] - p.count[Knight(opposite)]
        b.endgame += valueKnight.endgame * count
        b.midgame += valueKnight.midgame * count

        count = p.count[Bishop(color)] - p.count[Bishop(opposite)]
        b.endgame += valueBishop.endgame * count
        b.midgame += valueBishop.midgame * count

        count = p.count[Rook(color)] - p.count[Rook(opposite)]
        b.endgame += valueRook.endgame * count
        b.midgame += valueRook.midgame * count

        count = p.count[Queen(color)] - p.count[Queen(opposite)]
        b.endgame += valueQueen.endgame * count
        b.midgame += valueQueen.midgame * count
}

func (b *Brain) analyzeCoordination(p *Position) {
        var moves, attacks [2]int

        for square, piece := range p.pieces {
                if piece == 0 {
                        continue
                }
                color := piece.Color()

                // Mobility and agressivness.
                targets := p.targets[square]
                moves[color] += targets.Count()
                attacks[color] += targets.Intersect(p.board[color^1]).Count()

                // Piece/square adjustments.
                if color == WHITE {
                        square = flip[square]
                }
                switch piece.Kind() {
                case PAWN:
                        b.midgame += bonusPawn[square]
                        b.endgame += bonusPawn[square]
                case KNIGHT:
                        b.midgame += bonusKnight[square]
                        b.endgame += bonusKnight[square]
                case BISHOP:
                        b.midgame += bonusBishop[square]
                        b.endgame += bonusBishop[square]
                // case ROOK:
                //         bonus = bonusRook[square]
                // case QUEEN:
                //         bonus = bonusQueen[square]
                case KING:
                        b.midgame += bonusKing[square]
                        b.endgame += bonusKingEndgame[square]
                }
        }
        mobility := moves[b.color] - moves[b.color^1]
        if mobility != 0 {
                mobility = 25 * mobility / Abs(mobility)
        }
        aggression := attacks[b.color] - attacks[b.color^1]
        if aggression != 0 {
                aggression = 25 * aggression / Abs(aggression)
        }
        b.endgame += mobility + aggression
        b.midgame += mobility + aggression
}

func (b *Brain) analyzePawnStructure(p *Position) {
        // for color := WHITE; color <= BLACK; color++ {
        //     outposts = p->outposts(Pawn(color))
        //
        //     for outposts.IsNotEmpty() {
        //             square := outposts.FirstSet()
        //             outposts.Clear(target)
        //     }
        // }
}

func (b *Brain) analyzePassedPawns(p *Position) {
}

func (b *Brain) analyzeKingSafety(p *Position) {
}
