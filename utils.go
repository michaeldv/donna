// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (
	`fmt`
	`math/rand`
	`time`
)

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

// Returns a bitmask with light or dark squares set matching the color of the
// square.
func SameAs(square int) Bitmask {
	if bit[square] & maskDark != 0 {
		return maskDark
	}
	return ^maskDark
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

func Min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// Formats time duration in milliseconds in human readable form: MM:SS.XXX
func ms(duration int64) string {
	mm := duration / 1000 / 60
	ss := duration / 1000 % 60
	xx := duration - mm * 1000 * 60 - ss * 1000
	return fmt.Sprintf(`%02d:%02d.%03ds`, mm, ss, xx)
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

func Summary(metrics map[string]interface{}) {
	phase := metrics[`Phase`].(int)
	tally := metrics[`PST`].(Score)
	material := metrics[`Imbalance`].(Score)
	final := metrics[`Final`].(Score)
	units := float32(onePawn)

	fmt.Println()
	fmt.Printf("Metric              MidGame        |        EndGame        | Blended\n")
	fmt.Printf("                W      B     W-B   |    W      B     W-B   |  (%d)  \n", phase)
	fmt.Printf("-----------------------------------+-----------------------+--------\n")
	fmt.Printf("%-12s    -      -    %5.2f  |    -      -    %5.2f  >  %5.2f\n", `PST`,
		float32(tally.midgame)/units, float32(tally.endgame)/units, float32(tally.blended(phase))/units)
	fmt.Printf("%-12s    -      -    %5.2f  |    -      -    %5.2f  >  %5.2f\n", `Imbalance`,
		float32(material.midgame)/units, float32(material.endgame)/units, float32(material.blended(phase))/units)

	for _, tag := range([]string{`Tempo`, `Threats`, `Pawns`, `Passers`, `Mobility`, `+Pieces`, `-Knights`, `-Bishops`, `-Rooks`, `-Queens`, `+King`, `-Cover`, `-Safety`}) {
		white := metrics[tag].(Total).white
		black := metrics[tag].(Total).black

		var score Score
		score.add(white).subtract(black)

		if tag[0:1] == `+` {
			tag = tag[1:]
		} else if tag[0:1] == `-` {
			tag = `  ` + tag[1:]
		}

		fmt.Printf("%-12s  %5.2f  %5.2f  %5.2f  |  %5.2f  %5.2f  %5.2f  >  %5.2f\n", tag,
			float32(white.midgame)/units, float32(black.midgame)/units, float32(score.midgame)/units,
			float32(white.endgame)/units, float32(black.endgame)/units, float32(score.endgame)/units,
			float32(score.blended(phase))/units)
	}
	fmt.Printf("%-12s    -      -    %5.2f  |    -      -    %5.2f  >  %5.2f\n\n", `Final Score`,
		float32(final.midgame)/units, float32(final.endgame)/units, float32(final.blended(phase))/units)
}

// Logging wrapper around fmt.Printf() that could be turned on as needed. Typical
// usage is Log(true); defer Log(false) in tests.
func Log(args ...interface{}) {
	switch len(args) {
	case 0:
		// Calling Log() with no arguments flips the logging setting.
		engine.log = !engine.log
		engine.fancy = !engine.fancy
	case 1:
		switch args[0].(type) {
		case bool:
			engine.log = args[0].(bool)
			engine.fancy = args[0].(bool)
		default:
			if engine.log {
				fmt.Println(args...)
			}
		}
	default:
		if engine.log {
			fmt.Printf(args[0].(string), args[1:]...)
		}
	}
}
