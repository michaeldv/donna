package donna

import (`encoding/binary`; `os`; `fmt`)

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

func (b *Book) PickMove(position *Position) (move *Move) {
        entries := b.lookup(position)

        for i, entry := range entries {
                fmt.Printf(" move: %d\n", i+1)
                fmt.Printf("  key: 0x%016X\n", entry.Key)
                fmt.Printf(" move: 0x%04X\n", entry.Move)
                fmt.Printf("score: 0x%04X\n", entry.Score)
                fmt.Printf("learn: 0x%08X\n", entry.Learn)
                fmt.Printf("%016b: to c/r: %d/%d from c/r %d/%d promo %d\n",
                        entry.Move, entry.toCol(), entry.toRow(), entry.fromCol(), entry.fromRow(), entry.promoted())
        }

        if len(entries) == 0 {
                // TODO: set the "useless" flag after a few misses.
                return nil
        }

        return b.move(position, entries[Random(len(entries))])
}

func (b *Book) lookup(position *Position) (entries []Entry) {
        var entry Entry

	file, err := os.Open(b.fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

        key := b.polyglot(position)
        first, middle, last := int64(0), int64(0), b.entries

        for {
                if last - first <= 1 {
                        return entries // <-- Nothing was found.
                }

                middle = (first + last) / 2

                if _, err := file.Seek(middle * 16, 0); err == nil {
                        binary.Read(file, binary.BigEndian, &entry)
                        if key == entry.Key {
                                entries = append(entries, entry)
                                break // <-- Found it!
                        } else if key < entry.Key {
                                last = middle
                        } else {
                                first = middle
                        }
                } else {
                        return entries // <-- Nothing was found.
                }
        }
        //
        // Go up and down from the current spot to pick up remaining book
        // entries with the same polyglot hash key.
        //
        for offset := int64(0); ; offset += 16 {
                if _, err := file.Seek(offset, 1); err == nil {
                        binary.Read(file, binary.BigEndian, &entry)
                        if key == entry.Key {
                                entries = append(entries, entry)
                                continue
                        }
                }
                break
        }
        //
        // Go back to the middle and proceed backwards in 16-byte increments.
        //
        file.Seek(middle * 16, 0)
        for offset := int64(-16); ; offset -= 32 {
                if _, err := file.Seek(offset, 1); err == nil {
                        binary.Read(file, binary.BigEndian, &entry)
                        if key == entry.Key {
                                entries = append(entries, entry)
                                continue
                        }
                }
                break
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

        if position.enpassant.IsNotEmpty() {
                col := Col(position.enpassant.FirstSet())
                key ^= polyglotRandom[772 + col]
        }
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

// Poluglot encodes "promotion piece" as follows:
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
