package lape

import (
        `fmt`
        `math`
)

const CENTER4 = 0x0000001818000000 // 4 central squares
const CENTER12 = 0x00003C3C3C3C0000 // 12 central squares

type Brain struct {
        player *Player
        color  int
}

func (b *Brain) Initialize(player *Player) *Brain {
        b.player = player
        b.color = player.Color

        return b
}

func (b *Brain) Evaluate(p *Position) float64 {
        x1, x2, x3, x4 := b.material(p), b.mobility(p), b.aggressiveness(p), b.center(p)
        fmt.Printf("Score for %s is %.2f (mat: %.2f, mob: %.2f, agg: %.2f, ctr: %.2f)\n", C(b.color), math.Abs(x1 + x2 + x3 + x4), x1, x2, x3, x4)
        return math.Abs(x1 + x2 + x3 + x4)
}

func (b *Brain) material(p *Position) float64 {
        opposite := b.color^1

        score := 1000 * (p.count[King(b.color)] - p.count[King(opposite)]) +
                    9 * (p.count[Queen(b.color)] - p.count[Queen(opposite)]) +
                    5 * (p.count[Rook(b.color)] - p.count[Rook(opposite)]) +
                    3 * (p.count[Bishop(b.color)] - p.count[Bishop(opposite)]) +
                    3 * (p.count[Knight(b.color)] - p.count[Knight(opposite)]) +
                    1 * (p.count[Pawn(b.color)] - p.count[Pawn(opposite)])

        return float64(score) + 0.1 * float64(p.count[Bishop(b.color)] - p.count[Bishop(opposite)])
}

func (b *Brain) mobility(p *Position) float64 {
        return 0.25 * float64(b.movesAvailable(p, b.color) - b.movesAvailable(p, b.color^1))
}

func (b *Brain) aggressiveness(p *Position) float64 {
        return 0.20 * float64(b.attacksAvailable(p, b.color) - b.attacksAvailable(p, b.color^1))
}

func (b *Brain) center(p *Position) (center float64) {
        for side := p.board[b.color]; !side.IsEmpty(); {
                index := side.FirstSet()
                sq12 := p.Targets(index).Intersect(CENTER12)
                sq12_count, sq4_count := sq12.Count(), sq12.Intersect(CENTER4).Count()
                center += 0.05 * float64(sq12_count - sq4_count) + 0.3 * float64(sq4_count)
                side.Clear(index)
        }
        return
}

func (b *Brain) movesAvailable(p *Position, color int) (moves int) {
        for side := p.board[color]; !side.IsEmpty(); {
                index := side.FirstSet()
                moves += p.Targets(index).Count()
                side.Clear(index)
        }
        return
}

func (b *Brain) attacksAvailable(p *Position, color int) (attacks int) {
        for side := p.board[color]; !side.IsEmpty(); {
                index := side.FirstSet()
                attacks += p.Targets(index).Intersect(p.board[color^1]).Count()
                side.Clear(index)
        }
        return
}
