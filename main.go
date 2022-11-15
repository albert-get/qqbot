package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"qqbot/command"
	"qqbot/config"
	"qqbot/route"
	"time"
)

var conn *websocket.Conn
var beatS *int
var upChannel = make(chan int)

var Token = config.Token

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", route.WakeHandle)

	return r
}
func weakUp() {
	go startUp()

	for {

		select {
		case i := <-upChannel:
			if i == 0 {
				go startUp()
			}
		}
	}
}
func startUp() {

	c, _, errc := websocket.DefaultDialer.Dial("wss://sandbox.api.sgroup.qq.com/websocket", nil)
	if errc != nil {
		return
	}
	conn = c
	mapInstances := map[string]interface{}{}
	if beatS != nil {
		mapInstances["op"] = 6
		mapInstances["d"] = map[string]interface{}{
			"token":   Token,
			"intents": 513,
			"seq":     &beatS,
		}
	} else {
		mapInstances["op"] = 2
		mapInstances["d"] = map[string]interface{}{
			"token":   Token,
			"intents": 513,
		}
	}

	jsonStr, errJ := json.Marshal(mapInstances)
	if errJ != nil {
		return
	}
	errId := conn.WriteMessage(websocket.TextMessage, []byte(jsonStr))
	if errId != nil {
		return
	}

	go handleMess()
	go heartBeat()
}

func handleMess() {
	defer func() { upChannel <- 0 }()
	defer conn.Close()
	for {

		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("err:", err)

			return
		}

		var mess map[string]interface{}
		errRes := json.Unmarshal([]byte(message), &mess)
		if errRes != nil {
			return
		}
		d := mess["s"]
		switch d.(type) {
		case int:
			*beatS = d.(int)
		}

		messType := fmt.Sprintf("%v", mess["t"])
		switch messType {
		case "MESSAGE_CREATE":
			command.Cluctue(mess)
		}
	}

}

func heartBeat() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	send := map[string]interface{}{
		"op": 1,
		"d":  &beatS,
	}
	jsonSend, errJson := json.Marshal(send)
	if errJson != nil {
		return
	}
	for {
		select {
		case <-ticker.C:
			errHart := conn.WriteMessage(websocket.TextMessage, []byte(jsonSend))
			if errHart != nil {
				log.Println(errHart)
				return
			}
		}
	}

}

func main() {
	go weakUp()

	r := setupRouter()
	r.Run(":9501")
}
