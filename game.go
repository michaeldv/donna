// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`fmt`
	`strings`
	`time`
)

type History [14][64]int
type Killers [MaxPly][2]Move
type RootPv  []Move
type Pv      [MaxPly]RootPv

type Game struct {
	nodes       int 	// Number of regular nodes searched.
	qnodes      int 	// Number of quiescence nodes searched.
	token       uint8 	// Cache's expiration token.
	deepening   bool 	// True when searching first root move.
	improving   bool 	// True when root search score is not falling.
	volatility  float32 	// Root search stability count.
	initial     string   	// Initial position (FEN or algebraic).
	history     History  	// Good moves history.
	killers     Killers  	// Killer moves.
	rootpv      RootPv 	// Principal variation for root moves.
	pv          Pv 		// Principal variations for each ply.
	cache       Cache 	// Transposition table.
	pawnCache   PawnCache 	// Cache of pawn structures.
}

// Use single statically allocated variable.
var game Game

// We have two ways to initialize the game: 1) pass FEN string, and 2) specify
// white and black pieces using regular chess notation.
//
// In latter case we need to tell who gets to move first when starting the game.
// The second option is a bit less pricise (ex. no en-passant square) but it is
// much more useful when writing tests from memory.
func NewGame(args ...string) *Game {
	game = Game{}
	game.cache = NewCache(engine.cacheSize)
	game.pawnCache = PawnCache{}

	game.rootpv = make([]Move, 0, MaxPly)
	for ply := 0;  ply < MaxPly; ply++ {
		game.pv[ply] = make([]Move, 0, MaxPly)
	}

	switch len(args) {
	case 0: // Initial position.
		game.initial = `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`
	case 1: // Genuine FEN.
		game.initial = args[0]
	case 2: // Standard algebraic notation (white and black).
		game.initial = args[0] + `:` + args[1]
	}

	return &game
}

func (game *Game) start(args ...int) *Position {
	engine.clock.halt = false
	tree, node, rootNode = [1024]Position{}, 0, 0

	// Was the game started with FEN or algebraic notation?
	sides := strings.Split(game.initial, `:`)
	if len(sides) == 2 {
		return NewPosition(game, sides[White], sides[Black], args[0])
	}
	return NewPositionFromFEN(game, game.initial)
}

func (game *Game) position() *Position {
	return &tree[node]
}

// Resets principal variation as well as killer moves and move history. Cache
// entries get expired by incrementing cache token. Root node gets set to the
// current tree node to match the position.
func (game *Game) getReady() *Game {
	game.rootpv = game.rootpv[:0]
	for ply := 0;  ply < MaxPly; ply++ {
		game.pv[ply] = game.pv[ply][:0]
	}

	game.killers = Killers{}
	game.history = History{}
	game.deepening = false
	game.improving = true
	game.volatility = 0.0
	game.token++ // <-- Wraps around: ...254, 255, 0, 1...

	rootNode = node
	return game
}

func (game *Game) Think() Move {
	start := time.Now()
	position := game.position()
	game.nodes, game.qnodes = 0, 0

	if book, err := NewBook(); err == nil {
		if move := book.pickMove(position); move != 0 {
			game.printBestMove(move, time.Since(start).Seconds())
			return move
		}
	}

	game.getReady()
	score, move, status, alpha, beta := 0, Move(0), InProgress, -Checkmate, Checkmate

	if engine.uci {
		engine.debug(position.String())
	} else {
		fmt.Println(`Depth/Time     Nodes      QNodes     Nodes/s   Score   Best`)
	}

	if !engine.fixedDepth() {
		engine.startClock(); defer engine.stopClock();
	}

	for depth := 1; status == InProgress && game.keepThinking(depth, move); depth++ {
		// Save previous best score in case search gets interrupted.
		bestScore := score

		// Assume volatility decreases with each new iteration.
		game.volatility /= 2.0

		// At low depths do the search with full alpha/beta spread.
		// Aspiration window searches kick in at depth 5 and up.
		if depth < 5 {
			score = position.search(alpha, beta, depth)
			if score > alpha {
				bestScore = score
				game.rootpv = append(game.rootpv[:0], game.pv[0]...)
			}
		} else {
			aspiration := onePawn / 3
			alpha = max(score - aspiration, -Checkmate)
			beta = min(score + aspiration, Checkmate)

			// Do the search with smaller alpha/beta spread based on
			// previous iteration score, and re-search with the bigger
			// window as necessary.
			for {
				//Log("\tscore -> %d, searchRoot(%d, %d, %d)\n", score, alpha, beta, depth)
				score = position.search(alpha, beta, depth)
				if score > alpha {
					bestScore = score
					game.rootpv = append(game.rootpv[:0], game.pv[0]...)
				}

				if !engine.fixedDepth() && engine.clock.halt {
					break
				}

				if score <= alpha {
					game.improving = false
					alpha = max(score - aspiration, -Checkmate)
				} else if score >= beta {
					beta = min(score + aspiration, Checkmate)
				} else {
					break;
				}

				aspiration *= 2
			}
			// TBD: position.cache(game.rootpv[0], score, 0, 0)
		}
		if engine.clock.halt {
			//Log("\ttimed out pv => %v\n\ttimed out rv => %v\n", game.pv[0], game.rootpv)
			score = bestScore
		}

		move = game.rootpv[0]
		status = position.status(move, score)
		game.printPrincipal(depth, score, status, time.Since(start).Seconds())
	}

	game.printBestMove(move, time.Since(start).Seconds())

	return move
}

func (game *Game) keepThinking(depth int, move Move) bool {
	if depth == 1 {
		return true
	}

	if engine.fixedDepth() {
		return depth <= engine.options.maxDepth
	} else if engine.clock.halt {
		engine.debug("# Depth %02d Early out with %s\n", depth, move)
		return false
	}

	// Stop deepening if it's the only move.
	gen := NewRootGen(nil, depth)
	if gen.onlyMove() {
		engine.debug("# Depth %02d Only move %s\n", depth, move)
		return false
	}

	// Stop if the move seems to be obvious and we've searched deep enough.
	if gen.obvious != Move(0) {
		if move == gen.obvious {
			engine.debug("# Depth %02d Kinda obvious %s\n", depth, move)
			return depth < 9
		} else { // The move is no longer obvious.
			engine.debug("# Depth %02d No longer obvious %s\n", depth, move)
			gen.obvious = Move(0)
		}
	}

	// Stop if the time left is not enough to gets through the next iteration.
	if engine.varyingTime() {
		elapsed := engine.elapsed(time.Now())
		remaining := engine.factor(depth, game.volatility).remaining()

		engine.debug("# Depth %02d Volatility %.2f Elapsed %s Remaining %s\n", depth, game.volatility, ms(elapsed), ms(remaining))
		if elapsed > engine.factor(depth, game.volatility).remaining() {
			engine.debug("# Depth %02d Bailing out with %s\n", depth, move)
			return false
		}
	}

	return true
}

func (game *Game) printBestMove(move Move, duration float64) {
	if engine.uci {
		engine.uciBestMove(move, duration)
	} else {
		engine.replBestMove(move)
	}
}

// Prints principal variation. Note that in REPL advantage white is always +score
// and advantage black is -score whereas in UCI +score is advantage current side
// and -score is advantage opponent.
func (game *Game) printPrincipal(depth, score, status int, duration float64) {
	if engine.uci {
		engine.uciPrincipal(depth, score, duration)
	} else {
		if game.position().color == Black {
			score = -score
		}
		engine.replPrincipal(depth, score, status, duration)
	}
}

func (game *Game) saveBest(ply int, move Move) *Game {
	game.pv[ply] = append(game.pv[ply][0:ply], move)

	next := ply + 1
	if length := len(game.pv[next]); length > 0 {
		game.pv[ply] = append(game.pv[ply], game.pv[next][next : length]...)
	}

	return game
}

func (game *Game) saveGood(depth int, move Move) *Game {
	if ply := ply(); move.isQuiet() && move != game.killers[ply][0] {
		game.killers[ply][1] = game.killers[ply][0]
		game.killers[ply][0] = move
		game.history[move.piece()][move.to()] += depth * depth
	}

	return game
}

// Checks whether the move is among good moves captured so far and returns its
// history value.
func (game *Game) good(move Move) int {
	return game.history[move.piece()][move.to()]
}

func (game *Game) String() string {
	return game.position().String()
}
