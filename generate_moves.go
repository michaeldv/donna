// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) GenerateMoves() *MoveGen {
        color := gen.p.color
        return gen.pawnMoves(color).pieceMoves(color).kingMoves(color)
}

func (gen *MoveGen) pawnMoves(color int) *MoveGen {
        for pawns := gen.p.outposts[Pawn(color)]; pawns != 0; {
                square := pawns.pop()
                gen.movePawn(square, gen.p.targetsMask(square))
        }
        return gen
}

func (gen *MoveGen) pieceMoves(color int) *MoveGen {
        for _, kind := range [4]int{ KNIGHT, BISHOP, ROOK, QUEEN } {
                outposts := gen.p.outposts[Piece(kind|color)]
                for outposts != 0 {
                        square := outposts.pop()
                        gen.movePiece(square, gen.p.targetsMask(square))
                }
        }
        return gen
}

func (gen *MoveGen) kingMoves(color int) *MoveGen {
        if king := gen.p.outposts[King(color)]; king != 0 {
                square := king.pop()
                gen.moveKing(square, gen.p.targetsMask(square))
                if !gen.p.isInCheck(gen.p.color) {
                        if gen.p.can00(color) {
                                gen.moveKing(square, Bit(G1 + 56 * color))
                        }
                        if gen.p.can000(color) {
                                gen.moveKing(square, Bit(C1 + 56 * color))
                        }
                }
        }
        return gen
}

func (gen *MoveGen) movePawn(square int, targets Bitmask) *MoveGen {
        for targets != 0 {
                target := targets.pop()
                if target > H1 && target < A8 {
                        gen.add(gen.p.pawnMove(square, target))
                } else { // Promotion.
                        m1, m2, m3, m4 := gen.p.pawnPromotion(square, target)
                        gen.add(m1); gen.add(m2); gen.add(m3); gen.add(m4)
                }
        }
        return gen
}

func (gen *MoveGen) moveKing(square int, targets Bitmask) *MoveGen {
        for targets != 0 {
                target := targets.pop()
                if square == homeKing[gen.p.color] && Abs(square - target) == 2 {
                        gen.add(gen.p.NewCastle(square, target))
                } else {
                        gen.add(gen.p.NewMove(square, target))
                }
        }
        return gen
}

func (gen *MoveGen) movePiece(square int, targets Bitmask) *MoveGen {
        for targets != 0 {
                gen.add(gen.p.NewMove(square, targets.pop()))
        }
        return gen
}
