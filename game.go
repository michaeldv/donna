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
        bestLine     [MaxPly][MaxPly]Move
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
			piece = king(color)
		case `Q`:
			piece = queen(color)
		case `R`:
			piece = rook(color)
		case `B`:
			piece = bishop(color)
		case `N`:
			piece = knight(color)
		default:
			piece = pawn(color)
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

func (g *Game) Think(maxDepth int, position *Position) Move {
        if position == nil {
                position = g.Start(White)
        }

        book := NewBook("./books/gm2001.bin") // From http://www.chess2u.com/t5834-gm-polyglot-book
        if move := book.pickMove(position); move != 0 {
                fmt.Printf("Book move: %s\n", move)
                return move
        }

        fmt.Println(`Depth/Time     Nodes      QNodes     Nodes/s   Score   Best`)
        for depth := 1; depth <= maxDepth; depth++ {
                g.nodes, g.qnodes = 0, 0
                start := time.Now()
                score := g.Analyze(depth, position)
                finish := time.Since(start).Seconds()
                g.Print(depth, score, finish)
                if g.IsOver(score) {
                        return 0
                }
        }
        fmt.Printf("Best move: %s\n", g.bestLine[0][0])
        return g.bestLine[0][0]
}

func (g *Game) Print(depth, score int, finish float64) {
        if absScore := Abs(score); absScore > 32500 {
                movesLeft := (Checkmate - absScore) / 2
                if movesLeft > 0 {
                        fmt.Printf(" %d %02d:%02d    %8d    %8d   %9.1f   x%-4d   %v\n",
                                depth, int(finish) / 60, int(finish) % 60, g.nodes, g.qnodes,
                                float64(g.nodes + g.qnodes) / finish, movesLeft,
                                g.bestLine[0][0 : g.bestLength[0]])
                } else {
                        fmt.Printf(" %d %02d:%02d    %8d    %8d   %9.1f   Checkmate\n",
                                depth, int(finish) / 60, int(finish) % 60, g.nodes, g.qnodes,
                                float64(g.nodes + g.qnodes) / finish)
                }
        } else {
                fmt.Printf(" %d %02d:%02d    %8d    %8d   %9.1f   %5.2f   %v\n",
                        depth, int(finish) / 60, int(finish) % 60, g.nodes, g.qnodes,
                        float64(g.nodes + g.qnodes) / finish, float64(score) / 100.0,
                        g.bestLine[0][0 : g.bestLength[0]])
        }
}

func (g *Game) IsOver(score int) bool {
        return Abs(score) == Checkmate
}

func (g *Game) Analyze(depth int, position *Position) int {
        score := position.search(depth*2, 0, -Checkmate, Checkmate)
        if position.color == Black {
                return -score
        }
        return score
}

func (g *Game) Start(color int) *Position {
        tree = [1024]Position{}
        node = 0
        g.bestLine   = [MaxPly][MaxPly]Move{}
        g.bestLength = [MaxPly]int{}

        return NewPosition(g, g.pieces, color, Flags{})
}

func (g *Game) Search(depth int) Move {
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
