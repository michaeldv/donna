// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(
        `bytes`
)

var best [16][16]*Move // Assuming max depth = 4 which makes it 8 plies.
var bestlen [16]int

type Position struct {
        game      *Game
        pieces    [64]Piece     // Array of 64 squares with pieces on them.
        targets   [64]Bitmask   // Attack targets for each piece on the board.
        board     [3]Bitmask    // [0] white pieces only, [1] black pieces, and [2] all pieces.
        attacks   [3]Bitmask    // [0] all squares attacked by white, [1] by black, [2] either white or black.
        outposts  [16]Bitmask   // Bitmasks of each piece on the board, ex. white pawns, black king, etc.
        count     [16]int       // Counts of each piece on the board, ex. white pawns: 6, etc.
        enpassant Bitmask       // En-passant opportunity caused by previous move.
        color     int           // Side to make next move.
        stage     int           // Game stage (256 in the initial position).
        inCheck   bool          // Is our king under attack?
}

func NewPosition(game *Game, pieces [64]Piece, color int, enpassant Bitmask) *Position {
        position := new(Position)

        position.game = game
        position.pieces = pieces
        position.color = color
        position.enpassant = enpassant

        return position.setupPieces().setupAttacks()
}

func (p *Position) setupPieces() *Position {
        for i, piece := range p.pieces {
                if piece != 0 {
                        p.outposts[piece] = *new(Bitmask).Set(i)
                        p.board[piece.Color()].Set(i)
                        p.count[piece]++
                }
        }
        //
        // Combined board starts off with white pieces and adds black ones.
        //
        p.board[2] = p.board[WHITE]
        p.board[2].Combine(p.board[BLACK])
        //
        // Determine game stage.
        //
        p.stage = 2 * (p.count[Pawn(WHITE)]   + p.count[Pawn(BLACK)])   +
                  6 * (p.count[Knight(WHITE)] + p.count[Knight(BLACK)]) +
                 12 * (p.count[Bishop(WHITE)] + p.count[Bishop(BLACK)]) +
                 16 * (p.count[Rook(WHITE)]   + p.count[Rook(BLACK)])   +
                 44 * (p.count[Queen(WHITE)]  + p.count[Queen(BLACK)])
        return p
}

func (p *Position) setupAttacks() *Position {
        var king [2]int
        for i, piece := range p.pieces {
                if piece != 0 {
                        p.targets[i] = *p.Targets(i)
                        p.attacks[piece.Color()].Combine(p.targets[i])
                        if piece.IsKing() {
                                king[piece.Color()] = i
                        }
                }
        }
        //
        // Now that we have attack targets for both kings adjust them to make sure the
        // kings don't stomp on each other.
        //
        white_king_targets, black_king_targets := p.targets[king[WHITE]], p.targets[king[BLACK]]
        p.targets[king[WHITE]].Exclude(black_king_targets)
        p.targets[king[BLACK]].Exclude(white_king_targets)
        //
        // Combined board starts off with white pieces and adds black ones.
        //
        p.attacks[2] = p.attacks[0]
        p.attacks[2].Combine(p.attacks[1])
        //
        // Is our king being attacked?
        //
        p.inCheck = p.isCheck(p.color)

        //Log("\n%s\n", p)
        return p
}

func (p *Position) MakeMove(move *Move) *Position {
        pieces := p.pieces
        enpassant := Bitmask(0)

        pieces[move.From] = 0
        pieces[move.To] = move.Piece

        // Check if we need to update en-passant bitmask.
        if move.IsTwoSquarePawnAdvance() {
                if p.color == WHITE {
                        enpassant.Set(move.From + 8)
                } else {
                        enpassant.Set(move.From - 8)
                }
        } else if move.IsCrossing(p.enpassant) { // Take out the en-passant pawn.
                if p.color == WHITE {
                        pieces[move.To - 8] = Piece(NONE)
                } else {
                        pieces[move.To + 8] = Piece(NONE)
                }
        } else if move.IsCastle() {
                switch move.To {
                case G1:
                        pieces[H1], pieces[F1] = 0, Rook(WHITE)
                case C1:
                        pieces[A1], pieces[D1] = 0, Rook(WHITE)
                case G8:
                        pieces[H8], pieces[F8] = 0, Rook(BLACK)
                case C8:
                        pieces[A8], pieces[D8] = 0, Rook(BLACK)
                }
        } else if move.Promoted != 0 { // Replace pawn with the promoted piece.
                pieces[move.To] = move.Promoted
        }

        return NewPosition(p.game, pieces, p.color^1, enpassant)
}

func (p *Position) isCheck(color int) bool {
        king := p.outposts[King(color)]
        return king.Intersect(p.attacks[color^1]).IsNotEmpty()
}


func (p *Position) saveBest(ply int, move *Move) {
        best[ply][ply] = move
        bestlen[ply] = ply + 1
        for i := ply + 1; i < bestlen[ply + 1]; i++ {
                best[ply][i] = best[ply + 1][i]
                bestlen[ply]++
        }
}

func (p *Position) isPawnPromotion(piece Piece, target int) bool {
        return piece.IsPawn() && ((piece.IsWhite() && target >= A8) || (piece.IsBlack() && target <= H1))
}

func (p *Position) tryCastle() (castle *Move) {
        if p.isKingSideCastleAllowed() {
                if p.color == WHITE {
                        castle = NewMove(E1, G1, King(WHITE), 0)
                } else {
                        castle = NewMove(E8, G8, King(BLACK), 0)
                }
        } else if p.isQueenSideCastleAllowed() {
                if p.color == WHITE {
                        castle = NewMove(E1, C1, King(WHITE), 0)
                } else {
                        castle = NewMove(E8, C8, King(BLACK), 0)
                }
        }
        return
}

func (p *Position) isKingSideCastleAllowed() bool {
        white, black := CASTLE_KING_WHITE, CASTLE_KING_BLACK

        return p.game.players[p.color].Can00 &&
               ((p.color == WHITE &&
                       p.pieces[E1] == King(WHITE) &&
                       p.pieces[F1] == 0 &&
                       p.pieces[G1] == 0 &&
                       p.pieces[H1] == Rook(WHITE) &&
                       white.Intersect(p.attacks[BLACK]).IsEmpty()) ||
                (p.color == BLACK &&
                       p.pieces[E8] == King(BLACK) &&
                       p.pieces[F8] == 0 &&
                       p.pieces[G8] == 0 &&
                       p.pieces[H8] == Rook(BLACK) &&
                       black.Intersect(p.attacks[WHITE]).IsEmpty()))
}

func (p *Position) isQueenSideCastleAllowed() bool {
        white, black := CASTLE_QUEEN_WHITE, CASTLE_QUEEN_BLACK

        return p.game.players[p.color].Can000 &&
               ((p.color == WHITE &&
                       p.pieces[E1] == King(WHITE) &&
                       p.pieces[D1] == 0 &&
                       p.pieces[C1] == 0 &&
                       p.pieces[B1] == 0 &&
                       p.pieces[A1] == Rook(WHITE) &&
                       white.Intersect(p.attacks[BLACK]).IsEmpty()) ||
                (p.color == BLACK &&
                       p.pieces[E8] == King(BLACK) &&
                       p.pieces[D8] == 0 &&
                       p.pieces[C8] == 0 &&
                       p.pieces[B8] == 0 &&
                       p.pieces[A8] == Rook(BLACK) &&
                       black.Intersect(p.attacks[WHITE]).IsEmpty()))
}

func (p *Position) String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h")
        if !p.inCheck {
                buffer.WriteString("\n")
        } else {
                buffer.WriteString("  Check to " + C(p.color) + "\n")
        }
	for row := 7;  row >= 0;  row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0;  col <= 7; col++ {
			square := Square(row, col)
			buffer.WriteByte(' ')
			if piece := p.pieces[square]; piece != 0 {
				buffer.WriteString(piece.String())
			} else {
				buffer.WriteString("\u22C5")
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
