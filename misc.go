package donna

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
	A1 = iota
	    B1; C1; D1; E1; F1; G1; H1
	A2; B2; C2; D2; E2; F2; G2; H2
	A3; B3; C3; D3; E3; F3; G3; H3
	A4; B4; C4; D4; E4; F4; G4; H4
	A5; B5; C5; D5; E5; F5; G5; H5
	A6; B6; C6; D6; E6; F6; G6; H6
	A7; B7; C7; D7; E7; F7; G7; H7
	A8; B8; C8; D8; E8; F8; G8; H8
)

const (
        CASTLE_KING_WHITE  = Bitmask((1 << E1) | (1 << F1) | (1 << G1))
        CASTLE_KING_BLACK  = Bitmask((1 << E8) | (1 << F8) | (1 << G8))
        CASTLE_QUEEN_WHITE = Bitmask((1 << E1) | (1 << D1) | (1 << C1))
        CASTLE_QUEEN_BLACK = Bitmask((1 << E8) | (1 << D8) | (1 << C8))
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
func Col(n int) int {
	return n % 8 // n & 7
}

// Returns row and column numbers for the given bit index.
func Coordinate(n int) (int, int) {
        return Row(n), Col(n)
}

// Returns n for the given the given row/column coordinate.
func Square(row, column int) int {
	return (row << 3) + column
}

// Integer version of math/abs.
func Abs(n int) int {
        if n < 0 {
                return -n
        }
        return n
}

// Returns, as an integer, a non-negative pseudo-random number
// in [0, limit) range. It panics if limit <= 0.
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
