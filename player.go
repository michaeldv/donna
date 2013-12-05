package lape

type Player struct {
        game    *Game  // The game we're playing.
        brain   *Brain // Brain evaluate positions.
	Color   int    // 0: white, 1: black
        Can00   bool   // Can castle king's side?
        Can000  bool   // Can castle queen's side?
}

func NewPlayer(game *Game, color int) *Player {
        player := new(Player)

        player.game = game
        player.brain = NewBrain(player)
        player.Color = color
        player.Can00 = true
        player.Can000 = true

        return player
}

func (p *Player) IsWhite() bool {
	return p.Color == WHITE
}

func (p *Player) IsBlack() bool {
	return p.Color == BLACK
}
