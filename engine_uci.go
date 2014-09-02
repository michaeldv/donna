// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(
	`bufio`
	`io`
	`os`
	`strconv`
	`strings`
)

// Brain-damaged universal chess interface (UCI) protocol as described at
// http://wbec-ridderkerk.nl/html/UCIProtocol.html
func (e *Engine) Uci() *Engine {
	var game *Game
	var position *Position

	e.uci = true

	// "uci" command handler.
	doUci := func(args []string) {
		e.reply("Donna v1.0.0 Copyright (c) 2014 by Michael Dvorkin. All Rights Reserved.\n")
		e.reply("id name Donna v1.0.0\n")
		e.reply("id author Michael Dvorkin\n")
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
		default: return
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
		options := e.options

		for i, token := range args {
			// Boolen "infinite" and "ponder" commands have no arguments.
			if token == `infinite` {
				options = Options{ infinite: true }
			} else if token == `ponder` {
				options = Options{ ponder: true }
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
						options.movesToGo = n
					}
				}
			}
		}
		if options.timeLeft != 0 || options.timeInc != 0 || options.movesToGo != 0 {
			e.varyingLimits(options)
		} else {
			e.fixedLimit(options)
		}
		game.Think()
	}

	// Stop calculating as soon as possible.
	doStop := func(args []string) {
		e.clock.halt = true
	}

	var commands = map[string]func([]string){
		`isready`: doIsReady,      
		`uci`: doUci,
		`ucinewgame`: doUciNewGame,
		`position`: doPosition,
		`go`: doGo,
		`stop`: doStop,
	}

 	bio := bufio.NewReader(os.Stdin)
	for {
		command, err := bio.ReadString('\n')
		if err != io.EOF && len(command) > 0 {
			e.debug("> " + command)
			args := strings.Split(command[:len(command)-1], ` `)
			if args[0] == `quit` {
				break
			} else if args[0] == `debug` {
				engine.log = true
			}
			if handler, ok := commands[args[0]]; ok {
				handler(args[1:])
			}
		}
	}
	return e
}
