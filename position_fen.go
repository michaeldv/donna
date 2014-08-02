// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`fmt`; `strings`)

// Decodes FEN string and creates new position.
func NewPositionFromFEN(game *Game, fen string) *Position {
	tree[node] = Position{game: game}
	p := &tree[node]

	// Expected matches of interest are as follows:
	// [0] - Pieces (entire board).
	// [1] - Color of side to move.
	// [2] - Castle rights.
	// [3] - En-passant square.
	// [4] - Number of half-moves.
	// [5] - Number of full moves.
	matches := strings.Split(fen, ` `)
	if len(matches) == 5 {
		return nil
	}
	// fmt.Printf("%q\n", matches)

	// [0] - Pieces (entire board).
	square := A8
	for _, char := range(matches[0]) {
		piece := Piece(0)
		switch(char) {
		case 'P':
			piece = Pawn
		case 'p':
			piece = BlackPawn
		case 'N':
			piece = Knight
		case 'n':
			piece = BlackKnight
		case 'B':
			piece = Bishop
		case 'b':
			piece = BlackBishop
		case 'R':
			piece = Rook
		case 'r':
			piece = BlackRook
		case 'Q':
			piece = Queen
		case 'q':
			piece = BlackQueen
		case 'K':
			piece = King
			p.king[White] = square
		case 'k':
			piece = BlackKing
			p.king[Black] = square
		case '/':
			square -= 16
		case '1', '2', '3', '4', '5', '6', '7', '8':
			square += int(char - '0')
		}
		if piece != 0 {
			p.pieces[square] = piece
			p.outposts[piece].set(square)
			p.outposts[piece.color()].set(square)
			p.balance += materialBalance[piece]
			p.count[piece]++
			square++
		}
	}

	// [1] - Color of side to move.
	if matches[1] == `w` {
		p.color = White
	} else {
		p.color = Black
	}

	// [2] - Castle rights.
	for _, char := range(matches[2]) {
		switch(char) {
		case 'K':
			p.castles |= castleKingside[White]
		case 'Q':
			p.castles |= castleQueenside[White]
		case 'k':
			p.castles |= castleKingside[Black]
		case 'q':
			p.castles |= castleQueenside[Black]
		case '-':
			// No castling rights.
		}
	}

	// [3] - En-passant square.
	if matches[3] != `-` {
		p.enpassant = Square(int(matches[3][1] - '1'), int(matches[3][0] - 'a'))

	}

	p.reversible = true
	p.board = p.outposts[White] | p.outposts[Black]
	p.hash, p.hashPawns = p.polyglot()
	p.tally = p.valuation()

	return p
}

// Encodes position as FEN string.
func (p *Position) fen() (fen string) {
	fancy := Settings.Fancy
	Settings.Fancy = false; defer func() { Settings.Fancy = fancy }()

	// Board: start from A8->H8 going down to A1->H1.
	empty := 0
	for row := A8H8; row >= A1H1; row-- {
		for col := A1A8; col <= H1H8; col++ {
			square := Square(row, col)
			piece := p.pieces[square]

			if piece != 0 {
				if empty != 0 {
					fen += fmt.Sprintf(`%d`, empty)
					empty = 0
				}
				fen += piece.String()
			} else {
				empty++
			}

			if col == 7 {
				if empty != 0 {
					fen += fmt.Sprintf(`%d`, empty)
					empty = 0
				}
				if row != 0 {
					fen += `/`
				}
			}
		}
	}

	// Side to move.
	if p.color == White {
		fen += ` w`
	} else {
		fen += ` b`
	}

	// Castle rights for both sides, if any.
	if p.castles & 0x0F != 0 {
		fen += ` `
		if p.castles & castleKingside[White] != 0 {
			fen += `K`
		}
		if p.castles & castleQueenside[White] != 0 {
			fen += `Q`
		}
		if p.castles & castleKingside[Black] != 0 {
			fen += `k`
		}
		if p.castles & castleQueenside[Black] != 0 {
			fen += `q`
		}
	} else {
		fen += ` -`
	}

	// En-passant square, if any.
	if p.enpassant != 0 {
		row, col := Coordinate(p.enpassant)
		fen += fmt.Sprintf(` %c%d`, col + 'a', row + 1)
	} else {
		fen += ` -`
	}

	return
}
