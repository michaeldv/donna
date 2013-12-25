// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`testing`)

// See test key values at http://hardy.uhasselt.be/Toga/book_format.html
func TestBook000(t *testing.T) {
        game := NewGame().InitialPosition()
        position := game.Start()
        book := NewBook(``)
        expect(t, book.polyglot(position), uint64(0x463B96181691FC9C))
}

func TestBook010(t *testing.T) { // 1. e4
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e7,f7,g7,h7`)
        position := game.Start()
        position.color ^= 1
        book := NewBook(``)
        expect(t, book.polyglot(position), uint64(0x823C9B50FD114196))
}

func TestBook020(t *testing.T) { // 1. e4 d5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d5,e7,f7,g7,h7`)
        position := game.Start()
        book := NewBook(``)
        expect(t, book.polyglot(position), uint64(0x0756B94461C50FB0))
}

func TestBook030(t *testing.T) { // 1. e4 d5 2. e5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e5,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d5,e7,f7,g7,h7`)
        position := game.Start()
        position.color ^= 1
        book := NewBook(``)
        expect(t, book.polyglot(position), uint64(0x662FAFB965DB29D4))
}

// func TestBook040(t *testing.T) { // TODO: 1. e4 d5 2. e5 f5 <-- Enpassant
//         game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e5,f2,g2,h2`,
//                                 `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d5,e7,f5,g7,h7`)
//         position := game.Start()
//         book := NewBook(``)
//         expect(t, book.polyglot(position), uint64(0x22A48B5A8E47FF78))
// }
// 
// func TestBook050(t *testing.T) { // TODO: 1. e4 d5 2. e5 f5 3. Ke2 <-- White Castle
//         game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke2,Bf1,Ng1,Rh1,a2,b2,c2,d2,e5,f2,g2,h2`,
//                                 `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d5,e7,f5,g7,h7`)
//         position := game.Start()
//         position.color ^= 1
//         book := NewBook(``)
//         expect(t, book.polyglot(position), uint64(0x652A607CA3F242C1))
// }
// 
// func TestBook060(t *testing.T) { // TODO: 1. e4 d5 2. e5 f5 3. Ke2 Kf7 <-- Black Castle
//         game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke2,Bf1,Ng1,Rh1,a2,b2,c2,d2,e5,f2,g2,h2`,
//                                 `Ra8,Nb8,Bc8,Qd8,Kf7,Bf8,Ng8,Rh8,a7,b7,c7,d5,e7,f5,g7,h7`)
//         position := game.Start()
//         book := NewBook(``)
//         expect(t, book.polyglot(position), uint64(0x652A607CA3F242C1))
// }
// 
// func TestBook070(t *testing.T) { // TODO: 1. a2a4 b7b5 2. h2h4 b5b4 3. c2c4 <-- Enpassant
//         game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a4,b2,c4,d2,e2,f2,g2,h4`,
//                                 `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b4,c7,d7,e7,f7,g7,h7`)
//         position := game.Start()
//         position.color ^= 1
//         book := NewBook(``)
//         expect(t, book.polyglot(position), uint64(0x00FDD303C946BDD9))
// }
// 
// func TestBook080(t *testing.T) { // TODO: 1. a2a4 b7b5 2. h2h4 b5b4 3. c2c4 b4xc3 4. Ra1a3 <-- Enpassant/Castle
//         game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a4,b2,d2,e2,f2,g2,h4`,
//                                 `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,c3,c7,d7,e7,f7,g7,h7`)
//         position := game.Start()
//         position.color ^= 1
//         book := NewBook(``)
//         expect(t, book.polyglot(position), uint64(0x5C3F9B829B279560))
// }


func TestBook100(t *testing.T) { // 1. e4 e5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d2,e4,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d7,e5,f7,g7,h7`)
        position := game.Start()
        book := NewBook(``)
        expect(t, book.polyglot(position), uint64(0x0844931A6EF4B9A0))
}

func TestBook110(t *testing.T) { // 1. d4 d5
        game := NewGame().Setup(`Ra1,Nb1,Bc1,Qd1,Ke1,Bf1,Ng1,Rh1,a2,b2,c2,d4,e2,f2,g2,h2`,
                                `Ra8,Nb8,Bc8,Qd8,Ke8,Bf8,Ng8,Rh8,a7,b7,c7,d5,e7,f7,g7,h7`)
        position := game.Start()
        book := NewBook(``)
        expect(t, book.polyglot(position), uint64(0x06649BA69B8C9FF8))
}
