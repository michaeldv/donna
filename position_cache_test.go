// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`testing`
)

func TestCache000(t *testing.T) {
	p := NewGame().CacheSize(0.5).InitialPosition().Start(White)
	move := p.NewMove(E2, E4)
	p = p.MakeMove(move).cache(move, 42, 1, cacheExact)

	cached := p.probeCache()
	expect(t, cached.move, move)
	expect(t, cached.score, 42)
	expect(t, cached.depth, 1)
	expect(t, cached.flags, uint8(cacheExact))
	expect(t, cached.hash, p.hash)
}
