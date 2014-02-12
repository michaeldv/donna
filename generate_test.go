// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

func TestGenerate010(t *testing.T) {
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `Kg8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves().rank()

        // Moves should be sorted by relative strength.
        expect(t, gen.allMoves(), `[e6-e7 f5-f6 d2-d4 c4-c5 d2-d3 g4-g5 b3-b4 h3-h4 a2-a4 a2-a3]`)
}

// LVA/MVV capture ordering.
func TestGenerate110(t *testing.T) {
        game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5`)
        gen := game.Start(White).StartMoveGen(1).GenerateCaptures().rank()

        expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Qh5xd5 Kd4xd5]`)
}

func TestGenerate120(t *testing.T) {
        game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5`)
        gen := game.Start(White).StartMoveGen(1).GenerateCaptures().rank()

        expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5]`)
}

func TestGenerate130(t *testing.T) {
        game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6`)
        gen := game.Start(White).StartMoveGen(1).GenerateCaptures().rank()

        expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6]`)
}

func TestGenerate140(t *testing.T) {
        game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6,Nh3`)
        gen := game.Start(White).StartMoveGen(1).GenerateCaptures().rank()

        expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6 Nf4xh3 Qh5xh3]`)
}

func TestGenerate150(t *testing.T) {
        game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6,Nh3,e2`)
        gen := game.Start(White).StartMoveGen(1).GenerateCaptures().rank()

        expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6 Nf4xh3 Qh5xh3 Nf4xe2 Bc4xe2 Qh5xe2]`)
}
