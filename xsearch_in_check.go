// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`fmt`)

// Search for the node in check.
func (p *Position) xSearchInCheck(beta, depth int) int {
        beta = beta
        depth = depth

        fmt.Printf("%*schck/%s> depth: %d, ply: %d\n", depth*2, ` `, C(p.color), depth, Ply())

        return 0
}
