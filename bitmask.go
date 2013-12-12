package donna

import (
        `bytes`
        `fmt`
)

type Bitmask uint64

const (
        CASTLE_KING_WHITE  = Bitmask((1 << E1) | (1 << F1) | (1 << G1))
        CASTLE_KING_BLACK  = Bitmask((1 << E8) | (1 << F8) | (1 << G8))
        CASTLE_QUEEN_WHITE = Bitmask((1 << E1) | (1 << D1) | (1 << C1))
        CASTLE_QUEEN_BLACK = Bitmask((1 << E8) | (1 << D8) | (1 << C8))
)
var (
	DE_BRUJIN = [64]int{
		 0, 47,  1, 56, 48, 27,  2, 60,
		57, 49, 41, 37, 28, 16,  3, 61,
		54, 58, 35, 52, 50, 42, 21, 44,
		38, 32, 29, 23, 17, 11,  4, 62,
		46, 55, 26, 59, 40, 36, 15, 53,
		34, 51, 20, 43, 31, 22, 10, 45,
		25, 39, 14, 33, 19, 30,  9, 24,
		13, 18,  8, 12,  7,  6,  5, 63,
        }
)

// Returns true if all bitmask bits are clear.
func (b Bitmask) IsEmpty() bool {
	return b == Bitmask(0)
}

func (b Bitmask) IsNotEmpty() bool {
	return b != Bitmask(0)
}

// Returns true if a bit at given position is set.
func (b Bitmask) IsSet(position int) bool {
	return b & (1 << uint(position)) != Bitmask(0)
}

// Returns true if a bit at given position is clear.
func (b Bitmask) IsClear(position int) bool {
	return !b.IsSet(position)
}

// Sets a bit at given position.
func (b *Bitmask) Set(position int) *Bitmask {
	*b |= 1 << uint(position)
        return b
}

// Clears a bit at given position.
func (b *Bitmask) Clear(position int) *Bitmask {
	*b &= ^(1 << uint(position))
        return b
}

// Combines two bitmasks using bitwise OR operator.
func (b *Bitmask) Combine(bitmask Bitmask) *Bitmask {
	*b |= bitmask
        return b
}

// Intersects two bitmasks using bitwise AND operator.
func (b *Bitmask) Intersect(bitmask Bitmask) *Bitmask {
	*b &= bitmask
        return b
}

// Excludes bits of one bitmask from another using bitwise XOR operator.
func (b *Bitmask) Exclude(bitmask Bitmask) *Bitmask {
	*b ^= (bitmask & *b)
        return b
}

func (b *Bitmask) FirstSet() int {
        if *b == Bitmask(0) {
                return -1
        }
	return DE_BRUJIN[((*b ^ (*b-1)) * 0x03F79D71B4CB0A89) >> 58]
}

// Returns number of bits set.
func (b *Bitmask) Count() (count int) {
        if *b != Bitmask(0) {
                mask := *b                      // int count = 0;
                for ; mask != Bitmask(0); {     // while (x) {
                        count++                 //      count++;
                        mask &= mask -1         //      x &= x - 1; // reset LS1B
                }                               // }
        }                                       //
        return                                  // return count;
}

// ...
func (b *Bitmask) FirstSetFrom(square, direction int) int {
	rose := Rose(direction)
	for i, j := square, square+rose; Adjacent(i, j); i, j = j, j+rose {
		if b.IsSet(j) {
			return j
		}
	}

	return -1
}

func (b *Bitmask) ClearFrom (blocker, direction int) *Bitmask {
	rose := Rose(direction)
        for i, j := blocker, blocker; Adjacent(i, j); i, j = j, j+rose {
                b.Clear(j)
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
			if b.IsSet(position) {
				buffer.WriteString("\u2022") // Set
			} else {
				buffer.WriteString("\u22C5") // Clear
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}


// 0x0123456789ABCDEF
//   a b c d e f g h
// [0123456789ABCDEF]
// 8 X ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
// 7 X X ⋅ ⋅ ⋅ X ⋅ ⋅
// 6 X ⋅ X ⋅ ⋅ ⋅ X ⋅
// 5 X X X ⋅ ⋅ X X ⋅
// 4 X ⋅ ⋅ X ⋅ ⋅ ⋅ X
// 3 X X ⋅ X ⋅ X ⋅ X
// 2 X ⋅ X X ⋅ ⋅ X X
// 1 X X X X ⋅ X X X
//
// 7:  1   0   0   0   0   0   0   0
//     56  57  58  59  60  61  62  63
// 6:  1   1   0   0   0   1   0   0
//     48  49  50  51  52  53  54  55
// 5:  1   0   1   0   0   0   1   0
//     40  41  42  43  44  45  46  47
// 4:  1   1   1   0   0   1   1   0
//     32  33  34  35  36  37  38  39
// 3:  1   0   0   1   0   0   0   1
//     24  25  26  27  28  29  30  31
// 2:  1   1   0   1   0   1   0   1
//     16  17  18  19  20  21  22  23
// 1:  1   0   1   1   0   0   1   1
//     08  09  10  11  12  13  14  15
// 0:  1   1   1   1   0   1   1   1
//     00  01  02  03  04  05  06  07
//
