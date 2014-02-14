// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import()

type Magic struct {
	mask   Bitmask
	magic  Bitmask
}

var (
        kingMoves        [64]Bitmask
        knightMoves      [64]Bitmask
        pawnMoves        [2][64]Bitmask
	rookMagicMoves   [64][4096]Bitmask
	bishopMagicMoves [64][512]Bitmask

        maskPassed       [2][64]Bitmask
        maskInFront      [2][64]Bitmask

        // Complete file or rank mask if both squares reside on on the same file
        // or rank.
        maskStraight     [64][64]Bitmask

        // Complete diagonal mask if both squares reside on on the same diagonal.
        maskDiagonal     [64][64]Bitmask

        // If a king on square [x] gets checked from square [y] it can evade the
        // check from all squares except maskEvade[x][y]. For example, if white
        // king on B2 gets checked by black bishop on G7 the king can't step back
        // to A1 (despite not being attacked by black).
        maskEvade        [64][64]Bitmask

        // If a king on square [x] gets checked from square [y] the check can be
        // evaded by moving a piece to maskBlock[x][y]. For example, if white
        // king on B2 gets checked by black bishop on G7 the check can be evaded
        // by moving white piece onto C3-G7 diagonal (including capture on G7).
        maskBlock        [64][64]Bitmask

        // Bitmask to indicate pawn attacks for a square. For example, C3 is being
        // attacked by white pawns on B2 and D2, and black pawns on B4 and D4.
        maskPawn         [2][64]Bitmask
)

func init() {
	for square := A1; square <= H8; square++ {
                row, col := Coordinate(square)

		// Rooks.
		mask := createRookMask(square)
		bits := uint(mask.count())
		for i := 0; i < (1 << bits); i++ {
			bitmask := indexedBitmask(i, mask)
			index := (bitmask * rookMagic[square].magic) >> 52
			rookMagicMoves[square][index] = createRookAttacks(square, bitmask)
		}

		// Bishops.
		mask = createBishopMask(square)
                bits = uint(mask.count())
		for i := 0; i < (1 << bits); i++ {
			bitmask := indexedBitmask(i, mask)
			index := (bitmask * bishopMagic[square].magic) >> 55
			bishopMagicMoves[square][index] = createBishopAttacks(square, bitmask)
		}

                // Pawns.
                if row >= 1 && row <= 7 {
                        if col > 0 {
                                pawnMoves[White][square].set(Square(row+1, col-1))
                                pawnMoves[Black][square].set(Square(row-1, col-1))
                        }
                        if col < 7 {
                                pawnMoves[White][square].set(Square(row+1, col+1))
                                pawnMoves[Black][square].set(Square(row-1, col+1))
                        }
                }

                // Blocks, Evasions, Straight, Diagonals, Knights, and Kings.
                for i := A1; i <= H8; i++ {
			r, c := Coordinate(i)
			setMasks(square, i, row, col, r, c)

                        if i == square || Abs(i - square) > 17 {
                                continue // No king or knight can reach that far.
                        }
                        if (Abs(r - row) == 2 && Abs(c - col) == 1) || (Abs(r - row) == 1 && Abs(c - col) == 2) {
                                knightMoves[square].set(i)
                        }
                        if Abs(r - row) <= 1 && Abs(c - col) <= 1 {
                                kingMoves[square].set(i)
                        }
                }

                // Pawn attacks.
                if row > 1 { // White pawns can't attack first two ranks.
                        if col != 0 {
                                maskPawn[White][square] |= Bit(square - 9)
                        }
                        if col != 7 {
                                maskPawn[White][square] |= Bit(square - 7)
                        }
                }
                if row < 6 { // Black pawns can attack 7th and 8th ranks.
                        if col != 0 {
                                maskPawn[Black][square] |= Bit(square + 7)
                        }
                        if col != 7 {
                                maskPawn[Black][square] |= Bit(square + 9)
                        }
                }

                // Masks to check for passed pawns.
                if col > 0 {
                        maskPassed[White][square].fill(square - 1,  8, 0, 0x00FFFFFFFFFFFFFF)
                        maskPassed[Black][square].fill(square - 1, -8, 0, 0xFFFFFFFFFFFFFF00)
                }
                maskPassed[White][square].fill(square,  8, 0, 0x00FFFFFFFFFFFFFF)
                maskPassed[Black][square].fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
                if col < 7 {
                        maskPassed[White][square].fill(square + 1,  8, 0, 0x00FFFFFFFFFFFFFF)
                        maskPassed[Black][square].fill(square + 1, -8, 0, 0xFFFFFFFFFFFFFF00)
                }

                // Vertical squares in front of a pawn.
                maskInFront[White][square].fill(square,  8, 0, 0x00FFFFFFFFFFFFFF)
                maskInFront[Black][square].fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
	}
}

func indexedBitmask(index int, mask Bitmask) (bitmask Bitmask) {
	count := mask.count()

	for i, his := 0, mask; i < count; i++ {
		her := ((his - 1) & his) ^ his
		his &= his - 1
		if (1 << uint(i)) & index != 0 {
			bitmask |= her
		}
	}
	return
}

func createRookMask(square int) (bitmask Bitmask) {
	row, col := Coordinate(square)

	// North.
	for r := row + 1; r < 7; r++ {
		bitmask |= Bit(r * 8 + col)
	}
	// West.
	for c := col - 1; c > 0; c-- {
		bitmask |= Bit(row * 8 + c)
	}
	// South.
	for r := row - 1; r > 0; r-- {
		bitmask |= Bit(r * 8 + col)
	}
	// East.
	for c := col + 1; c < 7; c++ {
		bitmask |= Bit(row * 8 + c)
	}
	return
}

func createBishopMask(square int) (bitmask Bitmask) {
	row, col := Coordinate(square)

	// North West.
	for c, r := col - 1, row + 1; c > 0 && r < 7; c, r = c-1, r+1 {
		bitmask |= Bit(r * 8 + c)
	}
	// South West.
	for c, r := col - 1, row - 1; c > 0 && r > 0; c, r = c-1, r-1 {
		bitmask |= Bit(r * 8 + c)
	}
	// South East.
	for c, r := col + 1, row - 1; c < 7 && r > 0; c, r = c+1, r-1 {
		bitmask |= Bit(r * 8 + c)
	}
	// North East.
	for c, r := col + 1, row + 1; c < 7 && r < 7; c, r = c+1, r+1 {
		bitmask |= Bit(r * 8 + c)
	}
	return
}

func createRookAttacks(square int, mask Bitmask) (bitmask Bitmask) {
	row, col := Coordinate(square)

	// North.
	for c, r := col, row + 1; r <= 7; r++ {
                bit := Bit(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// East.
	for c, r := col + 1, row; c <= 7; c++ {
                bit := Bit(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// South.
	for c, r := col, row - 1; r >= 0; r-- {
                bit := Bit(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// West
	for c, r := col - 1, row; c >= 0; c-- {
                bit := Bit(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	return
}

func createBishopAttacks(square int, mask Bitmask) (bitmask Bitmask) {
	row, col := Coordinate(square)

	// North East.
	for c, r := col + 1, row + 1; c <= 7 && r <= 7; c, r = c+1, r+1 {
                bit := Bit(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// South East.
	for c, r := col + 1, row - 1; c <= 7 && r >= 0; c, r = c+1, r-1 {
                bit := Bit(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
        // South West.
	for c, r := col - 1, row - 1; c >= 0 && r >= 0; c, r = c-1, r-1 {
                bit := Bit(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
        // North West.
	for c, r := col - 1, row + 1; c >= 0 && r <= 7; c, r = c-1, r+1 {
                bit := Bit(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	return
}

func setMasks(square, target, row, col, r, c int) {
	if row == r {
		if col < c {
			maskBlock[square][target].fill(square, 1, Bit(target), maskFull)
			maskEvade[square][target] = ^((Bit(square) & ^maskFile[0]) >> 1)
		} else if col > c {
			maskBlock[square][target].fill(square, -1, Bit(target), maskFull)
			maskEvade[square][target] = ^((Bit(square) & ^maskFile[7]) << 1)
		}
                if col != c {
                        maskStraight[square][target] = maskRank[r]
                }
	} else if col == c {
		if row < r {
			maskBlock[square][target].fill(square, 8, Bit(target), maskFull)
			maskEvade[square][target] = ^((Bit(square) & ^maskRank[0]) >> 8)
		} else {
			maskBlock[square][target].fill(square, -8, Bit(target), maskFull)
			maskEvade[square][target] = ^((Bit(square) & ^maskRank[7]) << 8)
		}
                if row != r {
                        maskStraight[square][target] = maskFile[c]
                }
	} else if r + col == row + c { // Diagonals (A1->H8).
		if col < c {
			maskBlock[square][target].fill(square, 9, Bit(target), maskFull)
			maskEvade[square][target] = ^((Bit(square) & ^maskRank[0] & ^maskFile[0]) >> 9)
		} else {
			maskBlock[square][target].fill(square, -9, Bit(target), maskFull)
			maskEvade[square][target] = ^((Bit(square) & ^maskRank[7] & ^maskFile[7]) << 9)
		}
                if shift := (r - c) & 15; shift < 8 { // A1-A8-H8
                        maskDiagonal[square][target] = maskA1H8 << uint(8 * shift)
                } else { // B1-H1-H7
                        maskDiagonal[square][target] = maskA1H8 >> uint(8 * (16-shift))
                }
	} else if row + col == r + c { // AntiDiagonals (H1->A8).
		if col < c {
			maskBlock[square][target].fill(square, -7, Bit(target), maskFull)
			maskEvade[square][target] = ^((Bit(square) & ^maskRank[7] & ^maskFile[0]) << 7)
		} else {
			maskBlock[square][target].fill(square, 7, Bit(target), maskFull)
			maskEvade[square][target] = ^((Bit(square) & ^maskRank[0] & ^maskFile[7]) >> 7)
		}
                if shift := 7 ^ (r + c); shift < 8 { // A8-A1-H1
                        maskDiagonal[square][target] = maskH1A8 >> uint(8 * shift)
                } else { // B8-H8-H2
                        maskDiagonal[square][target] = maskH1A8 << uint(8 * (16-shift))
                }
	}
	//
	// Default values are all 0 for maskBlock[square][target] (Go sets it for us)
	// and all 1 for maskEvade[square][target].
	//
	if maskEvade[square][target] == 0 {
		maskEvade[square][target] = maskFull
	}
}
