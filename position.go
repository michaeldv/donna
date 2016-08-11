// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`bytes`
	`fmt`
	`strconv`
	`strings`
)

var tree [1024]Position
var node, rootNode int

type Position struct {		 // 224 bytes long.
	id           uint64	 // Polyglot hash value for the position.
	pawnId       uint64	 // Polyglot hash value for position's pawn structure.
	board        Bitmask	 // Bitmask of all pieces on the board.
	king         [2]int	 // King's square for both colors.
	pieces       [64]Piece	 // Array of 64 squares with pieces on them.
	outposts     [14]Bitmask // Bitmasks of each piece on the board; [0] all white, [1] all black.
	tally        Score	 // Positional valuation score based on PST.
	balance      int	 // Material balance index.
	score        int	 // Blended evaluation score.
	color        int	 // Side to make next move.
	enpassant    int	 // En-passant square caused by previous move.
	count50      int	 // 50 moves rule counter.
	reversible   bool	 // Is this position reversible?
	castles      uint8	 // Castle rights mask.
}

func NewPosition(game *Game, white, black string) *Position {
	tree[node] = Position{}
	p := &tree[node]

	p.setupSide(white, White).setupSide(black, Black)

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
		if piece.some() {
			p.outposts[piece].set(square)
			p.outposts[piece.color()].set(square)
			if piece.isKing() {
				p.king[piece.color()] = square
			}
			p.balance += materialBalance[piece]
		}
	}

	p.reversible = true
	p.board = p.outposts[White] | p.outposts[Black]
	p.id, p.pawnId = p.polyglot()
	p.tally = p.valuation()
	p.score = Unknown

	return p
}

// Parses Donna chess format string for one side. Besides [K]ing, [Q]ueen, [R]ook,
// [B]ishop, and k[N]ight the following pseudo pieces could be specified:
//
// [M]ove:      specifies the right to move along with the optional move number.
//              For example, "M42" for Black means the Black is making 42nd move.
//              Default value is "M1" for White.
//
// [C]astle:    specifies castle right squares. For example, "Cg1" and "Cc8" encode
//              allowed kingside castle for White, and queenside castle for Black.
//              By default all castles are allowed, i.e. defult value is "Cc1,Cg1"
//              for White and "Cc8,Cg8" for Black. The actual castle rights are
//              checked during position setup to make sure they do not violate
//              chess rules. If castle rights are specified incorrectly they are
//              quietly ignored.
//
// [E]npassant: specifies en-passant square if any. For example, "Ed3" marks D3
//              square as en-passant. Default value is no en-passant.
//
func (p *Position) setupSide(str string, color int) *Position {
	invalid := func (move string, color int) {
		// Don't panic.
		panic(fmt.Sprintf("Invalid notation '%s' for %s\n", move, C(color)))
	}

	for _, move := range strings.Split(str, `,`) {
		if move[0] == 'M' { // TODO: parse move number.
			p.color = color
		} else {
			arr := reMove.FindStringSubmatch(move)
			if len(arr) == 0 {
				invalid(move, color)
			}
			square := square(int(arr[3][0]-'1'), int(arr[2][0]-'a'))

			switch move[0] {
			case 'K':
				p.pieces[square] = king(color)
			case 'Q':
				p.pieces[square] = queen(color)
			case 'R':
				p.pieces[square] = rook(color)
			case 'B':
				p.pieces[square] = bishop(color)
			case 'N':
				p.pieces[square] = knight(color)
			case 'E':
				p.enpassant = square
			case 'C':
				if (square == C1 + int(color)) || (square == C8 + int(color)) {
					p.castles |= castleQueenside[color]
				} else if (square == G1 + int(color)) || (square == G8 + int(color)) {
					p.castles |= castleKingside[color]
				}
			default:
				// When everything else fails, read the instructions.
				p.pieces[square] = pawn(color)
			}
		}
	}

	return p
}

// Sets up initial chess position.
func NewInitialPosition(game *Game) *Position {
	return NewPositionFromFEN(game, `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`)
}

// Decodes FEN string and creates new position.
func NewPositionFromFEN(game *Game, fen string) *Position {
	tree[node] = Position{}
	p := &tree[node]

	// Expected matches of interest are as follows:
	// [0] - Pieces (entire board).
	// [1] - Color of side to move.
	// [2] - Castle rights.
	// [3] - En-passant square.
	// [4] - Number of half-moves.
	// [5] - Number of full moves.
	matches := strings.Split(fen, ` `)
	if len(matches) < 4 {
		return nil
	}

	// [0] - Pieces (entire board).
	sq := A8
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
			p.king[White] = sq
		case 'k':
			piece = BlackKing
			p.king[Black] = sq
		case '/':
			sq -= 16
		case '1', '2', '3', '4', '5', '6', '7', '8':
			sq += int(char - '0')
		}
		if piece.some() {
			p.pieces[sq] = piece
			p.outposts[piece].set(sq)
			p.outposts[piece.color()].set(sq)
			p.balance += materialBalance[piece]
			sq++
		}
	}

	// [1] - Color of side to move.
	p.color = let(matches[1] == `w`, White, Black)

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
		p.enpassant = square(int(matches[3][1] - '1'), int(matches[3][0] - 'a'))

	}

	// [4] - Number of half-moves.
	if n, err := strconv.Atoi(matches[4]); err == nil {
		p.count50 = n
	}

	p.reversible = true
	p.board = p.outposts[White] | p.outposts[Black]
	p.id, p.pawnId = p.polyglot()
	p.tally = p.valuation()
	p.score = Unknown

	return p
}

// Computes initial values of position's polyglot hash and pawn hash. When
// making a move these values get updated incrementally.
func (p *Position) polyglot() (hash, pawnHash uint64) {
	for board := p.board; board.any(); board = board.pop() {
		square := board.first()
		piece := p.pieces[square]
		random := piece.polyglot(square)
		hash ^= random
		if piece.isPawn() {
			pawnHash ^= random
		}
	}

	hash ^= hashCastle[p.castles]
	if p.enpassant != 0 {
		hash ^= hashEnpassant[p.enpassant & 7] // p.enpassant column.
	}
	if p.color == White {
		hash ^= polyglotRandomWhite
	}

	return hash, pawnHash
}

// Computes positional valuation score based on PST. When making a move the
// valuation tally gets updated incrementally.
func (p *Position) valuation() (score Score) {
	for bm := p.board; bm.any(); bm = bm.pop() {
		square := bm.first()
		piece := p.pieces[square]
		score.add(pst[piece][square])
	}

	return score
}

// Returns true if material balance is insufficient to win the game.
func (p *Position) insufficient() bool {
	return materialBase[p.balance].flags & materialDraw != 0
}

// Reports game status for current position or after the given move. The status
// helps to determine whether to continue with search or if the game is over.
func (p *Position) status(move Move, blendedScore int) int {
	if move.some() {
		p = p.makeMove(move)
		defer func() { p = p.undoLastMove() }()
	}

	switch ply, score := ply(), abs(blendedScore); score {
	case 0:
		if ply == 1 {
			if p.insufficient() {
				return Insufficient
			} else if p.thirdRepetition() {
				return Repetition
			} else if p.fifty() {
				return FiftyMoves
			}
		}
		if !NewGen(p, MaxPly).generateMoves().anyValid() {
			return Stalemate
		}
	case Checkmate - ply:
		if p.isInCheck(p.color) {
			return let(p.color == White, BlackWon, WhiteWon)
		}
		return Stalemate
	default:
		if score > Checkmate - MaxDepth && (score + ply) / 2 > 0 {
			return let(p.color == White, BlackWinning, WhiteWinning)
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
			square := square(row, col)
			piece := p.pieces[square]

			if piece.some() {
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
		row, col := coordinate(int(p.enpassant))
		fen += fmt.Sprintf(` %c%d`, col + 'a', row + 1)
	} else {
		fen += ` -`
	}

	// Number of half-moves (50 moves counter).
	fen += fmt.Sprintf(` %d`, p.count50)

	// TODO: Number of full moves.
	fen += ` 1`

	return
}

// Encodes position as DCF string (Donna Chess Format).
func (p *Position) dcf() string {
	fancy := engine.fancy
	engine.fancy = false; defer func() { engine.fancy = fancy }()

	encode := func (square int) string {
		var buffer bytes.Buffer

		buffer.WriteByte(byte(col(square)) + 'a')
		buffer.WriteByte(byte(row(square)) + '1')

		return buffer.String()
	}

	var pieces [2][]string

	for color := White; color <= Black; color++ {
		// Right to move and (TODO) move number.
		if color == p.color && color == Black {
			pieces[color] = append(pieces[color], `M`)
		}

		// King.
		pieces[color] = append(pieces[color], `K` + encode(int(p.king[color])))

		// Queens, Rooks, Bishops, and Knights.
		for outposts := p.outposts[queen(color)]; outposts.any(); outposts = outposts.pop() {
			pieces[color] = append(pieces[color], `Q` + encode(outposts.first()))
		}
		for outposts := p.outposts[rook(color)]; outposts.any(); outposts = outposts.pop() {
			pieces[color] = append(pieces[color], `R` + encode(outposts.first()))
		}
		for outposts := p.outposts[bishop(color)]; outposts.any(); outposts = outposts.pop() {
			pieces[color] = append(pieces[color], `B` + encode(outposts.first()))
		}
		for outposts := p.outposts[knight(color)]; outposts.any(); outposts = outposts.pop() {
			pieces[color] = append(pieces[color], `N` + encode(outposts.first()))
		}

		// Castle rights.
		if p.castles & castleQueenside[color] == 0 || p.castles & castleKingside[color] == 0 {
			if p.castles & castleQueenside[color] != 0 {
				pieces[color] = append(pieces[color], `C` + encode(C1 + 56 * int(color)))
			}
			if p.castles & castleKingside[color] != 0 {
				pieces[color] = append(pieces[color], `C` + encode(G1 + 56 * int(color)))
			}
		}

		// En-passant square if any. Note that this gets assigned to the
		// current side to move.
		if p.enpassant != 0 && color == p.color {
			pieces[color] = append(pieces[color], `E` + encode(int(p.enpassant)))
		}

		// Pawns.
		for outposts := p.outposts[pawn(color)]; outposts.any(); outposts = outposts.pop() {
			pieces[color] = append(pieces[color], encode(outposts.first()))
		}
	}

	return strings.Join(pieces[White], `,`) + ` : ` + strings.Join(pieces[Black], `,`)
}

func (p *Position) String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h  " + C(p.color) + " to move")
	if !p.isInCheck(p.color) {
		buffer.WriteString("\n")
	} else {
		buffer.WriteString(", check\n")
	}
	for row := 7; row >= 0; row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0; col <= 7; col++ {
			buffer.WriteByte(' ')
			if piece := p.pieces[square(row, col)]; piece.some() {
				buffer.WriteString(piece.String())
			} else {
				buffer.WriteString("\u22C5")
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
