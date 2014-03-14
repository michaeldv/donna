// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

func (p *Position) Perft(depth int) (total int64) {
        if depth == 0 {
                return 1
        }

        gen := p.StartMoveGen(depth)
        if p.isInCheck(p.color) {
                gen.GenerateEvasions()
        } else {
                gen.GenerateMoves() // TODO: GenerateNonEvasions()
        }
        for move := gen.NextMove(); move != 0; move = gen.NextMove() {
                if position := p.MakeMove(move); position != nil {
                        delta := position.Perft(depth - 1)
                        total += delta
                        position.TakeBack(move)
                }
        }
        return
}
