// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`sort`
)

type MoveWithScore struct {
	move  Move
	score int
}

type MoveGen struct {
	p      *Position
	list   [256]MoveWithScore
	head   int
	tail   int
	ply    int
}

// Pre-allocate move generator array (one entry per ply) to avoid garbage
// collection overhead.
var moveList [MaxPly]MoveGen

// Returns "new" move generator for the given ply. Since move generator array
// has been pre-allocated already we simply return a pointer to the existing
// array element re-initializing all its data.
func NewGen(p *Position, ply int) (gen *MoveGen) {
	gen = &moveList[ply]
	gen.p = p
	gen.list = [256]MoveWithScore{}
	gen.head, gen.tail = 0, 0
	gen.ply = ply
	return gen
}

// Returns new move generator for the initial step of iterative deepening
// (depth == 1) and existing one for subsequent iterations (depth > 1).
//
// This is used in iterative deepening search when all the moves are being
// generated at depth one, and reused later as the search deepens.
func NewRootGen(p *Position, depth int) (gen *MoveGen) {
	if depth > 1 {
		return moveList[0].reset().rank(p.cachedMove())
	}

	// 1) generate all moves or check evasions; 2) return if we've got the
	// only move; 3) and get rid of invalid moves so that we don't do it on
	// each iteration; 4) return sorted list.
	gen = NewGen(p, 0)
	if p.isInCheck(p.color) {
		gen.generateEvasions()
		if gen.onlyMove() {
			return gen
		}
		return gen.validOnly(p).quickRank()
	}

	gen.generateMoves()
	if gen.onlyMove() {
		return gen
	}
	return gen.validOnly(p).rank(p.cachedMove())
}

func (gen *MoveGen) reset() *MoveGen {
	gen.head = 0
	return gen
}

func (gen *MoveGen) size() int {
	return gen.tail - gen.head
}

func (gen *MoveGen) onlyMove() bool {
	return gen.size() == 1
}

func (gen *MoveGen) NextMove() (move Move) {
	if gen.head < gen.tail {
		move = gen.list[gen.head].move
		gen.head++
	}
	return
}

// Removes invalid moves from the generated list. We use in iterative deepening
// to avoid stumbling upon invalid moves on each iteration.
func (gen *MoveGen) validOnly(p *Position) *MoveGen {
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if position := p.MakeMove(move); position == nil {
			gen.remove()
		} else {
			position.UndoLastMove()
		}
	}
	return gen.reset()
}

// Probes a list of generated moves and returns true if it contains at least
// one valid move.
func (gen *MoveGen) anyValid(p *Position) bool {
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		if position := p.MakeMove(move); position != nil {
			position.UndoLastMove()
			return true
		}
	}
	return false
}

func (gen *MoveGen) rank(bestMove Move) *MoveGen {
	if gen.size() < 2 {
		return gen
	}

	for i, game := gen.head, gen.p.game; i < gen.tail; i++ {
		move := gen.list[i].move
		if move == bestMove {
			gen.list[i].score = 0xFFFF
		} else if move & isCapture != 0 {
			gen.list[i].score = 8192 + move.value()
		} else if move == game.killers[gen.ply][0] {
			gen.list[i].score = 4096
		} else if move == game.killers[gen.ply][1] {
			gen.list[i].score = 2048
		} else {
			gen.list[i].score = game.good(move)
		}
	}

	sort.Sort(byScore{gen.list[gen.head:gen.tail]})
	return gen
}

func (gen *MoveGen) quickRank() *MoveGen {
	if gen.size() < 2 {
		return gen
	}

	for i, game := gen.head, gen.p.game; i < gen.tail; i++ {
		if move := gen.list[i].move; move & isCapture != 0 {
			gen.list[i].score = 8192 + move.value()
		} else {
			gen.list[i].score = game.good(move)
		}
	}

	sort.Sort(byScore{gen.list[gen.head:gen.tail]})
	return gen
}

func (gen *MoveGen) add(move Move) *MoveGen {
	gen.list[gen.tail].move = move
	gen.tail++
	return gen
}

// Removes current move from the list by copying over the ramaining moves. Head and
// tail pointers get decremented so that calling NexMove() works as expected.
func (gen *MoveGen) remove() *MoveGen {
	copy(gen.list[gen.head-1:], gen.list[gen.head:])
	gen.head--
	gen.tail--
	return gen
}

// Returns an array of generated moves by continuously appending the NextMove()
// until the list is empty.
func (gen *MoveGen) allMoves() (moves []Move) {
	for move := gen.NextMove(); move != 0; move = gen.NextMove() {
		moves = append(moves, move)
	}
	return
}

// Sorting moves by their relative score based on piece/square for regular moves
// or least valuaeable attacker/most valueable victim for captures.
type byScore struct {
	list []MoveWithScore
}

func (her byScore) Len() int           { return len(her.list) }
func (her byScore) Swap(i, j int)      { her.list[i], her.list[j] = her.list[j], her.list[i] }
func (her byScore) Less(i, j int) bool { return her.list[i].score > her.list[j].score }
