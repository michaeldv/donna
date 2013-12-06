package donna

import (`testing`; `fmt`)

func TestPosition010(t *testing.T) {
        game := NewGame().Setup(`Ke1,e2`, `Kg8,d7`)
        initial := NewPosition(game, game.pieces, WHITE, Bitmask(0))
        expect(t, fmt.Sprintf("0x%016X", uint64(initial.enpassant)), `0x0000000000000000`)

        position := initial.MakeMove(game, NewMove(E2, E4, Pawn(WHITE), Piece(0)))
        expect(t, fmt.Sprintf("0x%016X", uint64(position.enpassant)), `0x0000000000100000`)

        position = initial.MakeMove(game, NewMove(D7, D5, Pawn(BLACK), Piece(0)))
        expect(t, fmt.Sprintf("0x%016X", uint64(position.enpassant)), `0x0000080000000000`)
}

func TestPosition020(t *testing.T) {
        game := NewGame().Setup(`Ke1,b7,e4`, `Kg8,Ra8,h2`)
        p := NewPosition(game, game.pieces, WHITE, Bitmask(0))

        expect(t, fmt.Sprintf("%v", p.isPawnPromotion(Pawn(WHITE), A8)), `true`)
        expect(t, fmt.Sprintf("%v", p.isPawnPromotion(Pawn(WHITE), B8)), `true`)
        expect(t, fmt.Sprintf("%v", p.isPawnPromotion(Pawn(WHITE), E5)), `false`)
        expect(t, fmt.Sprintf("%v", p.isPawnPromotion(Pawn(BLACK), H1)), `true`)
}
