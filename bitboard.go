package lape

import (
)

type Bitboard struct {
        Rook   [64]Bitmask
        Knight [64]Bitmask
        Bishop [64]Bitmask
        Queen  [64]Bitmask
        King   [64]Bitmask
}

func (b *Bitboard)Initialize() *Bitboard {
        for i := 0;  i < 64;  i++ {
                row, col := Row(i), Column(i)
                for j := 0;  j < 64;  j++ {
                        r, c := Row(j), Column(j)
                        if r == row && c == col {
                                continue
                        }
                        if c == col || r == row {
                                b.Rook[i].Set(Index(r, c))
                                b.Queen[i].Set(Index(r, c))
                        }
                        if (Abs(r - row) == 2 && Abs(c - col) == 1) || (Abs(r - row) == 1 && Abs(c - col) == 2) {
                                b.Knight[i].Set(Index(r, c))
                        }
                        if Abs(r - row) == Abs(c - col) {
                                b.Bishop[i].Set(Index(r, c))
                                b.Queen[i].Set(Index(r, c))
                        }
                        if Abs(r - row) <= 1 && Abs(c - col) <= 1 {
                                b.King[i].Set(Index(r, c))
                        }
                }
        }

        return b
}
