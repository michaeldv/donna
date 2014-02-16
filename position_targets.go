// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (p *Position) targetsMask(square int) Bitmask {
        if p.targets[square] == 0 {
                p.targets[square] = p.Targets(square, p.pieces[square])
        }
        return p.targets[square]
}

func (p *Position) attacksMask(color int) Bitmask {
        if p.attacks[color] == 0 {
                p.attacks[color] = p.pawnAttacks(color) | p.knightAttacks(color) | p.bishopAttacks(color) |
                                   p.rookAttacks(color) | p.queenAttacks(color) | p.kingAttacks(color)
        }
        return p.attacks[color]
}

func (p *Position) isAttacked(square, color int) bool {
        return (knightMoves[square] & p.outposts[Knight(color)]) != 0 ||
               (maskPawn[color][square] & p.outposts[Pawn(color)]) != 0 ||
               (kingMoves[square] & p.outposts[King(color)]) != 0 ||
               (p.rookMoves(square) & (p.outposts[Rook(color)] | p.outposts[Queen(color)])) != 0 ||
               (p.bishopMoves(square) & (p.outposts[Bishop(color)] | p.outposts[Queen(color)])) != 0
}

func (p *Position) rookMoves(square int) Bitmask {
        magic := ((rookMagic[square].mask & p.board[2]) * rookMagic[square].magic) >> 52
        return rookMagicMoves[square][magic]
}

func (p *Position) bishopMoves(square int) Bitmask {
        magic := ((bishopMagic[square].mask & p.board[2]) * bishopMagic[square].magic) >> 55
        return bishopMagicMoves[square][magic]
}

func (p *Position) Targets(square int, piece Piece) (bitmask Bitmask) {
        switch kind, color := piece.kind(), piece.color(); kind {
        case PAWN:
                bitmask = pawnMoves[color][square] & p.board[color^1]
                //
		// If the square in front of the pawn is empty then add it as possible
		// target.
		//
		// If the pawn is in its initial position and two squares in front of
		// the pawn are empty then add the second square as possible target.
                //
		row := RelRow(square, color)
		shift := eight[color]

		if p.board[2].isClear(square + shift) { // Can white pawn move up one square?
			bitmask.set(square + shift)
			if row == 1 && p.board[2].isClear(square + shift * 2) { // How about two squares?
				bitmask.set(square + shift * 2)
			}
		}
                //
                // If the last move set the en-passant square and it is diagonally adjacent
                // to the current pawn, then add en-passant to the pawn's attack targets.
                //
                if row == 4 && p.flags.enpassant != 0 {
                        if maskPawn[color][p.flags.enpassant] & Bit(square) != 0 {
                                bitmask |= Bit(p.flags.enpassant)
                        }
                }
        case KNIGHT:
                bitmask = knightMoves[square] & ^p.board[color]
        case BISHOP:
                bitmask = p.bishopMoves(square) & ^p.board[color]
        case ROOK:
                bitmask = p.rookMoves(square) & ^p.board[color]
        case QUEEN:
                bitmask = (p.bishopMoves(square) | p.rookMoves(square)) & ^p.board[color]
        case KING:
                bitmask = kingMoves[square] & ^p.board[color]
        }

        return bitmask
}

func (p *Position) pawnAttacks(color int) (bitmask Bitmask) {
        if color == White {
                bitmask  = (p.outposts[Pawn(White)] & ^maskFile[0]) << 7
                bitmask |= (p.outposts[Pawn(White)] & ^maskFile[7]) << 9
        } else {
                bitmask  = (p.outposts[Pawn(Black)] & ^maskFile[0]) >> 9
                bitmask |= (p.outposts[Pawn(Black)] & ^maskFile[7]) >> 7
        }
        return
}

func (p *Position) knightAttacks(color int) (bitmask Bitmask) {
        outposts := p.outposts[Knight(color)]
        for outposts != 0 {
                bitmask |= knightMoves[outposts.pop()]
        }
        return
}

func (p *Position) bishopAttacks(color int) (bitmask Bitmask) {
        outposts := p.outposts[Bishop(color)]
        for outposts != 0 {
                bitmask |= p.bishopMoves(outposts.pop())
        }
        return
}

func (p *Position) rookAttacks(color int) (bitmask Bitmask) {
        outposts := p.outposts[Rook(color)]
        for outposts != 0 {
                bitmask |= p.rookMoves(outposts.pop())
        }
        return
}

func (p *Position) queenAttacks(color int) (bitmask Bitmask) {
        outposts := p.outposts[Queen(color)]
        for outposts != 0 {
                square := outposts.pop()
                bitmask |= p.rookMoves(square)  | p.bishopMoves(square)
        }
        return
}

func (p *Position) kingAttacks(color int) Bitmask {
         return kingMoves[p.outposts[King(color)].first()]
}
