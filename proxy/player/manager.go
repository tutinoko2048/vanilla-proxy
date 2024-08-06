package player

type PlayerManager struct {
	// [uuid]: Player
	Players map[string]*Player
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		Players: make(map[string]*Player),
	}
}

func (manager *PlayerManager) AddPlayer(player *Player) {
 	uuid := player.GetUUID()
	manager.Players[uuid] = player
}

func (manager *PlayerManager) RemovePlayer(player *Player) {
	uuid := player.GetUUID()
	delete(manager.Players, uuid)
}

func (manager *PlayerManager) GetPlayers() []*Player {
	var players []*Player
	for _, player := range manager.Players {
		players = append(players, player)
	}
	return players
}

func (manager *PlayerManager) GetPlayer(uuid string) *Player {
	return manager.Players[uuid]
}

func (manager *PlayerManager) GetPlayerByXUID(xuid string) *Player {
	for _, player := range manager.Players {
		if player.GetXUID() == xuid {
			return player
		}
	}
	return nil
}

func (manager *PlayerManager) GetPlayerCount() int {
	return len(manager.Players)
}
