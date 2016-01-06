// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) generateEvasions() *MoveGen {
	p := gen.p
	color, enemy := p.color, p.color^1
	square := int(p.king[color])

	// Find out what pieces are checking the king. Usually it's a single
	// piece but double check is also a possibility.
	checkers := maskPawn[enemy][square] & p.outposts[pawn(enemy)]
	checkers |= p.targetsFor(square, knight(color)) & p.outposts[knight(enemy)]
	checkers |= p.targetsFor(square, bishop(color)) & (p.outposts[bishop(enemy)] | p.outposts[queen(enemy)])
	checkers |= p.targetsFor(square, rook(color)) & (p.outposts[rook(enemy)] | p.outposts[queen(enemy)])

	// Generate possible king retreats first, i.e. moves to squares not
	// occupied by friendly pieces and not attacked by the opponent.
	retreats := p.targets(square) & ^p.allAttacks(enemy)

	// If the attacking piece is bishop, rook, or queen then exclude the
	// square behind the king using evasion mask. Note that knight's
	// evasion mask is full board so we only check if the attacking piece
	// is not a pawn.
	attackSquare := checkers.pop()
	if p.pieces[attackSquare] != pawn(enemy) {
		retreats &= maskEvade[square][attackSquare]
	}

	// If checkers mask is not empty then we've got double check and
	// retreat is the only option.
	if checkers.any() {
		attackSquare = checkers.first()
		if p.pieces[attackSquare] != pawn(enemy) {
			retreats &= maskEvade[square][attackSquare]
		}
		return gen.movePiece(square, retreats)
	}

	// Generate king retreats. Since castle is not an option there is no
	// reason to use moveKing().
	gen.movePiece(square, retreats)

	// Pawn captures: do we have any pawns available that could capture
	// the attacking piece?
	pawns := maskPawn[color][attackSquare] & p.outposts[pawn(color)]
	for pawns.any() {
		move := NewMove(p, pawns.pop(), attackSquare)
		if attackSquare >= A8 || attackSquare <= H1 {
			move = move.promote(Queen)
		}
		gen.add(move)
	}

	// Rare case when the check could be avoided by en-passant capture.
	// For example: Ke4, c5, e5 vs. Ke8, d7. Black's d7-d5+ could be
	// evaded by c5xd6 or e5xd6 en-passant captures.
	if p.enpassant != 0 {
		if enpassant := attackSquare + up[color]; enpassant == int(p.enpassant) {
			pawns := maskPawn[color][enpassant] & p.outposts[pawn(color)]
			for pawns.any() {
				gen.add(NewMove(p, pawns.pop(), enpassant))
			}
		}
	}

	// See if the check could be blocked or the attacked piece captured.
	block := maskBlock[square][attackSquare] | bit[attackSquare]

	// Create masks for one-square pawn pushes and two-square jumps.
	jumps := ^p.board
	if color == White {
		pawns = (p.outposts[Pawn] << 8) & ^p.board
		jumps &= maskRank[3] & (pawns << 8)
	} else {
		pawns = (p.outposts[BlackPawn] >> 8) & ^p.board
		jumps &= maskRank[4] & (pawns >> 8)
	}
	pawns &= block; jumps &= block

	// Handle one-square pawn pushes: promote to Queen if reached last rank.
	for pawns.any() {
		to := pawns.pop()
		from := to - up[color]
		move := NewMove(p, from, to) // Can't cause en-passant.
		if to >= A8 || to <= H1 {
			move = move.promote(Queen)
		}
		gen.add(move)
	}

	// Handle two-square pawn jumps that can cause en-passant.
	for jumps.any() {
		to := jumps.pop()
		from := to - 2 * up[color]
		gen.add(NewPawnMove(p, from, to))
	}

	// What's left is to generate all possible knight, bishop, rook, and
	// queen moves that evade the check.
	outposts := p.outposts[color] & ^p.outposts[pawn(color)] & ^p.outposts[king(color)]
	for outposts.any() {
		from := outposts.pop()
		targets := p.targets(from) & block
		gen.movePiece(from, targets)
	}

	return gen
}
