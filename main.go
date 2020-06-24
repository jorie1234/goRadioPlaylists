package main

func main() {

	NewRadioPlaylistClient("radiohannover").MergePlaylistsAndSave()
	NewRadioPlaylistClient("radioffn").MergePlaylistsAndSave()
	NewRadioPlaylistClient("njoyradio").MergePlaylistsAndSave()
}
