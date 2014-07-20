// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `time`

const Ping = 125 // Check time 8 times a second.

type Clock struct {
	stopSearch  bool  // Stop search when set to true.
	msSoftStop  int   // Intermediate soft time limit.
	msHardStop  int   // Immediate hard time limit.
	ticker      *time.Ticker
}

func (game *Game) startClock() {
	game.clock.stopSearch = false

	if game.options.msMoveTime == 0 && game.options.msGameTime == 0 {
		return
	}

	if game.options.msMoveTime > 0 {
		start := time.Now()
		game.clock.ticker = time.NewTicker(time.Millisecond * Ping)
		go func() {
			if len(game.rootpv) == 0 {
				return // Haven't found the move yet.
			}
			for now := range game.clock.ticker.C {
				elapsed := now.Sub(start).Nanoseconds() / 1000000
				//Log("    ->clock %d limit %d left %d\n", elapsed, game.options.msMoveTime, (game.options.msMoveTime - elapsed))
				if elapsed >= (game.options.msMoveTime - Ping) {
					//Log("    <-CLOCK %d limit %d left %d\n", elapsed, game.options.msMoveTime, (game.options.msMoveTime - elapsed))
					game.clock.stopSearch = true
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

