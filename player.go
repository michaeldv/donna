package lape

type Player struct {
        game    *Game  // The game we're playing.
        brain   *Brain // Brain evaluate positions.
	Color   int    // 0: white, 1: black
        Can00   bool   // Can castle king's side?
        Can000  bool   // Can castle queen's side?
}

func (p *Player)Initialize(game *Game, color int) *Player {
        p.game = game
        p.brain = new(Brain).Initialize(p)
        p.Color = color
        p.Can00 = true
        p.Can000 = true

        return p
}

func (p *Player) IsWhite() bool {
	return p.Color == 0
}

func (p *Player) IsBlack() bool {
	return p.Color != 0
}
