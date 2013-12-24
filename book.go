// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package donna

import (`encoding/binary`; `os`)

type Book struct {
        fileName  string
        entries   int64
}

type Entry struct {
        Key    uint64
        Move   uint16 
        Score  uint16
        Learn  uint32
}

func NewBook(fileName string) *Book {
        book := new(Book)

        book.fileName = fileName
        if fi, err := os.Stat(book.fileName); err == nil {
                book.entries = fi.Size() / 16
        }

        return book
}

func (b *Book) pickMove(position *Position) (move *Move) {
        entries := b.lookup(position)
        if len(entries) == 0 {
                return nil // TODO: set the "useless book" flag after a few misses.
        }

        return b.move(position, entries[Random(len(entries))])
}

func (b *Book) lookup(position *Position) (entries []Entry) {
        var entry Entry

        file, err := os.Open(b.fileName)
        if err != nil {
                return
        }
        defer file.Close()

        key := b.polyglot(position)
        //
        // Since book entries are ordered by the polyglot key use binary
        // search to find *first* book entry that matches the position.
        //
        first, current, last := int64(-1), int64(0), b.entries
        for ; first < last; {
                current = (first + last) / 2
                file.Seek(current * 16, 0)
                binary.Read(file, binary.BigEndian, &entry)
                if key <= entry.Key {
                        last = current
                } else {
                        first = current + 1
                }
        }
        //
        // Read all book entries for the given position.
        //
        file.Seek(first * 16, 0)
        for ;; {
                binary.Read(file, binary.BigEndian, &entry)
                if key != entry.Key {
                        break
                } else {
                        entries = append(entries, entry)
                }
        }
        return
}

func (b *Book) move(p *Position, entry Entry) (move *Move) {
        from := Square(entry.fromRow(), entry.fromCol())
        to   := Square(entry.toRow(), entry.toCol())
        //
        // Check if this is a castle move. In Polyglot they are represented
        // as e1-h1, e1-a1, e8-h8, and e8-a8.
        //
        if from == E1 && to == H1 {
                to = G1
        } else if from == E1 && to == A1 {
                to = C1
        } else if from == E8 && to == H8 {
                to = G8
        } else if from == E8 && to == A8 {
                to = C8
        }

        move = NewMove(from, to, p.pieces[from], p.pieces[to])
        if promo := entry.promoted(); promo != 0 {
                move.Promote(promo)
        }
        return
}

func (b *Book) polyglot(position *Position) (key uint64) {
        for i, piece := range position.pieces {
                if piece != 0 {
                        key ^= polyglotRandom[0:768][64 * piece.Polyglot() + i]
                }
        }

	if position.game.players[WHITE].Can00 {
                key ^= polyglotRandom[768]
	}
	if position.game.players[WHITE].Can000 {
                key ^= polyglotRandom[769]
	}
	if position.game.players[BLACK].Can00 {
                key ^= polyglotRandom[770]
	}
	if position.game.players[BLACK].Can000 {
                key ^= polyglotRandom[771]
	}

        // TODO: 1. e2-e4 sets enpassant on e3 which causes invalid polyglot key.
        // if position.enpassant.IsNotEmpty() {
        //         col := Col(position.enpassant.FirstSet())
        //         key ^= polyglotRandom[772 + col]
        // }

	if position.color == WHITE {
                key ^= polyglotRandom[780]
	}

	return
}

func (e *Entry) toCol() int {
        return int(e.Move & 7)
}

func (e *Entry) toRow() int {
        return int((e.Move >> 3) & 7)
}

func (e *Entry) fromCol() int {
        return int((e.Move >> 6) & 7)
}

func (e *Entry) fromRow() int {
        return int((e.Move >> 9) & 7)
}

// Polyglot encodes "promotion piece" as follows:
//   knight  1 => 4
//   bishop  2 => 6
//   rook    3 => 8
//   queen   4 => 10
func (e *Entry) promoted() int {
        piece := int((e.Move >> 12) & 7)
        if piece == 0 {
                return piece
        }
        return piece * 2 + 2
}
