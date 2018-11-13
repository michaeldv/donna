// Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.
//
// I am making my contributions/submissions to this project solely in my
// personal capacity and am not conveying any rights to any intellectual
// property of any third parties.

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

	expect.Eq(t, hash, uint64(0x1A734D7E0DD57DD8))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0xC6918785471EC42C))
	expect.Eq(t, pawnHash, p.pid)
}

func TestBook010(t *testing.T) { // 1. e4
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x16B3AA798A244C5C))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0x328746286FC870A1))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, 0)
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook020(t *testing.T) { // 1. e4 d5
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0xFFEE236E9AB1D24B))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0x230CE995D07A6BBF))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, 0)
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook030(t *testing.T) { // 1. e4 d5 2. e5
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	p = p.makeMove(book.move(p, polyglotEntry(E4, E5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0xD0E81E00A8B40DBD))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0xF4DCF2514D583140))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, 0)
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook040(t *testing.T) { // 1. e4 d5 2. e5 f5 <-- Enpassant
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	p = p.makeMove(book.move(p, polyglotEntry(E4, E5)))
	p = p.makeMove(book.move(p, polyglotEntry(F7, F5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x9796B5F5D1654751))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0x9B903D74CEBA05D7))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, F6)
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

	expect.Eq(t, hash, uint64(0x6AC160CCBC8F24CC))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0x9B903D74CEBA05D7))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, 0)
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

	expect.Eq(t, hash, uint64(0x7CA4772D3AC3FEE0))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0x9B903D74CEBA05D7))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, 0)
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

	expect.Eq(t, hash, uint64(0xF5A62AE51DD0F3FB))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0xD1A8556C4ABCA664))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, C3)
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

	expect.Eq(t, hash, uint64(0xB5C8873B19358920))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0xA1AF43655164A4B8))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, 0)
	expect.Eq(t, p.castles, castleKingside[White]|castleKingside[Black]|castleQueenside[Black])
}

func TestBook100(t *testing.T) { // 1. e4 e5
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(E2, E4)))
	p = p.makeMove(book.move(p, polyglotEntry(E7, E5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0x3F33F2AEB33E6F5B))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0xE3D13855F9F5D6AF))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, 0)
	expect.Eq(t, p.castles, uint8(0x0F))
}

func TestBook110(t *testing.T) { // 1. d4 d5
	book, p := openBook()
	p = p.makeMove(book.move(p, polyglotEntry(D2, D4)))
	p = p.makeMove(book.move(p, polyglotEntry(D7, D5)))
	hash, pawnHash := p.polyglot()

	expect.Eq(t, hash, uint64(0xC22ADACEB0B462F2))
	expect.Eq(t, hash, p.id)
	expect.Eq(t, pawnHash, uint64(0x1EC81035FA7FDB06))
	expect.Eq(t, pawnHash, p.pid)
	expect.Eq(t, p.enpassant, 0)
	expect.Eq(t, p.castles, uint8(0x0F))
}
