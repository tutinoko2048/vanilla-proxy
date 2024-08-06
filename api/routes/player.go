package routes

import (
	"strconv"
	"github.com/HyPE-Network/vanilla-proxy/proxy"
)

type PlayerListEntry struct {
	Name string `json:"name"`
	Xuid string `json:"xuid"`
	Uuid string `json:"uuid"`
	UniqueId string `json:"uniqueId"`
}

type PlayerHandler struct {
	proxyInstance *proxy.Proxy
}

func NewPlayerHandler(proxyInstance *proxy.Proxy) *PlayerHandler {
	return &PlayerHandler{proxyInstance: proxyInstance}
}

// return array of PlayerListEntry
func (handler *PlayerHandler) GetPlayerList() (playerList []PlayerListEntry, err error) {
	for _, player := range handler.proxyInstance.PlayerManager.GetPlayers() {
		playerList = append(playerList, PlayerListEntry{
			Name: player.GetName(),
			Xuid: player.GetXUID(),
			Uuid: player.GetUUID(),
			UniqueId: strconv.FormatInt(player.GetUniqueID(), 10),
		})
	}
	return playerList, nil
}

func (handler *PlayerHandler) SendToast(uniqueId string, title string, body string) (nil, err error) {
	player, err := handler.proxyInstance.PlayerManager.GetPlayerByUniqueID(uniqueId)
	if err != nil {
		return nil, err
	}
	player.SendToast(title, body)
	return nil, nil
}