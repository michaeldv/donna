// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) GenerateCaptures() *MoveGen {
        color := gen.p.color
        return gen.pawnCaptures(color).pieceCaptures(color)
}

// Generates all pseudo-legal pawn captures and promotions.
func (gen *MoveGen) pawnCaptures(color int) *MoveGen {
        enemy := color^1
        pawns := gen.p.outposts[Pawn(color)]

        for pawns != 0 {
                square := pawns.pop()
                //
                // For pawns on files 2-6 the moves include captures only,
                // while for pawns on the 7th file the moves include captures
                // as well as promotion on empty square in front of the pawn.
                //
                if row := RelRow(square, color); row != 6 {
                        gen.movePawn(square, gen.p.targets(square) & gen.p.board[enemy])
                } else {
                        gen.movePawn(square, gen.p.targets(square))
                }
        }
        return gen
}

// Generates all pseudo-legal captures by pieces other than pawn.
func (gen *MoveGen) pieceCaptures(color int) *MoveGen {
        enemy := color^1
        outposts := gen.p.board[color] & ^gen.p.outposts[Pawn(color)] & ^gen.p.outposts[King(color)]
        for outposts != 0 {
                square := outposts.pop()
                gen.movePiece(square, gen.p.targets(square) & gen.p.board[enemy])
        }
        if king := gen.p.outposts[King(color)]; king != 0 {
                square := king.pop()
                gen.moveKing(square, gen.p.targets(square) & gen.p.board[enemy])
        }
        return gen
}
