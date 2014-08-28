// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package main

import (
	`github.com/michaeldv/donna`
	`os`
)

// Default engine settings are: 64MB transposition table, 5s per move.
func main() {
	if len(os.Args) > 1 && os.Args[1] == `-i` {
		donna.NewEngine(
			`fancy`, true,
			`cache`, 64,
			`movetime`, 5000,
		).Repl()
	} else {
		donna.NewEngine(
			`uci`, true,
			`fancy`, true,
			`cache`, 64,
			`movetime`, 5000,
		).Uci()
	}
}
