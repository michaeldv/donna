package lape

import()

type Position struct {
        pieces  [64]Piece
        board   [2]Bitmask
        attacks [2]Bitmask
        kings   [2]Bitmask
        queens  [2]Bitmask
        rooks   [2]Bitmask
        bishops [2]Bitmask
        knights [2]Bitmask
        pawns   [2]Bitmask
}

func (p *Position)Initialize(g *Game) *Position {
        p.pieces = g.pieces

        p.setupBoard()
        return p
}

func (p *Position)Moves() []*Move {
        var moves []*Move
        
        moves = append(moves, new(Move).Initialize(Index(1,4), Index(3,4), Pawn(0), 0))

        return moves
}

func (p *Position)setupBoard() *Position {
        for i, piece := range p.pieces {
                if piece == 0 {
                        continue
                }
                color := piece.Color()
                p.board[color].Set(i)

                switch piece.Type() {
                case KING:
                        p.kings[color].Set(i)
                case QUEEN:
                        p.queens[color].Set(i)
                case ROOK:
                        p.rooks[color].Set(i)
                case BISHOP:
                        p.bishops[color].Set(i)
                case KNIGHT:
                        p.knights[color].Set(i)
                case PAWN:
                        p.pawns[color].Set(i)
                }
        }

        return p
}
