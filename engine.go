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
	log         bool     // Enable logging.
	uci	    bool     // Use UCI protocol.
	trace       bool     // Trace evaluation scores.
	fancy       bool     // Represent pieces as UTF-8 characters.
	status      uint8    // Engine status.
	cacheSize   float64  // Default cache size.
	clock       Clock
	options     Options
}

// Use single statically allocated variable.
var engine Engine

func NewEngine(args ...interface{}) *Engine {
	engine = Engine{}
	for i := 0; i < len(args); i += 2 {
		key, value := args[i], args[i+1]
		//fmt.Printf("engine.Set(key `%s` value %v)\n", key, value)
		switch key {
		case `log`:
			engine.log = value.(bool)
		case `uci`:
			engine.uci = value.(bool)
		case `trace`:
			engine.trace = value.(bool)
		case `fancy`:
			engine.fancy = value.(bool)
		case `cache`:
			switch value.(type) {
			default: // :-)
				engine.cacheSize = value.(float64)
			case int:
				engine.cacheSize = float64(value.(int))
			}
		}
	}

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

func (e *Engine) fixedLimit(options Options) *Engine {
	e.options = options
	return e
}

func (e *Engine) variableLimits(options Options) *Engine {
	e.options = options
	options.ponder = false
	options.infinite = false
	options.maxDepth = 0
	options.maxNodes = 0
	options.moveTime = 0
	return e
}
