// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (gen *MoveGen) generateEvasions() *MoveGen {
	p := gen.p
	our, their := p.colors()
	square := p.pick(our).home

	// Find out what pieces are checking the king. Usually it's a single
	// piece but double check is also a possibility.
	checkers := maskPawn[their&1][square] & p.outposts[pawn(their)]
	checkers |= p.knightAttacksAt(square, our) & p.outposts[knight(their)]
	checkers |= p.bishopAttacksAt(square, our) & (p.outposts[bishop(their)] | p.outposts[queen(their)])
	checkers |= p.rookAttacksAt(square, our) & (p.outposts[rook(their)] | p.outposts[queen(their)])

	// Generate possible king retreats first, i.e. moves to squares not
	// occupied by friendly pieces and not attacked by the opponent.
	retreats := p.targets(square) & ^p.allAttacks(their)

	// If the attacking piece is bishop, rook, or queen then exclude the
	// square behind the king using evasion mask. Note that knight's
	// evasion mask is full board so we only check if the attacking piece
	// is not a pawn.
	attackSquare := checkers.first()
	if p.pieces[attackSquare] != pawn(their) {
		retreats &= maskEvade[square][attackSquare]
	}

	// If checkers mask is not empty then we've got double check and
	// retreat is the only option.
	if checkers = checkers.pop(); checkers.anyʔ() {
		attackSquare = checkers.first()
		if p.pieces[attackSquare] != pawn(their) {
			retreats &= maskEvade[square][attackSquare]
		}
		return gen.movePiece(square, retreats)
	}

	// Generate king retreats. Since castle is not an option there is no
	// reason to use moveKing().
	gen.movePiece(square, retreats)

	// Pawn captures: do we have any pawns available that could capture
	// the attacking piece?
	for bm := maskPawn[our&1][attackSquare] & p.outposts[pawn(our)]; bm.anyʔ(); bm = bm.pop() {
		move := NewMove(p, bm.first(), attackSquare)
		if attackSquare >= A8 || attackSquare <= H1 {
			move = move.promote(Queen)
		}
		gen.add(move)
	}

	// Rare case when the check could be avoided by en-passant capture.
	// For example: Ke4, c5, e5 vs. Ke8, d7. Black's d7-d5+ could be
	// evaded by c5xd6 or e5xd6 en-passant captures.
	if p.enpassant != 0 {
		if enpassant := attackSquare + up[our&1]; enpassant == p.enpassant {
			for bm := maskPawn[our&1][enpassant] & p.outposts[pawn(our)]; bm.anyʔ(); bm = bm.pop() {
				gen.add(NewMove(p, bm.first(), enpassant))
			}
		}
	}

	// See if the check could be blocked or the attacked piece captured.
	block := maskBlock[square][attackSquare] | bit(attackSquare)

	// Create masks for one-square pawn pushes and two-square jumps.
	pawns, jumps := maskNone, ^p.board
	if our == White {
		pawns = (p.outposts[Pawn] << 8) & ^p.board
		jumps &= maskRank[3] & (pawns << 8)
	} else {
		pawns = (p.outposts[BlackPawn] >> 8) & ^p.board
		jumps &= maskRank[4] & (pawns >> 8)
	}

	// Handle one-square pawn pushes: promote to Queen if reached last rank.
	for bm := pawns & block; bm.anyʔ(); bm = bm.pop() {
		to := bm.first()
		from := to - up[our&1]
		move := NewMove(p, from, to) // Can't cause en-passant.
		if to >= A8 || to <= H1 {
			move = move.promote(Queen)
		}
		gen.add(move)
	}

	// Handle two-square pawn jumps that can cause en-passant.
	for bm := jumps & block; bm.anyʔ(); bm = bm.pop() {
		to := bm.first()
		from := to - 2 * up[our&1]
		gen.add(NewPawnMove(p, from, to))
	}

	// What's left is to generate all possible knight, bishop, rook, and
	// queen moves that evade the check.
	for bm := p.outposts[our&1] & ^p.outposts[pawn(our)] & ^p.outposts[king(our)]; bm.anyʔ(); bm = bm.pop() {
		from := bm.first()
		gen.movePiece(from, p.targets(from) & block)
	}

	return gen
}
