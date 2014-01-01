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
        enpassant int           // En-passant square caused by previous move.
        color     int           // Side to make next move.
        stage     int           // Game stage (256 in the initial position).
        inCheck   bool          // Is our king under attack?
        can00     [2]bool       // Is king-side castle allowed?
        can000    [2]bool       // Is queen-side castle allowed?
}

func NewPosition(game *Game, pieces [64]Piece) *Position {
        p := new(Position)
        p.game = game
        p.color = p.game.current
        p.pieces = pieces
        p.can00[WHITE]  = p.pieces[E1] == King(WHITE) && p.pieces[H1] == Rook(WHITE)
        p.can00[BLACK]  = p.pieces[E8] == King(BLACK) && p.pieces[H8] == Rook(BLACK)
        p.can000[WHITE] = p.pieces[E1] == King(WHITE) && p.pieces[A1] == Rook(WHITE)
        p.can000[BLACK] = p.pieces[E8] == King(BLACK) && p.pieces[A8] == Rook(BLACK)

        return p.setupPieces().setupAttacks()
}

func (p *Position) setupPieces() *Position {
        for square, piece := range p.pieces {
                if piece != 0 {
                        p.outposts[piece].Set(square)
                        p.board[piece.Color()].Set(square)
                        p.count[piece]++
                }
        }
        //
        // Combined board starts off with white pieces and adds black ones.
        //
        p.board[2] = p.board[WHITE] | p.board[BLACK]
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
        var kingSquare [2]int

        board := p.board[2]
        for board.IsNotEmpty() {
                square := board.FirstSet()
                piece := p.pieces[square]
                p.targets[square] = p.Targets(square)
                p.attacks[piece.Color()].Combine(p.targets[square])
                if piece.IsKing() {
                        kingSquare[piece.Color()] = square
                }
                board.Clear(square)
        }
        //
        // Now that we have attack targets for both kings adjust them to make sure the
        // kings don't stomp on each other.
        //
        kingTargets := [2]Bitmask{ p.targets[kingSquare[WHITE]], p.targets[kingSquare[BLACK]] }
        p.targets[kingSquare[WHITE]].Exclude(kingTargets[BLACK])
        p.targets[kingSquare[BLACK]].Exclude(kingTargets[WHITE])
        //
        // Combined board starts off with white pieces and adds black ones.
        //
        p.attacks[2] = p.attacks[WHITE] | p.attacks[BLACK]
        //
        // Is our king being attacked?
        //
        p.inCheck = p.isCheck(p.color)

        //Log("\n%s\n", p)
        return p
}

func (p *Position) MakeMove(move *Move) (position *Position) {
        color := move.piece.Color()
        position = new(Position)
        position.game = p.game
        position.color = p.color^1
        position.can00 = p.can00
        position.can000 = p.can000
        position.outposts = p.outposts
        position.board = p.board
        position.count = p.count
        position.enpassant = 0

        position.pieces = p.pieces
        position.pieces[move.from] = 0
        position.pieces[move.to] = move.piece
        position.board[color].Clear(move.from)
        position.board[color].Set(move.to)
        position.outposts[move.piece].Clear(move.from)
        position.outposts[move.piece].Set(move.to)

        if kind := move.piece.Kind(); kind == KING {
                if move.isCastle() {
                        switch move.to {
                        case G1:
                                position.pieces[H1], position.pieces[F1] = 0, Rook(WHITE)
                                position.board[WHITE].Clear(H1)
                                position.board[WHITE].Set(F1)
                                position.outposts[Rook(WHITE)].Clear(H1)
                                position.outposts[Rook(WHITE)].Set(F1)
                        case C1:
                                position.pieces[A1], position.pieces[D1] = 0, Rook(WHITE)
                                position.board[WHITE].Clear(A1)
                                position.board[WHITE].Set(D1)
                                position.outposts[Rook(WHITE)].Clear(A1)
                                position.outposts[Rook(WHITE)].Set(D1)
                        case G8:
                                position.pieces[H8], position.pieces[F8] = 0, Rook(BLACK)
                                position.board[BLACK].Clear(H8)
                                position.board[BLACK].Set(F8)
                                position.outposts[Rook(BLACK)].Clear(H8)
                                position.outposts[Rook(BLACK)].Set(F8)
                        case C8:
                                position.pieces[A8], position.pieces[D8] = 0, Rook(BLACK)
                                position.board[BLACK].Clear(A8)
                                position.board[BLACK].Set(D8)
                                position.outposts[Rook(BLACK)].Clear(A8)
                                position.outposts[Rook(BLACK)].Set(D8)
                        }
                }
                position.can00[color], position.can000[color] = false, false
        } else {
                if kind == PAWN {
                        if move.isEnpassant(p.outposts[Pawn(color^1)]) {
                                if color == WHITE {
                                        position.enpassant = move.from + 8
                                } else {
                                        position.enpassant = move.from - 8
                                }
                        } else if move.isEnpassantCapture(p.enpassant) { // Take out the en-passant pawn.
                                position.count[Pawn(color^1)]--
                                if color == WHITE {
                                        position.pieces[move.to - 8] = 0
                                        position.board[color^1].Clear(move.to - 8)
                                        position.outposts[Pawn(color^1)].Clear(move.to - 8)
                                } else {
                                        position.pieces[move.to + 8] = 0
                                        position.board[color^1].Clear(move.to + 8)
                                        position.outposts[Pawn(color^1)].Clear(move.to + 8)
                                }
                        } else if move.promoted != 0 { // Replace pawn with the promoted piece.
                                position.pieces[move.to] = move.promoted
                                position.board[color].Set(move.to)
                                position.outposts[move.promoted].Set(move.to)
                                position.count[Pawn(color)]--
                                position.count[move.promoted]++
                        }
                }
                if position.can00[p.color] {
                        if p.color == WHITE {
                                position.can00[WHITE] = position.pieces[H1] == Rook(WHITE) //&& position.pieces[E1] == King(WHITE)
                        } else {
                                position.can00[BLACK] = position.pieces[H8] == Rook(BLACK) //&& position.pieces[E8] == King(BLACK)
                        }
                }
                if position.can000[p.color] {
                        if p.color == WHITE {
                                position.can000[WHITE] = position.pieces[A1] == Rook(WHITE) //&& position.pieces[E1] == King(WHITE)
                        } else {
                                position.can000[BLACK] = position.pieces[A8] == Rook(BLACK) //&& position.pieces[E8] == King(BLACK)
                        }
                }
        }
        if move.captured != 0 {
                position.board[color^1].Clear(move.to)
                position.outposts[move.captured].Clear(move.to)
                position.count[move.captured]--
        }

        position.board[2] = position.board[WHITE] | position.board[BLACK]
        p.stage = 2 * (p.count[Pawn(WHITE)]   + p.count[Pawn(BLACK)])   +
                  6 * (p.count[Knight(WHITE)] + p.count[Knight(BLACK)]) +
                 12 * (p.count[Bishop(WHITE)] + p.count[Bishop(BLACK)]) +
                 16 * (p.count[Rook(WHITE)]   + p.count[Rook(BLACK)])   +
                 44 * (p.count[Queen(WHITE)]  + p.count[Queen(BLACK)])

        position.setupAttacks()
        return
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

func (p *Position) isKingSideCastleAllowed() bool {
        if p.color == WHITE {
                return p.can00[WHITE] && p.pieces[F1] == 0 && p.pieces[G1] == 0 &&
                       CASTLE_KING_WHITE & p.attacks[BLACK] == 0
        } else {
                return p.can00[BLACK] && p.pieces[F8] == 0 && p.pieces[G8] == 0 &&
                       CASTLE_KING_BLACK & p.attacks[WHITE] == 0
        }
}

func (p *Position) isQueenSideCastleAllowed() bool {
        if p.color == WHITE {
                return p.can000[WHITE] && p.pieces[D1] == 0 && p.pieces[C1] == 0 && p.pieces[B1] == 0 &&
                       CASTLE_QUEEN_WHITE & p.attacks[BLACK] == 0
        } else {
                return p.can000[BLACK] && p.pieces[D8] == 0 && p.pieces[C8] == 0 && p.pieces[B8] == 0 &&
                       CASTLE_QUEEN_BLACK & p.attacks[WHITE] == 0
        }
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
