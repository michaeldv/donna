// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (p *Position) Targets(square int) (bitmask Bitmask) {
        piece := p.pieces[square]
        kind, color := piece.Kind(), piece.Color()

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
		row := Row(square)
		if color == WHITE {
			if p.board[2].IsClear(square + 8) { // Can white pawn move up one square?
				bitmask.Set(square + 8)
				if row == 1 && p.board[2].IsClear(square + 16) { // How about two squares?
					bitmask.Set(square + 16)
				}
			}
		} else if p.board[2].IsClear(square - 8) { // Can black pawn move up one square?
			bitmask.Set(square - 8)
			if row == 6 && p.board[2].IsClear(square - 16) { // How about two squares?
				bitmask.Set(square - 16)
			}
		}
                //
                // If the last move set the en-passant square and it is diagonally adjacent
                // to the current pawn, then add en-passant to the pawn's attack targets.
                //
                if target := p.enpassant; target != 0 {
                        if (color == WHITE && (target == square+7 || target == square+9)) || // Up/left or up/right a square.
                           (color == BLACK && (target == square-9 || target == square-7)) {  // Down/left or down/right a square.
                                bitmask |= Shift(target)
                        }
                }
        case KNIGHT:
                bitmask = knightMoves[square]
                bitmask.Exclude(p.board[color])
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
                bitmask.Combine(rookMagicMoves[square][magic] & ^p.board[color])
        case KING:
                bitmask = kingMoves[square]
                bitmask.Exclude(p.board[color])
        }

        return bitmask
}
