// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

import(`bytes`; `fmt`)

type Bitmask uint64

var deBruijn = [64]Square{
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
var msbLookup[256]Square

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

// Returns a bitmask with bit set at given square.
func bit(sq Square) Bitmask {
	return 1 << (uint(sq) & 63)
}

// Returns true if all bitmask bits are clear. Even if it's wrong, it's only
// off by a bit.
func (b Bitmask) noneʔ() bool {
	return b == 0
}

// Returns true if at least one bit is set.
func (b Bitmask) anyʔ() bool {
	return b != 0
}

// Returns true if a bit at given square is set.
func (b Bitmask) onʔ(sq Square) bool {
	return (b & bit(sq)).anyʔ()
}

// Returns true if a bit at given square is clear.
func (b Bitmask) offʔ(sq Square) bool {
	return !b.onʔ(sq)
}

// Returns true if a bitmask has single bit set.
func (b Bitmask) singleʔ() bool {
	return b.pop().noneʔ()
}

// Returns number of bits set.
func (b Bitmask) count() int {
	if b.noneʔ() {
		return 0
	}

	b -= ((b >> 1) & 0x5555555555555555)
	b =  ((b >> 2) & 0x3333333333333333) + (b & 0x3333333333333333)
	b =  ((b >> 4) + b) & 0x0F0F0F0F0F0F0F0F
	b += b >> 8
	b += b >> 16
	b += b >> 32

	return int(b) & 63
}

// Finds least significant bit set (LSB) in non-zero bitmask. Returns
// an integer in 0..63 range.
func (b Bitmask) first() Square {
	return deBruijn[((b ^ (b - 1)) * 0x03F79D71B4CB0A89) >> 58] & 63
}

// Eugene Nalimov's bitScanReverse: finds most significant bit set (MSB).
func (b Bitmask) last() (sq Square) {
	if b > 0xFFFFFFFF {
		b >>= 32; sq = 32
	}
	if b > 0xFFFF {
		b >>= 16; sq += 16
	}
	if b > 0xFF {
		b >>= 8; sq += 8
	}

	return sq + msbLookup[b]
}

func (b Bitmask) closest(color int) Square {
	if color == White {
		return b.first()
	}
	return b.last()
}

func (b Bitmask) farthest(color int) Square {
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

// Sets a bit at given square.
func (b *Bitmask) set(sq Square) *Bitmask {
	*b |= bit(sq)
	return b
}

// Clears a bit at given square.
func (b *Bitmask) clear(sq Square) *Bitmask {
	*b &= ^bit(sq)
	return b
}

func (b Bitmask) shift(offset int) Bitmask {
	if offset > 0 {
		return b << uint(offset)
	}

	return b >> -uint(offset)
}

func (b Bitmask) charm(sq Square) (bitmask Bitmask) {
	count := b.count()

	for i := 0; i < count; i++ {
		pop := b ^ b.pop()
		b = b.pop()
		if (bit(Square(i)) & Bitmask(sq)).anyʔ() {
			bitmask |= pop
		}
	}

	return bitmask
}

func (b *Bitmask) fill(sq Square, offset int, occupied, board Bitmask) *Bitmask {
	for bm := (bit(sq) & board).shift(offset); bm.anyʔ(); bm = bm.shift(offset) {
		*b |= bm
		if (bm & occupied).anyʔ() {
			break
		}
		bm &= board
	}

	return b
}

func (b *Bitmask) spot(sq Square, offset int, board Bitmask) *Bitmask {
	*b = ^((bit(sq) & board).shift(offset))
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
			sq := square(row, col)
			buffer.WriteByte(' ')
			if b.onʔ(sq) {
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
