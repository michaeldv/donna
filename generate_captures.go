// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

func (gen *MoveGen) generateCaptures(p *Position) *MoveGen {
	return gen.pawnCaptures(p).pieceCaptures(p)
}

// Generates all pseudo-legal pawn captures and promotions.
func (gen *MoveGen) pawnCaptures(p *Position) *MoveGen {
	our, their := p.colors()
	opponent := p.outposts[their]

	for pawns := p.outposts[pawn(our)]; pawns.anyʔ(); pawns = pawns.pop() {
		sq := pawns.first()

		// For pawns on files 2-6 the moves include captures only,
		// while for pawns on the 7th file the moves include captures
		// as well as queen promotion.
		if sq.rank(our) != A7H7 {
			gen.movePawn(p, sq, p.targets(sq) & opponent)
		} else {
			for bm := p.targets(sq); bm.anyʔ(); bm = bm.pop() {
				target := bm.first()
				//- mQ, _, _, _ := NewPromotion(p, sq, target)
				//- gen.add(mQ)
				gen.add(NewMove(p, sq, target).promote(Queen))
			}
		}
	}

	return gen
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (gen *MoveGen) pieceCaptures(p *Position) *MoveGen {
	our, their := p.colors()
	opponent := p.outposts[their]

	for bm := p.outposts[our] ^ p.outposts[pawn(our)] ^ p.outposts[king(our)]; bm.anyʔ(); bm = bm.pop() {
		sq := bm.first()
		gen.movePiece(p, sq, p.targets(sq) & opponent)
	}

	if p.outposts[king(our)].anyʔ() {
		sq := p.king[our]
		gen.moveKing(p, sq, kingMoves[sq] & opponent)
	}

	return gen
}
