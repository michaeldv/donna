// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (e *Evaluation) analyzePieces() {
	p := e.position
	var bonus, score Score
	var knight, bishop, rook, queen, mobility Total

	if engine.trace {
		defer func() {
			var our, their Score
			e.checkpoint(`Mobility`, mobility)
			e.checkpoint(`+Pieces`,  Total{*our.add(knight.white).add(bishop.white).add(rook.white).add(queen.white),
				*their.add(knight.black).add(bishop.black).add(rook.black).add(queen.black)})
			e.checkpoint(`-Knights`, knight)
			e.checkpoint(`-Bishops`, bishop)
			e.checkpoint(`-Rooks`,   rook)
			e.checkpoint(`-Queens`,  queen)
		}()
	}

	// Mobility masks for both sides exclude squares attacked by opponent's pawns,
	// king squares, pawns on first two ranks, and blocked pawns on other ranks.
	var pawnExclusions = [2]Bitmask {
		p.outposts[Pawn] & (maskRank[A2H2] | maskRank[A3H3] | p.board.up(Black)),
		p.outposts[BlackPawn] & (maskRank[A7H7] | maskRank[A6H6] | p.board.up(White)),
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
	if p.outposts[Knight].any() {
		knight.white, bonus = e.knights(White, maskSafe[White], isKingUnsafe[Black])
		e.score.add(knight.white)
		mobility.white.add(bonus)
	}
	if p.outposts[Bishop].any() {
		bishop.white, bonus = e.bishops(White, maskSafe[White], isKingUnsafe[Black])
		e.score.add(bishop.white)
		mobility.white.add(bonus)
	}
	if p.outposts[Rook].any() {
		rook.white, bonus  = e.rooks(White, maskSafe[White], isKingUnsafe[Black])
		e.score.add(rook.white)
		mobility.white.add(bonus)
	}

	// Evaluate black pieces except the queen.
	if p.outposts[BlackKnight].any() {
		knight.black, bonus = e.knights(Black, maskSafe[Black], isKingUnsafe[White])
		e.score.sub(knight.black)
		mobility.black.add(bonus)
	}
	if p.outposts[BlackBishop].any() {
		bishop.black, bonus = e.bishops(Black, maskSafe[Black], isKingUnsafe[White])
		e.score.sub(bishop.black)
		mobility.black.add(bonus)
	}
	if p.outposts[BlackRook].any() {
		rook.black, bonus = e.rooks(Black, maskSafe[Black], isKingUnsafe[White])
		e.score.sub(rook.black)
		mobility.black.add(bonus)
	}

	// Now that we've built all attack bitmasks we can adjust mobility to exclude
	// attacks by enemy's knights, bishops, and rooks and evaluate the queens.
	if p.outposts[Queen].any() {
		maskSafe[White] &= ^(e.attacks[BlackKnight] | e.attacks[BlackBishop] | e.attacks[BlackRook])
		queen.white, bonus = e.queens(White, maskSafe[White], isKingUnsafe[Black])
		e.score.add(queen.white)
		mobility.white.add(bonus)
	}
	if p.outposts[BlackQueen].any() {
		maskSafe[Black] &= ^(e.attacks[Knight] | e.attacks[Bishop] | e.attacks[Rook])
		queen.black, bonus = e.queens(Black, maskSafe[Black], isKingUnsafe[White])
		e.score.sub(queen.black)
		mobility.black.add(bonus)
	}

	// Update attack bitmasks for both sides.
	e.attacks[White] |= e.attacks[Knight] | e.attacks[Bishop] | e.attacks[Rook] | e.attacks[Queen]
	e.attacks[Black] |= e.attacks[BlackKnight] | e.attacks[BlackBishop] | e.attacks[BlackRook] | e.attacks[BlackQueen]

	// Calculate total mobility score applying mobility weight.
	score.add(mobility.white).sub(mobility.black).apply(weightMobility)
	e.score.add(score)
}

func (e *Evaluation) knights(our uint8, maskSafe Bitmask, unsafeKing bool) (score, mobility Score) {
	p, their := e.position, our^1
	outposts := p.outposts[knight(our)]

	for outposts.any() {
		square := outposts.pop()
		attacks := Bitmask(0)

		// Bonus for knight's mobility -- unless the knight is pinned.
		if e.pinned[our].off(square) {
			attacks = p.attacks(square)
			mobility.add(mobilityKnight[(attacks & maskSafe).count()])
		}

		// Penalty if knight is attacked by enemy's pawn.
		if (maskPawn[their][square] & p.outposts[pawn(their)]).any() {
			score.sub(penaltyPawnThreat[Knight/2])
		}

		// Bonus if knight is behind friendly pawn.
		if rank(our, square) < 4 && p.outposts[pawn(our)].on(square + up[our]) {
			score.add(behindPawn)
		}

		// Track if knight attacks squares around enemy's king.
		if unsafeKing {
			e.kingThreats(knight(our), attacks)
		}

		// Update attack bitmask for the knight.
		e.attacks[knight(our)] |= attacks
	}

	return
}

func (e *Evaluation) bishops(our uint8, maskSafe Bitmask, unsafeKing bool) (score, mobility Score) {
	p, their := e.position, our^1
	outposts := p.outposts[bishop(our)]

	for outposts.any() {
		square := outposts.pop()
		attacks := p.xrayAttacks(square)

		// Bonus for bishop's mobility: if the bishop is pinned then restrict the attacks.
		if e.pinned[our].on(square) {
			attacks &= maskLine[p.king[our]][square]
		}
		mobility.add(mobilityBishop[(attacks & maskSafe).count()])


		// Penalty for light/dark-colored pawns restricting a bishop.
		if count := (same(square) & p.outposts[pawn(our)]).count(); count > 0 {
			score.sub(bishopPawn.times(count))
		}

		// Penalty if bishop is attacked by enemy's pawn.
		if (maskPawn[their][square] & p.outposts[pawn(their)]).any() {
			score.sub(penaltyPawnThreat[Bishop/2])
		}

		// Bonus if bishop is behind friendly pawn.
		if rank(our, square) < 4 && p.outposts[pawn(our)].on(square + up[our]) {
			score.add(behindPawn)
		}

		// Middle game penalty for boxed bishop.
		if e.material.phase > 160 {
			if our == White {
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

		// Extra bonus if bishop is on central ranks.
		extra := Score{0, 0}
		if extra.midgame = extraBishop[flip(our, square)]; extra.midgame > 0 {
			extra.endgame = extra.midgame / 2
			score.add(extra)
		}

		// Track if bishop attacks squares around enemy's king.
		if unsafeKing {
			e.kingThreats(bishop(our), attacks)
		}

		// Update attack bitmask for the bishop.
		e.attacks[bishop(our)] |= attacks
	}

	return
}


func (e *Evaluation) rooks(our uint8, maskSafe Bitmask, unsafeKing bool) (score, mobility Score) {
	p, their := e.position, our^1
	ourPawns := p.outposts[pawn(our)]
	theirPawns := p.outposts[pawn(their)]
	outposts := p.outposts[rook(our)]

	// Bonus if rook is on 7th rank and enemy's king trapped on 8th.
	if count := (outposts & mask7th[our]).count(); count > 0 && p.outposts[king(their)] & mask8th[our] != 0 {
		score.add(rookOn7th.times(count))
	}
	for outposts.any() {
		square := outposts.pop()
		attacks := p.xrayAttacks(square)

		// Bonus for rook's mobility: if the rook is pinned then restrict the attacks.
		if e.pinned[our].on(square) {
			attacks &= maskLine[p.king[our]][square]
		}
		safeSquares := (attacks & maskSafe).count()
		mobility.add(mobilityRook[safeSquares])

		// Penalty if rook is attacked by enemy's pawn.
		if maskPawn[their][square] & theirPawns != 0 {
			score.sub(penaltyPawnThreat[Rook/2])
		}

		// Bonus if rook is attacking enemy's pawns.
		if rank(our, square) >= 4 {
			if count := (attacks & theirPawns).count(); count > 0 {
				score.add(rookOnPawn.times(count))
			}
		}

		// Bonuses if rook is on open or semi-open file.
		column := col(square)
		isFileAjar := (ourPawns & maskFile[column] == 0)
		if isFileAjar {
			if theirPawns & maskFile[column] == 0 {
				score.add(rookOnOpen)
			} else {
				score.add(rookOnSemiOpen)
			}
		}

		// Middle game penalty if a rook is boxed. Extra penalty if castle
		// rights have been lost.
		if safeSquares <= 3 || !isFileAjar {
			kingSquare := int(p.king[our])
			kingColumn := col(kingSquare)

			// Queenside box: king on D/C/B vs. rook on A/B/C files. Increase the
			// the penalty since no castle is possible.
			if column < kingColumn && rookBoxA[our].on(square) && kingBoxA[our].on(kingSquare) {
				score.midgame -= (rookBoxed.midgame - safeSquares * 10) * 2
			}

			// Kingside box: king on E/F/G vs. rook on H/G/F files.
			if column > kingColumn && rookBoxH[our].on(square) && kingBoxH[our].on(kingSquare) {
				score.midgame -= (rookBoxed.midgame - safeSquares * 10)
				if p.castles & castleKingside[our] == 0 {
					score.midgame -= (rookBoxed.midgame - safeSquares * 10)
				}
			}
		}

		// Track if rook attacks squares around enemy's king.
		if unsafeKing {
			e.kingThreats(rook(our), attacks)
		}

		// Update attack bitmask for the rook.
		e.attacks[rook(our)] |= attacks
	}

	return
}

func (e *Evaluation) queens(our uint8, maskSafe Bitmask, unsafeKing bool) (score, mobility Score) {
	p, their := e.position, our^1
	outposts := p.outposts[queen(our)]

	for outposts.any() {
		square := outposts.pop()
		attacks := p.attacks(square)

		// Bonus for queen's mobility: if the queen is pinned then restrict the attacks.
		if e.pinned[our].on(square) {
			attacks &= maskLine[p.king[our]][square]
		}
		mobility.add(mobilityQueen[min(15, (attacks & maskSafe).count())])

		// Penalty if queen is attacked by enemy's pawn.
		if (maskPawn[their][square] & p.outposts[pawn(their)]).any() {
			score.sub(penaltyPawnThreat[Queen/2])
		}

		// Track if queen attacks squares around enemy's king.
		if unsafeKing {
			e.kingThreats(queen(our), attacks)
		}

		// Update attack bitmask for the queen.
		e.attacks[queen(our)] |= attacks
	}

	return
}

// Updates safety data used later on when evaluating king safety.
func (e *Evaluation) kingThreats(piece Piece, attacks Bitmask) {
	their := piece.color()^1

	if (attacks & e.safety[their].fort).any() {
		e.safety[their].attackers++
		e.safety[their].threats += kingThreat[piece.id()]
		if bits := attacks & e.attacks[king(their)]; bits.any() {
			e.safety[their].attacks += bits.count()
		}
	}
}

// Initializes the fort bitmask around king's square. For example, for a king on
// G1 the bitmask covers F1,F2,F3, G2,G3, and H1,H2,H3. For a king on a corner
// square, say H1, the bitmask covers F1,F2, G1,G2,G3, and H2,H3.
func (e *Evaluation) setupFort(our uint8) (bitmask Bitmask) {
	bitmask = e.attacks[king(our)] | e.attacks[king(our)].up(our)
	switch e.position.king[our] {
	case A1, A8:
		bitmask |= e.attacks[king(our)] << 1
	case H1, H8:
		bitmask |= e.attacks[king(our)] >> 1
	}

	return bitmask
}
