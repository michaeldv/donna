// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

type MoveWithScore struct {
	move  Move
	score int
}

type MoveGen struct {
	p	*Position
	list	[128]MoveWithScore
	ply	int
	head	int
	tail	int
	pins	Bitmask
}

// Pre-allocate move generator array (one entry per ply) to avoid garbage
// collection overhead. Last entry serves for utility move generation, ex. when
// converting string notations or determining a stalemate.
var moveList [MaxPly+1]MoveGen

// Returns "new" move generator for the given ply. Since move generator array
// has been pre-allocated already we simply return a pointer to the existing
// array element re-initializing all its data.
func NewGen(p *Position, ply int) (gen *MoveGen) {
	gen = &moveList[ply]
	gen.p = p
	gen.list = [128]MoveWithScore{}
	gen.ply = ply
	gen.head, gen.tail = 0, 0
	gen.pins = p.pins(p.king[p.color])

	return gen
}

// Convenience method to return move generator for the current ply.
func NewMoveGen(p *Position) *MoveGen {
	return NewGen(p, ply())
}

// Returns new move generator for the initial step of iterative deepening
// (depth == 1) and existing one for subsequent iterations (depth > 1).
func NewRootGen(p *Position, depth int) *MoveGen {
	if depth == 1 {
		return NewGen(p, 0) // Zero ply.
	}

	return &moveList[0]
}

func (gen *MoveGen) reset() *MoveGen {
	gen.head = 0

	return gen
}

func (gen *MoveGen) size() int {
	return gen.tail
}

func (gen *MoveGen) onlyMove() bool {
	return gen.tail == 1
}

func (gen *MoveGen) nextMove() (move Move) {
	move = gen.list[gen.head].move
	gen.head++

	return move
}

// Removes invalid moves from the generated list. We use in iterative deepening
// to avoid filtering out invalid moves on each iteration.
func (gen *MoveGen) validOnly() *MoveGen {
	for move := gen.nextMove(); move.some(); move = gen.nextMove() {
		if !move.valid(gen.p, gen.pins) {
			gen.remove()
		}
	}

	return gen.reset()
}

// Probes a list of generated moves and returns true if it contains at least
// one valid move.
func (gen *MoveGen) anyValid() bool {
	for move := gen.nextMove(); move.some(); move = gen.nextMove() {
		if move.valid(gen.p, gen.pins) {
			return true
		}
	}

	return false
}

// Probes valid-only list of generated moves and returns true if the given move
// is one of them.
func (gen *MoveGen) amongValid(someMove Move) bool {
	for move := gen.nextMove(); move.some(); move = gen.nextMove() {
		if someMove == move {
			return true
		}
	}

	return false
}

// Assigns given score to the last move returned by the gen.nextMove().
func (gen *MoveGen) scoreMove(depth, score int) *MoveGen {
	current := &gen.list[gen.head - 1]

	if depth == 1 || current.score == score + 1 {
		current.score = score
	} else if score != -depth || (score == -depth && current.score != score) {
		current.score += score // Fix up aspiration search drop.
	}

	return gen
}

// Shell sort that is somewhat faster that standard Go sort. It also seems
// to outperform:
// 	loop {
// 		gen.shuffleRandomly()
// 		if gen.isSorted() {
// 			break
// 		}
// 	}
func (gen *MoveGen) sort() *MoveGen {
	total := gen.tail - gen.head
	count := total
	pocket := MoveWithScore{}
	ever := true

	for (ever) {
		count = (count + 1) / 2
		ever = count > 1
		for i := 0; i < total - count; i++ {
			if this := gen.list[i + count]; this.score > gen.list[i].score {
				pocket = this
				gen.list[i + count] = gen.list[i]
				gen.list[i] = pocket
				ever = true
			}
		}
	}

	return gen
}

func (gen *MoveGen) rank(bestMove Move) *MoveGen {
	if gen.size() < 2 {
		return gen
	}

	for i := gen.head; i < gen.tail; i++ {
		move := gen.list[i].move
		if move == bestMove {
			gen.list[i].score = 0xFFFF
		} else if !move.isQuiet() || move.isEnpassant() {
			gen.list[i].score = 8192 + move.value()
		} else if move == game.killers[gen.ply][0] {
			gen.list[i].score = 4096
		} else if move == game.killers[gen.ply][1] {
			gen.list[i].score = 2048
		} else {
			gen.list[i].score = game.good(move)
		}
	}

	return gen.sort()
}

func (gen *MoveGen) quickRank() *MoveGen {
	if gen.size() < 2 {
		return gen
	}

	for i := gen.head; i < gen.tail; i++ {
		if move := gen.list[i].move; !move.isQuiet() || move.isEnpassant() {
			gen.list[i].score = 8192 + move.value()
		} else {
			gen.list[i].score = game.good(move)
		}
	}

	return gen.sort()
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

// Returns an array of generated moves by continuously appending the nextMove()
// until the list is empty.
func (gen *MoveGen) allMoves() (moves []Move) {
	for move := gen.nextMove(); move.some(); move = gen.nextMove() {
		moves = append(moves, move)
	}
	gen.reset()

	return moves
}
