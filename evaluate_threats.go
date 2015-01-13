// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzeThreats() {
	var threats, center Total

	if engine.trace {
		defer func() {
			e.checkpoint(`Threats`, threats)
			if e.material.turf != 0 && e.material.flags & (whiteKingSafety | blackKingSafety) != 0 {
				e.checkpoint(`Center`, center)
			}
		}()
	}

	threats.white = e.threats(White, e.attacks[White], e.attacks[Black])
	threats.black = e.threats(Black, e.attacks[Black], e.attacks[White])
	e.score.add(threats.white).subtract(threats.black)

	if e.material.turf != 0 && e.material.flags & (whiteKingSafety | blackKingSafety) != 0 {
		center.white = e.center(White, e.attacks[White], e.attacks[Black], e.attacks[pawn(Black)])
		center.black = e.center(Black, e.attacks[Black], e.attacks[White], e.attacks[pawn(White)])
		e.score.add(center.white).subtract(center.black)
	}
}

func (e *Evaluation) threats(color uint8, hisAttacks, herAttacks Bitmask) (score Score) {
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

func (e *Evaluation) center(color uint8, hisAttacks, herAttacks, herPawnAttacks Bitmask) (score Score) {
	pawns := e.position.outposts[pawn(color)]
	safe := homeTurf[color] & ^pawns & ^herPawnAttacks & (hisAttacks | ^herAttacks)
	turf := safe & pawns

	if color == White {
		turf |= turf >> 8   // A4..H4 -> A3..H3
		turf |= turf >> 16  // A4..H4 | A3..H3 -> A2..H2 | A1..H1
		turf &= safe 	    // Keep safe squares only.
		safe <<= 32 	    // Move up to black's half of the board.
	} else {
		turf |= turf << 8   // A5..H5 -> A6..H6
		turf |= turf << 16  // A5..H5 | A6..H6 -> A7..H7 | A8..H8
		turf &= safe 	    // Keep safe squares only.
		safe >>= 32 	    // Move down to white's half of the board.
	}

	score.midgame = (safe | turf).count() * e.material.turf / 100
	return
}
