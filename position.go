// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`bytes`)

var tree [1024]Position
var node int

type Flags struct {
        enpassant     int       // En-passant square caused by previous move.
        irreversible  bool      // Is this position reversible?
}

type Position struct {
        game      *Game
        flags     Flags         // Flags set by last move leading to this position.
        pieces    [64]Piece     // Array of 64 squares with pieces on them.
        targets   [64]Bitmask   // Attack targets for pieces on each square of the board.
        board     [3]Bitmask    // [0] white pieces only, [1] black pieces, and [2] all pieces.
        outposts  [14]Bitmask   // Bitmasks of each piece on the board, ex. white pawns, black king, etc.
        count     [16]int       // counts of each piece on the board, ex. white pawns: 6, etc.
        color     int           // Side to make next move.
        stage     int           // Game stage (256 in the initial position).
        hash      uint64        // Polyglot hash value.
        castles   uint8         // Castle rights mask.
}

func NewPosition(game *Game, pieces [64]Piece, color int, flags Flags) *Position {
        tree[node] = Position{ game: game, pieces: pieces, color: color }
        p := &tree[node]

        p.castles = castleKingside[White] | castleQueenside[White] |
                    castleKingside[Black] | castleQueenside[Black]

        if p.pieces[E1] != King(White) || p.pieces[H1] != Rook(White) {
                p.castles &= ^castleKingside[White]
        }
        if p.pieces[E1] != King(White) || p.pieces[A1] != Rook(White) {
                p.castles &= ^castleQueenside[White]
        }

        if p.pieces[E8] != King(Black) || p.pieces[H8] != Rook(Black) {
                p.castles &= ^castleKingside[Black]
        }
        if p.pieces[E8] != King(Black) || p.pieces[A8] != Rook(Black) {
                p.castles &= ^castleQueenside[Black]
        }

        return p.setupPieces().computeStage()
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

func (p *Position) computeStage() *Position {
        p.hash  = p.polyglot()
        p.stage = 2 * (p.count[Pawn(White)]   + p.count[Pawn(Black)])   +
                  6 * (p.count[Knight(White)] + p.count[Knight(Black)]) +
                 12 * (p.count[Bishop(White)] + p.count[Bishop(Black)]) +
                 16 * (p.count[Rook(White)]   + p.count[Rook(Black)])   +
                 44 * (p.count[Queen(White)]  + p.count[Queen(Black)])
        return p
}

func (p *Position) MakeMove(move Move) *Position {
        color := move.color()
        flags := Flags{}

        from, to, piece, capture := move.split()
        flags.irreversible = capture != 0

        squares := []int{ from, to }
        delta := p.pieces
        delta[from], delta[to] = 0, piece

        // Castle rights for current node are based on the castle rights from
        // the previous node.
        castles := tree[node].castles & castleRights[from] & castleRights[to]

        if piece.isKing() {
                castles &= ^castleKingside[color]
                castles &= ^castleQueenside[color]
                if move.isCastle() {
                        flags.irreversible = true
                        switch to {
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
        } else if piece.isPawn() {
                flags.irreversible = true
                if move.isEnpassant() {
                        flags.enpassant = from + eight[color]           // Save the en-passant square.
                } else if to == p.flags.enpassant {
                        delta[to - eight[color]] = 0                    // Take out the en-passant pawn...
                        squares = append(squares, to - eight[color])    // and decrement opponent's pawn count.
                } else if promo := move.promo(); promo != 0 {
                        delta[to] = promo                               // Place the promoted piece.
                }
        }

	node++
	tree[node] = Position{
		game:     p.game,
		board:    p.board,
		count:    p.count,
		pieces:   p.pieces,
		outposts: p.outposts,
		color:    color^1,
		flags:    flags,
		castles:  castles,
	}

	tree[node].updatePieces(delta, squares)
	if tree[node].isInCheck(color) {
		node--
		return nil
	}
	return tree[node].computeStage()

}

func (p *Position) TakeBack(move Move) *Position {
        node--
        return &tree[node]
}

func (p *Position) isInCheck(color int) bool {
        return p.isAttacked(p.outposts[King(color)].first(), color^1)
}

func (p *Position) isRepetition() bool {
        if p.flags.irreversible {
                return false
        }

        for reps, prevNode := 1, node - 1; prevNode >= 0; prevNode-- {
                if tree[prevNode].flags.irreversible {
                        return false
                }
                if tree[prevNode].color == p.color && tree[prevNode].hash == p.hash {
                        reps++
                        if reps == 3 {
                                return true
                        }
                }
        }

        return false
}

func (p *Position) saveBest(ply int, move Move) {
        p.game.bestLine[ply][ply] = move
        p.game.bestLength[ply] = ply + 1
        for i := ply + 1; i < p.game.bestLength[ply + 1]; i++ {
                p.game.bestLine[ply][i] = p.game.bestLine[ply + 1][i]
                p.game.bestLength[ply]++
        }
}

func (p *Position) canCastle(color int) (kingside, queenside bool) {
        attacks := p.attacks(color^1)
        kingside = p.castles & castleKingside[color] != 0 &&
                   (gapKing[color] & p.board[2] == 0) &&
                   (castleKing[color] & attacks == 0)

        queenside = p.castles & castleQueenside[color] != 0 &&
                    (gapQueen[color] & p.board[2] == 0) &&
                    (castleQueen[color] & attacks == 0)
        return
}

// Compute position's polyglot hash.
func (p *Position) polyglot() (key uint64) {
        for i, piece := range p.pieces {
                if piece != 0 {
                        key ^= polyglotRandom[0:768][64 * piece.polyglot() + i]
                }
        }

	if p.castles & castleKingside[White] != 0 {
                key ^= polyglotRandom[768]
	}
	if p.castles & castleQueenside[White] != 0 {
                key ^= polyglotRandom[769]
	}
	if p.castles & castleKingside[Black] != 0 {
                key ^= polyglotRandom[770]
	}
	if p.castles & castleQueenside[Black] != 0 {
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
        if !p.isInCheck(p.color) {
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
