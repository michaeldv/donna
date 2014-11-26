// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) generateCaptures() *MoveGen {
	color := gen.p.color
	return gen.pawnCaptures(color).pieceCaptures(color)
}

// Generates all pseudo-legal pawn captures and promotions.
func (gen *MoveGen) pawnCaptures(color uint8) *MoveGen {
	enemy := gen.p.outposts[color^1]
	pawns := gen.p.outposts[pawn(color)]

	for pawns != 0 {
		square := pawns.pop()

		// For pawns on files 2-6 the moves include captures only,
		// while for pawns on the 7th file the moves include captures
		// as well as promotion on empty square in front of the pawn.
		if rank(color, square) != 6 {
			gen.movePawn(square, gen.p.targets(square) & enemy)
		} else {
			gen.movePawn(square, gen.p.targets(square))
		}
	}
	return gen
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (gen *MoveGen) pieceCaptures(color uint8) *MoveGen {
	enemy := gen.p.outposts[color^1]
	outposts := gen.p.outposts[color] & ^gen.p.outposts[pawn(color)] & ^gen.p.outposts[king(color)]

	for outposts != 0 {
		square := outposts.pop()
		gen.movePiece(square, gen.p.targets(square) & enemy)
	}
	if king := gen.p.outposts[king(color)]; king != 0 {
		square := king.pop()
		gen.moveKing(square, gen.p.targets(square) & enemy)
	}
	return gen
}
