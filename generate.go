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
        gen.list = [256]MoveWithScore{}
        gen.head, gen.tail = 0, 0
        gen.ply = ply
        return
}

func (p *Position) UseMoveGen(ply int) (gen *MoveGen) {
        return &moveList[ply]
}

func (gen *MoveGen) reset() *MoveGen {
        gen.head = 0
        return gen
}

func (gen *MoveGen) listSize() int {
        return gen.tail - gen.head
}

func (gen *MoveGen) theOnlyMove() bool {
        return gen.listSize() <= 1
}

func (gen *MoveGen) NextMove() (move Move) {
        if gen.head < gen.tail {
                move = gen.list[gen.head].move
                gen.head++
        }
        return
}

func (gen *MoveGen) quickRank() *MoveGen {
        if gen.listSize() < 2 {
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

func (gen *MoveGen) rank() *MoveGen {
        if gen.listSize() < 2 {
                return gen
        }
        for i := gen.head; i < gen.tail; i++ {
                move := gen.list[i].move
                if move == gen.p.game.bestLine[0][gen.ply] {
                        gen.list[i].score = 0xFFFF
                } else if move == gen.p.killers[0] {
                        gen.list[i].score = 0xFFFE
                } else if move == gen.p.killers[1] {
                        gen.list[i].score = 0xFFFD
                } else if move & isCapture != 0 {
                        gen.list[i].score = move.value()
                } else {
                        endgame, midgame := move.score()
                        gen.list[i].score = gen.p.score(midgame, endgame)
                        gen.list[i].score += gen.p.game.goodMoves[move.piece()][move.to()]
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

// Return a list of generated moves by continuously calling the next move
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
