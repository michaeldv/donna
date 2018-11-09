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
	color := p.color
	empty := ^p.board

	// Castles.
	if !p.inCheckʔ(color) {
		home := homeKing[color]
		kingside, queenside := p.canCastleʔ(color)
		if kingside {
			gen.addQuiet(NewCastle(p, home, home + 2))
		}
		if queenside {
			gen.addQuiet(NewCastle(p, home, home - 2))
		}
	}

	// Pawns.
	last := let(color == White, 7, 0)
	for bm := p.outposts[pawn(color)].up(color) & empty & ^maskRank[last]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		forward, backward := square + up[color], square - up[color]
		if rank(color, square) == 2 && p.pieces[forward].noneʔ() {
			gen.addQuiet(NewPawnMove(gen.p, backward, forward)) // Jump.
		}
		gen.addQuiet(NewPawnMove(gen.p, backward, square)) // Push.
	}

	// Knights.
	for bm := p.outposts[knight(color)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		for bm := knightMoves[square] & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, square, bm.first()))
		}
	}

	// Bishops.
	for bm := p.outposts[bishop(color)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		for bm := p.bishopMoves(square) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, square, bm.first()))
		}
	}

	// Rooks.
	for bm := p.outposts[rook(color)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		for bm := p.rookMoves(square) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, square, bm.first()))
		}
	}

	// Queens.
	for bm := p.outposts[queen(color)]; bm.anyʔ(); bm = bm.pop() {
		square := bm.first()
		for bm := (p.bishopMoves(square) | p.rookMoves(square)) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, square, bm.first()))
		}
	}

	// King.
	square := p.king[color]
	for bm := (kingMoves[square] & empty); bm.anyʔ(); bm = bm.pop() {
		gen.addQuiet(NewMove(p, square, bm.first()))
	}

	return gen
}
