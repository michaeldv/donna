// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`fmt`
	`math/rand`
	`time`
)

type Globals struct {
	Log   bool // Enable logging.
	Trace bool // Trace evaluation scores.
	Fancy bool // Represent pieces as UTF-8 characters.
}

var Settings Globals

// Returns row number for the given bit index.
func Row(n int) int {
	return n >> 3 // n / 8
}

// Returns column number for the given bit index.
func Col(n int) int {
	return n & 7 // n % 8
}

// Returns row and column numbers for the given bit index.
func Coordinate(n int) (int, int) {
	return Row(n), Col(n)
}

func RelRow(square, color int) int {
	return Row(square) ^ (color * 7)
}

// Returns 0..63 square number for the given row/column coordinate.
func Square(row, column int) int {
	return (row << 3) + column
}

func Flip(color, square int) int {
	if color == White {
		return square ^ 56
	}
	return square
}

// Returns bitmask with light or dark squares set, based on color of the square.
func Same(square int) Bitmask {
	return (bit[square] & maskDark) | (bit[square] & ^maskDark)
}

func IsBetween(from, to, between int) bool {
	return ((maskStraight[from][to] | maskDiagonal[from][to]) & bit[between]) != 0
}

func Ply() int {
	return node - rootNode
}

// Integer version of math/abs.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Returns, as an integer, a non-negative pseudo-random number
// in [0, limit) range. It panics if limit <= 0.
func Random(limit int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(limit)
}

func C(color int) string {
	return [2]string{`white`, `black`}[color]
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
	return [8]int{8, 9, 1, -7, -8, -9, -1, 7}[direction]
}


// Logging wrapper around fmt.Printf() that could be turned on as needed. Typical
// usage is Log(true); defer Log(false) in tests.
func Log(args ...interface{}) {
	switch len(args) {
	case 0:
		// Calling Log() with no arguments flips the logging setting.
		Settings.Log = !Settings.Log
		Settings.Fancy = !Settings.Fancy
	case 1:
		switch args[0].(type) {
		case bool:
			Settings.Log = args[0].(bool)
			Settings.Fancy = args[0].(bool)
		default:
			if Settings.Log {
				fmt.Println(args...)
			}
		}
	default:
		if Settings.Log {
			fmt.Printf(args[0].(string), args[1:]...)
		}
	}
}
