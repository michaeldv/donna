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

func (a *Attack) Targets(index int, piece Piece, board [3]Bitmask) (bitmask Bitmask) {
        kind, color := piece.Kind(), piece.Color()

        switch kind {
        case PAWN:
                // Not yet.
        case KNIGHT:
                bitmask = a.Knight[index]
                bitmask.Exclude(board[color])
        case BISHOP:
                bitmask = a.Bishop[index]
		x1, x2, x3, x4 := a.DiagonalBlockers(index, color^1, board)
		bitmask.ClearFrom(x1, NorthEast).ClearFrom(x2, SouthEast).ClearFrom(x3, SouthWest).ClearFrom(x4, NorthWest)
        case ROOK:
                bitmask = a.Rook[index]
		x1, x2, x3, x4 := a.LineBlockers(index, color^1, board)
		bitmask.ClearFrom(x1, North).ClearFrom(x2, East).ClearFrom(x3, South).ClearFrom(x4, West)
        case QUEEN:
                bitmask = a.Queen[index]
		x1, x2, x3, x4 := a.LineBlockers(index, color^1, board)
		bitmask.ClearFrom(x1, North).ClearFrom(x2, East).ClearFrom(x3, South).ClearFrom(x4, West)
		x1, x2, x3, x4 = a.DiagonalBlockers(index, color^1, board)
		bitmask.ClearFrom(x1, NorthEast).ClearFrom(x2, SouthEast).ClearFrom(x3, SouthWest).ClearFrom(x4, NorthWest)
        case KING:
                // Not yet.
        }

        return
}

func (a *Attack) LineBlockers(index, opposite int, board [3]Bitmask) (north, east, south, west int) {
	north = board[2].FirstSetFrom(index, North)
	if north >= 0 && board[opposite].IsSet(north) {
		north += Rose(North)
	}
	east = board[2].FirstSetFrom(index, East)
	if east >= 0 && board[opposite].IsSet(east) {
		east += Rose(East)
	}
	south = board[2].FirstSetFrom(index, South)
	if south >= 0 && board[opposite].IsSet(south) {
		south += Rose(South)
	}
	west = board[2].FirstSetFrom(index, West)
	if west >= 0 && board[opposite].IsSet(west) {
		west += Rose(West)
	}

	return
}

func (a *Attack) DiagonalBlockers(index, opposite int, board [3]Bitmask) (northEast, southEast, southWest, northWest int) {
	northEast = board[2].FirstSetFrom(index, NorthEast)
	if northEast >= 0 && board[opposite].IsSet(northEast) {
		northEast += Rose(NorthEast)
	}
	southEast = board[2].FirstSetFrom(index, SouthEast)
	if southEast >= 0 && board[opposite].IsSet(southEast) {
		southEast += Rose(SouthEast)
	}
	southWest = board[2].FirstSetFrom(index, SouthWest)
	if southWest >= 0 && board[opposite].IsSet(southWest) {
		southWest += Rose(SouthWest)
	}
	northWest = board[2].FirstSetFrom(index, NorthWest)
	if northWest >= 0 && board[opposite].IsSet(northWest) {
		northWest += Rose(West)
	}

	return
}
