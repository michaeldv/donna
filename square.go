// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

package donna

type Square int

// Returns 0..63 square number for the given row/column coordinate.
func square(row, col int) Square {
	return Square((row << 3) + col)
}

// Returns row number in 0..7 range for the given square.
func (sq Square) row() int {
	return int(sq) >> 3
}

// Returns column number in 0..7 range for the given square.
func (sq Square) col() int {
	return int(sq) & 7
}

// Returns both row and column numbers for the given square.
func (sq Square) coordinate() (int, int) {
	return sq.row(), sq.col()
}

// Returns relative rank for the square in 0..7 range. For example E2 is
// rank 1 for white and rank 6 for black.
func (sq Square) rank(color int) int {
	return sq.row() ^ (color * 7)
}

// Flips the square verically for white (ex. E2 becomes E7).
func (sq Square) flip(color int) Square {
	if color == White {
		return sq ^ 56
	}
	return sq
}

// Returns a bitmask with light or dark squares set matching the color of
// the given square.
func (sq Square) same() Bitmask {
	if maskDark.onʔ(sq) {
		return maskDark
	}

	return ^maskDark
}

// Resturns a distance between two squares.
func (sq Square) upto(square Square) int {
	return distance[sq][square]
}

// Returns a square in front of the given square.
func (sq Square) push(color int) Square {
	if color == White {
		return sq + 8
	}
	return sq - 8
}

// Returns a square two ranks in front of the given square.
func (sq Square) jump(color int) Square {
	if color == White {
		return sq + 16
	}
	return sq - 16
}

// Returns string representation of the square.
func (sq Square) str() string {
	return string(sq.col() + 'a') + string(sq.row() + '1')
}
