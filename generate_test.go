package donna

import (`testing`)//; `fmt`)

// Old tests for naive move genarator modernized to test the incremental one.
///////////////////////////////////////////////////////////////////////////////
func TestGenerate000(t *testing.T) {
        game := NewGame().InitialPosition()
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        // All possible moves in the initial position, pawn-to-queen, left-to right, unsorted.
        expect(t, gen.allMoves(), `[a2-a3 a2-a4 b2-b3 b2-b4 c2-c3 c2-c4 d2-d3 d2-d4 e2-e3 e2-e4 f2-f3 f2-f4 g2-g3 g2-g4 h2-h3 h2-h4 Nb1-a3 Nb1-c3 Ng1-f3 Ng1-h3]`)
}

func TestGenerate020(t *testing.T) {
        game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `a3,b4,c5,e7,f6,g5,h4,Kg8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        // All possible moves, left-to right, unsorted.
        expect(t, gen.allMoves(), `[d2-d3 d2-d4]`)
}

func TestGenerate030(t *testing.T) {
        game := NewGame().Setup(`a2,e4,g2`, `b3,f5,f3,h3,Kg8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        // All possible moves, left-to right, unsorted.
        expect(t, gen.allMoves(), `[a2-a3 a2xb3 a2-a4 g2xf3 g2-g3 g2xh3 g2-g4 e4-e5 e4xf5]`)
}

// Should not include castles when rook has moved.
func TestGenerate040(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rf1,g2`, `Ke8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

func TestGenerate050(t *testing.T) {
        game := NewGame().Setup(`Ke1,Rb1,b2`, `Ke8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when king has moved.
func TestGenerate060(t *testing.T) {
        game := NewGame().Setup(`Kf1,Ra1,a2,Rh1,h2`, `Ke8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when rooks are not there.
func TestGenerate070(t *testing.T) {
        game := NewGame().Setup(`Ke1`, `Ke8`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when king is in check.
func TestGenerate080(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Re7`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when target square is a capture.
func TestGenerate090(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Nc1,Ng1`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}

// Should not include castles when king is to jump over attacked square.
func TestGenerate100(t *testing.T) {
        game := NewGame().Setup(`Ke1,Ra1,Rf1`, `Ke8,Bc4,Bf4`)
        gen := game.Start(White).StartMoveGen(1).GenerateMoves()

        doesNotContain(t, gen.allMoves(), `0-0`)
}
// func TestGenerate010(t *testing.T) {
//         Settings.Log = false
//         Settings.Fancy = false
//         game := NewGame().Setup(`a2,b3,c4,d2,e6,f5,g4,h3`, `Kg8`)
//         moves := game.Start(White).Moves(0)
//
//         // Moves should be sorted by relative strength.
//         expect(t, moves, `[h3-h4 a2-a4 e6-e7 a2-a3 g4-g5 f5-f6 b3-b4 c4-c5 d2-d4 d2-d3]`)
// }
//
// // LVA/MVV capture ordering.
// func TestGenerate110(t *testing.T) {
//         game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5`)
//         captures := game.Start(White).Captures(0)
//
//         expect(t, captures, `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Qh5xd5 Kd4xd5]`)
// }
//
// func TestGenerate120(t *testing.T) {
//         game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5`)
//         captures := game.Start(White).Captures(0)
//
//         expect(t, captures, `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5]`)
// }
//
// func TestGenerate130(t *testing.T) {
//         game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6`)
//         captures := game.Start(White).Captures(0)
//
//         expect(t, captures, `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6]`)
// }
//
// func TestGenerate140(t *testing.T) {
//         game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6,Nh3`)
//         captures := game.Start(White).Captures(0)
//
//         expect(t, captures, `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6 Nf4xh3 Qh5xh3]`)
// }
//
// func TestGenerate150(t *testing.T) {
//         game := NewGame().Setup(`Kd4,e4,Nf4,Bc4,Ra5,Qh5`, `Kd8,Qd5,Rf5,Bg6,Nh3,e2`)
//         captures := game.Start(White).Captures(0)
//
//         expect(t, captures, `[e4xd5 Nf4xd5 Bc4xd5 Ra5xd5 Kd4xd5 e4xf5 Qh5xf5 Nf4xg6 Qh5xg6 Nf4xh3 Qh5xh3 Nf4xe2 Bc4xe2 Qh5xe2]`)
// }
