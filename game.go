// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
        `bytes`
        `fmt`
        `regexp`
        `time`
)

type Game struct {
	pieces	     [64]Piece
        nodes        int
        qnodes       int
        killers      [MaxPly][2]*Move
        bestLine     [MaxPly][MaxPly]*Move // Assuming max depth = 4 which makes it 8 plies.
        bestLength   [MaxPly]int
}

func NewGame() *Game {
        return new(Game)
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
        g.pieces[Square(row, col)] = piece

        return g
}

func (g *Game) InitialPosition() *Game {
        return g.Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e2,f2,g2,h2`,
                       `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`)
}

func (g *Game) Think(maxDepth int, position *Position) *Move {
        book := NewBook("./books/gm2001.bin") // From http://www.chess2u.com/t5834-gm-polyglot-book
        if position == nil {
                position = g.Start(White)
        }
        move := book.pickMove(position)
        if move != nil {
                fmt.Printf("Book move: %s\n", move)
                return move
        }

        // fmt.Printf("%s", position)
        fmt.Println(`Depth/Time     Nodes      QNodes     Nodes/s   Score   Best`)
        for depth := 1; depth <= maxDepth; depth++ {
                g.nodes, g.qnodes = 0, 0
                start := time.Now()
                score := g.Analyze(depth, position)
                finish := time.Since(start).Seconds()
                fmt.Printf(" %d %02d:%02d    %8d    %8d    %8.1f   %5s   %v\n",
                        depth, int(finish) / 60, int(finish) % 60, g.nodes, g.qnodes,
                        float64(g.nodes + g.qnodes) / finish, score,
                        g.bestLine[0][0 : g.bestLength[0]])
        }
        fmt.Printf("Best move: %s\n", g.bestLine[0][0])
        return g.bestLine[0][0]
}

func (g *Game) Analyze(depth int, position *Position) string {
        score := position.search(depth*2, 0, -Checkmate, Checkmate)
        if position.color == Black {
                score = -score
        }
        return fmt.Sprintf(`%.2f`, float64(score) / 100.0)
}

func (g *Game) Start(color int) *Position {
        tree = [1024]Position{}
        node = 0
        g.bestLine   = [MaxPly][MaxPly]*Move{}
        g.bestLength = [MaxPly]int{}
        g.killers    = [MaxPly][2]*Move{}

        return NewPosition(g, g.pieces, color, Flags{})
}

func (g *Game) Search(depth int) *Move {
        g.Analyze(depth, NewPosition(g, g.pieces, White, Flags{}))
        return g.bestLine[0][0]
}

func (g *Game)String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h\n")
	for row := 7;  row >= 0;  row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0;  col <= 7; col++ {
			square := Square(row, col)
			buffer.WriteByte(' ')
			if piece := g.pieces[square]; piece != 0 {
				buffer.WriteString(piece.String())
			} else {
				buffer.WriteString("\u22C5")
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
