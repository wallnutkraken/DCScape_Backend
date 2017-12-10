package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var worlds []World
var worldJson []byte

func main() {
	if err := getWorlds(); err != nil {
		fmt.Println("getWorlds:", err)
	}
	http.HandleFunc("/worlds", onGetWorlds)
	err := http.ListenAndServe(":10101", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err.Error())
	}
}

func getWorlds() error {
	doc, err := goquery.NewDocument("http://oldschool.runescape.com/slu")
	if err != nil {
		return err
	}

	worlds := make([]World, 0)
	doc.Find(".server-list__row").Each(func(index int, sel *goquery.Selection) {
		world := World{}
		worldID, _ := sel.Find(".server-list__world-link").Attr("id")
		world.Number, _ = strconv.Atoi(strings.Replace(worldID, "slu-world-", "", -1))
		world.Number -= 300
		world.Type = sel.Find(".server-list__row-cell--type").Text()
		worldUrl := fmt.Sprintf("oldschool%d.runescape.com", world.Number)
		ips, err := net.LookupIP(worldUrl)
		if err != nil {
			world.Hostname = worldUrl
		} else {
			world.Hostname = ips[0].String()
		}
		world.Activity = sel.Find(".server-list__row-cell").Last().Text()

		worlds = append(worlds, world)
	})

	worldJson, err = json.Marshal(worlds)
	return err
}

func onGetWorlds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(worldJson)
}
