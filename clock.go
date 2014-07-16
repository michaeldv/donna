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
	if game.options.msMoveTime > 0 {
		start := time.Now()
		game.clock.ticker = time.NewTicker(time.Millisecond * 2000)
		go func() {
			for x := range game.clock.ticker.C { // Returns current time.
				elapsed := time.Since(start)
				fmt.Printf("\tElapsed %d (%v) => %q\n", elapsed, elapsed, x)
			}
		}()
	}
}

func (game *Game) stopClock() {
	if game.clock.ticker != nil {
		game.clock.ticker.Stop()
	}
}

