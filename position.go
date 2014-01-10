// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`bytes`)

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
        history   int           // Index to the tip of the repetitions array.
        inCheck   bool          // Is our king under attack?
        can00     [2]bool       // Is king-side castle allowed?
        can000    [2]bool       // Is queen-side castle allowed?
}

func NewPosition(game *Game, pieces [64]Piece, color, enpassant int) *Position {
        p := new(Position)
        p.game = game
        p.pieces = pieces
        p.enpassant = enpassant
        p.color = color

        p.can00[White]  = p.pieces[E1] == King(White) && p.pieces[H1] == Rook(White)
        p.can00[Black]  = p.pieces[E8] == King(Black) && p.pieces[H8] == Rook(Black)
        p.can000[White] = p.pieces[E1] == King(White) && p.pieces[A1] == Rook(White)
        p.can000[Black] = p.pieces[E8] == King(Black) && p.pieces[A8] == Rook(Black)

        return p.setupPieces().computeStage().setupAttacks() //.saveHistory(nil)
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
        p.board[2] = p.board[White] | p.board[Black]

        p.stage = 2 * (p.count[Pawn(White)]   + p.count[Pawn(Black)])   +
                  6 * (p.count[Knight(White)] + p.count[Knight(Black)]) +
                 12 * (p.count[Bishop(White)] + p.count[Bishop(Black)]) +
                 16 * (p.count[Rook(White)]   + p.count[Rook(Black)])   +
                 44 * (p.count[Queen(White)]  + p.count[Queen(Black)])
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
        p.attacks[2] = p.attacks[White] | p.attacks[Black]
        p.inCheck = p.isCheck(p.color)

        return p
}

// Save position's polyglot hash in the repetitions array incrementing history to point
// to the next available spot. Positions resulted from the null move are not saved.
func (p *Position) saveHistory(prev *Position) *Position {
        spot := 0
        if prev != nil {
                if p.color == prev.color { // <-- Null move.
                        return p
                }
                spot = prev.history
        }

        p.game.repetitions[spot] = p.polyglot()
        p.history = spot + 1

        return p
}

func (p *Position) updateKingTargets(kingSquare [2]int) *Position {
        p.targets[kingSquare[White]].exclude(p.targets[kingSquare[Black]])
        p.targets[kingSquare[Black]].exclude(p.targets[kingSquare[White]])
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

func (p *Position) MakeMove(move *Move) *Position {
        eight := [2]int{ 8, -8 }
        color := move.piece.color()
        enpassant := 0

        pieces := p.pieces
        pieces[move.from] = 0
        pieces[move.to] = move.piece

        if kind := move.piece.kind(); kind == KING {
                if move.isCastle() {
                        switch move.to {
                        case G1:
                                pieces[H1], pieces[F1] = 0, Rook(White)
                        case C1:
                                pieces[A1], pieces[D1] = 0, Rook(White)
                        case G8:
                                pieces[H8], pieces[F8] = 0, Rook(Black)
                        case C8:
                                pieces[A8], pieces[D8] = 0, Rook(Black)
                        }
                }
                p.can00[color], p.can000[color] = false, false
        } else {
                if kind == PAWN {
                        if move.isEnpassant(p.outposts[Pawn(color^1)]) {
                                //
                                // Mark the en-passant square.
                                //
                                enpassant = move.from + eight[color]
                        } else if move.isEnpassantCapture(p.enpassant) {
                                //
                                // Take out the en-passant pawn and decrement opponent's pawn count.
                                //
                                pieces[move.to - eight[color]] = 0
                        } else if move.promoted != 0 {
                                //
                                // Replace a pawn on 8th rank with the promoted piece.
                                //
                                pieces[move.to] = move.promoted
                        }
                }
                if p.can00[color] {
                        rookSquare := [2]int{ H1, H8 }
                        p.can00[color] = p.pieces[rookSquare[color]] == Rook(color)
                }
                if p.can000[color] {
                        rookSquare := [2]int{ A1, A8 }
                        p.can000[color] = p.pieces[rookSquare[color]] == Rook(color)
                }
        }

        return NewPosition(p.game, pieces, color^1, enpassant)
}

func (p *Position) isCheck(color int) bool {
        king := p.outposts[King(color)]
        return king.intersect(p.attacks[color^1]).isNotEmpty()
}

func (p *Position) isRepetition() bool {
        if p.history > 6 {
                reps, hash := 0, p.game.repetitions[p.history - 1]
                for i := p.history - 1; i >= 0; i-- {
                        if p.game.repetitions[i] == hash {
                                reps++
                                if reps == 3 {
                                        return true
                                }
                        }
                }
        }
        return false
}

func (p *Position) saveBest(ply int, move *Move) {
        p.game.bestLine[ply][ply] = move
        p.game.bestLength[ply] = ply + 1
        for i := ply + 1; i < p.game.bestLength[ply + 1]; i++ {
                p.game.bestLine[ply][i] = p.game.bestLine[ply + 1][i]
                p.game.bestLength[ply]++
        }
}

func (p *Position) isPawnPromotion(piece Piece, target int) bool {
        return piece.isPawn() && ((piece.isWhite() && target >= A8) || (piece.isBlack() && target <= H1))
}

func (p *Position) isKingSideCastleAllowed(color int) bool {
        if color == White {
                return p.can00[White] && p.pieces[F1] == 0 && p.pieces[G1] == 0 && castleKingWhite & p.attacks[Black] == 0
        }
        return p.can00[Black] && p.pieces[F8] == 0 && p.pieces[G8] == 0 && castleKingBlack & p.attacks[White] == 0
}

func (p *Position) isQueenSideCastleAllowed(color int) bool {
        if p.color == White {
                return p.can000[White] && p.pieces[D1] == 0 && p.pieces[C1] == 0 && p.pieces[B1] == 0 && castleQueenWhite & p.attacks[Black] == 0
        }
        return p.can000[Black] && p.pieces[D8] == 0 && p.pieces[C8] == 0 && p.pieces[B8] == 0 && castleQueenBlack & p.attacks[White] == 0
}

// Compute position's polyglot hash.
func (p *Position) polyglot() (key uint64) {
        for i, piece := range p.pieces {
                if piece != 0 {
                        key ^= polyglotRandom[0:768][64 * piece.polyglot() + i]
                }
        }

	if p.can00[White] {
                key ^= polyglotRandom[768]
	}
	if p.can000[White] {
                key ^= polyglotRandom[769]
	}
	if p.can00[Black] {
                key ^= polyglotRandom[770]
	}
	if p.can000[Black] {
                key ^= polyglotRandom[771]
	}
        if p.enpassant != 0 {
                col := Col(p.enpassant)
                key ^= polyglotRandom[772 + col]
        }
	if p.color == White {
                key ^= polyglotRandom[780]
	}

	return
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
