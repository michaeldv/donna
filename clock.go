// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `time`

const Ping = 125 // Check time 8 times a second.

type Clock struct {
	halt        bool   // Stop search immediately when set to true.
	checkpoint  int64  // First time limit check.
	softStop    int64  // Intermediate soft time limit.
	hardStop    int64  // Immediate hard time limit.
	ticker      *time.Ticker
}

func (game *Game) startClock() {
	game.clock.halt = false

	if game.options.moveTime == 0 && game.options.gameTime == 0 {
		return
	}

	if game.options.moveTime > 0 {
		start := time.Now()
		game.clock.ticker = time.NewTicker(time.Millisecond * Ping)
		go func() {
			if len(game.rootpv) == 0 {
				return // Haven't found the move yet.
			}
			for now := range game.clock.ticker.C {
				elapsed := now.Sub(start).Nanoseconds() / 1000000
				//Log("    ->clock %d limit %d left %d\n", elapsed, game.options.moveTime, (game.options.moveTime - elapsed))
				if elapsed >= game.options.moveTime - Ping {
					//Log("    <-CLOCK %d limit %d left %d\n", elapsed, game.options.moveTime, (game.options.moveTime - elapsed))
					game.clock.halt = true
				}
			}
		}()
	}
}

func (game *Game) stopClock() {
	if game.clock.ticker != nil {
		game.clock.ticker.Stop()
	}
}

