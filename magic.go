// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

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
)

func init() {
	for square := A1; square <= H8; square++ {
                row, col := Coordinate(square)

		// Rooks.
		mask := createRookMask(square)
		bits := uint(mask.Count())
		for i := 0; i < (1 << bits); i++ {
			bitmask := indexedBitmask(i, mask)
			index := (bitmask * rookMagic[square].magic) >> 52
			rookMagicMoves[square][index] = createRookAttacks(square, bitmask)
		}

		// Bishops.
		mask = createBishopMask(square)
                bits = uint(mask.Count())
		for i := 0; i < (1 << bits); i++ {
			bitmask := indexedBitmask(i, mask)
			index := (bitmask * bishopMagic[square].magic) >> 55
			bishopMagicMoves[square][index] = createBishopAttacks(square, bitmask)
		}

                // Pawns.
                if row >= 1 && row <= 7 {
                        if col > 0 {
                                pawnMoves[WHITE][square].Set(Square(row+1, col-1))
                                pawnMoves[BLACK][square].Set(Square(row-1, col-1))
                        }
                        if col < 7 {
                                pawnMoves[WHITE][square].Set(Square(row+1, col+1))
                                pawnMoves[BLACK][square].Set(Square(row-1, col+1))
                        }
                }

                // Knights and Kings.
                for i := A1; i <= H8; i++ {
                        if i == square || Abs(i - square) > 17 {
                                continue
                        }

                        r, c := Coordinate(i)
                        if (Abs(r - row) == 2 && Abs(c - col) == 1) || (Abs(r - row) == 1 && Abs(c - col) == 2) {
                                knightMoves[square].Set(i)
                        }

                        if Abs(r - row) <= 1 && Abs(c - col) <= 1 {
                                kingMoves[square].Set(i)
                        }
                }

                // Masks to check for passed pawns.
                if col > 0 {
                        maskPassed[WHITE][square].Fill(square - 1,  8, 0, 0x00FFFFFFFFFFFFFF)
                        maskPassed[BLACK][square].Fill(square - 1, -8, 0, 0xFFFFFFFFFFFFFF00)
                }
                maskPassed[WHITE][square].Fill(square,  8, 0, 0x00FFFFFFFFFFFFFF)
                maskPassed[BLACK][square].Fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
                if col < 7 {
                        maskPassed[WHITE][square].Fill(square + 1,  8, 0, 0x00FFFFFFFFFFFFFF)
                        maskPassed[BLACK][square].Fill(square + 1, -8, 0, 0xFFFFFFFFFFFFFF00)
                }

                // Vertical squares in front of a pawn.
                maskInFront[WHITE][square].Fill(square,  8, 0, 0x00FFFFFFFFFFFFFF)
                maskInFront[BLACK][square].Fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
	}
}

func indexedBitmask(index int, mask Bitmask) (bitmask Bitmask) {
	count := mask.Count()

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
		bitmask |= Shift(r * 8 + col)
	}
	// West.
	for c := col - 1; c > 0; c-- {
		bitmask |= Shift(row * 8 + c)
	}
	// South.
	for r := row - 1; r > 0; r-- {
		bitmask |= Shift(r * 8 + col)
	}
	// East.
	for c := col + 1; c < 7; c++ {
		bitmask |= Shift(row * 8 + c)
	}
	return
}

func createBishopMask(square int) (bitmask Bitmask) {
	row, col := Coordinate(square)

	// North West.
	for c, r := col - 1, row + 1; c > 0 && r < 7; c, r = c-1, r+1 {
		bitmask |= Shift(r * 8 + c)
	}
	// South West.
	for c, r := col - 1, row - 1; c > 0 && r > 0; c, r = c-1, r-1 {
		bitmask |= Shift(r * 8 + c)
	}
	// South East.
	for c, r := col + 1, row - 1; c < 7 && r > 0; c, r = c+1, r-1 {
		bitmask |= Shift(r * 8 + c)
	}
	// North East.
	for c, r := col + 1, row + 1; c < 7 && r < 7; c, r = c+1, r+1 {
		bitmask |= Shift(r * 8 + c)
	}
	return
}

func createRookAttacks(square int, mask Bitmask) (bitmask Bitmask) {
	row, col := Coordinate(square)

	// North.
	for c, r := col, row + 1; r <= 7; r++ {
                bit := Shift(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// East.
	for c, r := col + 1, row; c <= 7; c++ {
                bit := Shift(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// South.
	for c, r := col, row - 1; r >= 0; r-- {
                bit := Shift(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// West
	for c, r := col - 1, row; c >= 0; c-- {
                bit := Shift(r * 8 + c)
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
                bit := Shift(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// South East.
	for c, r := col + 1, row - 1; c <= 7 && r >= 0; c, r = c+1, r-1 {
                bit := Shift(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
        // South West.
	for c, r := col - 1, row - 1; c >= 0 && r >= 0; c, r = c-1, r-1 {
                bit := Shift(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
        // North West.
	for c, r := col - 1, row + 1; c >= 0 && r <= 7; c, r = c-1, r+1 {
                bit := Shift(r * 8 + c)
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	return
}
