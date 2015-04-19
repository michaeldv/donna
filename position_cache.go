// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `unsafe`

const (
	cachedNone = uint8(0)
	cacheExact = uint8(1)
	cacheAlpha = uint8(2) // Upper bound.
	cacheBeta  = uint8(4) // Lower bound.
	cacheEntrySize = int(unsafe.Sizeof(CacheEntry{}))
)

type CacheEntry struct {
	id    uint32
	move  Move
	score int16
	depth int16
	flags uint8
	token uint8
}

type Cache []CacheEntry

func cacheUsage() (hits int) {
	for i := 0; i < len(game.cache); i++ {
		if game.cache[i].id != uint32(0) {
			hits++
		}
	}
	return
}

func NewCache(megaBytes float64) Cache {
	if megaBytes > 0.0 {
		cacheSize := int(1024*1024*megaBytes) / cacheEntrySize
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

func (p *Position) cache(move Move, score, depth, ply int, flags uint8) *Position {
	if cacheSize := len(game.cache); cacheSize > 0 {
		index := p.hash & uint64(cacheSize - 1)
		// fmt.Printf("cache size %d entries, index %d\n", len(game.cache), index)
		entry := &game.cache[index]

		if depth > int(entry.depth) || game.token != entry.token {
			if score > Checkmate-MaxPly && score <= Checkmate {
				entry.score = int16(score + ply)
			} else if score >= -Checkmate && score < -Checkmate+MaxPly {
				entry.score = int16(score - ply)
			} else {
				entry.score = int16(score)
			}
			if move != Move(0) || uint32(p.hash >> 32) != entry.id {
				entry.move = move
			}
			entry.depth = int16(depth)
			entry.flags = flags
			entry.token = game.token
			entry.id = uint32(p.hash >> 32)
		}
	}

	return p
}

func (p *Position) probeCache() *CacheEntry {
	if cacheSize := len(game.cache); cacheSize > 0 {
		index := p.hash & uint64(cacheSize - 1)
		if entry := &game.cache[index]; entry.id == uint32(p.hash >>32) {
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
