package lape

import (
        `fmt`
        `math`
)

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
        x1, x2, x3 := b.material(p), b.mobility(p), b.aggressiveness(p)
        fmt.Printf("Score for %s is %.2f (mat: %.2f, mob: %.2f, agg: %.2f)\n", math.Abs(x1 + x2 + x3), C(b.color), x1, x2, x3)
        return math.Abs(x1 + x2 + x3)
}

func (b *Brain) material(p *Position) float64 {
        opposite := b.color^1

        score := 1000 * (p.count[KING|b.color] - p.count[KING|opposite]) +
                9 * (p.count[QUEEN|b.color] - p.count[QUEEN|opposite]) +
                5 * (p.count[ROOK|b.color] - p.count[ROOK|opposite]) +
                3 * (p.count[BISHOP|b.color] - p.count[BISHOP|opposite]) +
                3 * (p.count[KNIGHT|b.color] - p.count[KNIGHT|opposite]) +
                1 * (p.count[PAWN|b.color] - p.count[PAWN|opposite])

        return float64(score) + 0.1 * float64(p.count[BISHOP|b.color] - p.count[BISHOP|opposite])
}

func (b *Brain) mobility(p *Position) float64 {
        return 0.25 * float64(p.count[b.color] - p.count[b.color^1]) // TODO
}

func (b *Brain) aggressiveness(p *Position) float64 {
        return 0.2 * float64(p.count[b.color+100] - p.count[(b.color^1)+100]) // TODO
}
