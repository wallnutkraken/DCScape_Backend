package main

import (
	"fmt"
	"net/http"
	"time"
)

var worlds []World
var worldJson []byte
var players []byte

func main() {
	var err error
	go asyncGetPlayers()
	if worldJson, err = getWorlds(); err != nil {
		fmt.Println("getWorlds:", err)
	}
	http.HandleFunc("/worlds", onGetWorlds)
	http.HandleFunc("/players", onGetPlayers)
	err = http.ListenAndServe(":10101", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err.Error())
	}
}

func onGetWorlds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(worldJson)
}

func asyncGetPlayers() {
	for {
		data, err := getPlayers()
		if err != nil {
			fmt.Println("asyncGetPlayers:", err)
			continue
		}
		players = data

		/* Sleep for 30 seconds */
		time.Sleep(time.Second * 30)
	}
}

func onGetPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json, err := getPlayers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(json)
	}
}
