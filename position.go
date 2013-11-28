package lape

import(
        `bytes`
        `fmt`
)

type Position struct {
        game      *Game
        pieces    [64]Piece // Position as an array of 64 squares with pieces on them.
        board     [3]Bitmask // Position as a bitmask: [0] white pieces only, [1] black pieces, and [2] all pieces.
        attacks   [3]Bitmask // [0] all squares attacked by white, [1] by black, [2] by either white or black.
        count     map[Piece]int // Counts of each piece on the board, ex. white pawns: 6, etc.
        outposts  map[Piece]*Bitmask // Bitmasks of each piece on the board, ex. white pawns, black king, etc.
        targets   map[Piece]*Bitmask // Bitmasks of target attack squares for each piece on the board.
}

func (p *Position) Initialize(game *Game, pieces [64]Piece) *Position {
        p.game = game
        p.pieces = pieces

        p.count = make(map[Piece]int)
        p.targets = make(map[Piece]*Bitmask)
        p.outposts = make(map[Piece]*Bitmask)
        for piece := Piece(PAWN); piece <= Piece(KING); piece++ {
                p.targets[piece] = new(Bitmask)
                p.targets[piece | BLACK] = new(Bitmask)
                p.outposts[piece] = new(Bitmask)
                p.outposts[piece | BLACK] = new(Bitmask)
        }

        return p.setupPosition().setupAttacks()
}

func (p *Position) MakeMove(game *Game, move *Move) *Position {
        fmt.Printf("Making move %s for %s\n", move, C(move.Piece.Color()))
        pieces := p.pieces
        pieces[move.From] = 0
        pieces[move.To] = move.Piece

        return new(Position).Initialize(game, pieces)
}

func (p *Position) Score(depth, color int, alpha, beta float64) float64 {
        //fmt.Printf("Score(depth: %d, color: %d, alpha: %.2f, beta: %.2f)\n", depth, color, alpha, beta)
        if depth == 0 {
                return p.Evaluate(color)
        }

        color ^= 1
        moves := p.Moves(color)
	for i, move := range moves {
	        score := -p.MakeMove(p.game, move).Score(depth-1, color, -beta, -alpha)
		if score >= beta {
                        fmt.Printf("\n  Done at depth %d after move %d out of %d for %s\n", depth, i+1, len(moves), C(color))
                        fmt.Printf("  Searched %v\n", moves[:i+1])
                        fmt.Printf("  Skipping %v\n", moves[i+1:])
                        fmt.Printf("  Picking %v\n\n", move)
			return score
		}
                if score > alpha {
                        alpha = score
                }
	}
	return alpha
}

func (p *Position) Evaluate(color int) float64 {
        return p.game.players[color].brain.Evaluate(p)
}

// Returns bitmask of attack targets for the piece at the index.
func (p *Position) Targets(index int) *Bitmask {
        return p.game.attacks.Targets(index, p.pieces[index], p.board)
}

// All moves.
func (p *Position) Moves(color int) (moves []*Move) {
        for side := p.board[color]; !side.IsEmpty(); {
                index := side.FirstSet()
                piece := p.pieces[index]
                moves = append(moves, p.PossibleMoves(index, piece)...)
                side.Clear(index)
        }
        fmt.Printf("%d candidates for %s: %v\n", len(moves), C(color), moves)

        return
}

// All moves for the piece in certain square.
func (p *Position) PossibleMoves(index int, piece Piece) (moves []*Move) {
        targets := p.Targets(index)
        for !targets.IsEmpty() {
                target := targets.FirstSet()
                moves = append(moves, new(Move).Initialize(index, target, piece, p.pieces[target]))
                targets.Clear(target)
        }

        return
}

func (p *Position) IsCheck(color int) bool {
        return false
}

func (p *Position) setupPosition() *Position {
        for i, piece := range p.pieces {
                if piece != 0 {
                        p.outposts[piece].Set(i)
                        p.board[piece.Color()].Set(i)
                        p.count[piece]++
                }
        }
        p.board[2] = p.board[0] // Combined board starts off with white pieces...
        p.board[2].Combine(p.board[1]) // ...and adds black ones.

        fmt.Printf("\n%s\n", p)
        return p
}

// attacks   [3]Bitmask // [0] all squares attacked by white, [1] by black, [2] by either white or black.
// targets   map[Piece]*Bitmask // Bitmasks of target attack squares for each piece on the board.
func (p *Position) setupAttacks() *Position {
        for i, piece := range p.pieces {
                if piece != 0 {
                        p.targets[piece] = p.Targets(i)
                        p.attacks[piece.Color()].Combine(*p.targets[piece])
                }
        }
        p.attacks[2] = p.attacks[0] // Combined board starts off with white pieces...
        p.attacks[2].Combine(p.attacks[1]) // ...and adds black ones.

        return p
}

func (p *Position) String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h\n")
	for row := 7;  row >= 0;  row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0;  col <= 7; col++ {
			index := Index(row, col)
			buffer.WriteByte(' ')
			if piece := p.pieces[index]; piece != 0 {
				buffer.WriteString(piece.String())
			} else {
				buffer.WriteString("\u22C5")
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
