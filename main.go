package main

import (
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	c.AddFunc("30 * * * *", func() { GetRadioData() })
	StartServer()
}

func GetRadioData() {
	NewRadioPlaylistClient("radiohannover").MergePlaylistsAndSave()
	NewRadioPlaylistClient("radioffn").MergePlaylistsAndSave()
	NewRadioPlaylistClient("njoyradio").MergePlaylistsAndSave()
}
