// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzePieces() {
	p := e.position
	var white, black [4]Score

	if Settings.Trace {
		defer func() {
			var his, her Score
			e.checkpoint(`+Pieces`,  Total{*his.add(white[0]).add(white[1]).add(white[2]).add(white[3]),
				*her.add(black[0]).add(black[1]).add(black[2]).add(black[3])})
			e.checkpoint(`-Knights`, Total{white[0], black[0]})
			e.checkpoint(`-Bishops`, Total{white[1], black[1]})
			e.checkpoint(`-Rooks`,   Total{white[2], black[2]})
			e.checkpoint(`-Queens`,  Total{white[3], black[3]})
		}()
	}

	// Mobility mask for both sides excludes a) squares attacked by enemy's
	// pawns and b) squares occupied by own pawns and king.
	maskMobile := [2]Bitmask{
		^(e.attacks[BlackPawn] | p.outposts[Pawn] | p.outposts[King]),
		^(e.attacks[Pawn] | p.outposts[BlackPawn] | p.outposts[BlackKing]),
	}

	// Evaluate white pieces except queen.
	if p.count[Knight] > 0 {
		white[0] = e.knights(White, maskMobile[White])
	}
	if p.count[Bishop] > 0 {
		white[1] = e.bishops(White, maskMobile[White])
	}
	if p.count[Rook] > 0 {
		white[2] = e.rooks(White, maskMobile[White])
	}

	// Evaluate black pieces except queen.
	if p.count[BlackKnight] > 0 {
		black[0] = e.knights(Black, maskMobile[Black])
	}
	if p.count[BlackBishop] > 0 {
		black[1] = e.bishops(Black, maskMobile[Black])
	}
	if p.count[BlackRook] > 0 {
		black[2] = e.rooks(Black, maskMobile[Black])
	}

	// Now that we've built all attack bitmasks we can adjust mobility to
	// exclude attacks by enemy's knights, bishops, and rooks and evaluate
	// the queens.
	if p.count[Queen] > 0 {
		maskMobile[White] &= ^(e.attacks[BlackKnight] | e.attacks[BlackBishop] | e.attacks[BlackRook])
		white[3] = e.queens(White, maskMobile[White])
	}
	if p.count[BlackQueen] > 0 {
		maskMobile[Black] &= ^(e.attacks[Knight] | e.attacks[Bishop] | e.attacks[Rook])
		black[3] = e.queens(Black, maskMobile[Black])
	}

	// Update attack bitmasks for both sides.
	e.attacks[White] |= e.attacks[Knight] | e.attacks[Bishop] | e.attacks[Rook] | e.attacks[Queen]
	e.attacks[Black] |= e.attacks[BlackKnight] | e.attacks[BlackBishop] | e.attacks[BlackRook] | e.attacks[BlackQueen]

	// Update cumulative score based on white vs. black delta.
	e.score.add(white[0]).add(white[1]).add(white[2]).add(white[3])
	e.score.subtract(black[0]).subtract(black[1]).subtract(black[2]).subtract(black[3])
}

func (e *Evaluation) knights(color int, maskMobile Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[knight(color)]

	for outposts != 0 {
		square := outposts.pop()
		attacks := p.attacks(square)

		// Bonus for knight's mobility.
		score.add(mobilityKnight[(attacks & maskMobile).count()])

		// Penalty if knight is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Knight/2])
		}

		// Bonus if knight is behind friendly pawn.
		if RelRow(color, square) < 4 && p.outposts[pawn(color)].isSet(square + eight[color]) {
			score.add(behindPawn)
		}

		// Extra bonus if knight is in the center. Increase the extra
		// bonus if the knight is supported by a pawn and can't be
		// exchanged.
		flip := Flip(color, square)
		if extra := extraKnight[flip]; extra > 0 {
			if p.pawnAttacks(color).isSet(square) {
				if p.count[knight(color^1)] == 0 {
					extra *= 2 // No knights to exchange.
				}
				extra += extra / 2 // Supported by a pawn.
			}
			score.adjust(extra)
		}

		// Track if knight attacks squares around enemy's king.
		e.enemyKingThreat(knight(color), attacks)
	}
	return
}

func (e *Evaluation) bishops(color int, maskMobile Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[bishop(color)]

	for outposts != 0 {
		square := outposts.pop()
		attacks := p.xrayAttacks(square)

		// Bonus for bishop's mobility
		score.add(mobilityBishop[(attacks & maskMobile).count()])

		// Penalty for light/dark square bishop and matching pawns.
		if count := (Same(square) & p.outposts[pawn(color)]).count(); count > 0 {
			score.subtract(bishopPawns)
		}

		// Penalty if bishop is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Bishop/2])
		}

		// Bonus if bishop is behind friendly pawn.
		if RelRow(color, square) < 4 && p.outposts[pawn(color)].isSet(square + eight[color]) {
			score.add(behindPawn)
		}

		// Middle game penalty for boxed bishop.
		if e.material.phase > 160 {
			if color == White {
				if (square == C1 && p.pieces[D2].isPawn() && p.pieces[D3] != 0) ||
				   (square == F1 && p.pieces[E2].isPawn() && p.pieces[E3] != 0) {
					score.midgame -= bishopBoxed.midgame
				}
			} else {
				if (square == C8 && p.pieces[D7].isPawn() && p.pieces[D6] != 0) ||
				   (square == F8 && p.pieces[E7].isPawn() && p.pieces[E6] != 0) {
					score.midgame -= bishopBoxed.midgame
				}
			}
		}

		// Extra bonus if bishop is in the center. Increase the extra
		// bonus if the bishop is supported by a pawn and can't be
		// exchanged.
		flip := Flip(color, square)
		if extra := extraBishop[flip]; extra > 0 {
			if p.pawnAttacks(color).isSet(square) {
				if p.count[bishop(color^1)] == 0 {
					extra *= 2 // No bishops to exchange.
				}
				extra += extra / 2 // Supported by a pawn.
			}
			score.adjust(extra)
		}

		// Track if bishop attacks squares around enemy's king.
		e.enemyKingThreat(bishop(color), attacks)
	}

	// Bonus for the pair of bishops.
	if bishops := p.count[bishop(color)]; bishops >= 2 {
		score.add(bishopPair)
	}
	return
}


func (e *Evaluation) rooks(color int, maskMobile Bitmask) (score Score) {
	p := e.position
	hisPawns := p.outposts[pawn(color)]
	herPawns := p.outposts[pawn(color^1)]
	outposts := p.outposts[rook(color)]

	// Bonus if rook is on 7th rank and enemy's king trapped on 8th.
	if count := (outposts & mask7th[color]).count(); count > 0 && p.outposts[king(color^1)] & mask8th[color] != 0 {
		score.add(rookOn7th.times(count))
	}
	for outposts != 0 {
		square := outposts.pop()
		attacks := p.xrayAttacks(square)

		// Bonus for rook's mobility
		mobility := (attacks & maskMobile).count()
		score.add(mobilityRook[mobility])

		// Penalty if rook is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Rook/2])
		}

		// Bonus if rook is attacking enemy's pawns.
		if count := (attacks & p.outposts[pawn(color^1)]).count(); count > 0 {
			score.add(rookOnPawn.times(count))
		}

		// Bonuses if rook is on open or semi-open file.
		column := Col(square)
		isFileAjar := (hisPawns & maskFile[column] == 0)
		if isFileAjar {
			if herPawns & maskFile[column] == 0 {
				score.add(rookOnOpen)
			} else {
				score.add(rookOnSemiOpen)
			}
		}

		// Middle game penalty if a rook is boxed. Extra penalty if castle
		// rights have been lost.
		if mobility <= 3 || !isFileAjar {
			kingSquare := p.king[color]
			kingColumn := Col(kingSquare)

			// Queenside box: king on D/C/B vs. rook on A/B/C files. Double the
			// the penalty since no castle is possible.
			if column < kingColumn && rookBoxA[color].isSet(square) && kingBoxA[color].isSet(kingSquare) {
				score.midgame -= rookBoxed.midgame * 2
			}

			// Kingside box: king on E/F/G vs. rook on H/G/F files.
			if column > kingColumn && rookBoxH[color].isSet(square) && kingBoxH[color].isSet(kingSquare) {
				score.midgame -= rookBoxed.midgame
				if p.castles & castleKingside[color] == 0 {
					score.midgame -= rookBoxed.midgame
				}
			}
		}

		// Track if rook attacks squares around enemy's king.
		e.enemyKingThreat(rook(color), attacks)
	}
	return
}

func (e *Evaluation) queens(color int, maskMobile Bitmask) (score Score) {
	p := e.position
	outposts := p.outposts[queen(color)]

	// Bonus if queen is on 7th rank and enemy's king trapped on 8th.
	if count := (outposts & mask7th[color]).count(); count > 0 && p.outposts[king(color^1)] & mask8th[color] != 0 {
		score.add(queenOn7th.times(count))
	}
	for outposts != 0 {
		square := outposts.pop()
		attacks := p.attacks(square)

		// Bonus for queen's mobility.
		score.add(mobilityQueen[Min(15, (attacks & maskMobile).count())])

		// Penalty if queen is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Queen/2])
		}

		// Bonus if queen is out and attacking enemy's pawns.
		if count := (attacks & p.outposts[pawn(color^1)]).count(); count > 0 && RelRow(color, square) > 3 {
			score.add(queenOnPawn.times(count))
		}

		// Track if queen attacks squares around enemy's king.
		e.enemyKingThreat(queen(color), attacks)
	}
	return
}

func (e *Evaluation) enemyKingThreat(piece Piece, attacks Bitmask) {
	color := piece.color() ^ 1

	if attacks & e.safety[color].fort != 0 {
		e.safety[color].attackers++
		e.safety[color].threats += bonusKingThreat[piece.kind()/2]
		if bits := attacks & e.attacks[king(color)]; bits != 0 {
			e.safety[color].attacks += bits.count()
		}
	}

	// Update attack bitmask for the given piece.
	e.attacks[piece] |= attacks
}
