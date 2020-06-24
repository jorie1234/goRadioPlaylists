package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"sort"
)

type RadioPlaylist struct {
	RadioName string
	BaseURL   *url.URL
	pd *PlayData
	Client *Client
}

type PlayData struct {
	Playlist       Playlist `json:"playlist"`
	CurrentTrackID string   `json:"currentTrackId"`
}
type PlaylistEntry struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int    `json:"created"`
}

type Playlist []PlaylistEntry

func NewRadioPlaylistClient(sender string) *RadioPlaylist {
	p:=RadioPlaylist{
		RadioName: sender,
	}
	p.BaseURL,_=url.Parse("https://onlineradiobox.com",)
	p.Client=NewClient(nil)
	p.Client.SetBaseURL(p.BaseURL.String())
	return &p
}
func (p *RadioPlaylist) GetFileName()string  {
	return fmt.Sprintf("%s.json", p.RadioName)

}
func (p *RadioPlaylist) ReadPlaylistFromFile() *Playlist {
	d, _ := ioutil.ReadFile(p.GetFileName())
	var data Playlist
	err := json.Unmarshal(d, &data)
	if err != nil {
		log.Printf("cannot unmarshal file %s -> %v", p.GetFileName(), err)
		return &data
	}
	return &data
}

func (p *RadioPlaylist) GetRadioData() (*Playlist, error) {
	var err error
	p.pd, err =p.Client.GetPlaylist(fmt.Sprintf("/json/de/%s/playlist/0?tz=-1000&rnd=0.36527829491055281", p.RadioName))
	if err != nil {
		log.Printf("Cannot get Playlist %v", err)
	}
	return &p.pd.Playlist, err
}

func (p *RadioPlaylist)  SavePlaylistToFile(data *Playlist) error {
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("cannot marshal playlist %v", err)
		return err
	}
	ioutil.WriteFile(p.GetFileName(), b, 0600)
	return nil
}

func (p *RadioPlaylist)MergePlaylistsAndSave()  {
	d, err := p.GetRadioData()
	if err != nil {
		log.Printf("cannot get Playlist %v", err)
	}
	log.Printf("Playlist %v", d)
	log.Printf("Playlist len %d", len(*d))
	oldPlaylist:=p.ReadPlaylistFromFile()
	log.Printf("Old list %v", oldPlaylist)
	log.Printf("Old list len %d", len(*oldPlaylist))

	for _,v:=range *oldPlaylist {
		found:=false
		for _,vv:=range *d {
			if v.Created==vv.Created{
				found=true
			}
		}
		if !found {
			*d=append(*d, v)
			log.Printf("Append %v to playlist", v)
		}
	}

	sort.Slice(*d, func(i, j int) bool {
		return (*d)[i].Created < (*d)[j].Created
	})
	p.SavePlaylistToFile(d)
}

