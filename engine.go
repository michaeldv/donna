// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`fmt`; `os`; `time`)

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
		switch value := args[i+1]; args[i] {
		case `log`:
			engine.log = value.(bool)
		case `uci`:
			engine.uci = value.(bool)
		case `trace`:
			engine.trace = value.(bool)
		case `fancy`:
			engine.fancy = value.(bool)
		case `depth`:
			engine.options.maxDepth = value.(int)
		case `movetime`:
			engine.options.moveTime = int64(value.(int))
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

// Dumps the string to standard output.
func (e *Engine) print(arg string) *Engine {
	os.Stdout.WriteString(arg)
	return e
}

// Appends the string to log file.
func (e *Engine) debug(arg string) *Engine {
	logFile, err := os.OpenFile("/tmp/donna.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err == nil {
		defer logFile.Close()
		logFile.WriteString(arg) // f.Write() and friends are unbuffered.
	}
	return e
}

// Dumps the string to standard output and logs it to file.
func (e *Engine) reply(args ...interface{}) *Engine {
	if len := len(args); len > 1 {
		data := fmt.Sprintf(args[0].(string), args[1:]...)
		e.print(data)
		e.debug(data)
	} else if len == 1 {
		e.print(args[0].(string))
		e.debug(args[0].(string))
	}
	return e
}

func (e *Engine) startClock() *Engine {
	e.clock.halt = false

	if e.options.moveTime == 0 && e.options.timeLeft == 0 {
		return e
	}

	if e.options.moveTime > 0 {
		return e.fixedMoveTime()
	}

	return e.varyingMoveTime()
}

func (e *Engine) stopClock() *Engine {
	if e.clock.ticker != nil {
		e.clock.ticker.Stop()
		e.clock.ticker = nil
	}
	return e
}

func (e *Engine) fixedMoveTime() *Engine {
	start := time.Now()
	e.clock.ticker = time.NewTicker(time.Millisecond * Ping)

	go func() {
		if e.clock.ticker == nil {
			return
		}
		for now := range e.clock.ticker.C {
			if len(game.rootpv) == 0 {
				continue // Haven't found the move yet.
			}
			elapsed := now.Sub(start).Nanoseconds() / 1000000
			if elapsed >= e.options.moveTime - Ping {
				e.clock.halt = true
				return
			}
		}
	}()

	return e
}

func (e *Engine) varyingMoveTime() *Engine {
	start := time.Now()
	e.clock.ticker = time.NewTicker(time.Millisecond * Ping)

	go func() {
		if e.clock.ticker == nil {
			return
		}
		for now := range e.clock.ticker.C {
			if len(game.rootpv) == 0 {
				continue // Haven't found the move yet.
			}
			elapsed := now.Sub(start).Nanoseconds() / 1000000
			// TODO:
			// - UCI info reporting
			// - better time management taking into account fail
			//   high/dropping scores and oft/hard time limits.
			if elapsed >= int64(e.clock.checkpoint - Ping) {
				e.debug(fmt.Sprintf("# halt %d limit %d left %d\n", elapsed, e.clock.checkpoint, e.clock.checkpoint - elapsed))
				e.clock.halt = true
				return
			}
		}
	}()

	return e
}

func (e *Engine) fixedLimit(options Options) *Engine {
	e.options = options
	return e
}

func (e *Engine) varyingLimits(options Options) *Engine {
	var moves, soft, hard int64

	e.options = options
	e.options.ponder = false
	e.options.infinite = false
	e.options.maxDepth = 0
	e.options.maxNodes = 0
	e.options.moveTime = 0

	// Use known number of moves till the end of the game or time control.
	moves = int64(e.options.movesToGo)
	if moves == 0 {
		moves = int64(40) // Default. TODO: calculate based on game phase.
	}

	// Calculate hard and soft stops.
	hard = options.timeLeft + options.timeInc * (moves - 1)
	soft = Max64(0, hard / moves * 120 / 100) * 4

	// Adjust hard stop to leave some emergency reserve plus account for
	// possible I/O lag.
	hard -= hard * (moves - 1) / 50
	hard -= Max64(50, hard * 5 / 100) // 5% or 50ms.
	hard = Max64(0, hard)

	e.clock.hardStop = hard
	e.clock.softStop = Min64(hard, soft)
	e.clock.checkpoint = Min64(hard, soft / 4)
	e.debug(fmt.Sprintf("# Make %d moves in %02d:%02ds\n", moves, e.options.timeLeft / 1000 / 60, e.options.timeLeft / 1000 % 60))
	e.debug(fmt.Sprintf("# checkpoint: %8d -> %02d:%02ds\n", e.clock.checkpoint, e.clock.checkpoint / 1000 / 60, e.clock.checkpoint / 1000 % 60))
	e.debug(fmt.Sprintf("#   softStop: %8d -> %02d:%02ds\n", e.clock.softStop, e.clock.softStop / 1000 / 60, e.clock.softStop / 1000 % 60))
	e.debug(fmt.Sprintf("#   hardStop: %8d -> %02d:%02ds\n", e.clock.hardStop, e.clock.hardStop / 1000 / 60, e.clock.hardStop / 1000 % 60))

	return e
}
