// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`bytes`
	`fmt`
)

type Bitmask uint64

// Returns true if all bitmask bits are clear.
func (b Bitmask) isEmpty() bool {
	return b == 0
}

func (b Bitmask) isNotEmpty() bool {
	return b != 0
}

// Returns true if a bit at given position is set.
func (b Bitmask) isSet(position int) bool {
	return b & (1 << uint(position)) != 0
}

// Returns true if a bit at given position is clear.
func (b Bitmask) isClear(position int) bool {
	return !b.isSet(position)
}

// Returns number of bits set.
func (b Bitmask) count() int {
	mask := b
	mask -= (mask >> 1) & 0x5555555555555555
	mask = ((mask >> 2) & 0x3333333333333333) + (mask & 0x3333333333333333)
	mask = ((mask >> 4) + mask) & 0x0F0F0F0F0F0F0F0F
	return int((mask * 0x0101010101010101) >> 56)
}

// Finds least significant bit set (LSB) in non-zero bitmask. Returns
// an integer in 0..63 range.
func (b Bitmask) first() int {
	return deBruijn[((b ^ (b-1)) * 0x03F79D71B4CB0A89) >> 58]
}

// MSB: Eugene Nalimov's bitScanReverse.
func (b Bitmask) last() int {
	position := 0
	if b > 0xFFFFFFFF {
		b >>= 32
		position = 32
	}

	if b > 0xFFFF {
		b >>= 16
		position += 16
	}

	if b > 0xFF {
		b >>= 8
		position += 8
	}

       return position + msbLookup[b]
}

func (b Bitmask) closest(color int) int {
	if color == White {
		return b.first()
	}
	return b.last()
}

func (b Bitmask) pushed(color int) Bitmask {
	if color == White {
		return b << 8
	}
	return b >> 8
}

// Finds *and clears* least significant bit set (LSB) in non-zero
// bitmask. Returns an integer in 0..63 range.
func (b *Bitmask) pop() int {
	magic := (*b - 1)
	index := deBruijn[((*b ^ magic) * 0x03F79D71B4CB0A89) >> 58]
	*b &= magic
	return index
}

// Sets a bit at given position.
func (b *Bitmask) set(position int) *Bitmask {
	*b |= 1 << uint(position)
	return b
}

// Clears a bit at given position.
func (b *Bitmask) clear(position int) *Bitmask {
	*b &= ^(1 << uint(position))
	return b
}

// Combines two bitmasks using bitwise OR operator.
func (b *Bitmask) combine(bitmask Bitmask) *Bitmask {
	*b |= bitmask
	return b
}

// Intersects two bitmasks using bitwise AND operator.
func (b *Bitmask) intersect(bitmask Bitmask) *Bitmask {
	*b &= bitmask
	return b
}

// Excludes bits of one bitmask from another using bitwise XOR operator.
func (b *Bitmask) exclude(bitmask Bitmask) *Bitmask {
	*b ^= (bitmask & *b)
	return b
}

func (b *Bitmask) shift(offset int) *Bitmask {
	if offset > 0 {
		*b <<= uint(offset)
	} else {
		*b >>= -uint(offset)
	}
	return b
}

func (b *Bitmask) fill(square, direction int, occupied, board Bitmask) *Bitmask {
	mask := bit[square] & board

	for mask.shift(direction); mask.isNotEmpty(); mask.shift(direction) {
		b.combine(mask)
		if (mask & occupied).isNotEmpty() {
			break
		}
		mask.intersect(board)
	}
	return b
}

func (b *Bitmask) spot(square, direction int, board Bitmask) *Bitmask {
	*b = bit[square] & board
	*b = ^*(b.shift(direction))
	return b
}

func (b *Bitmask) trim(row, col int) *Bitmask {
	if row > 0 {
		*b &= 0xFFFFFFFFFFFFFF00
	}
	if row < 7 {
		*b &= 0x00FFFFFFFFFFFFFF
	}
	if col > 0 {
		*b &= 0xFEFEFEFEFEFEFEFE
	}
	if col < 7 {
		*b &= 0x7F7F7F7F7F7F7F7F
	}

	return b
}

func (b Bitmask) String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h  ")
	buffer.WriteString(fmt.Sprintf("0x%016X\n", uint64(b)))
	for row := 7; row >= 0; row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0; col <= 7; col++ {
			position := row << 3 + col
			buffer.WriteByte(' ')
			if b.isSet(position) {
				buffer.WriteString("\u2022") // Set
			} else {
				buffer.WriteString("\u22C5") // Clear
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

//     A   B   C   D   E   F   G   H
// 7>  56  57  58  59  60  61  62  63
// 6>  48  49  50  51  52  53  54  55
// 5>  40  41  42  43  44  45  46  47
// 4>  32  33  34  35  36  37  38  39
// 3>  24  25  26  27  28  29  30  31
// 2>  16  17  18  19  20  21  22  23
// 1>  08  09  10  11  12  13  14  15
// 0>  00  01  02  03  04  05  06  07
//     ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//     0   1   2   3   4   5   6   7
