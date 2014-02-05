// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (gen *MoveGen) GenerateEvasions() *MoveGen {
        color := gen.p.color
        enemy := gen.p.color^1
        square := gen.p.outposts[King(color)].first()
        pawn, knight, bishop, rook, queen := Pawn(enemy), Knight(enemy), Bishop(enemy), Rook(enemy), Queen(enemy)
        //
        // Find out what pieces are checking the king. Usually it's a single
        // piece but double check is also a possibility.
        //
        checkers := maskPawn[enemy][square] & gen.p.outposts[pawn]
        checkers |= gen.p.Targets(square, Knight(color)) & gen.p.outposts[knight]
        checkers |= gen.p.Targets(square, Bishop(color)) & (gen.p.outposts[bishop] | gen.p.outposts[queen])
        checkers |= gen.p.Targets(square, Rook(color)) & (gen.p.outposts[rook] | gen.p.outposts[queen])
        //
        // Generate possible king retreats first, i.e. moves to squares not
        // occupied by friendly pieces and not attacked by the opponent.
        //
        retreats := gen.p.targets[square] & ^gen.p.attacks[enemy]
        //
        // If the attacking piece is bishop, rook, or queen then exclude the
        // square behind the king using evasion mask. Note that knight's
        // evasion mask is full board so we only check if the attacking piece
        // is not a pawn.
        //
        attackSquare := checkers.pop()
        if gen.p.pieces[attackSquare] != pawn {
                retreats &= maskEvade[square][attackSquare]
        }
        //
        // If checkers mask is not empty then we've got double check and
        // retreat is the only option.
        //
        if checkers != 0 {
                attackSquare = checkers.first()
                if gen.p.pieces[attackSquare] != pawn {
                        retreats &= maskEvade[square][attackSquare]
                }
                for retreats != 0 {
                        gen.list[gen.tail].move = gen.p.NewMove(square, retreats.pop())
                        gen.tail++
                }
                return gen
        }
        //
        // Generate king retreats.
        //
        for retreats != 0 {
                gen.list[gen.tail].move = gen.p.NewMove(square, retreats.pop())
                gen.tail++
        }
        //
        // Pawn captures: do we have any pawns available that could capture
        // the attacking piece?
        //
        pawns := maskPawn[color][attackSquare] & gen.p.outposts[Pawn(color)]
        for pawns != 0 {
                gen.list[gen.tail].move = gen.p.NewMove(pawns.pop(), attackSquare)
                gen.tail++
        }
        //
        // Rare case when the check could be avoided by en-passant capture.
        // For example: Ke4, c5, e5 vs. Ke8, d7. Black's d7-d5+ could be
        // evaded by c5xd6 or e5xd6 en-passant captures.
        //
        if enpassant := attackSquare + eight[color]; gen.p.flags.enpassant == enpassant {
                pawns := maskPawn[color][enpassant] & gen.p.outposts[Pawn(color)]
                for pawns != 0 {
                        gen.list[gen.tail].move = gen.p.NewEnpassant(pawns.pop(), attackSquare + eight[color])
                        gen.tail++
                }
        }
        //
        // See if the check could be blocked.
        //
        block := maskBlock[square][attackSquare]
        //
        // Handle one square pawn pushes: promote to Queen when reaching last rank.
        //
        pawns = gen.p.pawnMovesMask(color) & block
        for pawns != 0 {
                to := pawns.pop(); from := to - eight[color]
                gen.list[gen.tail].move = gen.p.NewMove(from, to)
                if to >= A8 || to <= H1 {
                        gen.list[gen.tail].move.promote(QUEEN)
                }
                gen.tail++
        }
        //
        // Handle two square pawn pushes.
        //
        pawns = gen.p.pawnJumpsMask(color) & block
        for pawns != 0 {
                to := pawns.pop(); from := to - 2 * eight[color]
                gen.list[gen.tail].move = gen.p.NewMove(from, to)
                gen.tail++
        }
        //
        // What's left is to generate all possible knight, bishop, rook, and
        // queen moves that evade the check.
        //
        for _, kind := range [4]int{ KNIGHT, BISHOP, ROOK, QUEEN } {
                gen.addEvasion(Piece(kind|color), block)
        }

        return gen
}

func (gen *MoveGen) addEvasion(piece Piece, block Bitmask) {
        outposts := gen.p.outposts[piece]
        for outposts != 0 {
                from := outposts.pop()
                targets := gen.p.targets[from] & block
                for targets != 0 {
                        gen.list[gen.tail].move = gen.p.NewMove(from, targets.pop())
                        gen.tail++
                }
        }
}
