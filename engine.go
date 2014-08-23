// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `time`

const Ping = 125 // Check time 8 times a second.

type Clock struct {
	halt        bool     // Stop search immediately when set to true.
	checkpoint  int64    // First time limit check.
	softStop    int64    // Intermediate soft time limit.
	hardStop    int64    // Immediate hard time limit.
	ticker      *time.Ticker
}

type Options struct {
	ponder      bool     // (-) Pondering mode.
	infinite    bool     // (-) Search until the "stop" command.
	maxDepth    int      // Search X plies only.
	maxNodes    int      // (-) Search X nodes only.
	movesToGo   int      // Number of moves to make till time control.
	moveTime    int64    // Search exactly X milliseconds per move.
	timeLeft    int64    // Time left for all remaining moves.
	timeInc     int64    // Time increment after the move is made.
}

type Engine struct {
	log          bool     // Enable logging.
	trace        bool     // Trace evaluation scores.
	fancy        bool     // Represent pieces as UTF-8 characters.
	status       uint8    // Engine status.
	cacheSize    float64  // Default cache size.
	clock        Clock
	options      Options
}

// Use single statically allocated variable.
var engine Engine

func Self() *Engine {
	engine = Engine{}
	return &engine
}

func (e *Engine) startClock() {
	e.clock.halt = false

	if e.options.moveTime == 0 && e.options.timeLeft == 0 {
		return
	}

	if e.options.moveTime > 0 {
		start := time.Now()
		e.clock.ticker = time.NewTicker(time.Millisecond * Ping)
		go func() {
			if len(game.rootpv) == 0 {
				return // Haven't found the move yet.
			}
			for now := range e.clock.ticker.C {
				elapsed := now.Sub(start).Nanoseconds() / 1000000
				//Log("    ->clock %d limit %d left %d\n", elapsed, e.options.moveTime, (e.options.moveTime - elapsed))
				if elapsed >= e.options.moveTime - Ping {
					//Log("    <-CLOCK %d limit %d left %d\n", elapsed, e.options.moveTime, (e.options.moveTime - elapsed))
					e.clock.halt = true
				}
			}
		}()
	}
}

func (e *Engine) stopClock() {
	if e.clock.ticker != nil {
		e.clock.ticker.Stop()
	}
}

func (e *Engine) Set(args ...interface{}) *Engine {
	for i := 0; i < len(args); i += 2 {
		key, value := args[i], args[i+1]
		//fmt.Printf("engine.Set(key `%s` value %v)\n", key, value)
		switch key {
		case `fancy`:
			e.fancy = value.(bool)
		case `log`:
			e.log = value.(bool)
		case `trace`:
			e.trace = value.(bool)
		case `cache`:
			switch value.(type) {
			default: // :-)
				e.cacheSize = value.(float64)
			case int:
				e.cacheSize = float64(value.(int))
			}
		case `ponder`:
			e.options = Options{}
			e.options.ponder = true
		case `infinite`:
			e.options = Options{}
			e.options.infinite = true
		case `depth`:
			e.options = Options{}
			e.options.maxDepth = value.(int)
		case `movetime`:
			e.options = Options{}
			e.options.moveTime = int64(value.(int))
		case `time`:
			e.options.ponder   = false
			e.options.infinite = false
			e.options.maxDepth = 0
			e.options.maxNodes = 0
			e.options.moveTime = 0
			e.options.timeLeft = int64(value.(int))
		case `timeinc`:
			e.options.ponder   = false
			e.options.infinite = false
			e.options.maxDepth = 0
			e.options.maxNodes = 0
			e.options.moveTime = 0
			e.options.timeInc  = int64(value.(int))
		case `movestogo`:
			e.options.ponder   = false
			e.options.infinite = false
			e.options.maxDepth = 0
			e.options.maxNodes = 0
			e.options.moveTime = 0
			e.options.movesToGo = value.(int)
		}
	}
	// fmt.Printf("Engine options: %v\n", e.options)

	return e
}
