// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`unsafe`
)

const (
	cacheExact = iota
	cacheAlpha // Upper bound.
	cacheBeta  // Lower bound.
)

type CacheEntry struct {
	move  Move
	score int
	depth int
	flags uint8
	token uint8
	hash  uint64
}

type Cache []CacheEntry

func NewCache(megaBytes float64) Cache {
	cacheSize := uint(1024*1024*megaBytes) / uint(unsafe.Sizeof(CacheEntry{}))
	return make(Cache, cacheSize)
}

func (p *Position) cache(move Move, score, depth int, flags uint8) *Position {
	if cacheSize := len(p.game.cache); cacheSize > 0 {
		index := p.hash % uint64(cacheSize)
		// fmt.Printf("cache size %d entries, index %d\n", len(p.game.cache), index)
		entry := &p.game.cache[index]

		if depth > entry.depth || p.game.token != entry.token {
			if score > Checkmate-MaxPly && score <= Checkmate {
				entry.score = score + Ply()
			} else if score >= -Checkmate && score < -Checkmate+MaxPly {
				entry.score = score - Ply()
			} else {
				entry.score = score
			}
			entry.move = move
			entry.depth = depth
			entry.flags = flags
			entry.token = p.game.token
			entry.hash = p.hash
		}
	}

	return p
}

func (p *Position) probeCache() *CacheEntry {
	if cacheSize := len(p.game.cache); cacheSize > 0 {
		index := p.hash % uint64(cacheSize)
		if entry := &p.game.cache[index]; entry.hash == p.hash {
			return entry
		}
	}
	return nil
}

func (p *Position) cachedMove() Move {
	if cached := p.probeCache(); cached != nil {
		return cached.move
	}
	return Move(0)
}
