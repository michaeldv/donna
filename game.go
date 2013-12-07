package donna

import (
        `bytes`
        `fmt`
        `math`
        `regexp`
        `time`
)

const CHECKMATE = float64(math.MaxInt16)

type Game struct {
	pieces	[64]Piece
	players	[2]*Player
        attacks *Attack
        current int
        nodes   int
}

func NewGame() *Game {
        game := new(Game)
        game.players[0] = NewPlayer(game, WHITE)
        game.players[1] = NewPlayer(game, BLACK)
        game.attacks = NewAttack()
        game.current = WHITE

        return game
}

func (g *Game) Setup(white, black string) *Game {
	re := regexp.MustCompile(`\W+`)
	whitePieces, blackPieces := re.Split(white, -1), re.Split(black, -1)
	return g.SetupSide(whitePieces, 0).SetupSide(blackPieces, 1)
}

func (g *Game) SetupSide(moves []string, color int) *Game {
	re := regexp.MustCompile(`([KQRBN]?)([a-h])([1-8])`)

	for _, move := range moves {
		arr := re.FindStringSubmatch(move)
		if len(arr) == 0 {
			fmt.Printf("Invalid move '%s' for %s\n", move, C(color))
			return g
		}
		name, col, row := arr[1], arr[2][0]-'a', arr[3][0]-'1'

		var piece Piece
		switch name {
		case `K`:
			piece = King(color)
		case `Q`:
			piece = Queen(color)
		case `R`:
			piece = Rook(color)
		case `B`:
			piece = Bishop(color)
		case `N`:
			piece = Knight(color)
		default:
			piece = Pawn(color)
		}
		g.Set(int(row), int(col), piece)
	}
	return g
}

func (g *Game) Set(row, col int, piece Piece) *Game {
        g.pieces[Index(row, col)] = piece

        return g
}

func (g *Game) InitialPosition() *Game {
        return g.Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e2,f2,g2,h2`,
                       `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`)
}

func (g *Game) Think(maxDepth int) *Move {
        fmt.Println(`Depth     Nodes     Nodes/s     Score     Best`)
        for depth := 1; depth <= maxDepth; depth++ {
                g.nodes = 0
                start := time.Now()
                score := g.Analyze(depth)
                fmt.Printf("  %d      %6d     %7.1f   %7.1f     %v\n",
                        depth, g.nodes, float64(g.nodes)/time.Since(start).Seconds(), score, best[0][0 : bestlen[0]])
        }
        fmt.Printf("Best move: %s\n", best[0][0])
        return best[0][0]
}

func (g *Game) Analyze(depth int) float64 {
        position := NewPosition(g, g.pieces, g.current, Bitmask(0))
        return position.AlphaBeta(depth*2, 0, -CHECKMATE, CHECKMATE)
}

func (g *Game) Search(depth int) *Move {
        g.Analyze(depth)
        return best[0][0]
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
