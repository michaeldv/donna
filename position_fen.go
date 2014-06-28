// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `regexp`
// import `fmt`

func NewPositionFromFEN(game *Game, fen string) *Position {
	const ( 
		row1   = `([rnbqkRNBQK1-8]+/`
		rows26 = `([rnbqkpRNBQKP1-8]+/){6}`
		row8   = `[rnbqkRNBQK1-8]+)`
		color  = `([bw])`
		castle = `(-|K?Q?k?q?)`
		enpass = `(-|[a-h][36])`
		number = `(\d+)`
		spx    = `\s*`
		sp     = `\s`
	)

	tree[node] = Position{game: game}
	p := &tree[node]

	// Expected matches of interest are as follows:
	// [1] - Pieces (entire board).
	// [3] - Color of side to move.
	// [4] - Castle rights.
	// [5] - En-passant square.
	// [6] - Number of half-moves.
	// [7] - Number of full moves.
	re := regexp.MustCompile(spx + row1 + rows26 + row8 + sp + color + sp + castle + sp + enpass + sp + number + sp + number + spx);
	matches := re.FindStringSubmatch(fen)
	if matches == nil {
		return nil
	}
	// fmt.Println(spx + row1 + rows26 + row8 + sp + color + sp + castle + sp + enpass + sp + number + sp + number + spx)
	// fmt.Printf("%q\n", matches)
	// fmt.Printf("%s\n", matches[1])

	// [1] - Pieces (entire board).
	square := A8
	for _, char := range(matches[1]) {
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

	// [3] - Color of side to move.
	if matches[3] == `w` {
		p.color = White
	} else {
		p.color = Black
	}

	// [4] - Castle rights.
	for _, char := range(matches[4]) {
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

	// [5] - En-passant square.
	if matches[5] != `-` {
		p.enpassant = Square(int(matches[5][1] - '1'), int(matches[5][0] - 'a'))

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
