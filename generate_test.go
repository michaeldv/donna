// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`testing`
)

// func TestGenerate010(t *testing.T) {
// 	game := NewGame(`Ka1,a2,b3,c4,d2,e6,f5,g4,h3`, `Kc1`)
// 	gen := NewMoveGen(game.Start(White)).generateMoves().rank(Move(0))

// 	// TODO: moves should be sorted by good moves history.
// 	expect(t, gen.allMoves(), `[e6-e7 f5-f6 d2-d4 c4-c5 g4-g5 b3-b4 d2-d3 a2-a4 h3-h4 a2-a3 Ka1-b2 Ka1-b1]`)
// }

// LVA/MVV capture ordering.
func TestGenerate110(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5`)
	gen := NewMoveGen(game.Start(White)).generateCaptures().rank(Move(0))

	expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Qh5xd5 Kd4xd5]`)
}

func TestGenerate120(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5`)
	gen := NewMoveGen(game.Start(White)).generateCaptures().rank(Move(0))

	expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5]`)
}

func TestGenerate130(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6`)
	gen := NewMoveGen(game.Start(White)).generateCaptures().rank(Move(0))

	expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6]`)
}

func TestGenerate140(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6,Nh3`)
	gen := NewMoveGen(game.Start(White)).generateCaptures().rank(Move(0))

	expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6 Nf4xh3 Qh5xh3]`)
}

func TestGenerate150(t *testing.T) {
	game := NewGame(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6,Nh3,e2`)
	gen := NewMoveGen(game.Start(White)).generateCaptures().rank(Move(0))

	expect(t, gen.allMoves(), `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6 Nf4xh3 Qh5xh3 Nf4xe2 Bc4xe2 Qh5xe2]`)
}
