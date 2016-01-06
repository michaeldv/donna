// Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import(`github.com/michaeldv/donna/expect`; `testing`)

func openBook() (*Book, *Position) {
	return &Book{}, NewGame().start()
}

func polyglotEntry(source, target int) Entry {
	return Entry{Move: uint16(row(source)<<9) | uint16(col(source)<<6) |
		uint16(row(target)<<3) | uint16(col(target))}
}

// See test key values at http://hardy.uhasselt.be/Toga/book_format.html
func TestBook000(t *testing.T) {
	p := NewGame().start()
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x463B96181691FC9C))
	expect.Eq(t, pawnHash, uint64(0x37FC40DA841E1692))
}

func TestBook010(t *testing.T) { // 1. e4
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x823C9B50FD114196))
	expect.Eq(t, pawnHash, uint64(0x0B2D6B38C0B92E91))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook020(t *testing.T) { // 1. e4 d5
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x0756B94461C50FB0))
	expect.Eq(t, pawnHash, uint64(0x76916F86F34AE5BE))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook030(t *testing.T) { // 1. e4 d5 2. e5
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	p = p.makeMove(book.move(p, polyglotEntry(E4, E5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x662FAFB965DB29D4))
	expect.Eq(t, pawnHash, uint64(0xEF3E5FD1587346D3))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook040(t *testing.T) { // 1. e4 d5 2. e5 f5 <-- Enpassant
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	p = p.makeMove(book.move(p, polyglotEntry(E4, E5)))
	p = p.makeMove(book.move(p, polyglotEntry(F7, F5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x22A48B5A8E47FF78))
	expect.Eq(t, pawnHash, uint64(0x83871FE249DCEE04))
	expect.Eq(t, p.enpassant, uint8(F6))
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook050(t *testing.T) { // 1. e4 d5 2. e5 f5 3. Ke2 <-- White Castle
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	p = p.makeMove(book.move(p, polyglotEntry(E4, E5)))
	p = p.makeMove(book.move(p, polyglotEntry(F7, F5)))
	p = p.makeMove(book.move(p, polyglotEntry(E1, E2)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x652A607CA3F242C1))
	expect.Eq(t, pawnHash, uint64(0x83871FE249DCEE04))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.castles, castleKingside[Black]|castleQueenside[Black])
}

func TestBook060(t *testing.T) { // 1. e4 d5 2. e5 f5 3. Ke2 Kf7 <-- Black Castle
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	p = p.makeMove(book.move(p, polyglotEntry(E4, E5)))
	p = p.makeMove(book.move(p, polyglotEntry(F7, F5)))
	p = p.makeMove(book.move(p, polyglotEntry(E1, E2)))
	p = p.makeMove(book.move(p, polyglotEntry(E8, F7)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x00FDD303C946BDD9))
	expect.Eq(t, pawnHash, uint64(0x83871FE249DCEE04))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.castles, uint8(0))
}

func TestBook070(t *testing.T) { // 1. a2a4 b7b5 2. h2h4 b5b4 3. c2c4 <-- Enpassant
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(A2, A4)))
	p = p.makeMove(book.move(p, polyglotEntry(B7, B5)))
	p = p.makeMove(book.move(p, polyglotEntry(H2, H4)))
	p = p.makeMove(book.move(p, polyglotEntry(B5, B4)))
	p = p.makeMove(book.move(p, polyglotEntry(C2, C4)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x3C8123EA7B067637))
	expect.Eq(t, pawnHash, uint64(0xB5AA405AF42E7052))
	expect.Eq(t, p.enpassant, uint8(C3))
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook080(t *testing.T) { // 1. a2a4 b7b5 2. h2h4 b5b4 3. c2c4 b4xc3 4. Ra1a3 <-- Enpassant/Castle
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(A2, A4)))
	p = p.makeMove(book.move(p, polyglotEntry(B7, B5)))
	p = p.makeMove(book.move(p, polyglotEntry(H2, H4)))
	p = p.makeMove(book.move(p, polyglotEntry(B5, B4)))
	p = p.makeMove(book.move(p, polyglotEntry(C2, C4)))
	p = p.makeMove(book.move(p, polyglotEntry(B4, C3)))
	p = p.makeMove(book.move(p, polyglotEntry(A1, A3)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x5C3F9B829B279560))
	expect.Eq(t, pawnHash, uint64(0xE214F040EAA135A0))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.castles, castleKingside[White]|castleKingside[Black]|castleQueenside[Black])
}

func TestBook100(t *testing.T) { // 1. e4 e5
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(E7, E5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x0844931A6EF4B9A0))
	expect.Eq(t, pawnHash, uint64(0x798345D8FC7B53AE))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook110(t *testing.T) { // 1. d4 d5
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(D2, D4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x06649BA69B8C9FF8))
	expect.Eq(t, pawnHash, uint64(0x77A34D64090375F6))
	expect.Eq(t, p.enpassant, uint8(0))
	expect.Eq(t, p.castles, uint8(0x0F))
}
