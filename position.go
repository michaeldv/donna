// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`bytes`)

type Flags struct {
        enpassant   int         // En-passant square caused by previous move.
        reversible  bool        // Is this position reversible?
        can00       [2]bool     // Is king-side castle allowed?
        can000      [2]bool     // Is queen-side castle allowed?
}

type Position struct {
        game      *Game
        previous  *Position     // Previous position.
        flags     *Flags        // Flags set by last move leading to this position.
        pieces    [64]Piece     // Array of 64 squares with pieces on them.
        targets   [64]Bitmask   // Attack targets for each piece on the board.
        board     [3]Bitmask    // [0] white pieces only, [1] black pieces, and [2] all pieces.
        attacks   [3]Bitmask    // [0] all squares attacked by white, [1] by black, [2] either white or black.
        outposts  [16]Bitmask   // Bitmasks of each piece on the board, ex. white pawns, black king, etc.
        count     [16]int       // counts of each piece on the board, ex. white pawns: 6, etc.
        color     int           // Side to make next move.
        stage     int           // Game stage (256 in the initial position).
        hash      uint64        // Polyglot hash value.
        inCheck   bool          // Is our king under attack?
}

func NewPosition(game *Game, pieces [64]Piece, color int, flags *Flags) *Position {
        p := &Position{ game: game, pieces: pieces, color: color }
        if flags != nil {
	        p.flags = flags
        } else {
	        p.flags = &Flags{ enpassant: 0, reversible: true }
		p.flags.can00[White]  = p.pieces[E1] == King(White) && p.pieces[H1] == Rook(White)
		p.flags.can00[Black]  = p.pieces[E8] == King(Black) && p.pieces[H8] == Rook(Black)
		p.flags.can000[White] = p.pieces[E1] == King(White) && p.pieces[A1] == Rook(White)
		p.flags.can000[Black] = p.pieces[E8] == King(Black) && p.pieces[A8] == Rook(Black)
	}

        return p.setupPieces().setupAttacks().computeStage()
}

func (p *Position) setupPieces() *Position {
        for square, piece := range p.pieces {
                if piece != 0 {
                        p.outposts[piece].set(square)
                        p.board[piece.color()].set(square)
                        p.count[piece]++
                }
        }
        p.board[2] = p.board[White] | p.board[Black]
        return p
}

func (p *Position) updatePieces(updates [64]Piece, squares []int) *Position {
        for _, square := range squares {
		if newer, older := updates[square], p.pieces[square]; newer != older {
			if older != 0 {
				p.count[older]--
			}
			p.outposts[older].clear(square)
			p.board[newer.color()^1].clear(square)
			p.pieces[square] = newer
	                if newer != 0 {
	                        p.outposts[newer].set(square)
	                        p.board[newer.color()].set(square)
				p.count[newer]++
	                } else {
	                        p.outposts[newer].clear(square)
	                        p.board[newer.color()].clear(square)
	                }
		}
        }
        p.board[2] = p.board[White] | p.board[Black]
        return p
}

func (p *Position) setupAttacks() *Position {
        kingSquare := [2]int{ -1, -1 }

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

func (p *Position) updateKingTargets(kingSquare [2]int) *Position {
	if kingSquare[White] >= 0 && kingSquare[Black] >= 0 {
                p.targets[kingSquare[White]].exclude(p.targets[kingSquare[Black]])
                p.targets[kingSquare[Black]].exclude(p.targets[kingSquare[White]])
	}
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

func (p *Position) computeStage() *Position {
        p.hash  = p.polyglot()
        p.stage = 2 * (p.count[Pawn(White)]   + p.count[Pawn(Black)])   +
                  6 * (p.count[Knight(White)] + p.count[Knight(Black)]) +
                 12 * (p.count[Bishop(White)] + p.count[Bishop(Black)]) +
                 16 * (p.count[Rook(White)]   + p.count[Rook(Black)])   +
                 44 * (p.count[Queen(White)]  + p.count[Queen(Black)])
        return p
}

func (p *Position) MakeMove(move *Move) *Position {
        eight := [2]int{ 8, -8 }
        color := move.piece.color()
	flags := &Flags{
		enpassant:  0,
		can00:      p.flags.can00,
		can000:     p.flags.can000,
		reversible: true,
	}

        delta := p.pieces
        delta[move.from] = 0
        delta[move.to] = move.piece
        squares := []int{ move.from, move.to }

        if kind := move.piece.kind(); kind == KING {
                if move.isCastle() {
                        switch move.to {
                        case G1:
                                delta[H1], delta[F1] = 0, Rook(White)
				squares = append(squares, H1, F1)
                        case C1:
                                delta[A1], delta[D1] = 0, Rook(White)
				squares = append(squares, A1, D1)
                        case G8:
                                delta[H8], delta[F8] = 0, Rook(Black)
				squares = append(squares, H8, F8)
                        case C8:
                                delta[A8], delta[D8] = 0, Rook(Black)
				squares = append(squares, A8, D8)
                        }
                }
                flags.can00[color], flags.can000[color] = false, false
        } else {
                if kind == PAWN {
                        if move.isEnpassant(p.outposts[Pawn(color^1)]) {
                                //
                                // Mark the en-passant square.
                                //
                                flags.enpassant = move.from + eight[color]
                        } else if move.isEnpassantCapture(p.flags.enpassant) {
                                //
                                // Take out the en-passant pawn and decrement opponent's pawn count.
                                //
                                delta[move.to - eight[color]] = 0
				squares = append(squares, move.to - eight[color])
                        } else if move.promoted != 0 {
                                //
                                // Replace a pawn on 8th rank with the promoted piece.
                                //
                                delta[move.to] = move.promoted
                        }
                }
                if p.flags.can00[color] {
                        rookSquare := [2]int{ H1, H8 }
                        flags.can00[color] = delta[rookSquare[color]] == Rook(color)
                }
                if p.flags.can000[color] {
                        rookSquare := [2]int{ A1, A8 }
                        flags.can000[color] = delta[rookSquare[color]] == Rook(color)
                }
        }

	position := &Position{
		game:     p.game,
		previous: p,
		board:    p.board,
		count:    p.count,
		pieces:   p.pieces,
		outposts: p.outposts,
		color:    color^1,
		flags:    flags,
	}

	position.updatePieces(delta, squares).setupAttacks()
	if position.isCheck(color) {
		return nil
	}
	return position.computeStage()
}

func (p *Position) isCheck(color int) bool {
        king := p.outposts[King(color)]
        return king.intersect(p.attacks[color^1]).isNotEmpty()
}

func (p *Position) isRepetition() bool {
        // if p.history > 6 {
        //         reps, hash := 0, p.game.repetitions[p.history - 1]
        //         for i := p.history - 1; i >= 0; i-- {
        //                 if p.game.repetitions[i] == hash {
        //                         reps++
        //                         if reps == 3 {
        //                                 return true
        //                         }
        //                 }
        //         }
        // }
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
                return p.flags.can00[White] && p.pieces[F1] == 0 && p.pieces[G1] == 0 && castleKingWhite & p.attacks[Black] == 0
        }
        return p.flags.can00[Black] && p.pieces[F8] == 0 && p.pieces[G8] == 0 && castleKingBlack & p.attacks[White] == 0
}

func (p *Position) isQueenSideCastleAllowed(color int) bool {
        if color == White {
                return p.flags.can000[White] && p.pieces[D1] == 0 && p.pieces[C1] == 0 && p.pieces[B1] == 0 && castleQueenWhite & p.attacks[Black] == 0
        }
        return p.flags.can000[Black] && p.pieces[D8] == 0 && p.pieces[C8] == 0 && p.pieces[B8] == 0 && castleQueenBlack & p.attacks[White] == 0
}

// Compute position's polyglot hash.
func (p *Position) polyglot() (key uint64) {
        for i, piece := range p.pieces {
                if piece != 0 {
                        key ^= polyglotRandom[0:768][64 * piece.polyglot() + i]
                }
        }

	if p.flags.can00[White] {
                key ^= polyglotRandom[768]
	}
	if p.flags.can000[White] {
                key ^= polyglotRandom[769]
	}
	if p.flags.can00[Black] {
                key ^= polyglotRandom[770]
	}
	if p.flags.can000[Black] {
                key ^= polyglotRandom[771]
	}
        if p.flags.enpassant != 0 {
                col := Col(p.flags.enpassant)
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
