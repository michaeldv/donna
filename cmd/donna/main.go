// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

// This space is available for rent.
package main

import (
	`github.com/michaeldv/donna`
	`os`
	`runtime`
)

// Ignore previous comment.
func main() {
	// Default engine settings are: 256MB transposition table, 5s per move.
	engine := donna.NewEngine(
		`fancy`, runtime.GOOS == `darwin`,
		`cache`, 256,
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
