// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Player struct {
        game    *Game  // The game we're playing.
	Color   int    // 0: white, 1: black
        Can00   bool   // Can castle king's side?
        Can000  bool   // Can castle queen's side?
}

func NewPlayer(game *Game, color int) *Player {
        player := new(Player)

        player.game = game
        player.Color = color
        player.Can00 = true
        player.Can000 = true

        return player
}

func (p *Player) IsWhite() bool {
	return p.Color == WHITE
}

func (p *Player) IsBlack() bool {
	return p.Color == BLACK
}
