// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`bytes`
	`fmt`
	`regexp`
	`strings`
)

var tree [1024]Position
var node, rootNode int

type Position struct {
	game         *Game
	enpassant    int         // En-passant square caused by previous move.
	color        int         // Side to make next move.
	balance      int 	 // Material balance index.
	hash         uint64      // Polyglot hash value for the position.
	hashPawns    uint64      // Polyglot hash value for position's pawn structure.
	board        Bitmask     // Bitmask of all pieces on the board.
	king         [2]int      // King's square for both colors.
	count        [14]int     // Counts of each piece on the board.
	pieces       [64]Piece   // Array of 64 squares with pieces on them.
	outposts     [14]Bitmask // Bitmasks of each piece on the board; [0] all white, [1] all black.
	tally        Score       // Positional valuation score based on PST.
	reversible   bool        // Is this position reversible?
	castles      uint8       // Castle rights mask.
}

func NewPosition(game *Game, white, black string, color int) *Position {
	tree[node] = Position{game: game, color: color}
	p := &tree[node]

	p.setupSide(strings.Split(white, `,`), White)
	p.setupSide(strings.Split(black, `,`), Black)

	p.castles = castleKingside[White] | castleQueenside[White] | castleKingside[Black] | castleQueenside[Black]
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
			if piece.isKing() {
				p.king[piece.color()] = square
			}
			p.balance += materialBalance[piece]
		}
	}

	p.reversible = true
	p.board = p.outposts[White] | p.outposts[Black]
	p.hash, p.hashPawns = p.polyglot()
	p.tally = p.valuation()

	return p
}

func (p *Position) setupSide(moves []string, color int) *Position {
	re := regexp.MustCompile(`([KQRBN]?)([a-h])([1-8])`)

	for _, move := range moves {
		arr := re.FindStringSubmatch(move)
		if len(arr) == 0 {
			panic(fmt.Sprintf("Invalid notation '%s' for %s\n", move, C(color)))
		}
		name, col, row := arr[1], int(arr[2][0]-'a'), int(arr[3][0]-'1')

		var piece Piece
		switch name {
		case `K`: piece = king(color)
		case `Q`: piece = queen(color)
		case `R`: piece = rook(color)
		case `B`: piece = bishop(color)
		case `N`: piece = knight(color)
		default:  piece = pawn(color)
		}
		p.pieces[Square(row, col)] = piece
	}

	return p
}

// Sets up initial chess position.
func NewInitialPosition(game *Game) *Position {
	return NewPositionFromFEN(game, `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`)
}

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
	// fmt.Printf("%q\n", matches)
	if len(matches) < 4 {
		return nil
	}

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

// Computes initial values of position's polyglot hash, pawn hash, and material
// hash. When making a move these values get updated incrementally.
func (p *Position) polyglot() (hash, hashPawns uint64) {
	board := p.board
	for board != 0 {
		square := board.pop()
		piece := p.pieces[square]
		random := piece.polyglot(square)
		hash ^= random
		if piece.isPawn() {
			hashPawns ^= random
		}
	}

	hash ^= hashCastle[p.castles]
	if p.enpassant != 0 {
		hash ^= hashEnpassant[Col(p.enpassant)]
	}
	if p.color == White {
		hash ^= polyglotRandomWhite
	}

	return
}

// Computes positional valuation score based on PST. When making a move the
// valuation tally gets updated incrementally.
func (p *Position) valuation() (score Score) {
	board := p.board
	for board != 0 {
		square := board.pop()
		piece := p.pieces[square]
		score.add(pst[piece][square])
	}
	return
}

// Stub.
func (p *Position) isInsufficient() bool {
	return false
}

// Reports game status for current position or after the given move. The status
// help to determine whether to continue with search or if the game is over.
func (p *Position) status(move Move, blendedScore int) int {
	if move != Move(0) {
		p = p.MakeMove(move)
		defer func() { p = p.UndoLastMove() }()
	}

	switch ply, score := Ply(), Abs(blendedScore); score {
	case 0:
		if ply == 1 {
			if p.thirdRepetition() {
				return Repetition
			} else if p.isInsufficient() {
				return Insufficient
			}
		}
		if !NewGen(p, MaxPly).generateMoves().anyValid(p) {
			return Stalemate
		}
	case Checkmate - ply:
		if p.isInCheck(p.color) {
			if p.color == White {
				return BlackWon
			}
			return WhiteWon
		}
		return Stalemate
	default:
		if score > Checkmate-MaxDepth && (score+ply)/2 > 0 {
			if p.color == White {
				return BlackWinning
			}
			return WhiteWinning
		}
	}
	return InProgress
}

// Encodes position as FEN string.
func (p *Position) fen() (fen string) {
	fancy := engine.fancy
	engine.fancy = false; defer func() { engine.fancy = fancy }()

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

func (p *Position) String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h")
	if !p.isInCheck(p.color) {
		buffer.WriteString("\n")
	} else {
		buffer.WriteString("  Check to " + C(p.color) + "\n")
	}
	for row := 7; row >= 0; row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0; col <= 7; col++ {
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
