// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package main

import (
	`github.com/michaeldv/donna`
	`os`
	`runtime`
)

func main() {
	// Default engine settings are: 128MB transposition table, 5s per move.
	engine := donna.NewEngine(
		`fancy`, runtime.GOOS == `darwin`,
		`cache`, 128,
		`movetime`, 5000,
		`logfile`, os.Getenv(`DONNA_LOG`),
		`bookfile`, os.Getenv(`DONNA_BOOK`),
	)

	if len(os.Args) > 1 && os.Args[1] == `-i` {
		engine.Repl()
	} else {
		engine.Uci()
	}
}
