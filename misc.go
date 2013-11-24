package lape

import (
        `math/rand`
        `time`
)

const (
	North = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

// Returns row number for the given bit index.
func Row(n int) int {
	return n / 8 // n >> 3
}

// Returns column number for the given bit index.
func Column(n int) int {
	return n % 8 // n & 7
}

// Returns row and column numbers for the given bit index.
func Coordinate(n int) (int, int) {
        return Row(n), Column(n)
}

// Returns n for the given the given row/column coordinate.
func Index(row, column int) int {
	return (row << 3) + column
}

// Integer version of math/abs.
func Abs(n int) int {
        if n < 0 {
                return -n
        }
        return n
}

func Random(limit int) int {
        rand.Seed(time.Now().Unix())
        return rand.Intn(limit)
}
//
//   noWe         nort         noEa
//           +7    +8    +9
//               \  |  /
//   west    -1 <-  0 -> +1    east
//               /  |  \
//           -9    -8    -7
//   soWe         sout         soEa
//
func Rose(direction int) int {
	return []int{ 8, 9, 1, -7, -8, -9, -1, 7 }[direction]
}

func Adjacent(index, target int) bool {
        if target < 0 || target > 63 {
                return false
        }
        row, col := Coordinate(index)
        x, y := Coordinate(target)

        return Abs(row-x) <= 1 && Abs(col-y) <= 1
}
