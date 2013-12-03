package lape

import (
        `fmt`
        `math/rand`
        `time`
)

const (
        WHITE = iota
        BLACK
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

type Globals struct {
        Log   bool // Enable logging.
        Fancy bool // Represent pieces as UTF-8 characters.
}

var Settings Globals

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

func C(color int) string {
        if color == 0 {
                return `white`
        } else if color == 1 {
                return `black`
        }
        return `Zebra?!`
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

// Returns index of the edge square following certain direction.
func Edge(index, direction int) (edge int) {
	if index < 0 || index > 63 {
		return -1
	}

	switch direction {
	case North:
		edge = (index % 8) + 56
	case East:
		edge = (Row(index) + 1) * 8 - 1
	case South:
		edge = index % 8
	case West:
		edge = Row(index) * 8
	case NorthEast:
		return NORTH_EAST_EDGE[index]
	case SouthEast:
		return SOUTH_EAST_EDGE[index]
	case SouthWest:
		return SOUTH_WEST_EDGE[index]
	case NorthWest:
		return NORTH_WEST_EDGE[index]
	}

	if index == edge {
		edge = -1
	}
	return
}

func Lop(args ...interface{}) {
        if Settings.Log {
                fmt.Println(args...)
        }
}

func Log(format string, args ...interface{}) {
        if Settings.Log {
                fmt.Printf(format, args...)
        }
}
