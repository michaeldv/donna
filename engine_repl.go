// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(
	`fmt`
	`io/ioutil`
	`regexp`
	`strconv`
	`strings`
	`time`
)

func (e *Engine) replBestMove(move Move) *Engine {
	fmt.Printf("\033[0;36mDonna's move: %s", move)
	if game.nodes == 0 {
		fmt.Printf(" (book)")
	}
	fmt.Println("\033[0m\n")

	return e
}

func (e *Engine) replPrincipal(depth, score, status int, duration int64) {
	fmt.Printf(`%2d %s %10d %10d %9d   `, depth, ms(duration), game.nodes, game.qnodes, nps(duration))
	switch status {
	case WhiteWon:
		fmt.Println(`1-0 White Checkmates`)
	case BlackWon:
		fmt.Println(`0-1 Black Checkmates`)
	case Stalemate:
		fmt.Println(`1/2 Stalemate`)
	case Repetition:
		fmt.Println(`1/2 Repetition`)
	case FiftyMoves:
		fmt.Println(`1/2 Fifty Moves`)
	case WhiteWinning, BlackWinning: // Show moves till checkmate.
		fmt.Printf("%4dX   %v Checkmate\n", (Checkmate - abs(score)) / 2, game.rootpv)
	default:
		fmt.Printf("%5.2f   %v\n", float32(score) / float32(onePawn), game.rootpv)
	}
}

func (e *Engine) Repl() *Engine {
	var game *Game
	var position *Position

	setup := func() {
		if game == nil || position == nil {
			game = NewGame()
			position = game.start()
			fmt.Printf("%s\n", position)
		}
	}

	think := func() {
		if move := game.Think(); move != 0 {
			position = position.makeMove(move)
			fmt.Printf("%s\n", position)
		}
	}

	benchmark := func(fileName string) {
		maxDepth, moveTime := e.options.maxDepth, e.options.moveTime
		e.options.maxDepth, e.options.moveTime = 0, 10000
		defer func() {
			e.options.maxDepth, e.options.moveTime = maxDepth, moveTime
		}()

		content, err := ioutil.ReadFile(fileName)
		if err == nil {
			total, solved := 0, 0
			lines := strings.Split(string(content), "\n")
			re := regexp.MustCompile(`[\+\?!]`)

			NextLine:
			for i, line := range lines {
				if len(line) > 0 && line[0] != '#' {
					total++
					game := NewGame(line)
					position := game.start()

					best := strings.Split(line, ` # `)[1]
					fmt.Printf("\033[0;36m%d) %s for %s\033[0m\n%s\n", i, best, C(position.color), position)
					move := game.Think()

					for _, theBest := range strings.Split(best, ` `) {
						theBest = re.ReplaceAllLiteralString(theBest, ``)
						if move == NewMoveFromString(position, theBest) {
							solved++
							fmt.Printf("\033[0;32m%d: solved (%d/%d %2.1f%%)\033[0m\n\n\n", total, solved, total - solved, float32(solved) * 100.0 / float32(total))
							continue NextLine
						}
					}
					fmt.Printf("\033[0;31m%d: not solved (%d/%d %2.1f%%)\033[0m\n\n\n", total, solved, total - solved, float32(solved) * 100.0 / float32(total))
				}
			}
		} else {
			fmt.Printf("Could not open benchmark file '%s'\n", fileName)
		}
	}

	perft := func(parameter string) {
		if parameter == `` {
			parameter = `5`
		}
		if depth, err := strconv.Atoi(parameter); err == nil {
			position := NewGame().start()
			start := time.Now()
			total := position.Perft(depth)
			finish := time.Since(start).Nanoseconds() / 1000000
			fmt.Printf("  Depth: %d\n", depth)
			fmt.Printf("  Nodes: %d\n", total)
			fmt.Printf("Elapsed: %s\n", ms(finish))
			fmt.Printf("Nodes/s: %dK\n", total / finish)
		}
	}

	fmt.Printf("Donna v%s Copyright (c) 2014 by Michael Dvorkin. All Rights Reserved.\nType ? for help.\n\n", Version)
	for command, parameter := ``, ``; ; command, parameter = ``, `` {
		fmt.Print(`donna> `)
		fmt.Scanln(&command, &parameter)

		switch command {
		case ``:
		case `bench`:
			benchmark(parameter)
		case `exit`, `quit`:
			return e
		case `go`:
			setup()
			think()
		case `help`, `?`:
			fmt.Println("The commands are:\n\n" +
				"   bench   Run benchmark tests\n" +
				"   exit    Exit the program\n" +
				"   go      Take side and make a move\n" +
				"   help    Display this help\n" +
				"   new     Start new game\n" +
				"   perft   Run perft test\n" +
				"   score   Show evaluation summary\n" +
				"   undo    Undo last move\n")
		case `new`:
			game, position = nil, nil
			setup()
		case `perft`:
			perft(parameter)
		case `score`:
			setup()
			_, metrics := position.EvaluateWithTrace()
			Summary(metrics)
		case `undo`:
			if position != nil {
				position = position.undoLastMove()
				fmt.Printf("%s\n", position)
			}
		default:
			setup()
			if move := NewMoveFromString(position, command); move != 0 {
				if advance := position.makeMove(move); advance != nil {
					position = advance
					think()
					continue
				}
			}
			// Invalid move (typo) or non-evasion on check.
			fmt.Printf("%s appears to be an invalid move.\n", command)
		}
	}
	return e
}
