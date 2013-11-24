package lape

import (
        `bytes`
        `fmt`
	`math`
)

type Game struct {
	pieces	[64]Piece
	players	[2]*Player
        attacks *Attack
        current int
}

func (g *Game)Initialize() *Game {
        g.players[0] = new(Player).Initialize(0)
        g.players[1] = new(Player).Initialize(1)
        g.attacks = new(Attack).Initialize()
        g.current = 0

        return g
}

func (g *Game)Set(row, col int, piece Piece) *Game {
        g.pieces[Index(row, col)] = piece

        return g
}

func (g *Game)Get(row, col int) Piece {
        return g.pieces[Index(row, col)]
}

func (g *Game)SetInitialPosition() *Game {

        for color := 0;  color < 2; color ++ {
                g.Set(color * 7, 0, Rook(color))
                g.Set(color * 7, 1, Knight(color))
                g.Set(color * 7, 2, Bishop(color))
                g.Set(color * 7, 3, Queen(color))
                g.Set(color * 7, 4, King(color))
                g.Set(color * 7, 5, Bishop(color))
                g.Set(color * 7, 6, Knight(color))
                g.Set(color * 7, 7, Rook(color))
                for col := 1; col <= 7; col++ {
                        g.Set(color * 5 + 1, col, Pawn(color))
                }
        }
        // g.Set(6, 0, Pawn(0))
        // 
        // g.Set(5, 3, Pawn(1))
        // g.Set(4, 3, Rook(1))
        // g.Set(3, 3, Pawn(1))
        // g.Set(4, 4, Pawn(1))
        // 
        // g.Set(3, 0, Rook(0))
        // g.Set(3, 1, Pawn(0))
        // g.Set(2, 0, Pawn(0))
        return g
}

func (g *Game)Search(depth int) (best *Move) {
        position := new(Position).Initialize(g, g.pieces)
        moves := position.Moves(g.current)
        estimate := float64(-math.MaxInt32)

	for i, move := range moves {
		score := -position.MakeMove(g, move).Score(depth*2-1, g.current, float64(-math.MaxInt32), float64(math.MaxInt32))

                fmt.Printf("  %d/%d: %s for %s, score is %.2f\n", i+1, len(moves), move, C(g.current), score)

                if score >= estimate {
                        estimate = score
                        best = move
                        fmt.Printf("  New best move for %s is %s\n\n", C(g.current), best)
                }
	}
	return
}

func (g *Game)String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h\n")
	for row := 7;  row >= 0;  row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0;  col <= 7; col++ {
			index := Index(row, col)
			buffer.WriteByte(' ')
			if piece := g.pieces[index]; piece != 0 {
				buffer.WriteString(piece.String())
			} else {
				buffer.WriteString("\u22C5")
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
