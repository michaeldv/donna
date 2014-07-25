// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Options struct {
	infinite     bool   // (-) Search until the "stop" command.
	maxDepth     int    // Search X plies only.
	maxNodes     int    // (-) Search X nodes only.
	gameTime     int64  // Time for all remaining moves is X milliseconds.
	moveTime     int64  // Search exactly X milliseconds per move.
	moveTimeInc  int64  // Time increment after the move is X milliseconds.
}

func (game *Game) Set(args ...interface{}) *Game {
	for i := 0; i < len(args); i += 2 {
		key, value := args[i], args[i+1]
		switch key {
		case `cache`:
			switch value.(type) {
			default: // :-)
				game.cache = NewCache(value.(float64))
			case int:
				game.cache = NewCache(float64(value.(int)))
			}
			game.warmUpMaterialCache()
		case `depth`:
			game.options = Options{}
			game.options.maxDepth = value.(int)
		case `time`:
			game.options.infinite = false
			game.options.maxDepth = 0
			game.options.maxNodes = 0
			game.options.moveTime = 0
			game.options.gameTime = int64(value.(int))
		case `timeinc`:
			game.options.infinite = false
			game.options.maxDepth = 0
			game.options.maxNodes = 0
			game.options.moveTime = 0
			game.options.moveTimeInc = int64(value.(int))
		case `movetime`:
			game.options = Options{}
			game.options.moveTime = int64(value.(int))
		}
	}

	return game
}
