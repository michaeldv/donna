// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package main

import (
	`github.com/michaeldv/donna`
	`github.com/michaeldv/donna/cli`
	`os`
)

func main() {
	engine := donna.Self().Set(`fancy`, true)

	if len(os.Args) > 1 && os.Args[1] == `-i` {
		cli.Repl(engine)
	} else {
		cli.Uci(engine)
	}
}
