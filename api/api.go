// host a http server
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HyPE-Network/vanilla-proxy/api/routes"
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Data interface{} `json:"data"`
	Error bool `json:"error"`
	Message string `json:"message"`
}

type ProxyAPI struct {
	proxyInstance *proxy.Proxy
	playerHandler *routes.PlayerHandler
}

func Init(proxyInstance *proxy.Proxy) *ProxyAPI {
	api := ProxyAPI{
			proxyInstance: proxyInstance,
			playerHandler: routes.NewPlayerHandler(proxyInstance),
	}
	router := httprouter.New()

	router.GET("/player/list", api.handleRequest("GetPlayerList"))

	// サーバーをゴルーチンで開始
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil && err != http.ErrServerClosed {
			fmt.Println("ListenAndServe(): " + err.Error())
		}
	}()

	fmt.Println("API server started on port 8080")
	return &api
}

func (api ProxyAPI) handleRequest(action string) httprouter.Handle {
	return func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		response := Response{
			Data: nil,
			Error: true,
			Message: "Invalid action",
		}

		var data interface{}
		var err error

		switch action {
			case "GetPlayerList": data, err = api.playerHandler.GetPlayerList()
		}

		if err != nil {
			response.Error = true
			response.Message = err.Error()
		} else {
			response.Data = data
			response.Error = false
			response.Message = "Success"
		}

		json.NewEncoder(w).Encode(response)
	}
}


