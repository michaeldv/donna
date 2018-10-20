// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

import (`fmt`; `os`; `time`)

const Ping = 250 // Check time 4 times a second.

type Clock struct {
	halt        bool     // Stop search immediately when set to true.
	softStop    int64    // Target soft time limit to make a move.
	hardStop    int64    // Immediate stop time limit.
	extra       float32  // Extra time factor based on search volatility.
	start       time.Time
	ticker      *time.Ticker
}

type Options struct {
	ponder      bool     // (-) Pondering mode.
	infinite    bool     // (-) Search until the "stop" command.
	maxDepth    int      // Search X plies only.
	maxNodes    int      // (-) Search X nodes only.
	moveTime    int64    // Search exactly X milliseconds per move.
	movesToGo   int64    // Number of moves to make till time control.
	timeLeft    int64    // Time left for all remaining moves.
	timeInc     int64    // Time increment after the move is made.
}

type Engine struct {
	log         bool     // Enable logging.
	uci	    bool     // Use UCI protocol.
	trace       bool     // Trace evaluation scores.
	fancy       bool     // Represent pieces as UTF-8 characters.
	status      uint8    // Engine status.
	logFile     string   // Log file name.
	bookFile    string   // Polyglot opening book file name.
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
		case `logfile`:
			engine.logFile = value.(string)
		case `bookfile`:
			engine.bookFile = value.(string)
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
	os.Stdout.Sync() // <-- Flush it.
	return e
}

// Appends the string to log file. No flush is required as f.Write() and friends
// are unbuffered.
func (e *Engine) debug(args ...interface{}) *Engine {
	if len(e.logFile) != 0 {
		logFile, err := os.OpenFile(e.logFile, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
		if err == nil {
			defer logFile.Close()
			if len := len(args); len > 1 {
				logFile.WriteString(fmt.Sprintf(args[0].(string), args[1:]...))
			} else {
				logFile.WriteString(args[0].(string))
			}
		}
	}
	return e
}

// Dumps the string to standard output and optionally logs it to file.
func (e *Engine) reply(args ...interface{}) *Engine {
	if len := len(args); len > 1 {
		data := fmt.Sprintf(args[0].(string), args[1:]...)
		e.print(data)
		//\\ e.debug(data)
	} else if len == 1 {
		e.print(args[0].(string))
		//\\ e.debug(args[0].(string))
	}
	return e
}

func (e *Engine) fixedDepth() bool {
	return e.options.maxDepth > 0
}

func (e *Engine) fixedTime() bool {
	return e.options.moveTime > 0
}

func (e *Engine) varyingTime() bool {
	return e.options.moveTime == 0
}


// Returns elapsed time in milliseconds.
func (e *Engine) elapsed(now time.Time) int64 {
	return now.Sub(e.clock.start).Nanoseconds() / 1000000 //int64(time.Millisecond)
}

// Returns remaining search time to make a move. The remaining time extends the
// soft stop estimate based on search volatility factor.
func (e *Engine) remaining() int64 {
	return int64(float32(e.clock.softStop) * e.clock.extra)
}

// Sets extra time factor. For depths 5+ we take into account search volatility,
// i.e. extra time is given for uncertain positions where the best move is not clear.
func (e *Engine) factor(depth int, volatility float32) *Engine {
	e.clock.extra = 0.75
	if depth >= 5 {
		e.clock.extra *= (volatility + 1.0)
	}

	return e
}

// Starts the clock setting ticker callback function. The callback function is
// different for fixed and variable time controls.
func (e *Engine) startClock() *Engine {
	e.clock.halt = false

	if e.options.moveTime == 0 && e.options.timeLeft == 0 {
		return e
	}

	e.clock.start = time.Now()
	e.clock.ticker = time.NewTicker(time.Millisecond * Ping)

	if e.fixedTime() {
		return e.fixedTimeTicker()
	}

	// How long a minute is depends on which side of the bathroom door you're on.
	return e.varyingTimeTicker()
}

// Stop the clock so that the ticker callback function is longer invoked.
func (e *Engine) stopClock() *Engine {
	if e.clock.ticker != nil {
		e.clock.ticker.Stop()
		e.clock.ticker = nil
	}
	return e
}

// Ticker callback for fixed time control (ex. 5s per move). Search gets terminated
// when we've got the move and the elapsed time approaches time-per-move limit.
func (e *Engine) fixedTimeTicker() *Engine {
	go func() {
		if e.clock.ticker == nil {
			return // Nothing to do if the clock has been stopped.
		}
		for now := range e.clock.ticker.C {
			if game.rootpv.size == 0 {
				continue // Haven't found the move yet.
			}
			if e.elapsed(now) >= e.options.moveTime - Ping {
				e.clock.halt = true
				return
			}
		}
	}()

	return e
}

// Ticker callback for the variable time control (ex. 40 moves in 5 minutes). Search
// termination depends on multiple factors with hard stop being the ultimate limit.
func (e *Engine) varyingTimeTicker() *Engine {
	go func() {
		if e.clock.ticker == nil {
			return // Nothing to do if the clock has been stopped.
		}
		for now := range e.clock.ticker.C {
			if game.rootpv.size == 0 {
				continue // Haven't found the move yet.
			}
			elapsed := e.elapsed(now)
			if (game.deepening && game.improving && elapsed > e.remaining() * 4 / 5) || elapsed > e.clock.hardStop {
				//\\ e.debug("# Halt: Flags %v Elapsed %s Remaining %s Hard stop %s\n",
				//\\	game.deepening && game.improving, ms(elapsed), ms(e.remaining() * 4 / 5), ms(e.clock.hardStop))
				e.clock.halt = true
				return
			}
		}
	}()

	return e
}

// Sets fixed search limits such as maximum depth or time to make a move.
func (e *Engine) fixedLimit(options Options) *Engine {
	e.options = options
	return e
}

// Sets variable time control options and calculates soft and hard stop estimates.
func (e *Engine) varyingLimits(options Options) *Engine {

	// Note if it's a new time control before saving the options.
	e.options = options
	e.options.ponder = false
	e.options.infinite = false
	e.options.maxDepth = 0
	e.options.maxNodes = 0
	e.options.moveTime = 0

	// Set default number of moves till the end of the game or time control.
	// TODO: calculate based on game phase.
	if e.options.movesToGo == 0 {
		e.options.movesToGo = 40
	}

	// Calculate hard and soft stop estimates.
	moves := e.options.movesToGo - 1
	hard := options.timeLeft + options.timeInc * moves
	soft := hard / e.options.movesToGo

	//\\ e.debug("#\n# Make %d moves in %s soft stop %s hard stop %s\n", e.options.movesToGo, ms(e.options.timeLeft), ms(soft), ms(hard))

	// Adjust hard stop to leave enough time reserve for the remaining moves. The time
	// reserve starts at 100% of soft stop for one remaining move, and goes down to 80%
	// in 1% decrement for 20+ remaining moves.
	if moves > 0 { // The last move gets all remaining time and doesn't need the reserve.
		percent := max64(80, 100 - moves)
		reserve := soft * moves * percent / 100
		//\\ e.debug("# Reserve %d%% = %s\n", percent, ms(reserve))
		if hard - reserve > soft {
			hard -= reserve
		}
		// Hard stop can't exceed optimal time to make 3 moves.
		hard = min64(hard, soft * 3)
		//\\ e.debug("# Hard stop %s\n", ms(hard))
	}

	// Set the final values for soft and hard stops making sure the soft stop
	// never exceeds the hard one.
	if soft < hard {
		e.clock.softStop, e.clock.hardStop = soft, hard
	} else {
		e.clock.softStop, e.clock.hardStop = hard, soft
	}

	// Keep two ping cycles available to avoid accidental time forefeit.
	e.clock.hardStop -= 2 * Ping
	if e.clock.hardStop < 0 {
		e.clock.hardStop = options.timeLeft // Oh well...
	}

	//\\ e.debug("# Final soft stop %s hard stop %s\n#\n", ms(e.clock.softStop), ms(e.clock.hardStop))

	return e
}
