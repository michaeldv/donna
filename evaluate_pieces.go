// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzePieces() {
	p := e.position
	var mobile Score
	var knight, bishop, rook, queen, mobility Total

	if engine.trace {
		defer func() {
			var his, her Score
			e.checkpoint(`Mobility`, mobility)
			e.checkpoint(`+Pieces`,  Total{*his.add(knight.white).add(bishop.white).add(rook.white).add(queen.white),
				*her.add(knight.black).add(bishop.black).add(rook.black).add(queen.black)})
			e.checkpoint(`-Knights`, knight)
			e.checkpoint(`-Bishops`, bishop)
			e.checkpoint(`-Rooks`,   rook)
			e.checkpoint(`-Queens`,  queen)
		}()
	}

	// Mobility masks for both sides exclude squares attacked by rival's pawns,
	// king squares, pawns on first two ranks, and blocked pawns on other ranks.
	var pawnExclusions = [2]Bitmask {
		p.outposts[Pawn] & (maskRank[A2H2] | maskRank[A3H3] | p.board.pushed(Black)),
		p.outposts[BlackPawn] & (maskRank[A7H7] | maskRank[A6H6] | p.board.pushed(White)),
	}

	// Initialize safe mobility zones for both sides.
	var maskSafe = [2]Bitmask {
		^(e.attacks[BlackPawn] | p.outposts[King] | pawnExclusions[White]),
		^(e.attacks[Pawn] | p.outposts[BlackKing] | pawnExclusions[Black]),
	}

	// Initialize flags to see if kings for both sides require safety evaluation.
	var isKingUnsafe = [2]bool { e.isKingUnsafe(White), e.isKingUnsafe(Black) }

	// Initialize king fort bitmasks only when we need them.
	if isKingUnsafe[White] {
		e.safety[White].fort = e.setupFort(White)
	}
	if isKingUnsafe[Black] {
		e.safety[Black].fort = e.setupFort(Black)
	}

	// Evaluate white pieces except the queen.
	if p.outposts[Knight] != 0 {
		knight.white, mobile = e.knights(White, maskSafe[White], isKingUnsafe[Black])
		mobility.white.add(mobile)
	}
	if p.outposts[Bishop] != 0 {
		bishop.white, mobile = e.bishops(White, maskSafe[White], isKingUnsafe[Black])
		mobility.white.add(mobile)
	}
	if p.outposts[Rook] != 0 {
		rook.white, mobile = e.rooks(White, maskSafe[White], isKingUnsafe[Black])
		mobility.white.add(mobile)
	}

	// Evaluate black pieces except the queen.
	if p.outposts[BlackKnight] != 0 {
		knight.black, mobile = e.knights(Black, maskSafe[Black], isKingUnsafe[White])
		mobility.black.add(mobile)
	}
	if p.outposts[BlackBishop] != 0 {
		bishop.black, mobile = e.bishops(Black, maskSafe[Black], isKingUnsafe[White])
		mobility.black.add(mobile)
	}
	if p.outposts[BlackRook] != 0 {
		rook.black, mobile = e.rooks(Black, maskSafe[Black], isKingUnsafe[White])
		mobility.black.add(mobile)
	}

	// Now that we've built all attack bitmasks we can adjust mobility to
	// exclude attacks by enemy's knights, bishops, and rooks and evaluate
	// the queens.
	if p.outposts[Queen] != 0 {
		maskSafe[White] &= ^(e.attacks[BlackKnight] | e.attacks[BlackBishop] | e.attacks[BlackRook])
		queen.white, mobile = e.queens(White, maskSafe[White], isKingUnsafe[Black])
		mobility.white.add(mobile)
	}
	if p.outposts[BlackQueen] != 0 {
		maskSafe[Black] &= ^(e.attacks[Knight] | e.attacks[Bishop] | e.attacks[Rook])
		queen.black, mobile = e.queens(Black, maskSafe[Black], isKingUnsafe[White])
		mobility.black.add(mobile)
	}

	// Update attack bitmasks for both sides.
	e.attacks[White] |= e.attacks[Knight] | e.attacks[Bishop] | e.attacks[Rook] | e.attacks[Queen]
	e.attacks[Black] |= e.attacks[BlackKnight] | e.attacks[BlackBishop] | e.attacks[BlackRook] | e.attacks[BlackQueen]

	// Apply weights to the mobility scores.
	mobility.white.apply(weights[0])
	mobility.black.apply(weights[0])

	// Update cumulative score based on white vs. black bonuses and mobility.
	e.score.add(knight.white).add(bishop.white).add(rook.white).add(queen.white).add(mobility.white)
	e.score.subtract(knight.black).subtract(bishop.black).subtract(rook.black).subtract(queen.black).subtract(mobility.black)
}

func (e *Evaluation) knights(color uint8, maskSafe Bitmask, unsafeKing bool) (score, mobility Score) {
	p := e.position
	outposts := p.outposts[knight(color)]

	for outposts != 0 {
		square := outposts.pop()
		attacks := Bitmask(0)

		// Bonus for knight's mobility -- unless the knight is pinned.
		if e.pinned[color].off(square) {
			attacks = p.attacks(square)
			mobility.add(mobilityKnight[(attacks & maskSafe).count()])
		}

		// Penalty if knight is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Knight/2])
		}

		// Bonus if knight is behind friendly pawn.
		if rank(color, square) < 4 && p.outposts[pawn(color)].on(square + eight[color]) {
			score.add(behindPawn)
		}

		// Extra bonus if knight is on central ranks. Increase the extra bonus
		// if the knight is supported by a pawn.
		extra := Score{0, 0}
		if extra.midgame = extraKnight[flip(color, square)]; extra.midgame > 0 {
			extra.endgame = extra.midgame / 4
			if p.pawnAttacks(color).on(square) {
				extra.scale(35) // Bump by 35% if supported by a pawn.
			}
			score.add(extra)
		}

		// Track if knight attacks squares around enemy's king.
		if unsafeKing {
			e.kingThreats(knight(color), attacks)
		}

		// Update attack bitmask for the knight.
		e.attacks[knight(color)] |= attacks
	}
	return
}

func (e *Evaluation) bishops(color uint8, maskSafe Bitmask, unsafeKing bool) (score, mobility Score) {
	p := e.position
	outposts := p.outposts[bishop(color)]

	for outposts != 0 {
		square := outposts.pop()
		attacks := p.xrayAttacks(square)

		// Bonus for bishop's mobility: if the bishop is pinned then restrict the attacks.
		if e.pinned[color].on(square) {
			king := p.king[color]
			attacks &= (maskBlock[king][square] | maskDiagonal[king][square])
		}
		mobility.add(mobilityBishop[(attacks & maskSafe).count()])


		// Penalty for light/dark-colored pawns restricting a bishop.
		if count := (same(square) & p.outposts[pawn(color)]).count(); count > 0 {
			score.subtract(bishopPawn.times(count))
		}

		// Penalty if bishop is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Bishop/2])
		}

		// Bonus if bishop is behind friendly pawn.
		if rank(color, square) < 4 && p.outposts[pawn(color)].on(square + eight[color]) {
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

		// Extra bonus if bishop is on central ranks. Increase the extra bonus
		// if the bishop is supported by a pawn.
		extra := Score{0, 0}
		if extra.midgame = extraBishop[flip(color, square)]; extra.midgame > 0 {
			extra.endgame = extra.midgame / 2
			if p.pawnAttacks(color).on(square) {
				extra.scale(35) // Bump by 35% if supported by a pawn.
			}
			score.add(extra)
		}

		// Track if bishop attacks squares around enemy's king.
		if unsafeKing {
			e.kingThreats(bishop(color), attacks)
		}

		// Update attack bitmask for the bishop.
		e.attacks[bishop(color)] |= attacks
	}
	return
}


func (e *Evaluation) rooks(color uint8, maskSafe Bitmask, unsafeKing bool) (score, mobility Score) {
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

		// Bonus for rook's mobility: if the rook is pinned then restrict the attacks.
		if e.pinned[color].on(square) {
			king := p.king[color]
			attacks &= (maskBlock[king][square] | maskStraight[king][square])
		}
		safeSquares := (attacks & maskSafe).count()
		mobility.add(mobilityRook[safeSquares])

		// Penalty if rook is attacked by enemy's pawn.
		if maskPawn[color^1][square] & herPawns != 0 {
			score.subtract(penaltyPawnThreat[Rook/2])
		}

		// Bonus if rook is attacking enemy's pawns.
		if rank(color, square) >= 4 {
			if count := (attacks & herPawns).count(); count > 0 {
				score.add(rookOnPawn.times(count))
			}
		}

		// Bonuses if rook is on open or semi-open file.
		column := col(square)
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
		if safeSquares <= 3 || !isFileAjar {
			kingSquare := int(p.king[color])
			kingColumn := col(kingSquare)

			// Queenside box: king on D/C/B vs. rook on A/B/C files. Increase the
			// the penalty since no castle is possible.
			if column < kingColumn && rookBoxA[color].on(square) && kingBoxA[color].on(kingSquare) {
				score.midgame -= (rookBoxed.midgame - safeSquares * 10) * 2
			}

			// Kingside box: king on E/F/G vs. rook on H/G/F files.
			if column > kingColumn && rookBoxH[color].on(square) && kingBoxH[color].on(kingSquare) {
				score.midgame -= (rookBoxed.midgame - safeSquares * 10)
				if p.castles & castleKingside[color] == 0 {
					score.midgame -= (rookBoxed.midgame - safeSquares * 10)
				}
			}
		}

		// Track if rook attacks squares around enemy's king.
		if unsafeKing {
			e.kingThreats(rook(color), attacks)
		}

		// Update attack bitmask for the rook.
		e.attacks[rook(color)] |= attacks
	}
	return
}

func (e *Evaluation) queens(color uint8, maskSafe Bitmask, unsafeKing bool) (score, mobility Score) {
	p := e.position
	outposts := p.outposts[queen(color)]

	// Bonus if queen is on 7th rank and enemy's king trapped on 8th.
	if count := (outposts & mask7th[color]).count(); count > 0 && p.outposts[king(color^1)] & mask8th[color] != 0 {
		score.add(queenOn7th.times(count))
	}
	for outposts != 0 {
		square := outposts.pop()
		attacks := p.attacks(square)

		// Bonus for queen's mobility: if the queen is pinned then restrict the attacks.
		if e.pinned[color].on(square) {
			king := p.king[color]
			attacks &= (maskBlock[king][square] | maskDiagonal[king][square] | maskStraight[king][square])
		}
		mobility.add(mobilityQueen[min(15, (attacks & maskSafe).count())])

		// Penalty if queen is attacked by enemy's pawn.
		if maskPawn[color^1][square] & p.outposts[pawn(color^1)] != 0 {
			score.subtract(penaltyPawnThreat[Queen/2])
		}

		// Bonus if queen is out and attacking enemy's pawns.
		if count := (attacks & p.outposts[pawn(color^1)]).count(); count > 0 && rank(color, square) > 3 {
			score.add(queenOnPawn.times(count))
		}

		// Track if queen attacks squares around enemy's king.
		if unsafeKing {
			e.kingThreats(queen(color), attacks)
		}

		// Update attack bitmask for the queen.
		e.attacks[queen(color)] |= attacks
	}
	return
}

// Updates safety data used later on when evaluating king safety.
func (e *Evaluation) kingThreats(piece Piece, attacks Bitmask) {
	color := piece.color() ^ 1

	if attacks & e.safety[color].fort != 0 {
		e.safety[color].attackers++
		e.safety[color].threats += bonusKingThreat[piece.kind()/2]
		if bits := attacks & e.attacks[king(color)]; bits != 0 {
			e.safety[color].attacks += bits.count()
		}
	}
}

// Initializes the fort bitmask around king's square. For example, for a king on
// G1 the bitmask covers F1,F2,F3, G2,G3, and H1,H2,H3. For a king on a corner
// square, say H1, the bitmask covers F1,F2, G1,G2,G3, and H2,H3.
func (e *Evaluation) setupFort(color uint8) (bitmask Bitmask) {
	bitmask = e.attacks[king(color)] | e.attacks[king(color)].pushed(color)
	switch e.position.king[color] {
	case A1, A8:
		bitmask |= e.attacks[king(color)] << 1
	case H1, H8:
		bitmask |= e.attacks[king(color)] >> 1
	}
	return
}
