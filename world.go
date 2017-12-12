package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net"
	"strconv"
	"strings"
)

type World struct {
	Number   int    `json:"number"`
	Hostname string `json:"host"`
	Type     string `json:"type"`
	Activity string `json:"activity"`
	Players  int    `json:"players"`
}

type WorldPlayers struct {
	Number      int `json:"number"`
	PlayerCount int `json:"players"`
}

func (w *World) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}

func getWorlds() ([]byte, error) {
	doc, err := goquery.NewDocument("http://oldschool.runescape.com/slu")
	if err != nil {
		return nil, err
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
		playersText := sel.Find(".server-list__row-cell").First().Next().Text()
		world.Players, _ = strconv.Atoi(strings.Replace(playersText, " players", "", -1))

		worlds = append(worlds, world)
	})

	return json.Marshal(worlds)
}

func getPlayers() ([]byte, error) {
	doc, err := goquery.NewDocument("http://oldschool.runescape.com/slu")
	if err != nil {
		return nil, err
	}

	players := make([]WorldPlayers, 0)
	doc.Find(".server-list__row").Each(func(index int, sel *goquery.Selection) {
		wPlayers := WorldPlayers{}
		worldID, _ := sel.Find(".server-list__world-link").Attr("id")
		wPlayers.Number, _ = strconv.Atoi(strings.Replace(worldID, "slu-world-", "", -1))
		wPlayers.Number -= 300
		playersText := sel.Find(".server-list__row-cell").First().Next().Text()
		wPlayers.PlayerCount, _ = strconv.Atoi(strings.Replace(playersText, " players", "", -1))

		players = append(players, wPlayers)
	})

	return json.Marshal(players)
}
