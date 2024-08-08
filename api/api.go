// host a http server
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/HyPE-Network/vanilla-proxy/api/routes"
	"github.com/HyPE-Network/vanilla-proxy/proxy"
	"github.com/julienschmidt/httprouter"
)

type Request struct {
	Type string `json:"type"`
	Action string `json:"action"`
	Args map[string]interface{} `json:"args"`
}

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

	router.POST("/", api.handleRequest)

	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil && err != http.ErrServerClosed {
			fmt.Println("ListenAndServe(): " + err.Error())
		}
	}()

	fmt.Println("API server has started on port 8080")
	return &api
}

func (api ProxyAPI) handleRequest(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// read body
	body := make([]byte, r.ContentLength)
	length, parseErr := r.Body.Read(body)
  if parseErr != nil && parseErr != io.EOF {
		fmt.Println("error reading body", parseErr)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  //parse json
  var request Request
  parseErr = json.Unmarshal(body[:length], &request)
  if parseErr != nil {
		fmt.Println("error parsing json", parseErr)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

	response := Response{
		Data: nil,
		Error: true,
		Message: "invalid type",
	}

	var data interface{}
	var err error

	switch request.Type {
		case "player": data, err = api.playerHandler.Handle(request.Action, request.Args)
		default: err = errors.New("invalid type")
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


