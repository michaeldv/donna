package donna

import (`fmt`; `testing`)

func TestBook000(t *testing.T) {
        game := NewGame().InitialPosition()
        position := NewPosition(game, game.pieces, WHITE, Bitmask(0))
        book := NewBook("")
        expect(t, fmt.Sprintf(`0x%016X`, book.polyglot(position)), `0x463B96181691FC9C`)
}

func TestBook010(t *testing.T) { // 1. e4 e5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        position := NewPosition(game, game.pieces, WHITE, Bitmask(0))
        book := NewBook("")
        expect(t, fmt.Sprintf(`0x%016X`, book.polyglot(position)), `0x0844931A6EF4B9A0`)
}

func TestBook020(t *testing.T) { // 1. d4 d5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d4,e2,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d5,e7,f7,g7,h7`)
        position := NewPosition(game, game.pieces, WHITE, Bitmask(0))
        book := NewBook("")
        expect(t, fmt.Sprintf(`0x%016X`, book.polyglot(position)), `0x06649BA69B8C9FF8`)
}
