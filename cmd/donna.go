// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package main

import (
	`github.com/michaeldv/donna`
	`os`
)

func main() {
	// Default engine settings: 64MB transposition table, 5s per move.
	engine := donna.NewEngine(`fancy`, true, `cache`, 64, `movetime`, 5000)

	if len(os.Args) > 1 && os.Args[1] == `-i` {
		engine.Repl()
	} else {
		engine.Uci()
	}
}
