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
        board     Bitmask       // Bitmask of all pieces on the board.
        outposts  [14]Bitmask   // Bitmasks of each piece on the board; [0] all white, [1] all black.
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

        if p.pieces[E1] != King || p.pieces[H1] != Rook {
                p.castles &= ^castleKingside[White]
        }
        if p.pieces[E1] != King || p.pieces[A1] != Rook {
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
                        p.outposts[piece.color()].set(square)
                        p.count[piece]++
                }
        }

        p.hash = p.polyglot()
        p.board = p.outposts[White] | p.outposts[Black]

        return p
}

func (p *Position) movePiece(piece Piece, from, to int) *Position {
        p.pieces[from], p.pieces[to] = 0, piece
        p.outposts[piece] ^= bit[from] | bit[to]
        p.outposts[piece.color()] ^= bit[from] | bit[to]

        return p
}

func (p *Position) promotePawn(piece Piece, from, to int, promo Piece) *Position {
        p.pieces[from], p.pieces[to] = 0, promo
        p.outposts[piece] ^= bit[from]
        p.outposts[promo] ^= bit[to]
        p.outposts[piece.color()] ^= bit[from] | bit[to]
        p.count[piece]--
        p.count[promo]++

        return p
}

func (p *Position) capturePiece(capture Piece, from, to int) *Position {
        p.outposts[capture] ^= bit[to]
        p.outposts[capture.color()] ^= bit[to]
        p.count[capture]--

        return p
}

func (p *Position) captureEnpassant(capture Piece, from, to int) *Position {
        enpassant := to - eight[capture.color()^1]

        p.pieces[enpassant] = 0
        p.outposts[capture] ^= bit[enpassant]
        p.outposts[capture.color()] ^= bit[enpassant]
        p.count[capture]--

        return p
}

func (p *Position) MakeMove(move Move) *Position {
        color := move.color()
        from, to, piece, capture := move.split()
        //
        // Copy over the contents of previous tree node to the current one.
        //
        node++
        tree[node] = *p // => tree[node] = tree[node - 1]
        pp := &tree[node]

        pp.hash ^= hashCastle[pp.castles]
        if pp.flags.enpassant != 0 {
                pp.hash ^= hashEnpassant[Col(pp.flags.enpassant)]
        }
        pp.flags.enpassant, pp.flags.irreversible = 0, false
        //
        // Castle rights for current node are based on the castle rights from
        // the previous node.
        //
        pp.castles &= castleRights[from] & castleRights[to]
        pp.hash ^= hashCastle[pp.castles]

        if capture != 0 {
                pp.flags.irreversible = true
                if to != 0 && to == p.flags.enpassant {
                        pp.hash ^= polyglotRandom[64 * pawn(color^1).polyglot() + to - eight[color]]
                        pp.captureEnpassant(pawn(color^1), from, to)
                } else {
                        pp.hash ^= polyglotRandom[64 * p.pieces[to].polyglot() + to]
                        pp.capturePiece(capture, from, to)
                }
        }

        if promo := move.promo(); promo == 0 {
                poly := 64 * p.pieces[from].polyglot()
                pp.hash ^= polyglotRandom[poly + from] ^ polyglotRandom[poly + to]
                pp.movePiece(piece, from, to)
                if move.isCastle() {
                        pp.flags.irreversible = true
                        switch to {
                        case G1:
                                poly = 64 * Piece(Rook).polyglot()
                                pp.hash ^= polyglotRandom[poly + H1] ^ polyglotRandom[poly + F1]
                                pp.movePiece(Rook, H1, F1)
                        case C1:
                                poly = 64 * Piece(Rook).polyglot()
                                pp.hash ^= polyglotRandom[poly + A1] ^ polyglotRandom[poly + D1]
                                pp.movePiece(Rook, A1, D1)
                        case G8:
                                poly = 64 * Piece(BlackRook).polyglot()
                                pp.hash ^= polyglotRandom[poly + H8] ^ polyglotRandom[poly + F8]
                                pp.movePiece(BlackRook, H8, F8)
                        case C8:
                                poly = 64 * Piece(BlackRook).polyglot()
                                pp.hash ^= polyglotRandom[poly + A8] ^ polyglotRandom[poly + D8]
                                pp.movePiece(BlackRook, A8, D8)
                        }
                } else if piece.isPawn() {
                        pp.flags.irreversible = true
                        if move.isEnpassant() {
                                pp.flags.enpassant = from + eight[color] // Save the en-passant square.
                                pp.hash ^= hashEnpassant[Col(pp.flags.enpassant)]
                        }
                }
        } else {
                pp.flags.irreversible = true
                pp.hash ^= polyglotRandom[64 * pawn(color).polyglot() + from]
                pp.hash ^= polyglotRandom[64 * promo.polyglot() + to]
                pp.promotePawn(piece, from, to, promo)
        }

	pp.board = pp.outposts[White] | pp.outposts[Black]
	if pp.isInCheck(color) {
		node--
		return nil
	}

	if color == White {
                pp.hash ^= polyglotRandomWhite
	}
	pp.color = color^1

	return pp // => &tree[node]

}

func (p *Position) TakeBack(move Move) *Position {
        node--
        return &tree[node]
}

func (p *Position) isInCheck(color int) bool {
        return p.isAttacked(p.outposts[king(color)].first(), color^1)
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
                   (gapKing[color] & p.board == 0) &&
                   (castleKing[color] & attacks == 0)

        queenside = p.castles & castleQueenside[color] != 0 &&
                    (gapQueen[color] & p.board == 0) &&
                    (castleQueen[color] & attacks == 0)
        return
}

// Calculates position stage based on what pieces are on the board (256 for
// the initial position, 0 for bare kings).
func (p *Position) stage() int {
        return 2 * (p.count[Pawn]   + p.count[BlackPawn])   +
               6 * (p.count[Knight] + p.count[BlackKnight]) +
              12 * (p.count[Bishop] + p.count[BlackBishop]) +
              16 * (p.count[Rook]   + p.count[BlackRook])   +
              44 * (p.count[Queen]  + p.count[BlackQueen])
}

// Calculates normalized position score based on position stage and given
// midgame/endgame values.
func (p *Position) score(midgame, endgame int) int {
        stage := p.stage()
        return (midgame * stage + endgame * (256 - stage)) / 256
}

// Compute position's polyglot hash.
func (p *Position) polyglot() (key uint64) {
        board := p.board
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
