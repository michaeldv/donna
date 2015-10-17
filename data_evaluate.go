// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

const onePawn = 100
const unstoppablePawn = onePawn * 10
var (
	valuePawn      = Score{ onePawn *  1 +  0, onePawn *  1 + 29 }  //  100,  129
	valueKnight    = Score{ onePawn *  4 +  8, onePawn *  4 + 23 }  //  408,  423
	valueBishop    = Score{ onePawn *  4 + 18, onePawn *  4 + 28 }  //  418,  428
	valueRook      = Score{ onePawn *  6 + 35, onePawn *  6 + 40 }  //  635,  640
	valueQueen     = Score{ onePawn * 12 + 60, onePawn * 12 + 79 }  // 1260, 1279

	rightToMove    = Score{ 12,  5 }  // Tempo bonus.
	bishopPawn     = Score{  4,  6 }  // Penalty for each pawn on the same colored square as a bishop.
	bishopBoxed    = Score{ 73,  0 }  // Penalty for patterns like Bc1,d2,Nd3.
	bishopDanger   = Score{ 35,  0 }  // Bonus when king is under attack and sides have opposite-colored bishops.
	rookOnPawn     = Score{  6, 14 }  // Bonus for rook attacking a pawn.
	rookOnOpen     = Score{ 22, 10 }  // Bonus for rook on open file.
	rookOnSemiOpen = Score{ 10,  5 }  // Bonus for rook on semi-open file.
	rookOn7th      = Score{  5, 10 }  // Bonus for rook on 7th file.
	rookBoxed      = Score{ 45,  0 }  // Penalty for rook boxed by king.
	queenOnPawn    = Score{  2, 10 }  // Bonus for queen attacking a pawn.
	queenOn7th     = Score{  1,  4 }  // Bonus for queen on 7th rank.
	behindPawn     = Score{  8,  0 }  // Bonus for knight and bishop being behind friendly pawn.
	hangingAttack  = Score{  8, 12 }  // Bonus for attacking enemy pieces that are hanging.
	kingByPawn     = Score{  0,  8 }  // Penalty king being too far from friendly pawns.
	coverMissing   = Score{ 50,  0 }  // Penalty for missing cover pawn.
)

// Weight percentages applied to evaluation scores before computing the overall
// blended score.
var weights = []Score{
	{ 112, 134 }, 	// [0] Mobility.
	{  91,  79 }, 	// [1] Pawn structure.
	{  86, 107 }, 	// [2] Passed pawns.
	{ 126, 100 }, 	// [3] King safety.
	{ 126, 100 }, 	// [4] Enemy's king safety.
}

// Piece values for calculating most valueable victim/least valueable attacker,
// indexed by piece.
var pieceValue = [14]int{
	0, 0,
	valuePawn.midgame,   valuePawn.midgame,
	valueKnight.midgame, valueKnight.midgame,
	valueBishop.midgame, valueBishop.midgame,
	valueRook.midgame,   valueRook.midgame,
	valueQueen.midgame,  valueQueen.midgame,
	0, 0,
}

// Piece/square table: gets initilized on startup from the bonus arrays below.
var pst = [14][64]Score{{},}

var materialBalance = [14]int{
	0, 0,
	2*2*3*3*3*3*9,	  // Pawn
	2*2*3*3*3*3*9*9,  // Black Pawn
	2*2*3*3*3*3,	  // Knight
	2*2*3*3*3*3*3,	  // Black Knight
	2*2*3*3,	  // Bishop
	2*2*3*3*3,	  // Black Bishop
	2*2,	          // Rook
	2*2*3,	          // Black Rook
	1,	          // Queen
	1*2,	          // Black Queen
	0, 0,	          // Kings
}

// Piece/square bonus points, visually arranged from White's point of view. The
// square index is used directly for Black and requires a flip for White.
var bonusPawn = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	    0,    0,    0,    0,    0,    0,    0,    0,
	   -6,    8,   -4,   -2,   -2,   -4,    8,   -6,
	   -7,   -7,   -5,   -3,   -3,   -5,   -7,   -7,
	   -7,    0,   -1,    9,    9,   -1,    0,   -7,
	  -13,   -7,    8,   16,   16,    8,   -7,  -13,
	  -13,   -4,   10,   12,   12,   10,   -4,  -13,
	  -10,    1,    4,    2,    2,    4,    1,  -10,
	    0,    0,    0,    0,    0,    0,    0,    0,
	}, {
	    0,    0,    0,    0,    0,    0,    0,    0,
	    1,   -5,    1,    9,    9,    1,   -5,    1,
	    3,   -3,    1,    2,    2,    1,   -3,    3,
	    3,    5,    4,   -3,   -3,    4,    5,    3,
	    1,    2,   -4,   -2,   -2,   -4,    2,    1,
	   -3,   -3,    3,    2,    2,    3,   -3,   -3,
	    3,   -2,    4,   -1,   -1,    4,   -2,    3,
	    0,    0,    0,    0,    0,    0,    0,    0,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusKnight = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	  -98,  -33,  -21,  -15,  -15,  -21,  -33,  -98,
	  -31,   -9,    3,    7,    7,    3,   -9,  -31,
	   -6,   19,   28,   36,   36,   28,   19,   -6,
	  -13,    8,   19,   25,   25,   19,    8,  -13,
	  -13,    9,   22,   24,   24,   22,    9,  -13,
	  -36,  -11,    0,    5,    5,    0,  -11,  -36,
	  -42,  -22,  -11,   -5,   -5,  -11,  -22,  -42,
	  -72,  -48,  -40,  -37,  -37,  -40,  -48,  -72,
	}, {
	  -55,  -45,  -25,   -7,   -7,  -25,  -45,  -55,
	  -32,  -25,  -12,    7,    7,  -12,  -25,  -32,
	  -28,  -19,   -4,   14,   14,   -4,  -19,  -28,
	  -23,  -13,    1,   21,   21,    1,  -13,  -23,
	  -21,  -13,    4,   19,   19,    4,  -13,  -21,
	  -25,  -20,   -4,   14,   14,   -4,  -20,  -25,
	  -35,  -28,   -9,    5,    5,   -9,  -28,  -35,
	  -49,  -41,  -23,   -7,   -7,  -23,  -41,  -49,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusBishop = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	  -23,  -11,  -15,  -20,  -20,  -15,  -11,  -23,
	  -17,    4,   -2,   -6,   -6,   -2,    4,  -17,
	  -14,    3,    1,   -4,   -4,    1,    3,  -14,
	  -11,    7,    3,   -1,   -1,    3,    7,  -11,
	  -11,    9,    6,    0,    0,    6,    9,  -11,
	  -10,    9,    6,    1,    1,    6,    9,  -10,
	  -15,    5,    1,   -5,   -5,    1,    5,  -15,
	  -27,  -12,  -18,  -22,  -22,  -18,  -12,  -27,
	}, {
	  -33,  -21,  -23,  -14,  -14,  -23,  -21,  -33,
	  -22,  -11,  -11,   -2,   -2,  -11,  -11,  -22,
	  -18,   -7,   -5,    1,    1,   -5,   -7,  -18,
	  -18,   -7,   -9,    2,    2,   -9,   -7,  -18,
	  -18,   -7,   -8,    4,    4,   -8,   -7,  -18,
	  -16,   -5,   -7,    4,    4,   -7,   -5,  -16,
	  -22,   -9,  -12,   -3,   -3,  -12,   -9,  -22,
	  -34,  -20,  -23,  -14,  -14,  -23,  -20,  -34,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusRook = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	  -12,   -8,   -6,   -3,   -3,   -6,   -8,  -12,
	   -6,    2,    4,    6,    6,    4,    2,   -6,
	  -11,   -4,    0,    1,    1,    0,   -4,  -11,
	  -11,   -4,    0,    1,    1,    0,   -4,  -11,
	  -11,   -3,   -1,    1,    1,   -1,   -3,  -11,
	  -11,   -5,   -2,    1,    1,   -2,   -5,  -11,
	  -11,   -4,   -2,    0,    0,   -2,   -4,  -11,
	  -13,   -8,   -8,   -5,   -5,   -8,   -8,  -13,
	}, {
	    0,    0,    0,    0,    0,    0,    0,    0,
	    0,    0,    0,    0,    0,    0,    0,    0,
	    0,    0,    0,    0,    0,    0,    0,    0,
	    0,    0,    0,    0,    0,    0,    0,    0,
	    0,    0,    0,    0,    0,    0,    0,    0,
	    0,    0,    0,    0,    0,    0,    0,    0,
	    0,    0,    0,    0,    0,    0,    0,    0,
	    0,    0,    0,    0,    0,    0,    0,    0,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusQueen = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	   -1,   -2,   -1,    0,    0,   -1,   -2,   -1,
	   -1,    4,    4,    3,    3,    4,    4,   -1,
	   -1,    3,    4,    5,    5,    4,    3,   -1,
	   -2,    5,    4,    4,    4,    4,    5,   -2,
	   -1,    4,    5,    4,    4,    5,    4,   -1,
	   -1,    3,    5,    5,    5,    5,    3,   -1,
	   -2,    3,    5,    4,    4,    5,    3,   -2,
	    0,   -2,   -2,   -1,   -1,   -2,   -2,    0,
	}, {
	  -38,  -27,  -22,  -15,  -15,  -22,  -27,  -38,
	  -27,  -15,  -11,   -4,   -4,  -11,  -15,  -27,
	  -20,   -8,   -6,    2,    2,   -6,   -8,  -20,
	  -14,   -3,    5,   12,   12,    5,   -3,  -14,
	  -15,   -3,    5,    9,    9,    5,   -3,  -15,
	  -20,   -9,   -4,    3,    3,   -4,   -9,  -20,
	  -29,  -15,  -11,   -2,   -2,  -11,  -15,  -29,
	  -35,  -29,  -21,  -15,  -15,  -21,  -29,  -35,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusKing = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	   47,   61,   39,   16,   16,   39,   61,   47,
	   59,   80,   47,   24,   24,   47,   80,   59,
	   74,   95,   57,   35,   35,   57,   95,   74,
	   89,  104,   72,   47,   47,   72,  104,   89,
	  103,  107,   88,   69,   69,   88,  107,  103,
	  114,  137,  102,   69,   69,  102,  137,  114,
	  146,  166,  133,  104,  104,  133,  166,  146,
	  147,  174,  148,  111,  111,  148,  174,  147,
	}, {
	   15,   38,   51,   56,   56,   51,   38,   15,
	   36,   61,   72,   81,   81,   72,   61,   36,
	   59,   90,  100,   99,   99,  100,   90,   59,
	   67,   94,  113,  114,  114,  113,   94,   67,
	   66,   98,   98,  103,  103,   98,   98,   66,
	   55,   83,   98,   96,   96,   98,   83,   55,
	   35,   60,   86,   80,   80,   86,   60,   35,
	   14,   38,   52,   56,   56,   52,   38,   14,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusPassedPawn = [8]Score{
	{0, 0}, {0, 3}, {0, 7}, {17, 17}, {51, 35}, {102, 59}, {170, 91}, {0, 0},
}

var bonusSemiPassedPawn = [8]Score{
	{0, 0}, {3, 6}, {3, 6}, {7, 14}, {17, 34}, {41, 83}, {0, 0}, {0, 0},
}

var extraPassedPawn = [8]int{
	0, 0, 0, 1, 3, 6, 10, 0,
}

var extraKnight = [64]int{
     //vvvvvvvvvvvv Black vvvvvvvvvvvv
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
       21, 21, 21, 21, 21, 21, 21, 21,
       21, 21, 21, 21, 21, 21, 21, 21,
       21, 21, 21, 21, 21, 21, 21, 21,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
     //^^^^^^^^^^^^ White ^^^^^^^^^^^^
}

var extraBishop = [64]int{
     //vvvvvvvvvvvv Black vvvvvvvvvvvv
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	9,  9,  9,  9,  9,  9,  9,  9,
	9,  9,  9,  9,  9,  9,  9,  9,
	9,  9,  9,  9,  9,  9,  9,  9,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
     //^^^^^^^^^^^^ White ^^^^^^^^^^^^
}

// [1] Pawn, [2] Knight, [3] Bishop, [4] Rook, [5] Queen
var bonusMinorThreat = [6]Score{
	{0, 0}, {3, 18}, {12, 24}, {12, 24}, {20, 50}, {20, 50},
}

// [1] Pawn, [2] Knight, [3] Bishop, [4] Rook, [5] Queen
var bonusMajorThreat = [6]Score{
	{0, 0}, {7, 18}, {7, 22}, {7, 22}, {7, 22}, {12, 24},
}

// [1] Pawn, [2] Knight, [3] Bishop, [4] Rook, [5] Queen
var bonusKingThreat = [6]int {
	0, 0, 2, 2, 3, 5,
}

// [1] Pawn, [2] Knight, [3] Bishop, [4] Rook, [5] Queen
var bonusCloseCheck = [6]int {
	0, 0, 0, 0, 8, 12,
}

// [1] Pawn, [2] Knight, [3] Bishop, [4] Rook, [5] Queen
var bonusDistanceCheck = [6]int {
	0, 0, 1, 1, 4, 6,
}

var kingSafety = [64]int {
	  0,   0,   1,   2,   3,   5,   7,  10,
	 13,  16,  20,  24,  29,  34,  39,  45,
	 51,  58,  65,  72,  80,  88,  97, 106,
	115, 125, 135, 146, 157, 168, 180, 192,
	205, 218, 231, 245, 259, 274, 289, 304,
	319, 334, 349, 364, 379, 394, 409, 424,
	439, 454, 469, 484, 499, 514, 529, 544,
	559, 574, 589, 604, 619, 634, 640, 640,
}

// Supported pawn bonus arranged from White point of view. The actual score
// uses the same values for midgame and endgame.
var bonusSupportedPawn = [64]int{
      //vvvvvvvvvvvvv Black vvvvvvvvvvvv
	  0,  0,  0,  0,  0,  0,  0,  0,
	 62, 66, 66, 68, 68, 66, 66, 62,
	 31, 34, 34, 36, 36, 34, 34, 31,
	 13, 16, 16, 18, 18, 16, 16, 13,
	  4,  6,  6,  7,  7,  6,  6,  4,
	  1,  3,  3,  4,  4,  3,  3,  1,
	  0,  1,  1,  2,  2,  1,  1,  0,
	  0,  0,  0,  0,  0,  0,  0,  0,
     //^^^^^^^^^^^^^^ White ^^^^^^^^^^^^
}

// [1] Pawn, [2] Knight, [3] Bishop, [4] Rook, [5] Queen
var penaltyPawnThreat = [6]Score {
	{0, 0}, {0, 0}, {26, 35}, {26, 35}, {38, 49}, {43, 59},
}

// Penalty for doubled pawn: A to H, midgame/endgame.
var penaltyDoubledPawn = [8]Score{
	{7, 21}, {10, 24}, {12, 24}, {12, 24}, {12, 24}, {12, 24}, {10, 24}, {7, 21},
}

// Penalty for isolated pawn that is *not* exposed: A to H, midgame/endgame.
var penaltyIsolatedPawn = [8]Score{
	{12, 15}, {18, 17}, {20, 17}, {20, 17}, {20, 17}, {20, 17}, {18, 17}, {12, 15},
}

// Penalty for isolated pawn that is exposed: A to H, midgame/endgame.
var penaltyWeakIsolatedPawn = [8]Score{
	{18, 22}, {27, 26}, {30, 26}, {30, 26}, {30, 26}, {30, 26}, {27, 26}, {18, 22},
}

// Penalty for backward pawn that is *not* exposed: A to H, midgame/endgame.
var penaltyBackwardPawn = [8]Score{
	{10, 14}, {15, 16}, {17, 16}, {17, 16}, {17, 16}, {17, 16}, {15, 16}, {10, 14},
}

// Penalty for backward pawn that is exposed: A to H, midgame/endgame.
var penaltyWeakBackwardPawn = [8]Score{
	{15, 21}, {22, 23}, {25, 23}, {25, 23}, {25, 23}, {25, 23}, {22, 23}, {15, 21},
}

// Penalty for the weak king cover indexed by rank, midgame only.
var penaltyCover = [8]int {
	0, 0, 14, 38, 46, coverMissing.midgame, coverMissing.midgame, coverMissing.midgame,
}

var mobilityKnight = [9]Score{
	{-32, -25}, {-21, -15}, {-4, -5}, {1, 0}, {7, 5}, {13, 10}, {18, 14}, {21, 15}, {22, 16},
}

var mobilityBishop = [16]Score{
	{-26, -23}, {-14, -11}, { 3,  0}, {10,  7}, {17, 14}, {24, 21}, {30, 27}, {34, 31},
	{ 37,  34}, { 38,  36}, {40, 37}, {41, 38}, {42, 39}, {43, 40}, {43, 40}, {43, 40},
}

var mobilityRook = [16]Score{
	{-23, -26}, {-15, -13}, {-2,  0}, { 0,  8}, { 3, 16}, { 6, 24}, { 9, 32}, {11, 40},
	{ 13,  48}, { 14,  54}, {15, 57}, {16, 59}, {17, 61}, {18, 61}, {18, 62},
}

var mobilityQueen = [16]Score{
	{-21, -20}, {-14, -12}, {-2, -3}, { 0,  0}, { 3,  5}, { 5,  9}, { 6, 14}, { 9, 19},
	{ 10,  20}, { 10,  20}, {11, 20}, {11, 20}, {11, 20}, {12, 20}, {12, 20}, {12, 20},
}

// Boxed rooks.
var kingBoxA = [2]Bitmask{
	bit[D1]|bit[C1]|bit[B1], bit[D8]|bit[C8]|bit[B8],
}

var kingBoxH = [2]Bitmask{
	bit[E1]|bit[F1]|bit[G1], bit[E8]|bit[F8]|bit[G8],
}

var rookBoxA = [2]Bitmask{
	bit[A1]|bit[B1]|bit[C1], bit[A8]|bit[B8]|bit[C8],
}

var rookBoxH = [2]Bitmask{
	bit[H1]|bit[G1]|bit[F1], bit[H8]|bit[G8]|bit[F8],
}
