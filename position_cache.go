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
	if megaBytes > 0.0 {
		cacheSize := int(1024*1024*megaBytes) / int(unsafe.Sizeof(CacheEntry{}))
		// If cache size has changed then create a new cache; otherwise
		// simply clear the existing one.
		if cacheSize != len(game.cache) {
			return make(Cache, cacheSize)
		}
		game.cache = Cache{}
		return game.cache
	}
	return nil
}

func (p *Position) cache(move Move, score, depth int, flags uint8) *Position {
	if cacheSize := len(game.cache); cacheSize > 0 {
		index := p.hash % uint64(cacheSize)
		// fmt.Printf("cache size %d entries, index %d\n", len(game.cache), index)
		entry := &game.cache[index]

		if depth > entry.depth || game.token != entry.token {
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
			entry.token = game.token
			entry.hash = p.hash
		}
	}

	return p
}

func (p *Position) probeCache() *CacheEntry {
	if cacheSize := len(game.cache); cacheSize > 0 {
		index := p.hash % uint64(cacheSize)
		if entry := &game.cache[index]; entry.hash == p.hash {
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
