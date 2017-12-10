package main

import "encoding/json"

type World struct {
	Number   int    `json:"number"`
	Hostname string `json:"host"`
	Type     string `json:"type"`
	Activity string `json:"activity"`
}

func (w *World) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}
