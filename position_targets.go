// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (p *Position) rookMoves(square int) Bitmask {
        magic := ((rookMagic[square].mask & p.board) * rookMagic[square].magic) >> 52
        return rookMagicMoves[square][magic]
}

func (p *Position) bishopMoves(square int) Bitmask {
        magic := ((bishopMagic[square].mask & p.board) * bishopMagic[square].magic) >> 55
        return bishopMagicMoves[square][magic]
}

func (p *Position) targets(square int) Bitmask {
        return p.targetsFor(square, p.pieces[square])
}

func (p *Position) targetsFor(square int, piece Piece) (bitmask Bitmask) {
        switch kind, color := piece.kind(), piece.color(); kind {
        case Pawn:
                bitmask = pawnMoves[color][square] & p.outposts[color^1]
                //
		// If the square in front of the pawn is empty then add it as possible
		// target.
		//
		// If the pawn is in its initial position and two squares in front of
		// the pawn are empty then add the second square as possible target.
                //
		row := RelRow(square, color)
		shift := eight[color]

		if p.board.isClear(square + shift) { // Can white pawn move up one square?
			bitmask.set(square + shift)
			if row == 1 && p.board.isClear(square + shift * 2) { // How about two squares?
				bitmask.set(square + shift * 2)
			}
		}
                //
                // If the last move set the en-passant square and it is diagonally adjacent
                // to the current pawn, then add en-passant to the pawn's attack targets.
                //
                if row == 4 && p.enpassant != 0 {
                        if maskPawn[color][p.enpassant] & bit[square] != 0 {
                                bitmask |= bit[p.enpassant]
                        }
                }
        case Knight:
                bitmask = knightMoves[square] & ^p.outposts[color]
        case Bishop:
                bitmask = p.bishopMoves(square) & ^p.outposts[color]
        case Rook:
                bitmask = p.rookMoves(square) & ^p.outposts[color]
        case Queen:
                bitmask = (p.bishopMoves(square) | p.rookMoves(square)) & ^p.outposts[color]
        case King:
                bitmask = kingMoves[square] & ^p.outposts[color]
        }

        return bitmask
}

func (p *Position) attacks(color int) (bitmask Bitmask) {
        bitmask = p.pawnAttacks(color) | p.knightAttacks(color) | p.kingAttacks(color)

        outposts := p.outposts[bishop(color)] | p.outposts[queen(color)]
        for outposts != 0 {
                bitmask |= p.bishopMoves(outposts.pop())
        }

        outposts = p.outposts[rook(color)] | p.outposts[queen(color)]
        for outposts != 0 {
                bitmask |= p.rookMoves(outposts.pop())
        }
        return
}

func (p *Position) isAttacked(square, color int) bool {
        return (knightMoves[square] & p.outposts[knight(color)]) != 0 ||
               (maskPawn[color][square] & p.outposts[pawn(color)]) != 0 ||
               (kingMoves[square] & p.outposts[king(color)]) != 0 ||
               (p.rookMoves(square) & (p.outposts[rook(color)] | p.outposts[queen(color)])) != 0 ||
               (p.bishopMoves(square) & (p.outposts[bishop(color)] | p.outposts[queen(color)])) != 0
}

func (p *Position) pawnAttacks(color int) (bitmask Bitmask) {
        if color == White {
                bitmask  = (p.outposts[Pawn] & ^maskFile[0]) << 7
                bitmask |= (p.outposts[Pawn] & ^maskFile[7]) << 9
        } else {
                bitmask  = (p.outposts[BlackPawn] & ^maskFile[0]) >> 9
                bitmask |= (p.outposts[BlackPawn] & ^maskFile[7]) >> 7
        }
        return
}

func (p *Position) knightAttacks(color int) (bitmask Bitmask) {
        outposts := p.outposts[knight(color)]
        for outposts != 0 {
                bitmask |= knightMoves[outposts.pop()]
        }
        return
}

func (p *Position) bishopAttacks(color int) (bitmask Bitmask) {
        outposts := p.outposts[bishop(color)]
        for outposts != 0 {
                bitmask |= p.bishopMoves(outposts.pop())
        }
        return
}

func (p *Position) rookAttacks(color int) (bitmask Bitmask) {
        outposts := p.outposts[rook(color)]
        for outposts != 0 {
                bitmask |= p.rookMoves(outposts.pop())
        }
        return
}

func (p *Position) queenAttacks(color int) (bitmask Bitmask) {
        outposts := p.outposts[queen(color)]
        for outposts != 0 {
                square := outposts.pop()
                bitmask |= p.rookMoves(square)  | p.bishopMoves(square)
        }
        return
}

func (p *Position) kingAttacks(color int) Bitmask {
         return kingMoves[p.king[color]]
}
