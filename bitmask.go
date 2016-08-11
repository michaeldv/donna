// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`bytes`; `fmt`)

type Bitmask uint64

// One man's constant is another man's variable.
var bit = [64]Bitmask{
	1<<A1, 1<<B1, 1<<C1, 1<<D1, 1<<E1, 1<<F1, 1<<G1, 1<<H1,
	1<<A2, 1<<B2, 1<<C2, 1<<D2, 1<<E2, 1<<F2, 1<<G2, 1<<H2,
	1<<A3, 1<<B3, 1<<C3, 1<<D3, 1<<E3, 1<<F3, 1<<G3, 1<<H3,
	1<<A4, 1<<B4, 1<<C4, 1<<D4, 1<<E4, 1<<F4, 1<<G4, 1<<H4,
	1<<A5, 1<<B5, 1<<C5, 1<<D5, 1<<E5, 1<<F5, 1<<G5, 1<<H5,
	1<<A6, 1<<B6, 1<<C6, 1<<D6, 1<<E6, 1<<F6, 1<<G6, 1<<H6,
	1<<A7, 1<<B7, 1<<C7, 1<<D7, 1<<E7, 1<<F7, 1<<G7, 1<<H7,
	1<<A8, 1<<B8, 1<<C8, 1<<D8, 1<<E8, 1<<F8, 1<<G8, 1<<H8,
}

var deBruijn = [64]int{
	 0, 47,  1, 56, 48, 27,  2, 60,
	57, 49, 41, 37, 28, 16,  3, 61,
	54, 58, 35, 52, 50, 42, 21, 44,
	38, 32, 29, 23, 17, 11,  4, 62,
	46, 55, 26, 59, 40, 36, 15, 53,
	34, 51, 20, 43, 31, 22, 10, 45,
	25, 39, 14, 33, 19, 30,  9, 24,
	13, 18,  8, 12,  7,  6,  5, 63,
}

// Most-significant bit (MSB) lookup table.
var msbLookup[256]int

func init() {
	for i := 0; i < len(msbLookup); i++ {
		if i > 127 {
			msbLookup[i] = 7
		} else if i > 63 {
			msbLookup[i] = 6
		} else if i > 31 {
			msbLookup[i] = 5
		} else if i > 15 {
			msbLookup[i] = 4
		} else if i > 7 {
			msbLookup[i] = 3
		} else if i > 3 {
			msbLookup[i] = 2
		} else if i > 1 {
			msbLookup[i] = 1
		}
	}
}

// Returns true if all bitmask bits are clear. Even if it's wrong, it's only
// off by a bit.
func (b Bitmask) empty() bool {
	return b == 0
}

// Returns true if at least one bit is set.
func (b Bitmask) any() bool {
	return b != 0
}

// Returns true if a bit at given offset is set.
func (b Bitmask) on(offset int) bool {
	return (b & bit[offset]).any()
}

// Returns true if a bit at given offset is clear.
func (b Bitmask) off(offset int) bool {
	return !b.on(offset)
}

// Returns true if a bitmask has single bit set.
func (b Bitmask) single() bool {
	return b.pop().empty()
}

// Returns number of bits set.
func (b Bitmask) count() int {
	if b.empty() {
		return 0
	}
	b -= (b >> 1) & 0x5555555555555555
	b = ((b >> 2) & 0x3333333333333333) + (b & 0x3333333333333333)
	b = ((b >> 4) + b) & 0x0F0F0F0F0F0F0F0F
	return int((b * 0x0101010101010101) >> 56)
}

// Finds least significant bit set (LSB) in non-zero bitmask. Returns
// an integer in 0..63 range.
func (b Bitmask) first() int {
	return deBruijn[((b ^ (b - 1)) * 0x03F79D71B4CB0A89) >> 58]
}

// Eugene Nalimov's bitScanReverse: finds most significant bit set (MSB).
func (b Bitmask) last() (offset int) {
	if b > 0xFFFFFFFF {
		b >>= 32; offset = 32
	}
	if b > 0xFFFF {
		b >>= 16; offset += 16
	}
	if b > 0xFF {
		b >>= 8; offset += 8
	}

	return offset + msbLookup[b]
}

func (b Bitmask) closest(color int) int {
	if color == White {
		return b.first()
	}
	return b.last()
}

func (b Bitmask) farthest(color int) int {
    if color == White {
        return b.last()
    }
    return b.first()
}

func (b Bitmask) up(color int) Bitmask {
	if color == White {
		return b << 8
	}
	return b >> 8
}

// Returns bitmask with least significant bit off.
func (b Bitmask) pop() Bitmask {
	return b & (b - 1)
}

// Sets a bit at given offset.
func (b *Bitmask) set(offset int) *Bitmask {
	*b |= bit[offset]
	return b
}

// Clears a bit at given offset.
func (b *Bitmask) clear(offset int) *Bitmask {
	*b &= ^bit[offset]
	return b
}

func (b Bitmask) shift(offset int) Bitmask {
	if offset > 0 {
		return b << uint(offset)
	}

	return b >> -uint(offset)
}

func (b Bitmask) charm(offset int) (bitmask Bitmask) {
	count := b.count()

	for i := 0; i < count; i++ {
		pop := b ^ b.pop()
		b = b.pop()
		if (bit[i] & Bitmask(offset)).any() {
			bitmask |= pop
		}
	}

	return bitmask
}

func (b *Bitmask) fill(square, direction int, occupied, board Bitmask) *Bitmask {
	for bm := (bit[square] & board).shift(direction); bm.any(); bm = bm.shift(direction) {
		*b |= bm
		if (bm & occupied).any() {
			break
		}
		bm &= board
	}

	return b
}

func (b *Bitmask) spot(square, direction int, board Bitmask) *Bitmask {
	*b = ^((bit[square] & board).shift(direction))
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
			offset := row << 3 + col
			buffer.WriteByte(' ')
			if b.on(offset) {
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
