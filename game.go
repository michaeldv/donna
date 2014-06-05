// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
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
	nodes      int
	qnodes     int
	token      uint8
	cache      Cache
	pieces     [64]Piece
	bestLine   [MaxPly][MaxPly]Move
	bestLength [MaxPly]int
	killers    [MaxPly][2]Move
	goodMoves  [14][64]int
}

func NewGame() *Game {
	return new(Game)
}

func (game *Game) Setup(white, black string) *Game {
	re := regexp.MustCompile(`\W+`)
	whiteSide, blackSide := re.Split(white, -1), re.Split(black, -1)
	return game.SetupSide(whiteSide, 0).SetupSide(blackSide, 1)
}

func (game *Game) CacheSize(megaBytes float32) *Game {
	game.cache = NewCache(megaBytes)
	return game
}

func (game *Game) SetupSide(moves []string, color int) *Game {
	re := regexp.MustCompile(`([KQRBN]?)([a-h])([1-8])`)

	for _, move := range moves {
		arr := re.FindStringSubmatch(move)
		if len(arr) == 0 {
			fmt.Printf("Invalid move '%s' for %s\n", move, C(color))
			return game
		}
		name, col, row := arr[1], int(arr[2][0]-'a'), int(arr[3][0]-'1')

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
		game.pieces[Square(row, col)] = piece
	}
	return game
}

func (game *Game) InitialPosition() *Game {
	return game.Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e2,f2,g2,h2`,
		`Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`)
}

func (game *Game) getReady() *Game {
	rootNode = node
	game.token++ // <-- Wraps around: ...254, 255, 0, 1...
	game.bestLine = [MaxPly][MaxPly]Move{}
	game.bestLength = [MaxPly]int{}
	game.goodMoves = [14][64]int{}
	game.killers = [MaxPly][2]Move{}

	return game
}

func (game *Game) Start(color int) *Position {
	tree, node = [1024]Position{}, 0
	game.getReady()

	return NewPosition(game, game.pieces, color)
}

func (game *Game) Think(requestedDepth int, position *Position) Move {
	if position == nil {
		position = game.Start(White)
	}

	book := NewBook("./books/gm2001.bin") // From http://www.chess2u.com/t5834-gm-polyglot-book
	if move := book.pickMove(position); move != 0 {
		fmt.Printf("Book move: %s\n", move)
		return move
	}

	game.getReady()
	move, score, status := Move(0), 0, InProgress

	fmt.Println(`Depth/Time     Nodes      QNodes     Nodes/s   Score   Best`)
	for depth := 1; depth <= Min(MaxDepth, requestedDepth); depth++ {
		game.nodes, game.qnodes = 0, 0
		start := time.Now()
		move, score = position.searchRoot(depth)
		finish := time.Since(start).Seconds()
		if position.color == Black {
			score = -score
		}

		status = position.status(move, score)
		game.printBestLine(depth, score, status, finish)

		// No reason to search deeper if no moves are available at current depth.
		if move == Move(0) {
			return move
		}

		// No reason to search deeper if the game is over or mate in X moves was
		// found at current depth.
		if status != InProgress {
			break
		}

	}
	fmt.Printf("\nDonna's move: %s\n\n", move)
	return move
}

func (game *Game) printBestLine(depth, score, status int, finish float64) {
	switch status {
	case WhiteWon:
		fmt.Printf("%2d %02d:%02d    %8d    %8d   %9.1f   1-0 White Checkmates\n",
			depth, int(finish)/60, int(finish)%60, game.nodes, game.qnodes,
			float64(game.nodes+game.qnodes)/finish)
	case BlackWon:
		fmt.Printf("%2d %02d:%02d    %8d    %8d   %9.1f   0-1 Black Checkmates\n",
			depth, int(finish)/60, int(finish)%60, game.nodes, game.qnodes,
			float64(game.nodes+game.qnodes)/finish)
	case Stalemate:
		fmt.Printf("%2d %02d:%02d    %8d    %8d   %9.1f   1/2 Stalemate\n",
			depth, int(finish)/60, int(finish)%60, game.nodes, game.qnodes,
			float64(game.nodes+game.qnodes)/finish)
	case Repetition:
		fmt.Printf("%2d %02d:%02d    %8d    %8d   %9.1f   1/2 Repetition\n",
			depth, int(finish)/60, int(finish)%60, game.nodes, game.qnodes,
			float64(game.nodes+game.qnodes)/finish)
	case WhiteWinning, BlackWinning:
		movesLeft := Checkmate - Abs(score)
		fmt.Printf("%2d %02d:%02d    %8d    %8d   %9.1f   %4dX   %v Checkmate\n",
			depth, int(finish)/60, int(finish)%60, game.nodes, game.qnodes,
			float64(game.nodes+game.qnodes)/finish, movesLeft/2,
			game.bestLine[0][0:Min(movesLeft, game.bestLength[0])])
	default:
		fmt.Printf("%2d %02d:%02d    %8d    %8d   %9.1f   %5.2f   %v\n",
			depth, int(finish)/60, int(finish)%60, game.nodes, game.qnodes,
			float64(game.nodes+game.qnodes)/finish, float32(score)/float32(valuePawn.endgame),
			game.bestLine[0][0:game.bestLength[0]])
	}
}

func (game *Game) saveBest(ply int, move Move) *Game {
	game.bestLine[ply][ply] = move
	game.bestLength[ply] = ply + 1

	if length := game.bestLength[ply+1]; length > 0 {
		copy(game.bestLine[ply][ply+1:length],
			game.bestLine[ply+1][ply+1:length])
		game.bestLength[ply] = length
	}
	return game
}

func (game *Game) saveGood(depth int, move Move) *Game {
	if ply := Ply(); move&(isCapture|isPromo) == 0 && move != game.killers[ply][0] {
		game.killers[ply][1] = game.killers[ply][0]
		game.killers[ply][0] = move
		game.goodMoves[move.piece()][move.to()] += depth * depth
	}
	return game
}

// Checks whether the move is among good moves captured so far and returns its
// history value.
func (game *Game) good(move Move) int {
	return game.goodMoves[move.piece()][move.to()]
}

func (game *Game) String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h\n")
	for row := 7; row >= 0; row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0; col <= 7; col++ {
			square := Square(row, col)
			buffer.WriteByte(' ')
			if piece := game.pieces[square]; piece != 0 {
				buffer.WriteString(piece.String())
			} else {
				buffer.WriteString("\u22C5")
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}
