// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import `math`

type Magic struct {
	mask  Bitmask
	magic Bitmask
}

var (
	kingMoves [64]Bitmask
	knightMoves [64]Bitmask
	pawnAttacks [2][64]Bitmask
	rookMagicMoves [64][4096]Bitmask
	bishopMagicMoves [64][512]Bitmask

	maskPassed [2][64]Bitmask
	maskInFront [2][64]Bitmask

	// Complete file, rank or diagonal mask if both squares reside on on the
	// same file, rank, or diagonal. For example, maskLine[C3][F6] has bits
	// set for the entire A1-H8 diagonal.
	maskLine [64][64]Bitmask

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

	// Bitmasks to detect unstoppable passers (pawn square rule).
	maskSquare[2][64]Bitmask   // King doesn't have the right to move.
	maskSquareEx[2][64]Bitmask // King has the right to move (border bits added).

	// Two arrays to simplify incremental polyglot hash computation.
	hashCastle [16]uint64
	hashEnpassant [8]uint64

	// Distance between two squares.
	distance [64][64]int

	// Late move reductions indexed by depth and move number.
	lateMoveReductions [64][64]int

	// Precomputed database of material imbalance scores, evaluation flags,
	// and endgame handlers. I wish they all could be California girls.
	materialBase [2*2*3*3*3*3*3*3*9*9]MaterialEntry
)

func init() {
	initMasks()
	initArrays()
	initPST()
	initMaterial()
}

func initMasks() {
	for sq := A1; sq <= H8; sq++ {
		row, col := coordinate(sq)

		// Distance, Blocks, Evasions, Lines, Knights, and Kings.
		for i := A1; i <= H8; i++ {
			r, c := coordinate(i)

			distance[sq][i] = max(abs(row - r), abs(col - c))
			setupMasks(sq, i, row, col, r, c)

			if i == sq || abs(i-sq) > 17 {
				continue // No king or knight can reach that far.
			}
			if (abs(r-row) == 2 && abs(c-col) == 1) || (abs(r-row) == 1 && abs(c-col) == 2) {
				knightMoves[sq].set(i)
			}
			if abs(r-row) <= 1 && abs(c-col) <= 1 {
				kingMoves[sq].set(i)
			}
		}

		// Rooks.
		mask := createRookMask(sq)
		bits := uint(mask.count())
		for i := 0; i < (1 << bits); i++ {
			bitmask := mask.magicify(i)
			index := (bitmask * rookMagic[sq].magic) >> 52
			rookMagicMoves[sq][index] = createRookAttacks(sq, bitmask)
		}

		// Bishops.
		mask = createBishopMask(sq)
		bits = uint(mask.count())
		for i := 0; i < (1 << bits); i++ {
			bitmask := mask.magicify(i)
			index := (bitmask * bishopMagic[sq].magic) >> 55
			bishopMagicMoves[sq][index] = createBishopAttacks(sq, bitmask)
		}

		// Pawns.
		if row >= A2H2 && row <= A7H7 {
			if col > 0 {
				pawnAttacks[White][sq].set(square(row + 1, col - 1))
				pawnAttacks[Black][sq].set(square(row - 1, col - 1))
			}
			if col < 7 {
				pawnAttacks[White][sq].set(square(row + 1, col + 1))
				pawnAttacks[Black][sq].set(square(row - 1, col + 1))
			}
		}

		// Pawn attacks.
		if row > 1 { // White pawns can't attack first two ranks.
			if col != 0 {
				maskPawn[White][sq] |= bit[sq-9]
			}
			if col != 7 {
				maskPawn[White][sq] |= bit[sq-7]
			}
		}
		if row < 6 { // Black pawns can attack 7th and 8th ranks.
			if col != 0 {
				maskPawn[Black][sq] |= bit[sq+7]
			}
			if col != 7 {
				maskPawn[Black][sq] |= bit[sq+9]
			}
		}

		// Vertical squares in front of a pawn.
		maskInFront[White][sq] = (maskBlock[sq][A8+col] | bit[A8+col]) & ^bit[sq]
		maskInFront[Black][sq] = (maskBlock[A1+col][sq] | bit[A1+col]) & ^bit[sq]

		// Masks to check for passed pawns.
		if col > 0 {
			maskPassed[White][sq] |= maskInFront[White][sq-1]
			maskPassed[Black][sq] |= maskInFront[Black][sq-1]
			maskPassed[White][sq-1] |= maskInFront[White][sq]
			maskPassed[Black][sq-1] |= maskInFront[Black][sq]
		}
		maskPassed[White][sq] |= maskInFront[White][sq]
		maskPassed[Black][sq] |= maskInFront[Black][sq]
	}
}

func initArrays() {

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

	// Late move reductions.
	for i := 0; i < 64; i++ {
		for j := 0; j < 64; j++ {
			value := math.Log1p(float64(i)) * math.Log1p(float64(j)) / 1.25 - 3.25  //\\1.6 - 2.3
			if value < 0.0 {
				value = 0.0
			}
			lateMoveReductions[i][j] = int(math.Floor(value))
		}
	}
}

func initPST() {
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
		pst[BlackPawn]  [square].sub(Score{bonusPawn  [0][square], bonusPawn  [1][square]}).sub(valuePawn)
		pst[BlackKnight][square].sub(Score{bonusKnight[0][square], bonusKnight[1][square]}).sub(valueKnight)
		pst[BlackBishop][square].sub(Score{bonusBishop[0][square], bonusBishop[1][square]}).sub(valueBishop)
		pst[BlackRook]  [square].sub(Score{bonusRook  [0][square], bonusRook  [1][square]}).sub(valueRook)
		pst[BlackQueen] [square].sub(Score{bonusQueen [0][square], bonusQueen [1][square]}).sub(valueQueen)
		pst[BlackKing]  [square].sub(Score{bonusKing  [0][square], bonusKing  [1][square]})
	}
}

func initMaterial() {
	var index int

	for wQ := 0; wQ < 2; wQ++ {
		for bQ := 0; bQ < 2; bQ++ {
			for wR := 0; wR < 3; wR++ {
				for bR := 0; bR < 3; bR++ {
					for wB := 0; wB < 3; wB++ {
						for bB := 0; bB < 3; bB++ {
							for wN := 0; wN < 3; wN++ {
								for bN := 0; bN < 3; bN++ {
									for wP := 0; wP < 9; wP++ {
										for bP := 0; bP < 9; bP++ {
		index = wQ * materialBalance[Queen]       +
			bQ * materialBalance[BlackQueen]  +
			wR * materialBalance[Rook]        +
			bR * materialBalance[BlackRook]   +
			wB * materialBalance[Bishop]      +
			bB * materialBalance[BlackBishop] +
			wN * materialBalance[Knight]      +
			bN * materialBalance[BlackKnight] +
			wP * materialBalance[Pawn]        +
			bP * materialBalance[BlackPawn]

		// Compute game phase and home turf values.
		materialBase[index].phase = 12 * (wN + bN + wB + bB) + 18 * (wR + bR) + 44 * (wQ + bQ)
		materialBase[index].turf = (wN + bN + wB + bB) * (wN + bN + wB + bB)

		// Set up evaluation flags and endgame handlers.
		materialBase[index].flags,
		materialBase[index].endgame = endgames(wP, wN, wB, wR, wQ, bP, bN, bB, bR, bQ)

		// Compute material imbalance scores.
		if wQ != bQ || wR != bR || wB != bB || wN != bN || wP != bP {
			white := imbalance(wB/2, wP, wN, wB, wR, wQ,  bB/2, bP, bN, bB, bR, bQ)
			black := imbalance(bB/2, bP, bN, bB, bR, bQ,  wB/2, wP, wN, wB, wR, wQ)

			adjustment := (white - black) / 32
			materialBase[index].score.midgame += adjustment
			materialBase[index].score.endgame += adjustment
		}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

// Simplified second-degree polynomial material imbalance by Tord Romstad.
func imbalance(w2, wP, wN, wB, wR, wQ, b2, bP, bN, bB, bR, bQ int) int {
	polynom := func(x, a, b, c int) int {
		return a * (x * x) + (b + c) * x
	}

	return polynom(w2,    0, (   0                                                                                    ),  1756) +
	       polynom(wP,    2, (  39*w2 +                                      37*b2                                    ),  -164) +
	       polynom(wN,   -4, (  35*w2 + 271*wP +                             10*b2 +  62*bP                           ), -1067) +
	       polynom(wB,    0, (   0*w2 + 105*wP +   4*wN +                    57*b2 +  64*bP +  39*bN                  ),  -160) +
	       polynom(wR, -141, ( -27*w2 +  -2*wP +  46*wN + 100*wB +           50*b2 +  40*bP +  23*bN + -22*bB         ),   234) +
	       polynom(wQ,    0, (-177*w2 +  25*wP + 129*wN + 142*wB + -137*wR + 98*b2 + 105*bP + -39*bN + 141*bB + 274*bR),  -137)
}

func endgames(wP, wN, wB, wR, wQ, bP, bN, bB, bR, bQ int) (flags uint8, endgame Function) {
	wMinor, wMajor := wN + wB, wR + wQ
	bMinor, bMajor := bN + bB, bR + bQ
	allMinor, allMajor := wMinor + bMinor, wMajor + bMajor

	noPawns := (wP + bP == 0)
	bareKing := ((wP + wMinor + wMajor) * (bP + bMinor + bMajor) == 0) // Bare king (white, black or both).

	// Set king safety flags if the opposing side has a queen and at least one piece.
	if wQ > 0 && (wN + wB + wR) > 0 {
		flags |= blackKingSafety
	}
	if bQ > 0 && (bN + bB + bR) > 0 {
		flags |= whiteKingSafety
	}

	// Insufficient material endgames that don't require further evaluation:
	// 1) Two bare kings.
	if wP + bP + allMinor + allMajor == 0 {
		flags |= materialDraw

	// 2) No pawns and king with a minor.
	} else if noPawns && allMajor == 0 && wMinor < 2 && bMinor < 2 {
		flags |= materialDraw

	// 3) No pawns and king with two knights.
	} else if noPawns && allMajor == 0 && allMinor == 2 && (wN == 2 || bN == 2) {
		flags |= materialDraw

	// Known endgame: king and a pawn vs. bare king.
	} else if wP + bP == 1 && allMinor == 0 && allMajor == 0 {
		flags |= knownEndgame
		endgame = (*Evaluation).kingAndPawnVsBareKing

	// Known endgame: king with a knight and a bishop vs. bare king.
	} else if noPawns && allMajor == 0 && ((wN == 1 && wB == 1) || (bN == 1 && bB == 1)) {
		flags |= knownEndgame
		endgame = (*Evaluation).knightAndBishopVsBareKing

	// Known endgame: two bishops vs. bare king.
	} else if noPawns && allMajor == 0 && ((wN == 0 && wB == 2) || (bN == 0 && bB == 2)) {
		flags |= knownEndgame
		endgame = (*Evaluation).twoBishopsVsBareKing

	// Known endgame: king with some winning material vs. bare king.
	} else if bareKing && allMajor > 0 {
		flags |= knownEndgame
		endgame = (*Evaluation).winAgainstBareKing

	// Lesser known endgame: king and two or more pawns vs. bare king.
	} else if bareKing && allMinor + allMajor == 0 && wP + bP > 1 {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).kingAndPawnsVsBareKing

	// Lesser known endgame: queen vs. rook with pawn(s)
	} else if (wP + wMinor + wR == 0 && wQ == 1 && bMinor + bQ == 0 && bP > 0 && bR == 1) ||
	          (bP + bMinor + bR == 0 && bQ == 1 && wMinor + wQ == 0 && wP > 0 && wR == 1) {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).queenVsRookAndPawns

	// Lesser known endgame: king and pawn vs. king and pawn.
	} else if allMinor + allMajor == 0 && wP == 1 && bP == 1 {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).kingAndPawnVsKingAndPawn

	// Lesser known endgame: bishop and pawn vs. bare king.
	} else if bareKing && allMajor == 0 && wN + bN == 0 && (wB * wP == 1 || bB * bP == 1) {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).bishopAndPawnVsBareKing

	// Lesser known endgame: rook and pawn vs. rook.
	} else if allMinor == 0 && wQ + bQ == 0 && wR + bR == 2 && wP + bP == 1 {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).rookAndPawnVsRook

	// Lesser known endgame: no pawns left.
	} else if (wP == 0 || bP == 0) && wMajor - bMajor == 0 && abs(wMinor - bMinor) <= 1 {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).noPawnsLeft

	// Lesser known endgame: single pawn with not a lot of material.
	} else if (wP == 1 || bP == 1) && wMajor - bMajor == 0 && abs(wMinor - bMinor) <= 1 {
		flags |= lesserKnownEndgame
		endgame = (*Evaluation).lastPawnLeft

	// Check for potential opposite-colored bishops.
	} else if wB * bB == 1 {
		flags |= singleBishops
		if allMajor == 0 && allMinor == 2 {
			flags |= lesserKnownEndgame
			endgame = (*Evaluation).bishopsAndPawns
		} else if flags & (whiteKingSafety | blackKingSafety) == 0 {
			flags |= lesserKnownEndgame
			endgame = (*Evaluation).drawishBishops
		}
	}

	return
}

func createRookMask(square int) (bitmask Bitmask) {
	r, c := coordinate(square)
	bitmask = (maskRank[r] | maskFile[c]) ^ bit[square]

	return *bitmask.trim(r, c)
}

func createBishopMask(square int) (bitmask Bitmask) {
	r, c := coordinate(square)

	if sq := square + 7; sq <= H8 && col(sq) == c - 1 {
		bitmask = maskLine[square][sq]
	} else if sq := square - 7; sq >= A1 && col(sq) == c + 1 {
		bitmask = maskLine[square][sq]
	}

	if sq := square + 9; sq <= H8 && col(sq) == c + 1 {
		bitmask |= maskLine[square][sq]
	} else if sq := square - 9; sq >= A1 && col(sq) == c - 1 {
		bitmask |= maskLine[square][sq]
	}
	bitmask ^= bit[square]

	return *bitmask.trim(r, c)
}

func createRookAttacks(sq int, mask Bitmask) (bitmask Bitmask) {
	row, col := coordinate(sq)

	// North.
	for r := row + 1; r <= 7; r++ {
		sq := square(r, col)
		bitmask.set(sq)
		if mask.on(sq) {
			break
		}
	}
	// East.
	for c := col + 1; c <= 7; c++ {
		sq := square(row, c)
		bitmask.set(sq)
		if mask.on(sq) {
			break
		}
	}
	// South.
	for r := row - 1; r >= 0; r-- {
		sq := square(r, col)
		bitmask.set(sq)
		if mask.on(sq) {
			break
		}
	}
	// West
	for c := col - 1; c >= 0; c-- {
		sq := square(row, c)
		bitmask.set(sq)
		if mask.on(sq) {
			break
		}
	}
	return
}

func createBishopAttacks(sq int, mask Bitmask) (bitmask Bitmask) {
	row, col := coordinate(sq)

	// North East.
	for c, r := col + 1, row + 1; c <= 7 && r <= 7; c, r = c+1, r+1 {
		sq := square(r, c)
		bitmask.set(sq)
		if mask.on(sq) {
			break
		}
	}
	// South East.
	for c, r := col + 1, row - 1; c <= 7 && r >= 0; c, r = c+1, r-1 {
		sq := square(r, c)
		bitmask.set(sq)
		if mask.on(sq) {
			break
		}
	}
	// South West.
	for c, r := col - 1, row - 1; c >= 0 && r >= 0; c, r = c-1, r-1 {
		sq := square(r, c)
		bitmask.set(sq)
		if mask.on(sq) {
			break
		}
	}
	// North West.
	for c, r := col - 1, row + 1; c >= 0 && r <= 7; c, r = c-1, r+1 {
		sq := square(r, c)
		bitmask.set(sq)
		if mask.on(sq) {
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
			maskLine[square][target] = maskRank[r]
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
			maskLine[square][target] = maskFile[c]
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
			maskLine[square][target] = maskA1H8 << uint(8*shift)
		} else { // B1-H1-H7
			maskLine[square][target] = maskA1H8 >> uint(8*(16-shift))
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
			maskLine[square][target] = maskH1A8 >> uint(8*shift)
		} else { // B8-H8-H2
			maskLine[square][target] = maskH1A8 << uint(8*(16-shift))
		}
	}

	// Default values are all 0 for maskBlock[square][target] (Go sets it for us)
	// and all 1 for maskEvade[square][target].
	if maskEvade[square][target] == 0 {
		maskEvade[square][target] = maskFull
	}

	// Pawn square rule masks.
	if square != target {

		// White king chasing black pawn.
		if row > 1 {
			if row <= r && abs(col - c) <= 7 - row {
				maskSquare[White][square].set(target)
			}
			if row <= r + 1 && abs(col - c) <= 8 - row {
				maskSquareEx[White][square].set(target)
			}
		} else if row == 1 {
			if row < r && abs(col - c) < 7 - row {
				maskSquare[White][square].set(target)
			}
			if row <= r && abs(col - c) < 8 - row {
				maskSquareEx[White][square].set(target)
			}
		}

		// Black king chasing white pawn.
		if row < 6 {
			if row >= r && abs(col - c) <= row {
				maskSquare[Black][square].set(target)
			}
			if row + 1 >= r && abs(col - c) <= row + 1 {
				maskSquareEx[Black][square].set(target)
			}
		} else if row == 6 {
			if row > r && abs(col - c) < row {
				maskSquare[Black][square].set(target)
			}
			if row >= r && abs(col - c) <= row {
				maskSquareEx[Black][square].set(target)
			}
		}
	}
}
