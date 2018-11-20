// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (gen *MoveGen) addQuiet(move Move) *MoveGen {
	gen.list = append(gen.list, MoveWithScore{move, game.good(move)})

	return gen
}

// Generates pseudo-legal moves that preserve material balance, i.e.
// no captures or pawn promotions are allowed.
func (gen *MoveGen) generateQuiets(p *Position) *MoveGen {
	our := p.color; their := our^1
	empty := ^p.board

	// Castles.
	if !p.inCheckʔ(our) {
		home := homeKing[our]
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
		sq := bm.first()
		forward, backward := sq.push(our), sq.push(their)
		if sq.rank(our) == 2 && p.pieces[forward].noneʔ() {
			gen.addQuiet(NewMove(p, backward, forward)) // Jump.
		}
		gen.addQuiet(NewMove(p, backward, sq)) // Push.
	}

	// Knights.
	for bm := p.outposts[knight(our)]; bm.anyʔ(); bm = bm.pop() {
		sq := bm.first()
		for bm := knightMoves[sq] & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, sq, bm.first()))
		}
	}

	// Bishops.
	for bm := p.outposts[bishop(our)]; bm.anyʔ(); bm = bm.pop() {
		sq := bm.first()
		for bm := p.bishopMoves(sq) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, sq, bm.first()))
		}
	}

	// Rooks.
	for bm := p.outposts[rook(our)]; bm.anyʔ(); bm = bm.pop() {
		sq := bm.first()
		for bm := p.rookMoves(sq) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, sq, bm.first()))
		}
	}

	// Queens.
	for bm := p.outposts[queen(our)]; bm.anyʔ(); bm = bm.pop() {
		sq := bm.first()
		for bm := (p.bishopMoves(sq) | p.rookMoves(sq)) & empty; bm.anyʔ(); bm = bm.pop() {
			gen.addQuiet(NewMove(p, sq, bm.first()))
		}
	}

	// King.
	sq := p.king[our]
	for bm := (kingMoves[sq] & empty); bm.anyʔ(); bm = bm.pop() {
		gen.addQuiet(NewMove(p, sq, bm.first()))
	}

	return gen
}
