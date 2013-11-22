package lape

import (
        `bytes`
        //`fmt`
)

type Bitmask uint64

// Returns true if all bitmask bits are clear.
func (b Bitmask) IsEmpty() bool {
	return b == Bitmask(0)
}

// Returns true if all bitmask bits are set.
func (b Bitmask) IsFull() bool {
	return b == 0xFFFFFFFFFFFFFFFF
}

// Returns true if a bit at given position is set.
func (b Bitmask) IsSet(position int) bool {
	return b & (1 << uint(position)) != Bitmask(0)
}

// Returns true if a bit at given position is clear.
func (b Bitmask) IsClear(position int) bool {
	return !b.IsSet(position)
}

func (b Bitmask) String() string {
	buffer := bytes.NewBufferString("  a b c d e f g h\n")
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

// Sets a bit at given position.
func (b *Bitmask) Set(position int) {
	*b |= 1 << uint(position)
}

// Clears a bit at given position.
func (b *Bitmask) Clear(position int) {
	*b &= ^(1 << uint(position))
}

func (b *Bitmask) FirstSet() int {
        if *b == Bitmask(0) {
                return -1
        }
	return [64]int{ // 64-bit De Bruijn sequence.
		 0, 47,  1, 56, 48, 27,  2, 60,
		57, 49, 41, 37, 28, 16,  3, 61,
		54, 58, 35, 52, 50, 42, 21, 44,
		38, 32, 29, 23, 17, 11,  4, 62,
		46, 55, 26, 59, 40, 36, 15, 53,
		34, 51, 20, 43, 31, 22, 10, 45,
		25, 39, 14, 33, 19, 30,  9, 24,
		13, 18,  8, 12,  7,  6,  5, 63,
	}[((*b ^ (*b-1)) * 0x03F79D71B4CB0A89) >> 58]
}

// Combines two bitmasks using bitwise OR operator.
func (b *Bitmask) Combine(bitmask Bitmask) {
	*b |= bitmask
}

// Intersects two bitmasks using bitwise AND operator.
func (b *Bitmask) Intersect(bitmask Bitmask) {
	*b &= bitmask
}

// Mulitplies two bitmasks.
func (b *Bitmask) Multiply(bitmask Bitmask) {
	*b *= bitmask
}

// Excludes bits of one bitmask from another using bitwise XOR operator.
func (b *Bitmask) Exclude(bitmask Bitmask) {
	*b ^= (bitmask & *b)
}

// ...
func (b *Bitmask) Trim(index int, piece Piece, sides [2]Bitmask) {
        row, col := Row(index), Column(index)
        kind, color := piece.Kind(), piece.Color()
        same, opposite := sides[color], sides[color^1]
        if kind == ROOK || kind == QUEEN {
                b.trimLines(row, col, same, opposite)
        }
        if kind == BISHOP || kind == QUEEN {
                b.trimDiagonals(row, col, same, opposite)
        }
}

func (b *Bitmask) trimLines(row, col int, same, opposite Bitmask) {
        for c, clear := col-1, false;  c >= 0; c-- { // East.
                this, prev := Index(row, c), Index(row, c+1)
                if clear || same.IsSet(this) || (col < 7 && opposite.IsSet(prev)) {
                        b.Clear(this)
                        clear = true
                }
        }
        for r, clear := row+1, false;  r <= 7; r++ { // North.
                this, prev := Index(r, col), Index(r-1, col)
                if clear || same.IsSet(this) || (row > 0 && opposite.IsSet(prev)) {
                        b.Clear(this)
                        clear = true
                }
        }
        for c, clear := col+1, false;  c <= 7; c++ { // West.
                this, prev := Index(row, c), Index(row, c-1)
                if clear || same.IsSet(this) || (col > 0 && opposite.IsSet(prev)) {
                        b.Clear(this)
                        clear = true
                }
        }
        for r, clear := row-1, false;  r >= 0; r-- { // South.
                this, prev := Index(r, col), Index(r+1, col)
                if clear || same.IsSet(this) || (row < 7 && opposite.IsSet(prev)) {
                        b.Clear(this)
                        clear = true
                }
        }
}

func (b *Bitmask) trimDiagonals(row, col int, same, opposite Bitmask) {
        for r, c, clear := row-1, col-1, false;  r >= 0 && c >= 0; r, c = r-1, c-1 { // SW.
                this, prev := Index(r, c), Index(r+1, c+1)
                if clear || same.IsSet(this) || (row < 7 && col < 7 && opposite.IsSet(prev)) {
                        b.Clear(this)
                        clear = true
                }
        }
        for r, c, clear := row+1, col-1, false; r <= 7 && c >= 0; r, c = r+1, c-1 { // NW.
                this, prev := Index(r, c), Index(r-1, c+1)
                if clear || same.IsSet(this) || (row > 0 && col < 7 && opposite.IsSet(prev)) {
                        b.Clear(this)
                        clear = true
                }
        }
        for r, c, clear := row+1, col+1, false;  r <= 7 && c <= 7; r, c = r+1, c+1 { // NE.
                this, prev := Index(r, c), Index(r-1, c-1)
                if clear || same.IsSet(this) || (row > 0 && col > 0 && opposite.IsSet(prev)) {
                        b.Clear(this)
                        clear = true
                }
        }
        for r, c, clear := row-1, col+1, false;  r >= 0 && c <= 7; r, c = r-1, c+1 { // SE.
                this, prev := Index(r, c), Index(r+1, c-1)
                if clear || same.IsSet(this) || (row < 7 && col > 0 && opposite.IsSet(prev)) {
                        b.Clear(this)
                        clear = true
                }
        }
}

// Finds out row number of bit position.
func Row(position int) int {
	return position / 8 // position >> 3
}

// Finds out column number of bit position.
func Column(position int) int {
	return position % 8 // position & 7
}

// Finds out bit position for given row and column.
func Index(row, column int) int {
	return (row << 3) + column
}

func Abs(i int) int {
        if i >= 0 {
                return i
        } else {
                return -i
        }
}
