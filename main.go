package main

import (
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	c.AddFunc("@every 30m", func() { GetRadioData() })
	StartServer()
}

func GetRadioData() {
	NewRadioPlaylistClient("radiohannover").MergePlaylistsAndSave()
	NewRadioPlaylistClient("radioffn").MergePlaylistsAndSave()
	NewRadioPlaylistClient("njoyradio").MergePlaylistsAndSave()
}
