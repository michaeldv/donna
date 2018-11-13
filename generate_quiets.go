// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (gen *MoveGen) addQuiet(move Move) *MoveGen {
	gen.list[gen.tail].move = move
	gen.list[gen.tail].score = game.good(move)
	gen.tail++

	return gen
}

// Generates pseudo-legal moves that preserve material balance, i.e.
// no captures or pawn promotions are allowed.
func (gen *MoveGen) generateQuiets() *MoveGen {
	p := gen.p
	our := p.color
	empty := ^p.board

	// Castles.
	if !p.inCheckʔ(our) {
		home := homeKing[our&1]
		kingside, queenside := p.canCastleʔ(our)
		if kingside {
			gen.addQuiet(NewCastle(p, home, home + 2))
		}
		if queenside {
			gen.addQuiet(NewCastle(p, home, home - 2))
		}
	}

	// Pawns.
	last := let(our == White, 7, 0)
	for bm := p.outposts[pawn(our)].up(our) & empty & ^maskRank[last]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		forward, backward := square + up[our&1], square - up[our&1]
		if rank(our, square) == 2 && p.pieces[forward].noneʔ() {
			gen.addQuiet(NewMove(gen.p, backward, forward)) // Jump.
		}
		gen.addQuiet(NewMove(gen.p, backward, square)) // Push.
	}

	// Knights.
	for bm := p.outposts[knight(our)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		for bm := knightMoves[square] & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, square, bm.first()))
		}
	}

	// Bishops.
	for bm := p.outposts[bishop(our)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		for bm := p.bishopMoves(square) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, square, bm.first()))
		}
	}

	// Rooks.
	for bm := p.outposts[rook(our)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		for bm := p.rookMoves(square) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, square, bm.first()))
		}
	}

	// Queens.
	for bm := p.outposts[queen(our)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		for bm := (p.bishopMoves(square) | p.rookMoves(square)) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, square, bm.first()))
		}
	}

	// King.
	square := p.king[our&1]
	for bm := (kingMoves[square] & empty); bm.anyʔ(); bm = bm.pop() {
		gen.addQuiet(NewMove(p, square, bm.first()))
	}

	return gen
}
