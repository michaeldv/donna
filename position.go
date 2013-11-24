package lape

import(
        `bytes`
        `fmt`
        `math`
)

type Position struct {
        game      *Game
        pieces    [64]Piece
        board     [3]Bitmask // 0: white, 1: black, 2: both
        layout    map[Piece]*Bitmask
        count     map[int]int // number of white/black pieces
}

func (p *Position) Initialize(game *Game, pieces [64]Piece) *Position {
        p.game = game
        p.pieces = pieces

        p.count = make(map[int]int)
        p.layout = make(map[Piece]*Bitmask)
        for piece := Piece(PAWN); piece <= Piece(KING); piece++ {
                p.layout[piece] = new(Bitmask)
                p.layout[piece | 1] = new(Bitmask)
        }

        return p.setupBoard()
}

func (p *Position) ConsiderMove(game *Game, move *Move) *Position {
        fmt.Printf("ConsiderMove(*game, move: %s) color: %d\n", move, move.Piece.Color())
        pieces := p.pieces
        pieces[move.From] = 0
        pieces[move.To] = move.Piece

        return new(Position).Initialize(game, pieces)
}

func (p *Position) Score(depth, color int, alpha, beta float64) float64 {
        fmt.Printf("Score(depth: %d, color: %d, alpha: %f, beta: %f)\n", depth, color, alpha, beta)
        if depth == 0 {
                x1, x2, x3 := p.material(color), p.mobility(color), p.aggressiveness(color)
                fmt.Printf("Score: %f (mat: %f, mob: %f, agg: %f)\n", x1 + x2 + x3, x1, x2, x3)
                return x1 + x2 + x3
        }

        color ^= 1
	for _, move := range p.Moves(color) {
	        score := p.ConsiderMove(p.game, move).Score(depth-1, color, alpha, beta)
		if score >= beta {
			return beta
		}
                alpha = math.Max(alpha, score)
	}

	return alpha
}

// All moves.
func (p *Position) Moves(color int) (moves []*Move) {
        for side := p.board[color]; !side.IsEmpty(); {
                index := side.FirstSet()
                piece := p.pieces[index]
                moves = append(moves, p.PossibleMoves(index, piece)...)
                side.Clear(index)
        }
        fmt.Printf("%d candidates for %d: %v\n", len(moves), color, moves)

        return
}

// All moves for the piece in certain square.
func (p *Position) PossibleMoves(index int, piece Piece) (moves []*Move) {
        color := piece.Color()
        targets := p.game.attacks.Targets(index, piece, p.board)
        for !targets.IsEmpty() {
                target := targets.FirstSet()
                if p.board[color^1].IsSet(target) { // Target square is occupied by opposite color?
                        p.count[color+100]++ // Increment attacks count by white(100) or black(101)
                }
                moves = append(moves, new(Move).Initialize(index, target, piece, p.pieces[target]))
                targets.Clear(target)
        }

        return
}

func (p *Position) material(color int) float64 {
        opponent := color^1
        score := 1000 * (p.count[KING|color] - p.count[KING|opponent]) +
                9 * (p.count[QUEEN|color] - p.count[QUEEN|opponent]) +
                5 * (p.count[ROOK|color] - p.count[ROOK|opponent]) +
                3 * (p.count[BISHOP|color] - p.count[BISHOP|opponent]) +
                3 * (p.count[KNIGHT|color] - p.count[KNIGHT|opponent]) +
                1 * (p.count[PAWN|color] - p.count[PAWN|opponent])
        return float64(score) + 0.1 * float64(p.count[BISHOP|color] - p.count[BISHOP|opponent])
}

func (p *Position) mobility(color int) float64 {
        return 0.25 * float64(p.count[color] - p.count[color^1]) // TODO
}

func (p *Position) aggressiveness(color int) float64 {
        return 0.2 * float64(p.count[color+100] - p.count[(color^1)+100]) // TODO
}

func (p *Position) setupBoard() *Position {
        for i, piece := range p.pieces {
                if piece != 0 {
                        p.layout[piece].Set(i)
                        p.board[piece.Color()].Set(i)
                        p.count[int(piece)]++
                }
        }
        p.board[2] = p.board[0]
        p.board[2].Combine(p.board[1])
        p.count[0] = len(p.Moves(0))
        p.count[1] = len(p.Moves(1))

        fmt.Printf("\n%s\n", p)
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
        //buffer.WriteString(fmt.Sprintf("%v", p.count))
	return buffer.String()
}
