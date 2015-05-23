// Copyright (c) 2014-2015 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

const onePawn = 100
var (
	valuePawn      = Score{ onePawn *  1 +  0, onePawn *  1 + 29 }  //  100,  129
	valueKnight    = Score{ onePawn *  4 +  8, onePawn *  4 + 23 }  //  408,  423
	valueBishop    = Score{ onePawn *  4 + 18, onePawn *  4 + 28 }  //  418,  428
	valueRook      = Score{ onePawn *  6 + 35, onePawn *  6 + 39 }  //  635,  639
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
	{ 100, 100 }, 	// [0] Mobility.
	{ 100, 100 }, 	// [1] Pawn structure.
	{ 100, 100 }, 	// [2] Passed pawns.
	{ 100, 100 }, 	// [3] King safety.
	{ 100, 100 }, 	// [4] Enemy's king safety.
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
	  -10,    0,    0,    0,    0,    0,    0,  -10,
	  -10,    0,    0,    0,    0,    0,    0,  -10,
	  -10,    0,    5,   10,   10,    5,    0,  -10,
	  -10,    0,   10,   20,   20,   10,    0,  -10,
	  -10,    0,    5,   10,   10,    5,    0,  -10,
	  -10,    0,    0,    0,    0,    0,    0,  -10,
	    0,    0,    0,    0,    0,    0,    0,    0,
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

var bonusKnight = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	 -101,  -33,  -21,  -15,  -15,  -21,  -33, -101,
	  -32,  -10,    3,    9,    9,    3,  -10,  -32,
	   -5,   18,   30,   36,   36,   30,   18,   -5,
	  -15,    8,   20,   26,   26,   20,    8,  -15,
	  -14,    9,   21,   27,   27,   21,    9,  -14,
	  -35,  -12,    0,    6,    6,    0,  -12,  -35,
	  -44,  -22,  -10,   -4,   -4,  -10,  -22,  -44,
	  -73,  -55,  -43,  -37,  -37,  -43,  -55,  -73,
	}, {
	  -49,  -42,  -26,   -8,   -8,  -26,  -42,  -49,
	  -34,  -27,  -11,    7,    7,  -11,  -27,  -34,
	  -27,  -19,   -3,   15,   15,   -3,  -19,  -27,
	  -21,  -14,    3,   20,   20,    3,  -14,  -21,
	  -21,  -14,    3,   20,   20,    3,  -14,  -21,
	  -27,  -19,   -3,   15,   15,   -3,  -19,  -27,
	  -34,  -27,  -11,    7,    7,  -11,  -27,  -34,
	  -49,  -42,  -26,   -8,   -8,  -26,  -42,  -49,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusBishop = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	  -25,  -11,  -15,  -19,  -19,  -15,  -11,  -25,
	  -16,    3,   -1,   -6,   -6,   -1,    3,  -16,
	  -14,    5,    1,   -4,   -4,    1,    5,  -14,
	  -11,    8,    4,   -1,   -1,    4,    8,  -11,
	  -10,    9,    6,    1,    1,    6,    9,  -10,
	  -10,    9,    5,    1,    1,    5,    9,  -10,
	  -15,    4,    1,   -4,   -4,    1,    4,  -15,
	  -27,  -14,  -17,  -22,  -22,  -17,  -14,  -27,
	}, {
	  -33,  -21,  -22,  -13,  -13,  -22,  -21,  -33,
	  -22,  -10,  -11,   -2,   -2,  -11,  -10,  -22,
	  -17,   -5,   -6,    3,    3,   -6,   -5,  -17,
	  -18,   -6,   -7,    2,    2,   -7,   -6,  -18,
	  -18,   -6,   -7,    2,    2,   -7,   -6,  -18,
	  -17,   -5,   -6,    3,    3,   -6,   -5,  -17,
	  -22,  -10,  -11,   -2,   -2,  -11,  -10,  -22,
	  -33,  -21,  -22,  -13,  -13,  -22,  -21,  -33,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusRook = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	  -11,   -9,   -6,   -4,   -4,   -6,   -9,  -11,
	   -6,    2,    5,    7,    7,    5,    2,   -6,
	  -11,   -4,   -1,    1,    1,   -1,   -4,  -11,
	  -11,   -4,   -1,    1,    1,   -1,   -4,  -11,
	  -11,   -4,   -1,    1,    1,   -1,   -4,  -11,
	  -11,   -4,   -1,    1,    1,   -1,   -4,  -11,
	  -11,   -4,   -1,    1,    1,   -1,   -4,  -11,
	  -11,   -9,   -6,   -4,   -4,   -6,   -9,  -11,
	}, {
	    2,    2,    2,    2,    2,    2,    2,    2,
	    2,    2,    2,    2,    2,    2,    2,    2,
	    2,    2,    2,    2,    2,    2,    2,    2,
	    2,    2,    2,    2,    2,    2,    2,    2,
	    2,    2,    2,    2,    2,    2,    2,    2,
	    2,    2,    2,    2,    2,    2,    2,    2,
	    2,    2,    2,    2,    2,    2,    2,    2,
	    2,    2,    2,    2,    2,    2,    2,    2,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusQueen = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	   -1,   -1,   -1,   -1,   -1,   -1,   -1,   -1,
	   -1,    4,    4,    4,    4,    4,    4,   -1,
	   -1,    4,    4,    4,    4,    4,    4,   -1,
	   -1,    4,    4,    4,    4,    4,    4,   -1,
	   -1,    4,    4,    4,    4,    4,    4,   -1,
	   -1,    4,    4,    4,    4,    4,    4,   -1,
	   -1,    4,    4,    4,    4,    4,    4,   -1,
	   -1,   -1,   -1,   -1,   -1,   -1,   -1,   -1,
	}, {
	  -40,  -27,  -21,  -15,  -15,  -21,  -27,  -40,
	  -27,  -15,   -9,   -3,   -3,   -9,  -15,  -27,
	  -21,   -9,   -3,    3,    3,   -3,   -9,  -21,
	  -15,   -3,    3,    9,    9,    3,   -3,  -15,
	  -15,   -3,    3,    9,    9,    3,   -3,  -15,
	  -21,   -9,   -3,    3,    3,   -3,   -9,  -21,
	  -27,  -15,   -9,   -3,   -3,   -9,  -15,  -27,
	  -40,  -27,  -21,  -15,  -15,  -21,  -27,  -40,
	}, //^^^^^^^^^^^^^^^^^^ White ^^^^^^^^^^^^^^^^^^
}

var bonusKing = [2][64]int{
	{  //vvvvvvvvvvvvvvvvvv Black vvvvvvvvvvvvvvvvvv
	   49,   67,   37,   13,   13,   37,   67,   49,
	   60,   77,   47,   23,   23,   47,   77,   60,
	   74,   91,   61,   37,   37,   61,   91,   74,
	   87,  105,   75,   51,   51,   75,  105,   87,
	   99,  116,   86,   62,   62,   86,  116,   99,
	  113,  130,  101,   76,   76,  101,  130,  113,
	  145,  162,  132,  108,  108,  132,  162,  145,
	  151,  168,  138,  114,  114,  138,  168,  151,
	}, {
	   14,   41,   54,   58,   58,   54,   41,   14,
	   37,   64,   78,   82,   82,   78,   64,   37,
	   56,   83,   97,  101,  101,   97,   83,   56,
	   68,   95,  109,  113,  113,  109,   95,   68,
	   68,   95,  109,  113,  113,  109,   95,   68,
	   56,   83,   97,  101,  101,   97,   83,   56,
	   37,   64,   78,   82,   82,   78,   64,   37,
	   14,   41,   54,   58,   58,   54,   41,   14,
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
	0,  2,  8,  8,  8,  8,  2,  0,
	0,  4, 13, 17, 17, 13,  4,  0,
	0,  2,  8, 13, 13,  8,  2,  0,
	0,  0,  2,  4,  4,  2,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
     //^^^^^^^^^^^^ White ^^^^^^^^^^^^
}

var extraBishop = [64]int{
     //vvvvvvvvvvvv Black vvvvvvvvvvvv
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  2,  4,  4,  4,  4,  2,  0,
	0,  5, 10, 10, 10, 10,  5,  0,
	0,  2,  5,  5,  5,  5,  2,  0,
	0,  0,  2,  2,  2,  2,  0,  0,
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
