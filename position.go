// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`bytes`)

var killer [16][2]*Move
var best [16][16]*Move // Assuming max depth = 4 which makes it 8 plies.
var bestlen [16]int

type Position struct {
        game      *Game
        pieces    [64]Piece     // Array of 64 squares with pieces on them.
        targets   [64]Bitmask   // Attack targets for each piece on the board.
        board     [3]Bitmask    // [0] white pieces only, [1] black pieces, and [2] all pieces.
        attacks   [3]Bitmask    // [0] all squares attacked by white, [1] by black, [2] either white or black.
        outposts  [16]Bitmask   // Bitmasks of each piece on the board, ex. white pawns, black king, etc.
        count     [16]int       // counts of each piece on the board, ex. white pawns: 6, etc.
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

        return p.setupPieces().computeStage().setupAttacks()
}

func (p *Position) setupPieces() *Position {
        for square, piece := range p.pieces {
                if piece != 0 {
                        p.outposts[piece].set(square)
                        p.board[piece.color()].set(square)
                        p.count[piece]++
                }
        }
        return p
}

func (p *Position) computeStage() *Position {
        p.board[2] = p.board[WHITE] | p.board[BLACK]

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
        for board.isNotEmpty() {
                square := board.firstSet()
                piece := p.pieces[square]
                p.targets[square] = p.Targets(square)
                p.attacks[piece.color()].combine(p.targets[square])
                if piece.isKing() {
                        kingSquare[piece.color()] = square
                }
                board.clear(square)
        }
        //
        // Now that we have attack targets for both kings adjust them to make sure the
        // kings don't stomp on each other. Also, combine attacks bitmask and set the
        // flag is the king is being attacked.
        //
        p.updateKingTargets(kingSquare)
        p.attacks[2] = p.attacks[WHITE] | p.attacks[BLACK]
        p.inCheck = p.isCheck(p.color)

        return p
}

func (p *Position) updateKingTargets(kingSquare [2]int) *Position {
        p.targets[kingSquare[WHITE]].exclude(p.targets[kingSquare[BLACK]])
        p.targets[kingSquare[BLACK]].exclude(p.targets[kingSquare[WHITE]])
        //
        // Add castle jump targets if castles are allowed.
        //
        if kingSquare[p.color] == initialKingSquare[p.color] {
                if p.isKingSideCastleAllowed(p.color) {
                        p.targets[kingSquare[p.color]].set(kingSquare[p.color] + 2)
                }
                if p.isQueenSideCastleAllowed(p.color) {
                        p.targets[kingSquare[p.color]].set(kingSquare[p.color] - 2)
                }
        }

        return p
}

func (p *Position) clone() *Position {
        position := new(Position)

        position.game = p.game
        position.color = p.color
        position.board = p.board
        position.pieces = p.pieces
        position.outposts = p.outposts
        position.count = p.count
        position.can00 = p.can00
        position.can000 = p.can000
        position.enpassant = 0

        return position
}

func (p *Position) lift(move *Move) *Position {
        color := p.color
        if move.piece != 0 {
                color = move.piece.color()
        }

        p.pieces[move.from] = 0
        p.board[color].clear(move.from)
        p.outposts[move.piece].clear(move.from)
        return p
}

func (p *Position) put(move *Move) *Position {
        piece := move.piece
        if move.promoted != 0 {
                piece = move.promoted
        }

        p.pieces[move.to] = piece
        p.board[p.color].set(move.to)
        p.outposts[piece].set(move.to)

        if move.captured != 0 {
                p.board[p.color^1].clear(move.to)
                p.outposts[move.captured].clear(move.to)
                p.count[move.captured]--
        }
        return p
}

func (p *Position) make(move *Move) *Position {
        return p.lift(move).put(move)
}

func (p *Position) MakeMove(move *Move) *Position {
        eight := [2]int{ 8, -8 }
        color := move.piece.color()
        position := p.clone().make(move)

        if kind := move.piece.kind(); kind == KING {
                if move.isCastle() {
                        switch move.to {
                        case G1:
                                position.make(NewMove(position, H1, F1))
                        case C1:
                                position.make(NewMove(position, A1, D1))
                        case G8:
                                position.make(NewMove(position, H8, F8))
                        case C8:
                                position.make(NewMove(position, A8, D8))
                        }
                }
                position.can00[color], position.can000[color] = false, false
        } else {
                if kind == PAWN {
                        if move.isEnpassant(p.outposts[Pawn(color^1)]) {
                                //
                                // Mark the en-passant square.
                                //
                                position.enpassant = move.from + eight[color]
                        } else if move.isEnpassantCapture(p.enpassant) {
                                //
                                // Take out the en-passant pawn and decrement opponent's pawn count.
                                //
                                position.lift(NewMove(position, move.to + eight[color^1], move.to + eight[color^1]))
                                position.count[Pawn(color^1)]--
                        } else if move.promoted != 0 {
                                //
                                // Replace a pawn on 8th rank with the promoted piece.
                                //
                                position.put(move)
                                position.count[Pawn(color)]--
                                position.count[move.promoted]++
                        }
                }
                if position.can00[color] {
                        rookSquare := [2]int{ H1, H8 }
                        position.can00[color] = position.pieces[rookSquare[color]] == Rook(color)
                }
                if position.can000[color] {
                        rookSquare := [2]int{ A1, A8 }
                        position.can000[color] = position.pieces[rookSquare[color]] == Rook(color)
                }
        }

        position.color ^= 1 // <-- Switch side to move.
        return position.computeStage().setupAttacks()
}

func (p *Position) isCheck(color int) bool {
        king := p.outposts[King(color)]
        return king.intersect(p.attacks[color^1]).isNotEmpty()
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
        return piece.isPawn() && ((piece.isWhite() && target >= A8) || (piece.isBlack() && target <= H1))
}

func (p *Position) isKingSideCastleAllowed(color int) bool {
        if color == WHITE {
                return p.can00[WHITE] && p.pieces[F1] == 0 && p.pieces[G1] == 0 && castleKingWhite & p.attacks[BLACK] == 0
        }
        return p.can00[BLACK] && p.pieces[F8] == 0 && p.pieces[G8] == 0 && castleKingBlack & p.attacks[WHITE] == 0
}

func (p *Position) isQueenSideCastleAllowed(color int) bool {
        if p.color == WHITE {
                return p.can000[WHITE] && p.pieces[D1] == 0 && p.pieces[C1] == 0 && p.pieces[B1] == 0 && castleQueenWhite & p.attacks[BLACK] == 0
        }
        return p.can000[BLACK] && p.pieces[D8] == 0 && p.pieces[C8] == 0 && p.pieces[B8] == 0 && castleQueenBlack & p.attacks[WHITE] == 0
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
