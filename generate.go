// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`sort`)

const (
        stepPrincipal = iota
        stepCaptures
        stepPromotions
        stepKillers
        stepRemaining
)

type MoveWithScore struct {
        move   Move
        score  int
}

type MoveGen struct {
        p     *Position
        game  *Game
        list  [256]MoveWithScore
        head  int
        tail  int
        step  int
        ply   int
}

var moveList [MaxPly]MoveGen

func (p *Position) StartMoveGen(ply int) (gen *MoveGen) {
        gen = &moveList[ply]
        gen.p = p
        gen.game = p.game
        gen.list = [256]MoveWithScore{}
        gen.head, gen.tail = 0, 0
        gen.ply = ply
        return
}

func (p *Position) UseMoveGen(ply int) (gen *MoveGen) {
        return moveList[ply].reset()
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
                        position.TakeBack(move)
                }
        }
        return gen.reset()
}

func (gen *MoveGen) rank() *MoveGen {
        if gen.size() < 2 {
                return gen
        }
        for i := gen.head; i < gen.tail; i++ {
                move := gen.list[i].move
                if move == gen.game.bestLine[0][gen.ply] {
                        gen.list[i].score = 0xFFFF
                } else if move == gen.game.killers[gen.ply][0] {
                        gen.list[i].score = 0xFFFE
                } else if move == gen.game.killers[gen.ply][1] {
                        gen.list[i].score = 0xFFFD
                } else if move & isCapture != 0 {
                        gen.list[i].score = move.value()
                } else {
                        endgame, midgame := move.score()
                        gen.list[i].score = gen.p.score(midgame, endgame)
                        gen.list[i].score += gen.game.goodMoves[move.piece()][move.to()]
                }
        }
        sort.Sort(byScore{ gen.list[gen.head : gen.tail] })
        return gen
}

func (gen *MoveGen) quickRank() *MoveGen {
        if gen.size() < 2 {
                return gen
        }
        for i := gen.head; i < gen.tail; i++ {
                if move := gen.list[i].move; move & isCapture != 0 {
                        gen.list[i].score = move.value()
                } else {
                        endgame, midgame := move.score()
                        gen.list[i].score = gen.p.score(midgame, endgame)
                }
        }
        sort.Sort(byScore{ gen.list[gen.head : gen.tail] })
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
        gen.head--; gen.tail--
        return gen;
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
func (her byScore) Len() int           { return len(her.list)}
func (her byScore) Swap(i, j int)      { her.list[i], her.list[j] = her.list[j], her.list[i] }
func (her byScore) Less(i, j int) bool { return her.list[i].score > her.list[j].score }
