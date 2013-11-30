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

func (a *Attack) Targets(index int, p *Position)*Bitmask {
        var bitmask Bitmask
        piece := p.pieces[index]
        kind, color := piece.Kind(), piece.Color()

        switch kind {
        case PAWN:
                // Not yet.
        case KNIGHT:
                bitmask = a.Knight[index]
                bitmask.Exclude(p.board[color])
        case BISHOP:
                bitmask = a.Bishop[index]
                x1, x2, x3, x4 := a.DiagonalBlockers(index, p)
                bitmask.ClearFrom(x1, NorthEast).ClearFrom(x2, SouthEast).ClearFrom(x3, SouthWest).ClearFrom(x4, NorthWest)
        case ROOK:
                bitmask = a.Rook[index]
                x1, x2, x3, x4 := a.LineBlockers(index, p)
                bitmask.ClearFrom(x1, North).ClearFrom(x2, East).ClearFrom(x3, South).ClearFrom(x4, West)
        case QUEEN:
                bitmask = a.Queen[index]
		x1, x2, x3, x4 := a.LineBlockers(index, p)
		bitmask.ClearFrom(x1, North).ClearFrom(x2, East).ClearFrom(x3, South).ClearFrom(x4, West)
                x1, x2, x3, x4 = a.DiagonalBlockers(index, p)
                bitmask.ClearFrom(x1, NorthEast).ClearFrom(x2, SouthEast).ClearFrom(x3, SouthWest).ClearFrom(x4, NorthWest)
        case KING:
                bitmask = a.King[index]
                bitmask.Exclude(p.board[color]) // Exclude all squares occupied by the same color pieces.
        }

        return &bitmask
}

func (a *Attack) LineBlockers(index int, p *Position) (north, east, south, west int) {
        opposite := p.pieces[index].Color() ^ 1

	north = p.board[2].FirstSetFrom(index, North)
	if north >= 0 && p.board[opposite].IsSet(north) {
                if Row(north) != 7 {
                        north += Rose(North)
                } else {
                        north = -1
                }
	}
	east = p.board[2].FirstSetFrom(index, East)
	if east >= 0 && p.board[opposite].IsSet(east) {
                if Column(east) != 7 {
                        east += Rose(East)
                } else {
                        east = -1
                }
	}
	south = p.board[2].FirstSetFrom(index, South)
	if south >= 0 && p.board[opposite].IsSet(south) {
                if Row(south) != 0 {
		        south += Rose(South)
                } else {
                        south = -1
                }
	}
	west = p.board[2].FirstSetFrom(index, West)
	if west >= 0 && p.board[opposite].IsSet(west) {
                if Column(west) != 0 {
		        west += Rose(West)
                } else {
                        west = -1
                }
	}

	return
}

func (a *Attack) DiagonalBlockers(index int, p *Position) (northEast, southEast, southWest, northWest int) {
        opposite := p.pieces[index].Color() ^ 1

	northEast = p.board[2].FirstSetFrom(index, NorthEast)
	if northEast >= 0 && p.board[opposite].IsSet(northEast) {
                if Row(northEast) != 7 && Column(northEast) != 7 {
		        northEast += Rose(NorthEast)
                } else {
                        northEast = -1
                }
	}
	southEast = p.board[2].FirstSetFrom(index, SouthEast)
	if southEast >= 0 && p.board[opposite].IsSet(southEast) {
                if Row(southEast) != 0 && Column(southEast) != 7 {
		        southEast += Rose(SouthEast)
                } else {
                        southEast = -1
                }
	}
	southWest = p.board[2].FirstSetFrom(index, SouthWest)
	if southWest >= 0 && p.board[opposite].IsSet(southWest) {
                if Row(southWest) != 0 && Column(southWest) != 0 {
		        southWest += Rose(SouthWest)
                } else {
		        southWest = -1
                }
	}
	northWest = p.board[2].FirstSetFrom(index, NorthWest)
	if northWest >= 0 && p.board[opposite].IsSet(northWest) {
                if Row(northWest) != 7 && Column(northWest) != 0 {
		        northWest += Rose(NorthWest)
                } else {
		        northWest = -1
                }
	}

	return
}