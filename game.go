// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`fmt`
	`strings`
	`time`
)

type Options struct {
	msRemaining   int // (-) Remaining time for the rest of the game.
	msIncrement   int // (-) Time increment after the move.
	msToMakeMove  int // Time limit to make a move.
	maxDepth      int // Search depth limit.
	maxNodes      int // (-) Search nodes limit.
	msSoftStop    int // (-) Soft time limit stop.
	msHardStop    int // (-) Hard time limit stop.
}

type History [14][64]int
type Killers [MaxPly][2]Move
type Pv      [MaxPly][MaxPly]Move
type PvSize  [MaxPly]int

type Game struct {
	nodes    int 	  // Number of regular nodes searched.
	qnodes   int 	  // Number of quiescence nodes searched.
	token    uint8 	  // Expiration token for cache.
	initial  string   // Initial position (FEN or algebraic).
	cache    Cache 	  // Transposition table.
	history  History  // Good moves history.
	killers  Killers  // Killer moves.
	pv       Pv 	  // Principal variation.
	pvsize   PvSize   // Number of moves in principal variation.
	options  Options  // Game options.
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

func (game *Game) CacheSize(megaBytes float32) *Game {
	game.cache = NewCache(megaBytes)
	game.warmUpMaterialCache()

	return game
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

func (game *Game) Set(option string, value int) {
	switch option {
	case `depth`:
		game.options = Options{}
		game.options.maxDepth = value
	case `movetime`:
		game.options = Options{}
		game.options.msRemaining = value
		game.options.msToMakeMove = value
		game.options.msSoftStop = value
		game.options.msHardStop = value
	}
}

func (game *Game) Position() *Position {
	return &tree[node]
}

func (game *Game) Think(requestedDepth int) Move {
	position := game.Position()

	book := NewBook("./books/gm2001.bin") // From http://www.chess2u.com/t5834-gm-polyglot-book
	if move := book.pickMove(position); move != 0 {
		fmt.Printf("Book move: %s\n", move)
		return move
	}

	// Reset principal variation, killer moves and move history, and update
	// cache token to ignore existing cache entries.
	rootNode = node
	game.pv = Pv{}
	game.pvsize = PvSize{}
	game.killers = Killers{}
	game.history = History{}
	game.token++ // <-- Wraps around: ...254, 255, 0, 1...

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
			game.pv[0][0:Min(movesLeft, game.pvsize[0])])
	default:
		fmt.Printf("%2d %02d:%02d    %8d    %8d   %9.1f   %5.2f   %v\n",
			depth, int(finish)/60, int(finish)%60, game.nodes, game.qnodes,
			float64(game.nodes+game.qnodes)/finish, float32(score)/float32(valuePawn.endgame),
			game.pv[0][0:game.pvsize[0]])
	}
}

func (game *Game) saveBest(ply int, move Move) *Game {
	game.pv[ply][ply] = move
	game.pvsize[ply] = ply + 1

	if length := game.pvsize[ply+1]; length > 0 {
		copy(game.pv[ply][ply+1:length],
			game.pv[ply+1][ply+1:length])
		game.pvsize[ply] = length
	}
	return game
}

func (game *Game) saveGood(depth int, move Move) *Game {
	if ply := Ply(); move&(isCapture|isPromo) == 0 && move != game.killers[ply][0] {
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
