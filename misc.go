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

// Returns n for the given the given row/column coordinate.
func Square(row, column int) int {
	return (row << 3) + column
}

func IsBetween(from, to, between int) bool {
        return ((maskStraight[from][to] | maskDiagonal[from][to]) & bit[between]) != 0
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
	return [8]int{ 8, 9, 1, -7, -8, -9, -1, 7 }[direction]
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
