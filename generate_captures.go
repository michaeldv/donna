// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) generateCaptures() *MoveGen {
	color := gen.p.color
	return gen.pawnCaptures(color).pieceCaptures(color)
}

// Generates all pseudo-legal pawn captures and promotions.
func (gen *MoveGen) pawnCaptures(color int) *MoveGen {
	enemy := gen.p.outposts[color^1]

	for pawns := gen.p.outposts[pawn(color)]; pawns.any(); pawns = pawns.pop() {
		square := pawns.first()

		// For pawns on files 2-6 the moves include captures only,
		// while for pawns on the 7th file the moves include captures
		// as well as queen promotion.
		if rank(color, square) != A7H7 {
			gen.movePawn(square, gen.p.targets(square) & enemy)
		} else {
			for bm := gen.p.targets(square); bm.any(); bm = bm.pop() {
				target := bm.first()
				mQ, _, _, _ := NewPromotion(gen.p, square, target)
				gen.add(mQ)
			}
		}
	}

	return gen
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (gen *MoveGen) pieceCaptures(color int) *MoveGen {
	enemy := gen.p.outposts[color^1]

	for bm := gen.p.outposts[color] ^ gen.p.outposts[pawn(color)] ^ gen.p.outposts[king(color)]; bm.any(); bm = bm.pop() {
		square := bm.first()
		gen.movePiece(square, gen.p.targets(square) & enemy)
	}
	if gen.p.outposts[king(color)].any() {
		square := int(gen.p.king[color])
		gen.moveKing(square, gen.p.targets(square) & enemy)
	}

	return gen
}
