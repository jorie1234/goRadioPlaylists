package main

import (
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	m, _ := NewMongoDB()
	c.AddFunc("@every 30m", func() { GetRadioData(m) })
	c.Start()
	GetRadioData(m)
	StartServer()
}

func GetRadioData(m *mongodb) {
	NewRadioPlaylistClient("radiohannover", m).MergePlaylistsAndSave()
	NewRadioPlaylistClient("radioffn", m).MergePlaylistsAndSave()
	NewRadioPlaylistClient("njoyradio", m).MergePlaylistsAndSave()
}
