package lape

import (`testing`; `fmt`)

func TestPosition010(t *testing.T) {
        game := new(Game).Initialize().Setup(`Ke1,e2`, `Kg8,d7`)
        initial := new(Position).Initialize(game, game.pieces, WHITE, Bitmask(0))
        expect(t, fmt.Sprintf("%v", initial.enpassant),
`  a b c d e f g h  0x0000000000000000
8 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
7 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
6 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
5 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
4 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
3 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
2 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
1 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
`)

        position := initial.MakeMove(game, new(Move).Initialize(12, 28, Pawn(WHITE), Piece(0))) // e2-e4
        expect(t, fmt.Sprintf("%v", position.enpassant),
`  a b c d e f g h  0x0000000000100000
8 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
7 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
6 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
5 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
4 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
3 ⋅ ⋅ ⋅ ⋅ • ⋅ ⋅ ⋅
2 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
1 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
`)

        position = initial.MakeMove(game, new(Move).Initialize(51, 35, Pawn(BLACK), Piece(0))) // d7-d5
        // expect(t, fmt.Sprintf("0x%b", position.enpassant), `0x10000000000000000000000000000000000000000000`)
        expect(t, fmt.Sprintf("%v", position.enpassant),
`  a b c d e f g h  0x0000080000000000
8 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
7 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
6 ⋅ ⋅ ⋅ • ⋅ ⋅ ⋅ ⋅
5 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
4 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
3 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
2 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
1 ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅ ⋅
`)
}
