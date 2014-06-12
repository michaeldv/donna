// Copyright (c) 2013-2014 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import ()

type Magic struct {
	mask  Bitmask
	magic Bitmask
}

var (
	kingMoves        [64]Bitmask
	knightMoves      [64]Bitmask
	pawnMoves        [2][64]Bitmask
	rookMagicMoves   [64][4096]Bitmask
	bishopMagicMoves [64][512]Bitmask

	maskPassed  [2][64]Bitmask
	maskInFront [2][64]Bitmask

	// Complete file or rank mask if both squares reside on on the same file
	// or rank.
	maskStraight [64][64]Bitmask

	// Complete diagonal mask if both squares reside on on the same diagonal.
	maskDiagonal [64][64]Bitmask

	// If a king on square [x] gets checked from square [y] it can evade the
	// check from all squares except maskEvade[x][y]. For example, if white
	// king on B2 gets checked by black bishop on G7 the king can't step back
	// to A1 (despite not being attacked by black).
	maskEvade [64][64]Bitmask

	// If a king on square [x] gets checked from square [y] the check can be
	// evaded by moving a piece to maskBlock[x][y]. For example, if white
	// king on B2 gets checked by black bishop on G7 the check can be evaded
	// by moving white piece onto C3-G7 diagonal (including capture on G7).
	maskBlock [64][64]Bitmask

	// Bitmask to indicate pawn attacks for a square. For example, C3 is being
	// attacked by white pawns on B2 and D2, and black pawns on B4 and D4.
	maskPawn [2][64]Bitmask

	// Two arrays to simplify incremental polyglot hash computation.
	hashCastle    [16]uint64
	hashEnpassant [8]uint64

	// Distance between two squares.
	distance [64][64]int
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
				pawnMoves[White][square].set(Square(row + 1, col - 1))
				pawnMoves[Black][square].set(Square(row - 1, col - 1))
			}
			if col < 7 {
				pawnMoves[White][square].set(Square(row + 1, col + 1))
				pawnMoves[Black][square].set(Square(row - 1, col + 1))
			}
		}

		// Distance, Blocks, Evasions, Straight, Diagonals, Knights, and Kings.
		for i := A1; i <= H8; i++ {
			r, c := Coordinate(i)

			distance[square][i] = Max(Abs(row - r), Abs(col - c))
			setupMasks(square, i, row, col, r, c)

			if i == square || Abs(i-square) > 17 {
				continue // No king or knight can reach that far.
			}
			if (Abs(r-row) == 2 && Abs(c-col) == 1) || (Abs(r-row) == 1 && Abs(c-col) == 2) {
				knightMoves[square].set(i)
			}
			if Abs(r-row) <= 1 && Abs(c-col) <= 1 {
				kingMoves[square].set(i)
			}
		}

		// Pawn attacks.
		if row > 1 { // White pawns can't attack first two ranks.
			if col != 0 {
				maskPawn[White][square] |= bit[square-9]
			}
			if col != 7 {
				maskPawn[White][square] |= bit[square-7]
			}
		}
		if row < 6 { // Black pawns can attack 7th and 8th ranks.
			if col != 0 {
				maskPawn[Black][square] |= bit[square+7]
			}
			if col != 7 {
				maskPawn[Black][square] |= bit[square+9]
			}
		}

		// Masks to check for passed pawns.
		if col > 0 {
			maskPassed[White][square].fill(square-1, 8, 0, 0x00FFFFFFFFFFFFFF)
			maskPassed[Black][square].fill(square-1, -8, 0, 0xFFFFFFFFFFFFFF00)
		}
		maskPassed[White][square].fill(square, 8, 0, 0x00FFFFFFFFFFFFFF)
		maskPassed[Black][square].fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
		if col < 7 {
			maskPassed[White][square].fill(square+1, 8, 0, 0x00FFFFFFFFFFFFFF)
			maskPassed[Black][square].fill(square+1, -8, 0, 0xFFFFFFFFFFFFFF00)
		}

		// Vertical squares in front of a pawn.
		maskInFront[White][square].fill(square, 8, 0, 0x00FFFFFFFFFFFFFF)
		maskInFront[Black][square].fill(square, -8, 0, 0xFFFFFFFFFFFFFF00)
	}

	// Castle hash values.
	for mask := uint8(0); mask < 16; mask++ {
		if mask & castleKingside[White] != 0 {
			hashCastle[mask] ^= polyglotRandomCastle[0]
		}
		if mask & castleQueenside[White] != 0 {
			hashCastle[mask] ^= polyglotRandomCastle[1]
		}
		if mask & castleKingside[Black] != 0 {
			hashCastle[mask] ^= polyglotRandomCastle[2]
		}
		if mask & castleQueenside[Black] != 0 {
			hashCastle[mask] ^= polyglotRandomCastle[3]
		}
	}

	// Enpassant hash values.
	for col := A1; col <= H1; col++ {
		hashEnpassant[col] = polyglotRandomEnpassant[col]
	}

	// Initialize PST.
	for square := A1; square <= H8; square++ {

		// White pieces: flip square index since bonus points have been
		// set up from black's point of view.
		flip := square ^ A8
		pst[Pawn]  [square].add(Score{bonusPawn  [0][flip], bonusPawn  [1][flip]}).add(valuePawn)
		pst[Knight][square].add(Score{bonusKnight[0][flip], bonusKnight[1][flip]}).add(valueKnight)
		pst[Bishop][square].add(Score{bonusBishop[0][flip], bonusBishop[1][flip]}).add(valueBishop)
		pst[Rook]  [square].add(Score{bonusRook  [0][flip], bonusRook  [1][flip]}).add(valueRook)
		pst[Queen] [square].add(Score{bonusQueen [0][flip], bonusQueen [1][flip]}).add(valueQueen)
		pst[King]  [square].add(Score{bonusKing  [0][flip], bonusKing  [1][flip]})

		// Black pieces: use square index as is, and assign negative
		// values so we could use white + black without extra condition.
		pst[BlackPawn]  [square].subtract(Score{bonusPawn  [0][square], bonusPawn  [1][square]}).subtract(valuePawn)
		pst[BlackKnight][square].subtract(Score{bonusKnight[0][square], bonusKnight[1][square]}).subtract(valueKnight)
		pst[BlackBishop][square].subtract(Score{bonusBishop[0][square], bonusBishop[1][square]}).subtract(valueBishop)
		pst[BlackRook]  [square].subtract(Score{bonusRook  [0][square], bonusRook  [1][square]}).subtract(valueRook)
		pst[BlackQueen] [square].subtract(Score{bonusQueen [0][square], bonusQueen [1][square]}).subtract(valueQueen)
		pst[BlackKing]  [square].subtract(Score{bonusKing  [0][square], bonusKing  [1][square]})
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
		bitmask |= bit[r * 8 + col]
	}
	// West.
	for c := col - 1; c > 0; c-- {
		bitmask |= bit[row * 8 + c]
	}
	// South.
	for r := row - 1; r > 0; r-- {
		bitmask |= bit[r * 8 + col]
	}
	// East.
	for c := col + 1; c < 7; c++ {
		bitmask |= bit[row * 8 + c]
	}
	return
}

func createBishopMask(square int) (bitmask Bitmask) {
	row, col := Coordinate(square)

	// North West.
	for c, r := col - 1, row + 1; c > 0 && r < 7; c, r = c-1, r+1 {
		bitmask |= bit[r * 8 + c]
	}
	// South West.
	for c, r := col - 1, row - 1; c > 0 && r > 0; c, r = c-1, r-1 {
		bitmask |= bit[r * 8 + c]
	}
	// South East.
	for c, r := col + 1, row - 1; c < 7 && r > 0; c, r = c+1, r-1 {
		bitmask |= bit[r * 8 + c]
	}
	// North East.
	for c, r := col + 1, row + 1; c < 7 && r < 7; c, r = c+1, r+1 {
		bitmask |= bit[r * 8 + c]
	}
	return
}

func createRookAttacks(square int, mask Bitmask) (bitmask Bitmask) {
	row, col := Coordinate(square)

	// North.
	for c, r := col, row + 1; r <= 7; r++ {
		bit := bit[r * 8 + c]
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// East.
	for c, r := col + 1, row; c <= 7; c++ {
		bit := bit[r * 8 + c]
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// South.
	for c, r := col, row - 1; r >= 0; r-- {
		bit := bit[r * 8 + c]
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// West
	for c, r := col - 1, row; c >= 0; c-- {
		bit := bit[r * 8 + c]
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
		bit := bit[r * 8 + c]
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// South East.
	for c, r := col + 1, row - 1; c <= 7 && r >= 0; c, r = c+1, r-1 {
		bit := bit[r * 8 + c]
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// South West.
	for c, r := col - 1, row - 1; c >= 0 && r >= 0; c, r = c-1, r-1 {
		bit := bit[r * 8 + c]
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	// North West.
	for c, r := col - 1, row + 1; c >= 0 && r <= 7; c, r = c-1, r+1 {
		bit := bit[r * 8 + c]
		bitmask |= bit
		if mask & bit != 0 {
			break
		}
	}
	return
}

func setupMasks(square, target, row, col, r, c int) {
	if row == r {
		if col < c {
			maskBlock[square][target].fill(square, 1, bit[target], maskFull)
			maskEvade[square][target].spot(square, -1, ^maskFile[0])
		} else if col > c {
			maskBlock[square][target].fill(square, -1, bit[target], maskFull)
			maskEvade[square][target].spot(square, 1, ^maskFile[7])
		}
		if col != c {
			maskStraight[square][target] = maskRank[r]
		}
	} else if col == c {
		if row < r {
			maskBlock[square][target].fill(square, 8, bit[target], maskFull)
			maskEvade[square][target].spot(square, -8, ^maskRank[0])
		} else {
			maskBlock[square][target].fill(square, -8, bit[target], maskFull)
			maskEvade[square][target].spot(square, 8, ^maskRank[7])
		}
		if row != r {
			maskStraight[square][target] = maskFile[c]
		}
	} else if r+col == row+c { // Diagonals (A1->H8).
		if col < c {
			maskBlock[square][target].fill(square, 9, bit[target], maskFull)
			maskEvade[square][target].spot(square, -9,  ^maskRank[0] & ^maskFile[0])
		} else {
			maskBlock[square][target].fill(square, -9, bit[target], maskFull)
			maskEvade[square][target].spot(square, 9, ^maskRank[7] & ^maskFile[7])
		}
		if shift := (r - c) & 15; shift < 8 { // A1-A8-H8
			maskDiagonal[square][target] = maskA1H8 << uint(8*shift)
		} else { // B1-H1-H7
			maskDiagonal[square][target] = maskA1H8 >> uint(8*(16-shift))
		}
	} else if row+col == r+c { // AntiDiagonals (H1->A8).
		if col < c {
			maskBlock[square][target].fill(square, -7, bit[target], maskFull)
			maskEvade[square][target].spot(square, 7, ^maskRank[7] & ^maskFile[0])
		} else {
			maskBlock[square][target].fill(square, 7, bit[target], maskFull)
			maskEvade[square][target].spot(square, -7, ^maskRank[0] & ^maskFile[7])
		}
		if shift := 7 ^ (r + c); shift < 8 { // A8-A1-H1
			maskDiagonal[square][target] = maskH1A8 >> uint(8*shift)
		} else { // B8-H8-H2
			maskDiagonal[square][target] = maskH1A8 << uint(8*(16-shift))
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
