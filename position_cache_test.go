// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

func TestCache000(t *testing.T) {
	engine.cacheSize = 0.5
	p := NewGame().start()
	move := NewMove(p, E2, E4)
	p = p.makeMove(move).cache(move, 42, 1, 0, cacheExact)

	cached := p.probeCache()
	expect.Eq(t, cached.move, move)
	expect.Eq(t, cached.score, int16(42))
	expect.Eq(t, cached.depth, int16(1))
	expect.Eq(t, cached.flags, uint8(cacheExact))
	expect.Eq(t, cached.id, uint32(p.id >> 32))
}
