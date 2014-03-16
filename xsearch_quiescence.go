// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`fmt`)

// Quiescence search.
func (p *Position) xSearchQuiescence(alpha, beta, depth int) int {
        alpha = alpha
        beta = beta
        depth = depth

        fmt.Printf("%*squie/%s> depth: %d, ply: %d\n", depth*2, ` `, C(p.color), depth, Ply())

        return 0
}
