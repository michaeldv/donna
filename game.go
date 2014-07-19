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
type Pv      []Move
type PvTable [MaxPly]Pv

type Game struct {
	nodes    int 	  // Number of regular nodes searched.
	qnodes   int 	  // Number of quiescence nodes searched.
	token    uint8 	  // Expiration token for cache.
	initial  string   // Initial position (FEN or algebraic).
	cache    Cache 	  // Transposition table.
	history  History  // Good moves history.
	killers  Killers  // Killer moves.
	rootpv   Pv 	  // Principal variation for root moves.
	pv       PvTable  // Principal variations for each ply.
	options  Options  // Game options, might be set by REPL or UCI.
	clock    Clock    // Time controls.
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
	pawnCache = [8192]PawnEntry{}
	materialCache = [8192]MaterialEntry{}

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

// The color parameter is optional.
func (game *Game) Start(args ...int) *Position {
	tree, node, rootNode = [1024]Position{}, 0, 0

	// Was the game started with FEN or algebraic notation?
	sides := strings.Split(game.initial, `:`)
	if len(sides) == 2 {
		return NewPosition(game, sides[White], sides[Black], args[0])
	}
	return NewPositionFromFEN(game, game.initial)
}

func (game *Game) Position() *Position {
	return &tree[node]
}

func (game *Game) Think() Move {
	position := game.Position()

	book := NewBook("./books/gm2001.bin") // From http://www.chess2u.com/t5834-gm-polyglot-book
	if move := book.pickMove(position); move != 0 {
		fmt.Printf("Book move: %s\n", move)
		return move
	}

	// Reset principal variation, killer moves and move history, and update
	// cache token to ignore existing cache entries.
	rootNode = node
	game.killers = Killers{}
	game.history = History{}
	game.token++ // <-- Wraps around: ...254, 255, 0, 1...

	for ply := 0;  ply < MaxPly; ply++ {
		game.pv[ply] = game.pv[ply][:0]
	}

	move, score, status := Move(0), 0, InProgress
	alpha, beta := -Checkmate, Checkmate

	done := func(depth int) bool {
		return game.clock.stopSearch || (game.options.maxDepth > 0 && depth > game.options.maxDepth)
	}

	fmt.Println(`Depth/Time     Nodes      QNodes     Nodes/s   Score   Best`)

	game.startClock(); defer game.stopClock();
	for depth := 1; !done(depth); depth++ {
		game.nodes, game.qnodes = 0, 0

		// Save previous best score in case search gets interrupted.
		previousBest := score

		// At low depths do the search with full alpha/beta spread.
		// Aspiration window searches kick in at depth 5 and up.
		start := time.Now()
		if depth < 5 {
			move, score = position.search(alpha, beta, depth)
		} else {
			aspiration := valuePawn.midgame / 3
			alpha = Max(score - aspiration, -Checkmate)
			beta = Min(score + aspiration, Checkmate)

			// Do the search with smaller alpha/beta spread based on
			// previous iteration score, and re-search with the bigger
			// window as necessary.
			for {
				Log("\tscore -> %d, searchRoot(%d, %d, %d)\n", score, alpha, beta, depth)
				previousBest = score
				move, score = position.search(alpha, beta, depth)

				if game.clock.stopSearch {
					break
				}

				Log("\tscore => %d, pv => %v\n", score, game.pv[0])
				if score <= alpha {
					Log("\tscore %d <= alpha %d, new alpha %d\n", score, alpha, score - aspiration)
					alpha = Max(score - aspiration, -Checkmate)
				} else if score >= beta {
					Log("\tscore %d >= beta %d, new beta %d\n", score, beta, score + aspiration)
					beta = Min(score + aspiration, Checkmate)
				} else {
					break;
				}
				aspiration *= 2
			}
			// TBD: position.cache(move, score, 0, 0)
		}
		finish := time.Since(start).Seconds()

		if game.clock.stopSearch {
			Log("\ttimed out score %d previousBest %d move %s\n", score, previousBest, move)
			score = previousBest
		}

		status = position.status(move, score)
		game.printBestLine(depth, score, status, finish)

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
			float64(game.nodes+game.qnodes)/finish, movesLeft/2, game.pv[0])
	default:
		if game.Position().color == Black {
			score = -score
		}
		fmt.Printf("%2d %02d:%02d    %8d    %8d   %9.1f   %5.2f   %v\n",
			depth, int(finish)/60, int(finish)%60, game.nodes, game.qnodes,
			float64(game.nodes+game.qnodes)/finish, float32(score)/float32(valuePawn.endgame), game.pv[0])
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
	if ply := Ply(); move.isQuiet() && move != game.killers[ply][0] {
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
	return game.Position().String()
}
