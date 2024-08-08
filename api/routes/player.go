package routes

import (
	"errors"
	"strconv"

	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type PlayerListEntry struct {
	Name string `json:"name"`
	Xuid string `json:"xuid"`
	Uuid string `json:"uuid"`
	UniqueId string `json:"uniqueId"`
	DeviceId string `json:"deviceId"`
	DeviceOS protocol.DeviceOS `json:"deviceOS"`
	DeviceModel string `json:"deviceModel"`
}

type PlayerHandler struct {
	proxyInstance *proxy.Proxy
}

func NewPlayerHandler(proxyInstance *proxy.Proxy) *PlayerHandler {
	return &PlayerHandler{proxyInstance: proxyInstance}
}

func (handler *PlayerHandler) Handle(action string, args map[string]interface{}) (data interface{}, err error) {
	switch action {
		case "getList": data, err = handler.GetPlayerList()
		case "sendToast":	data, err = handler.SendToast(args)
		case "sendBossbar": data, err = handler.SendBossbar(args)
		default: return nil, errors.New("invalid action")
	}
	return
}

// return array of PlayerListEntry
func (handler *PlayerHandler) GetPlayerList() (playerList []PlayerListEntry, err error) {
	for _, player := range handler.proxyInstance.PlayerManager.GetPlayers() {
		playerList = append(playerList, PlayerListEntry{
			Name: player.GetName(),
			Xuid: player.GetXUID(),
			Uuid: player.GetUUID(),
			UniqueId: strconv.FormatInt(player.GetUniqueID(), 10),
			DeviceId: player.GetDeviceID(),
			DeviceOS: player.GetSession().ClientData.DeviceOS,
			DeviceModel: player.GetSession().ClientData.DeviceModel,
		})
	}
	return playerList, nil
}

func (handler *PlayerHandler) SendToast(args map[string]interface{}) (nil, err error) {
	uniqueId, ok1 := args["uniqueId"].(string)
	title, ok2 := args["title"].(string)
	body, ok3 := args["body"].(string)
	if !ok1 || !ok2 || !ok3 {
		return nil, errors.New("wrong arguments")
	}

	player, err := handler.proxyInstance.PlayerManager.GetPlayerByUniqueID(uniqueId)
	if err != nil {
		return nil, err
	}
	player.SendToast(title, body)
	return nil, nil
}

func (handler * PlayerHandler) SendBossbar(args map[string]interface{}) (nil, err error) {
	uniqueId, ok1 := args["uniqueId"].(string)
	title, ok2 := args["title"].(string)
	_percentage, ok3 := args["percentage"].(float64)
	_color, ok4 := args["color"].(float64)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, errors.New("wrong arguments")
	}
	percentage := float32(_percentage)
	color := uint32(_color)

	player, err := handler.proxyInstance.PlayerManager.GetPlayerByUniqueID(uniqueId)
	if err != nil {
		return nil, err
	}
	player.SetBossbar(title, percentage, color)
	return nil, nil
}

func (handler *PlayerHandler) ClearBossbar(args map[string]interface{}) (nil, err error) {
	uniqueId, ok := args["uniqueId"].(string)
	if !ok {
		return nil, errors.New("wrong arguments")
	}
	player, err := handler.proxyInstance.PlayerManager.GetPlayerByUniqueID(uniqueId)
	if err != nil {
		return nil, err
	}
	player.ClearBossbar()
	return nil, nil
}
