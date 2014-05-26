// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeThreats() {
	white := e.threats(White, e.attacks[White], e.attacks[Black])
	black := e.threats(Black, e.attacks[Black], e.attacks[White])

	if Settings.Trace {
		defer func() {
			e.checkpoint(`Threats`, Total{white, black})
		}()
	}

	e.score.add(white).subtract(black)
}

func (e *Evaluation) threats(color int, hisAttacks, herAttacks Bitmask) (score Score) {
	p := e.position

	// Find weak enemy pieces: the ones under attack and not defended by
	// pawns (excluding a king).
	weak := p.outposts[color^1] & hisAttacks & ^e.attacks[pawn(color^1)]
	weak &= ^p.outposts[king(color^1)]

	if weak != 0 {

		// Threat bonus for strongest enemy piece attacked by our pawns,
		// knights, or bishops.
		targets := weak & (e.attacks[pawn(color)] | e.attacks[knight(color)] | e.attacks[bishop(color)])
		if targets != 0 {
			piece := p.strongestPiece(color^1, targets)
			score.add(bonusMinorThreat[piece.kind()/2])
		}

		// Threat bonus for strongest enemy piece attacked by our rooks
		// or queen.
		targets = weak & (e.attacks[rook(color)] | e.attacks[queen(color)])
		if targets != 0 {
			piece := p.strongestPiece(color^1, targets)
			score.add(bonusMajorThreat[piece.kind()/2])
		}

		// Extra bonus when attacking enemy pieces that are hanging. Side
		// having the right to move gets bigger bonus.
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
