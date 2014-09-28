// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`bufio`
	`fmt`
	`io`
	`os`
	`strconv`
	`strings`
)

func (e *Engine) uciScore(depth, score, alpha, beta int) *Engine {
	str := fmt.Sprintf("info depth %d score", depth)

	if Abs(score) < Checkmate-MaxPly {
		str += fmt.Sprintf(" cp %d", score*100/onePawn)
	} else {
		mate := -Checkmate - score
		if score > 0 {
			mate = Checkmate - score + 1
		}
		str += fmt.Sprintf(" mate %d", mate/2)
	}
	if score <= alpha {
		str += " upperbound"
	} else if score >= beta {
		str += " lowerbound"
	}

	return engine.reply(str + "\n")
}

func (e *Engine) uciMove(move Move, moveno, depth, score int) *Engine {
	return engine.reply("info depth %d currmove %s currmovenumber %d # %d\n", depth, move.notation(), moveno, score)
}

func (e *Engine) uciBestMove(move Move, duration float64) *Engine {
	return engine.reply("info nodes %d time %d\nbestmove %s\n", game.nodes+game.qnodes, int64(duration*1000), move.notation())
}

func (e *Engine) uciPrincipal(depth, score int, duration float64) *Engine {
	str := fmt.Sprintf("info depth %d score", depth)

	if Abs(score) < Checkmate-MaxPly {
		str += fmt.Sprintf(" cp %d", score*100/onePawn)
	} else {
		mate := -Checkmate - score
		if score > 0 {
			mate = Checkmate - score + 1
		}
		str += fmt.Sprintf(" mate %d", mate/2)
	}
	nodes := game.nodes + game.qnodes
	str += fmt.Sprintf(" nodes %d nps %d time %d pv", nodes, int(float64(nodes)/duration), int64(duration*1000))

	for i := 0; i < len(game.rootpv); i++ {
		str += " " + game.rootpv[i].notation()
	}

	return engine.reply(str + "\n")
}

// Brain-damaged universal chess interface (UCI) protocol as described at
// http://wbec-ridderkerk.nl/html/UCIProtocol.html
func (e *Engine) Uci() *Engine {
	var game *Game
	var position *Position

	e.uci = true

	// "uci" command handler.
	doUci := func(args []string) {
		e.reply("Donna v%s Copyright (c) 2014 by Michael Dvorkin. All Rights Reserved.\n", Version)
		e.reply("id name Donna %s\n", Version)
		e.reply("id author Michael Dvorkin\n")
		e.reply("option name Hash type spin default %d min 1 max 1024\n", 64)
		e.reply("option name Mobility type spin default %d min 0 max 100\n", weights[0].midgame)
		e.reply("option name PawnStructure type spin default %d min 0 max 100\n", weights[1].midgame)
		e.reply("option name PassedPawns type spin default %d min 0 max 100\n", weights[2].midgame)
		e.reply("option name KingSafety type spin default %d min 0 max 100\n", weights[3].midgame)
		e.reply("option name EnemyKingSafety type spin default %d min 0 max 100\n", weights[4].midgame)
		e.reply("uciok\n")
	}

	// "ucinewgame" command handler.
	doUciNewGame := func(args []string) {
		game, position = nil, nil
	}

	// "isready" command handler.
	doIsReady := func(args []string) {
		e.reply("readyok\n")
	}

	// "position [startpos | fen ] [ moves ... ]" command handler.
	doPosition := func(args []string) {
		// Make sure we've started the game since "ucinewgame" is optional.
		if game == nil || position == nil {
			game = NewGame()
			position = game.Start()
		}

		switch args[0] {
		case `startpos`:
			args = args[1:]
			position = game.Start()
		case `fen`:
			fen := []string{}
			for _, token := range args[1:] {
				args = args[1:] // Shift the token.
				if token == `moves` {
					break
				}
				fen = append(fen, token)
			}
			position = NewPositionFromFEN(game, strings.Join(fen, ` `))
		default:
			return
		}

		if position != nil && len(args) > 0 && args[0] == `moves` {
			for _, move := range args[1:] {
				args = args[1:] // Shift the move.
				position = position.MakeMove(NewMoveFromNotation(position, move))
			}
		}
	}

	// "go [[wtime winc | btime binc ] movestogo] | depth | nodes | movetime"
	doGo := func(args []string) {
		think := true
		options := e.options

		for i, token := range args {
			// Boolen "infinite" and "ponder" commands have no arguments.
			if token == `infinite` {
				options = Options{infinite: true}
			} else if token == `ponder` {
				options = Options{ponder: true}
			} else if token == `test` { // <-- Custom token for use in tests.
				think = false
			} else if len(args) > i+1 {
				switch token {
				case `depth`:
					if n, err := strconv.Atoi(args[i+1]); err == nil {
						options = Options{ maxDepth: n }
					}
				case `nodes`:
					if n, err := strconv.Atoi(args[i+1]); err == nil {
						options = Options{ maxNodes: n }
					}
				case `movetime`:
					if n, err := strconv.Atoi(args[i+1]); err == nil {
						options = Options{ moveTime: int64(n) }
					}
				case `wtime`:
					if position.color == White {
						if n, err := strconv.Atoi(args[i+1]); err == nil {
							options.timeLeft = int64(n)
						}
					}
				case `btime`:
					if position.color == Black {
						if n, err := strconv.Atoi(args[i+1]); err == nil {
							options.timeLeft = int64(n)
						}
					}
				case `winc`:
					if position.color == White {
						if n, err := strconv.Atoi(args[i+1]); err == nil {
							options.timeInc = int64(n)
						}
					}
				case `binc`:
					if position.color == Black {
						if n, err := strconv.Atoi(args[i+1]); err == nil {
							options.timeInc = int64(n)
						}
					}
				case `movestogo`:
					if n, err := strconv.Atoi(args[i+1]); err == nil {
						options.movesToGo = int64(n)
					}
				}
			}
		}
		if options.timeLeft != 0 || options.timeInc != 0 || options.movesToGo != 0 {
			e.varyingLimits(options)
		} else {
			e.fixedLimit(options)
		}

		// Start "thinking" and come up with best move unless when running
		// tests where we verify argument parsing only.
		if think {
			game.Think()
		}
	}

	// Stop calculating as soon as possible.
	doStop := func(args []string) {
		e.clock.halt = true
	}

	var commands = map[string]func([]string){
		`isready`:    doIsReady,
		`uci`:        doUci,
		`ucinewgame`: doUciNewGame,
		`position`:   doPosition,
		`go`:         doGo,
		`stop`:       doStop,
	}

	bio := bufio.NewReader(os.Stdin)
	for {
		command, err := bio.ReadString('\n')
		if err != io.EOF && len(command) > 0 {
			e.debug("> " + command)
			args := strings.Split(strings.Trim(command, " \t\r\n"), ` `)
			if args[0] == `quit` {
				break
			}
			if handler, ok := commands[args[0]]; ok {
				handler(args[1:])
			}
		}
	}
	return e
}
