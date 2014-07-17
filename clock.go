// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`fmt`
	`time`
)

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
		game.clock.ticker = time.NewTicker(time.Millisecond * 125) // 8 times 1a second.
		go func() {
			for now := range game.clock.ticker.C {
				elapsed := now.Sub(start).Nanoseconds() / 1000000
				fmt.Printf("    ->clock %d limit %d left %d\n", elapsed, game.options.msMoveTime, (game.options.msMoveTime - elapsed))
				if elapsed >= game.options.msMoveTime {
					fmt.Printf("    <-CLOCK %d limit %d left %d\n", elapsed, game.options.msMoveTime, (game.options.msMoveTime - elapsed))
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

