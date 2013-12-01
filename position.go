package lape

import(
        `bytes`
        `fmt`
        `math`
)

type Position struct {
        game      *Game
        pieces    [64]Piece // Position as an array of 64 squares with pieces on them.
        targets   [64]Bitmask // Attack targets for each piece on the board.
        board     [3]Bitmask // Position as a bitmask: [0] white pieces only, [1] black pieces, and [2] all pieces.
        attacks   [3]Bitmask // [0] all squares attacked by white, [1] by black, [2] by either white or black.
        count     map[Piece]int // Counts of each piece on the board, ex. white pawns: 6, etc.
        outposts  map[Piece]*Bitmask // Bitmasks of each piece on the board, ex. white pawns, black king, etc.
        check     bool // Is there a check?
        next      int // Side to make next move.
}

func (p *Position) Initialize(game *Game, pieces [64]Piece, color int) *Position {
        p.game = game
        p.pieces = pieces
        p.next = color

        p.count = make(map[Piece]int)
        p.outposts = make(map[Piece]*Bitmask)
        for piece := Piece(PAWN); piece <= Piece(KING); piece++ {
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

        return new(Position).Initialize(game, pieces, move.Piece.Color()^1)
}

func (p *Position) Score(depth, color int, alpha, beta float64) float64 {
        fmt.Printf("Score(depth: %d, color: %s, alpha: %.2f, beta: %.2f)\n", depth, C(color), alpha, beta)

        if depth == 0 {
                return p.Evaluate(color)
        }

        color ^= 1
        moves := p.Moves(color)
        if len(moves) > 0 {
                for i, move := range moves {
                        score := -p.MakeMove(p.game, move).Score(depth-1, color, -beta, -alpha)
                        fmt.Printf("Move %d/%d: %s (%d): score: %.2f, alpha: %.2f, beta: %.2f\n", i+1, len(moves), C(color), depth, score, alpha, beta)
                        //if score >= beta || math.Abs(score) == math.Abs(float64(math.MinInt32)) {
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
        } else if p.IsCheck(color) {
                return float64(math.MinInt32) // <-- Checkmate value.
        } else {
                fmt.Printf("Stalemate\n") // TODO
        }

        fmt.Printf("End of Score(depth: %d, color: %s, alpha: %.2f, beta: %.2f) => %.2f\n", depth, C(color), alpha, beta, alpha)
	return alpha
}

func (p *Position) Evaluate(color int) float64 {
        return p.game.players[color].brain.Evaluate(p)
}

// Returns bitmask of attack targets for the piece at the index.
func (p *Position) Targets(index int) *Bitmask {
        return p.game.attacks.Targets(index, p)
}

// All moves.
func (p *Position) Moves(color int) (moves []*Move) {
        for i, piece := range p.pieces {
                if piece != 0 && piece.Color() == color {
                        moves = append(moves, p.PossibleMoves(i, piece)...)
                }
        }
        fmt.Printf("%d candidates for %s: %v\n", len(moves), C(color), moves)
        if len(moves) > 1 {
                moves = p.Reorder(moves)
                fmt.Printf("%d candidates for %s (reordered): %v\n", len(moves), C(color), moves)
        }

        return
}

// All moves for the piece in certain square.
func (p *Position) PossibleMoves(index int, piece Piece) (moves []*Move) {
        targets := p.targets[index]
        for !targets.IsEmpty() {
                target := targets.FirstSet()
                candidate := new(Move).Initialize(index, target, piece, p.pieces[target])
                if !p.MakeMove(p.game, candidate).IsCheck(piece.Color()) {
                        moves = append(moves, candidate)
                }
                targets.Clear(target)
        }

        return
}

func (p *Position) Reorder(moves []*Move) []*Move {
        var checks, captures, remaining []*Move

        for _, move := range moves {
                if p.MakeMove(p.game, move).IsCheck(move.Piece.Color()^1) {
                        checks = append(checks, move)
                } else if move.Captured != 0 {
                        captures = append(captures, move)
                } else {
                        remaining = append(remaining, move)
                }
        }

        return append(append(checks, captures...), remaining...)
}

func (p *Position) IsCheck(color int) bool {
        king := *p.outposts[King(color)]
        return king.Intersect(p.attacks[color^1]).IsNotEmpty()
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

        return p
}

func (p *Position) setupAttacks() *Position {
        var king [2]int
        for i, piece := range p.pieces {
                if piece != 0 {
                        p.targets[i] = *p.Targets(i)
                        p.attacks[piece.Color()].Combine(p.targets[i])
                        if piece.Kind() == KING {
                                king[piece.Color()] = i
                        }
                }
        }
        // Now that we have attack targets for both kings adjust them to make sure the
        // kings don't stomp on each other.
        white_king_targets, black_king_targets := p.targets[king[WHITE]], p.targets[king[BLACK]]
        p.targets[king[WHITE]].Exclude(black_king_targets)
        p.targets[king[BLACK]].Exclude(white_king_targets)

        // Combined board starts off with white pieces and adds black ones.
        p.attacks[2] = p.attacks[0]
        p.attacks[2].Combine(p.attacks[1])

        p.check = p.IsCheck(p.next)

        //fmt.Printf("\n%s\n", p)
        return p
}

func (p *Position) String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h")
        if !p.check {
                buffer.WriteString("\n")
        } else {
                buffer.WriteString("  Check to " + C(p.next) + "\n")
        }
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
