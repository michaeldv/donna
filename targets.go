// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (p *Position) Targets(square int, piece Piece) (bitmask Bitmask) {
        kind, color := piece.kind(), piece.color()

        switch kind {
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
		eight := [2]int{ 8, -8 }[color]

		if p.board[2].isClear(square + eight) { // Can white pawn move up one square?
			bitmask.set(square + eight)
			if row == 1 && p.board[2].isClear(square + eight * 2) { // How about two squares?
				bitmask.set(square + eight * 2)
			}
		}
                //
                // If the last move set the en-passant square and it is diagonally adjacent
                // to the current pawn, then add en-passant to the pawn's attack targets.
                //
                if target := p.flags.enpassant; row == 4 && target != 0 {
                        if target == square + (eight - 1) || target == square + (eight + 1) {
                                bitmask |= Bit(target)
                        }
                }
        case KNIGHT:
                bitmask = knightMoves[square]
                bitmask.exclude(p.board[color])
        case BISHOP:
                magic := ((bishopMagic[square].mask & p.board[2]) * bishopMagic[square].magic) >> 55
                bitmask = bishopMagicMoves[square][magic] & ^p.board[color]
        case ROOK:
                magic := ((rookMagic[square].mask & p.board[2]) * rookMagic[square].magic) >> 52
                bitmask = rookMagicMoves[square][magic] & ^p.board[color]
        case QUEEN:
                magic := ((bishopMagic[square].mask & p.board[2]) * bishopMagic[square].magic) >> 55
                bitmask = bishopMagicMoves[square][magic] & ^p.board[color]
                magic = ((rookMagic[square].mask & p.board[2]) * rookMagic[square].magic) >> 52
                bitmask.combine(rookMagicMoves[square][magic] & ^p.board[color])
        case KING:
                bitmask = kingMoves[square]
                bitmask.exclude(p.board[color])
        }

        return bitmask
}
