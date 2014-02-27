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
        board     [3]Bitmask    // [0] white pieces only, [1] black pieces, and [2] all pieces.
        outposts  [14]Bitmask   // Bitmasks of each piece on the board, ex. white pawns, black king, etc.
        count     [16]int       // counts of each piece on the board, ex. white pawns: 6, etc.
        color     int           // Side to make next move.
        hash      uint64        // Polyglot hash value.
        castles   uint8         // Castle rights mask.
}

func NewPosition(game *Game, pieces [64]Piece, color int, flags Flags) *Position {
        tree[node] = Position{ game: game, pieces: pieces, color: color }
        p := &tree[node]

        p.castles = castleKingside[White] | castleQueenside[White] |
                    castleKingside[Black] | castleQueenside[Black]

        if p.pieces[E1] != WhiteKing || p.pieces[H1] != WhiteRook {
                p.castles &= ^castleKingside[White]
        }
        if p.pieces[E1] != WhiteKing || p.pieces[A1] != WhiteRook {
                p.castles &= ^castleQueenside[White]
        }

        if p.pieces[E8] != BlackKing || p.pieces[H8] != BlackRook {
                p.castles &= ^castleKingside[Black]
        }
        if p.pieces[E8] != BlackKing || p.pieces[A8] != BlackRook {
                p.castles &= ^castleQueenside[Black]
        }

        for square, piece := range p.pieces {
                if piece != 0 {
                        p.outposts[piece].set(square)
                        p.board[piece.color()].set(square)
                        p.count[piece]++
                }
        }

        p.hash = p.polyglot()
        p.board[2] = p.board[White] | p.board[Black]

        return p
}

func (p *Position) movePiece(from, to int) *Position {
        piece, color := p.pieces[from], p.pieces[from].color()

        p.pieces[from], p.pieces[to] = 0, piece
        p.outposts[piece] ^= bit[from] | bit[to]
        p.board[color] ^= bit[from] | bit[to]

        return p
}

func (p *Position) promotePawn(from, to int, promo Piece) *Position {
        piece, color := p.pieces[from], p.pieces[from].color()

        p.pieces[from], p.pieces[to] = 0, promo
        p.outposts[piece] ^= bit[from]
        p.outposts[promo] ^= bit[to]
        p.board[color] ^= bit[from] | bit[to]
        p.count[piece]--
        p.count[promo]++

        return p
}

func (p *Position) capturePiece(from, to int) *Position {
        capture, enemy := p.pieces[to], p.pieces[to].color()

        p.outposts[capture] ^= bit[to]
        p.board[enemy] ^= bit[to]
        p.count[capture]--

        return p
}

func (p *Position) captureEnpassant(from, to int) *Position {
        color := p.pieces[from].color()
        enemy := color^1
        capture := Pawn(enemy)
        enpassant := to - eight[color]

        p.pieces[enpassant] = 0
        p.outposts[capture] ^= bit[enpassant]
        p.board[enemy] ^= bit[enpassant]
        p.count[capture]--

        return p
}

func (p *Position) MakeMove(move Move) *Position {
        color := move.color()
        flags := Flags{}

        from, to, piece, capture := move.split()
        //
        // Copy over the contents of previous tree node to the current one.
        //
        node++
        tree[node] = tree[node - 1] // Faster that tree[node] = *p ?!

        hash := tree[node].hash ^ hashCastle[tree[node].castles]
        if tree[node].flags.enpassant != 0 {
                hash ^= hashEnpassant[Col(tree[node].flags.enpassant)]
        }
        //
        // Castle rights for current node are based on the castle rights from
        // the previous node.
        //
        tree[node].castles &= castleRights[from] & castleRights[to]
        hash ^= hashCastle[tree[node].castles]

        if capture != 0 {
                flags.irreversible = true
                if to != 0 && to == tree[node].flags.enpassant {
                        hash ^= polyglotRandom[64 * Pawn(color^1).polyglot() + to - eight[color]]
                        tree[node].captureEnpassant(from, to)
                } else {
                        hash ^= polyglotRandom[64 * p.pieces[to].polyglot() + to]
                        tree[node].capturePiece(from, to)
                }
        }

        if promo := move.promo(); promo == 0 {
                poly := 64 * p.pieces[from].polyglot()
                hash ^= polyglotRandom[poly + from] ^ polyglotRandom[poly + to]
                tree[node].movePiece(from, to)
                if move.isCastle() {
                        flags.irreversible = true
                        switch to {
                        case G1:
                                poly = 64 * Piece(WhiteRook).polyglot()
                                hash ^= polyglotRandom[poly + H1] ^ polyglotRandom[poly + F1]
                                tree[node].movePiece(H1, F1)
                        case C1:
                                poly = 64 * Piece(WhiteRook).polyglot()
                                hash ^= polyglotRandom[poly + A1] ^ polyglotRandom[poly + D1]
                                tree[node].movePiece(A1, D1)
                        case G8:
                                poly = 64 * Piece(BlackRook).polyglot()
                                hash ^= polyglotRandom[poly + H8] ^ polyglotRandom[poly + F8]
                                tree[node].movePiece(H8, F8)
                        case C8:
                                poly = 64 * Piece(BlackRook).polyglot()
                                hash ^= polyglotRandom[poly + A8] ^ polyglotRandom[poly + D8]
                                tree[node].movePiece(A8, D8)
                        }
                } else if piece.isPawn() {
                        flags.irreversible = true
                        if move.isEnpassant() {
                                flags.enpassant = from + eight[color] // Save the en-passant square.
                                hash ^= hashEnpassant[Col(flags.enpassant)]
                        }
                }
        } else {
                flags.irreversible = true
                hash ^= polyglotRandom[64 * Pawn(color).polyglot() + from]
                hash ^= polyglotRandom[64 * promo.polyglot() + to]
                tree[node].promotePawn(from, to, promo)
        }

	if color == White {
                hash ^= polyglotRandomWhite
	}

	tree[node].color = color^1
	tree[node].flags = flags
	tree[node].hash = hash
	tree[node].board[2] = tree[node].board[White] | tree[node].board[Black]

	if tree[node].isInCheck(color) {
		node--
		return nil
	}
	return &tree[node]

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

// Calculates position stage based on what pieces are on the board (256 for
// the initial position, 0 for bare kings).
func (p *Position) stage() int {
        return 2 * (p.count[WhitePawn]   + p.count[BlackPawn])   +
               6 * (p.count[WhiteKnight] + p.count[BlackKnight]) +
              12 * (p.count[WhiteBishop] + p.count[BlackBishop]) +
              16 * (p.count[WhiteRook]   + p.count[BlackRook])   +
              44 * (p.count[WhiteQueen]  + p.count[BlackQueen])
}

// Calculates normalized position score based on position stage and given
// midgame/endgame values.
func (p *Position) score(midgame, endgame int) int {
        stage := p.stage()
        return (midgame * stage + endgame * (256 - stage)) / 256
}

// Compute position's polyglot hash.
func (p *Position) polyglot() (key uint64) {
        board := p.board[2]
        for board != 0 {
                square := board.pop() // Inline polyhash() is at lest 10% faster.
                key ^= polyglotRandom[64 * p.pieces[square].polyglot() + square]
        }

	key ^= hashCastle[p.castles]

	if p.flags.enpassant != 0 {
                key ^= hashEnpassant[Col(p.flags.enpassant)]
	}
	if p.color == White {
                key ^= polyglotRandomWhite
	}

	return
}

func (p *Position) polyhash (square int) uint64 {
       return polyglotRandom[64 * p.pieces[square].polyglot() + square]
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
