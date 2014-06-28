// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`testing`
)

// Piece captures.
func TestGenCaptures000(t *testing.T) {
	game := NewGame(`Ka1,Qd1,Rh1,Bb3,Ne5`, `Ka8,Qd8,Rh8,Be6,Ng6`)
	white := NewGen(game.Start(White), 0).generateCaptures()
	expect(t, white.allMoves(), `[Qd1xd8 Rh1xh8 Bb3xe6 Ne5xg6]`)

	black := NewGen(game.Start(Black), 0).generateCaptures()
	expect(t, black.allMoves(), `[Be6xb3 Ng6xe5 Qd8xd1 Rh8xh1]`)
}

// Pawn captures on rows 2-6 (no promotions).
func TestGenCaptures010(t *testing.T) {
	game := NewGame(`Ka1,b3,c4,d5`, `Ka8,a4,b5,c6,e6`)
	white := NewGen(game.Start(White), 0).generateCaptures()
	expect(t, white.allMoves(), `[b3xa4 c4xb5 d5xc6 d5xe6]`)

	black := NewGen(game.Start(Black), 0).generateCaptures()
	expect(t, black.allMoves(), `[a4xb3 b5xc4 c6xd5 e6xd5]`)
}

// Pawn captures with promotion, rows 1-7.
func TestGenCaptures020(t *testing.T) {
	game := NewGame(`Ka1,Bh1,Ng1,a2,b7,e7`, `Kb8,Rd8,Be8,Rf8,h2`)
	white := NewGen(game.Start(White), 0).generateCaptures()
	expect(t, white.allMoves(), `[e7xd8Q e7xd8R e7xd8B e7xd8N e7xf8Q e7xf8R e7xf8B e7xf8N]`)

	black := NewGen(game.Start(Black), 0).generateCaptures()
	expect(t, black.allMoves(), `[h2xg1Q h2xg1R h2xg1B h2xg1N Kb8xb7]`)
}

// Pawn promotions without capture, rows 1-7.
func TestGenCaptures030(t *testing.T) {
	game := NewGame(`Ka1,a2,e7`, `Ka8,Rd8,a7,h2`)
	white := NewGen(game.Start(White), 0).generateCaptures()
	expect(t, white.allMoves(), `[e7xd8Q e7xd8R e7xd8B e7xd8N e7-e8Q e7-e8R e7-e8B e7-e8N]`)

	black := NewGen(game.Start(Black), 0).generateCaptures()
	expect(t, black.allMoves(), `[h2-h1Q h2-h1R h2-h1B h2-h1N]`)
}
