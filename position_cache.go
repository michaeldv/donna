// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

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
	id	uint32	// 4
	move	Move	// +4 = 8
	xscore	int16	// +2 = 10
	xdepth	int8	// +1 = 11
	flags	uint8	// +1 = 12
	padding	uint8	// +1 = 13
}

type Cache []CacheEntry

func cacheUsage() (hits int) {
	for i := 0; i < len(game.cache); i++ {
		if game.cache[i].id != uint32(0) {
			hits++
		}
	}

	return hits
}

func (ce *CacheEntry) score(ply int) int {
	score := int(ce.xscore)

	if score >= matingIn(MaxPly) {
		return score - ply
	} else if score <= matedIn(MaxPly) {
		return score + ply
	}

	return score
}

func (ce *CacheEntry) depth() int {
	return int(ce.xdepth)
}

func (ce *CacheEntry) token() uint8 {
	return ce.flags & 0xFC
}

func (ce *CacheEntry) bounds() uint8 {
	return ce.flags & 3
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
		entry := &game.cache[index]

		if depth > entry.depth() || game.token != entry.token() {
			if score >= matingIn(MaxPly) {
				entry.xscore = int16(score + ply)
			} else if score <= matedIn(MaxPly) {
				entry.xscore = int16(score - ply)
			} else {
				entry.xscore = int16(score)
			}
			id := uint32(p.id >> 32)
			if move.some() || id != entry.id {
				entry.move = move
			}
			entry.xdepth = int8(depth)
			entry.flags = flags | game.token
			entry.id = id
		}
	}

	return p
}

func (p *Position) probeCache() *CacheEntry {
	if cacheSize := len(game.cache); cacheSize > 0 {
		index := p.id & uint64(cacheSize - 1)
		if entry := &game.cache[index]; entry.id == uint32(p.id >> 32) {
			return entry
		}
	}

	return nil
}
