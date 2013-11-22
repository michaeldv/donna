package lape

import ()

type Attack struct {
        Knight  [64]Bitmask
        Bishop  [64]Bitmask
        Rook    [64]Bitmask
        Queen   [64]Bitmask
        King    [64]Bitmask
}

func (b *Attack)Initialize() *Attack {
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

func (a *Attack) Targets(index int, piece Piece, sides [2]Bitmask) Bitmask {
        var bitmask Bitmask

        kind, color := piece.Kind(), piece.Color()

        switch kind {
        case PAWN:
                // Not yet.
        case KNIGHT:
                bitmask = a.Knight[index]
                bitmask.Exclude(sides[color])
        case BISHOP:
                bitmask = a.Bishop[index]
                bitmask.Trim(index, piece, sides)
        case ROOK:
                bitmask = a.Rook[index]
                bitmask.Trim(index, piece, sides)
        case QUEEN:
                bitmask = a.Queen[index]
                bitmask.Trim(index, piece, sides)
        case KING:
                // Not yet.
        }

        return bitmask
}
