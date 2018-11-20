// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (gen *MoveGen) generateRootMoves(p *Position) *MoveGen {
	gen.generateAllMoves(p)

	if !gen.onlyMoveʔ() {
		gen.validOnly().rank(Move(0))
	}

	return gen
}

// Copies last move returned by nextMove() to the top of the list shifting
// remaining moves down. Head/tail pointers remain unchanged.
func (gen *MoveGen) rearrangeRootMoves() *MoveGen {
	if gen.head > 0 {
		best := gen.list[gen.head - 1]
		copy(gen.list[1:], gen.list[0:gen.head - 1])
		gen.list[0] = best
	}

	return gen
}

func (gen *MoveGen) generateAllMoves(p *Position) *MoveGen {
	if p.inCheckʔ(p.color) {
		return gen.generateEvasions(p)
	}

	return gen.generateMoves(p)
}

func (gen *MoveGen) generateMoves(p *Position) *MoveGen {
	return gen.pawnMoves(p).pieceMoves(p).kingMoves(p)
}

func (gen *MoveGen) pawnMoves(p *Position) *MoveGen {
	our := p.color

	for bm := p.outposts[pawn(our)]; bm.anyʔ(); bm = bm.pop() {
		sq := bm.first()
		gen.movePawn(p, sq, p.targets(sq))
	}

	return gen
}

// Go over all pieces except pawns and the king.
func (gen *MoveGen) pieceMoves(p *Position) *MoveGen {
	our := p.color

	for bm := p.outposts[our] ^ p.outposts[pawn(our)] ^ p.outposts[king(our)]; bm.anyʔ(); bm = bm.pop() {
		sq := bm.first()
		gen.movePiece(p, sq, p.targets(sq))
	}

	return gen
}

func (gen *MoveGen) kingMoves(p *Position) *MoveGen {
	our := p.color

	if p.outposts[king(our)].anyʔ() {
		sq := p.king[our]
		gen.moveKing(p, sq, p.targets(sq))

		kside, qside := p.canCastleʔ(our)
		if kside {
			gen.moveKing(p, sq, bit(Square(G1 + 56 * our)))
		}
		if qside {
			gen.moveKing(p, sq, bit(Square(C1 + 56 * our)))
		}
	}

	return gen
}

func (gen *MoveGen) movePawn(p *Position, sq Square, targets Bitmask) *MoveGen {
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		target := bm.first()
		if target > H1 && target < A8 {
			gen.add(NewMove(p, sq, target))
		} else { // Promotion.
			mQ, mR, mB, mN := NewPromotion(p, sq, target)
			gen.add(mQ).add(mN).add(mR).add(mB)
		}
	}

	return gen
}

func (gen *MoveGen) moveKing(p *Position, sq Square, targets Bitmask) *MoveGen {
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		target := bm.first()
		if sq.upto(target) == 2 {
			gen.add(NewCastle(p, sq, target))
		} else {
			gen.add(NewMove(p, sq, target))
		}
	}

	return gen
}

func (gen *MoveGen) movePiece(p *Position, sq Square, targets Bitmask) *MoveGen {
	for bm := targets; bm.anyʔ(); bm = bm.pop() {
		gen.add(NewMove(p, sq, bm.first()))
	}

	return gen
}
