// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `unsafe`

const (
	cacheNone  = uint8(0)
	cacheAlpha = uint8(1) // Upper bound.
	cacheBeta  = uint8(2) // Lower bound.
	cacheExact = uint8(cacheAlpha | cacheBeta)
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

// Creates new or resets existing game cache (aka transposition table).
func NewCache(megaBytes float64) Cache {
	if megaBytes > 0.0 {
		cacheSize := int(1024 * 1024 * megaBytes) / cacheEntrySize
		// Cache size has changed.
		if existing := len(game.cache); cacheSize != existing {
			if existing == 0 {
				// Create brand new zero-initialized cache.
				return make(Cache, cacheSize)
			} else {
				// Reallocate existing cache (shrink or expand).
				game.cache = append([]CacheEntry{}, game.cache[:cacheSize]...)
			}
		}
		// Make sure the cache is all clear.
		for i := 0; i < len(game.cache); i++ {
			game.cache[i] = CacheEntry{}
		}
		return game.cache
	}
	return nil
}

func (p *Position) cache(move Move, score, depth, ply int, flags uint8) *Position {
	if cacheSize := len(game.cache); cacheSize > 0 {
		index := p.id & uint64(cacheSize - 1)
		// fmt.Printf("cache size %d entries, index %d\n", len(game.cache), index)
		entry := &game.cache[index]

		if depth > int(entry.depth) || game.token != entry.token {
			if score > Checkmate - MaxPly && score <= Checkmate {
				entry.score = int16(score + ply)
			} else if score < MaxPly - Checkmate && score >= -Checkmate {
				entry.score = int16(score - ply)
			} else {
				entry.score = int16(score)
			}
			id := uint32(p.id >> 32)
			if move != Move(0) || id != entry.id {
				entry.move = move
			}
			entry.depth = int16(depth)
			entry.flags = flags
			entry.token = game.token
			entry.id = id
		}
	}

	return p
}

func (p *Position) cacheDelta(move Move, score, depth, ply, alpha, beta int) *Position {
	cacheFlags := cacheNone

	if score > alpha {
		cacheFlags |= cacheBeta
	}
	if score < beta {
		cacheFlags |= cacheAlpha
	}

	return p.cache(move, score, depth, ply, cacheFlags)
}

func (p *Position) probeCache() *CacheEntry {
	if cacheSize := len(game.cache); cacheSize > 0 {
		index := p.id & uint64(cacheSize - 1)
		if entry := &game.cache[index]; entry.id == uint32(p.id >>32) {
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
