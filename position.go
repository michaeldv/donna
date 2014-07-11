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
	reversible   bool        // Is this position reversible?
	castles      uint8       // Castle rights mask.
	hash         uint64      // Polyglot hash value for the position.
	hashPawns    uint64      // Polyglot hash value for position's pawn structure.
	hashMaterial uint64      // Polyglot hash for position's playing material.
	board        Bitmask     // Bitmask of all pieces on the board.
	king         [2]int      // King's square for both colors.
	count        [14]int     // Counts of each piece on the board.
	pieces       [64]Piece   // Array of 64 squares with pieces on them.
	outposts     [14]Bitmask // Bitmasks of each piece on the board; [0] all white, [1] all black.
	tally        Score       // Positional valuation score based on PST.
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
		}
	}

	p.reversible = true
	p.board = p.outposts[White] | p.outposts[Black]
	p.hash, p.hashPawns, p.hashMaterial = p.polyglot()
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

func (p *Position) movePiece(piece Piece, from, to int) *Position {
	p.pieces[from], p.pieces[to] = 0, piece
	p.outposts[piece] ^= bit[from] | bit[to]
	p.outposts[piece.color()] ^= bit[from] | bit[to]

	// Update position's hash values.
	random := piece.polyglot(from) ^ piece.polyglot(to)
	p.hash ^= random
	if piece.isPawn() {
		p.hashPawns ^= random
	}

	// Update positional score.
	p.tally.subtract(pst[piece][from]).add(pst[piece][to])

	return p
}

func (p *Position) promotePawn(piece Piece, from, to int, promo Piece) *Position {
	p.pieces[from], p.pieces[to] = 0, promo
	p.outposts[piece] ^= bit[from]
	p.outposts[promo] ^= bit[to]
	p.outposts[piece.color()] ^= bit[from] | bit[to]
	p.count[piece]--

	// Update position's hash values.
	random := piece.polyglot(from)
	p.hash ^= random ^ promo.polyglot(to)
	p.hashPawns ^= random
	p.hashMaterial ^= piece.polyglot(p.count[piece]) ^ promo.polyglot(p.count[promo])
	p.count[promo]++

	// Update positional score.
	p.tally.subtract(pst[piece][from]).add(pst[promo][to])

	return p
}

func (p *Position) capturePiece(capture Piece, from, to int) *Position {
	p.outposts[capture] ^= bit[to]
	p.outposts[capture.color()] ^= bit[to]
	p.count[capture]--

	// Update position's hash values and count.
	random := capture.polyglot(to)
	p.hash ^= random
	if capture.isPawn() {
		p.hashPawns ^= random
	}
	p.hashMaterial ^= capture.polyglot(p.count[capture])

	// Update positional score.
	p.tally.subtract(pst[capture][to])

	return p
}

func (p *Position) captureEnpassant(capture Piece, from, to int) *Position {
	enpassant := to - eight[capture.color()^1]

	p.pieces[enpassant] = 0
	p.outposts[capture] ^= bit[enpassant]
	p.outposts[capture.color()] ^= bit[enpassant]
	p.count[capture]--

	// Update position's hash values and count.
	random := capture.polyglot(enpassant)
	p.hash ^= random
	p.hashPawns ^= random
	p.hashMaterial ^= capture.polyglot(p.count[capture])

	// Update positional score.
	p.tally.subtract(pst[capture][enpassant])

	return p
}

func (p *Position) MakeMove(move Move) *Position {
	color := move.color()
	from, to, piece, capture := move.split()

	// Copy over the contents of previous tree node to the current one.
	node++
	tree[node] = *p // => tree[node] = tree[node - 1]
	pp := &tree[node]

	pp.enpassant, pp.reversible = 0, true

	if capture != 0 {
		pp.reversible = false
		if to != 0 && to == p.enpassant {
			pp.captureEnpassant(pawn(color^1), from, to)
		} else {
			pp.capturePiece(capture, from, to)
		}
	}

	if promo := move.promo(); promo == 0 {
		pp.movePiece(piece, from, to)

		if piece.isKing() {
			pp.king[color] = to
			if move.isCastle() {
				pp.reversible = false
				switch to {
				case G1:
					pp.movePiece(Rook, H1, F1)
				case C1:
					pp.movePiece(Rook, A1, D1)
				case G8:
					pp.movePiece(BlackRook, H8, F8)
				case C8:
					pp.movePiece(BlackRook, A8, D8)
				}
			}
		} else if piece.isPawn() {
			pp.reversible = false
			if move.isEnpassant() {
				pp.enpassant = from + eight[color] // Save the en-passant square.
				pp.hash ^= hashEnpassant[Col(pp.enpassant)]
			}
		}
	} else {
		pp.reversible = false
		pp.promotePawn(piece, from, to, promo)
	}

	pp.board = pp.outposts[White] | pp.outposts[Black]

	// Ready to validate new position we have after making the move: if it is not
	// valid then revert back the node pointer and return nil.
	if pp.isInCheck(color) {
		node--
		return nil
	}

	// OK, the position after making the move is valid: all that's left is updating
	// castle rights, finishing off incremental hash value, and flipping the color.
	pp.castles &= castleRights[from] & castleRights[to]
	pp.hash ^= hashCastle[p.castles] ^ hashCastle[pp.castles]

	if p.enpassant != 0 {
		pp.hash ^= hashEnpassant[Col(p.enpassant)]
	}

	pp.hash ^= polyglotRandomWhite
	pp.color ^= 1 // <-- Flip side to move.

	return &tree[node] // pp
}

// Makes "null" move by copying over previous node position (i.e. preserving all pieces
// intact) and flipping the color.
func (p *Position) MakeNullMove() *Position {
	node++
	tree[node] = *p // => tree[node] = tree[node - 1]
	pp := &tree[node]

	// Flipping side to move obviously invalidates the enpassant square.
	if pp.enpassant != 0 {
		pp.hash ^= hashEnpassant[Col(pp.enpassant)]
		pp.enpassant = 0
	}
	pp.hash ^= polyglotRandomWhite
	pp.color ^= 1 // <-- Flip side to move.

	return &tree[node] // pp
}

// Restores previous position effectively taking back the last move made.
func (p *Position) TakeBack(move Move) *Position {
	node--
	return &tree[node]
}

func (p *Position) TakeBackNullMove() *Position {
	p.hash ^= polyglotRandomWhite
	p.color ^= 1

	return p.TakeBack(Move(0))
}

func (p *Position) isInCheck(color int) bool {
	return p.isAttacked(p.king[color], color^1)
}

func (p *Position) isNull() bool {
	return node > 0 && tree[node].board == tree[node-1].board
}

func (p *Position) isRepetition() bool {
	if !p.reversible {
		return false
	}

	for reps, prevNode := 1, node-1; prevNode >= 0; prevNode-- {
		if !tree[prevNode].reversible {
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

func (p *Position) isInsufficient() bool {
	return false
}

func (p *Position) canCastle(color int) (kingside, queenside bool) {
	attacks := p.allAttacks(color ^ 1)
	kingside = p.castles & castleKingside[color] != 0 &&
		(gapKing[color] & p.board == 0) &&
		(castleKing[color] & attacks == 0)

	queenside = p.castles & castleQueenside[color] != 0 &&
		(gapQueen[color] & p.board == 0) &&
		(castleQueen[color] & attacks == 0)
	return
}

// Reports game status for current position or after the given move. The status
// help to determine whether to continue with search or if the game is over.
func (p *Position) status(move Move, blendedScore int) int {
	if move != Move(0) {
		p = p.MakeMove(move)
		defer func() { p = p.TakeBack(move) }()
	}

	switch ply, score := Ply(), Abs(blendedScore); score {
	case 0:
		if ply == 1 {
			if p.isRepetition() {
				return Repetition
			} else if p.isInsufficient() {
				return Insufficient
			}
		}
		if !NewGen(p, ply+1).generateMoves().anyValid(p) {
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

// Computes initial values of position's polyglot hash, pawn hash, and material
// hash. When making a move these values get updated incrementally.
func (p *Position) polyglot() (hash, hashPawns, hashMaterial uint64) {
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

	for piece := Pawn; piece <= BlackQueen; piece++ {
		for count := 0; count < p.count[piece]; count++ {
			hashMaterial ^= Piece(piece).polyglot(count)
		}
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
