// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeThreats() {
	whiteAttacks := e.position.attacks(White)
	blackAttacks := e.position.attacks(Black)

	white := e.threats(White, whiteAttacks, blackAttacks)
	black := e.threats(Black, blackAttacks, whiteAttacks)

	if Settings.Trace {
		defer func() {
			e.checkpoint(`Threats`, Total{white, black})
		}()
	}

	e.score.add(white).subtract(black)
}

func (e *Evaluation) threats(color int, hisAttacks, herAttacks Bitmask) (score Score) {
	p := e.position

	// Find enemy pieces (excludes king) under attack that are not defended
	// by pawns.
	weak := p.outposts[color^1] & hisAttacks & ^p.pawnAttacks(color^1)
	weak &= ^p.outposts[king(color^1)]
	if weak != 0 {

		// Attacks by pawns, knights, and bishops.
		targets := weak & (p.pawnAttacks(color) | p.knightAttacks(color) | p.bishopAttacks(color))
		if targets != 0 {
			piece := p.strongestPiece(color^1, targets)
			score.add(bonusMinorThreat[piece.kind()/2])
		}

		// Attacks by rooks and queens.
		targets = weak & (p.rookAttacks(color) | p.queenAttacks(color))
		if targets != 0 {
			piece := p.strongestPiece(color^1, targets)
			score.add(bonusMajorThreat[piece.kind()/2])
		}

		// Bonus when pieces under attack are hanging. Whoever has the
		// right to move gets a bit extra.
		hanging := (weak & ^herAttacks).count()
		if hanging > 0 {
			if p.color == color {
				hanging++
			}
			score.add(hangingAttack.times(hanging))
		}
	}
	return
}