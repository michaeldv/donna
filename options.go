// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

type Options struct {
	infinite    bool  // (-) Search until the "stop" command.
	maxDepth    int   // Search X plies only.
	maxNodes    int   // (-) Search X nodes only.
	msMoveTime  int   // Search exactly X milliseconds per move.
	msGameTime  int   // Time for all remaining moves is X milliseconds.
	msTimeInc   int   // Time increment after the move is X milliseconds.
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
		case `movetime`:
			game.options = Options{}
			game.options.msMoveTime = value.(int)
		}
	}

	return game
}
