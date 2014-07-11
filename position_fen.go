// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `strings`

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
		// fmt.Printf("%c ", char)
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
			// fmt.Printf("%02d: %s\n", square, piece)
			p.pieces[square] = piece
			p.outposts[piece].set(square)
			p.outposts[piece.color()].set(square)
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
	p.hash, p.hashPawns, p.hashMaterial = p.polyglot()
	p.tally = p.valuation()

	return p
}

func (p *Position) fen() string {
	return ":-)"
}
