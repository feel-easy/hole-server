package game

type playerController struct {
	player Player
}

func newPlayerController(player Player) *playerController {
	return &playerController{
		player: player,
	}
}

func (c *playerController) Name() string {
	return c.player.NickName()
}

func (c *playerController) ID() int {
	return c.player.PlayerID()
}
func (c *playerController) Player() *Player {
	return &c.player
}
