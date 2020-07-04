package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type RadioPlaylist struct {
	RadioName string
	BaseURL   *url.URL
	pd        *PlayData
	Client    *Client
	mongo     *mongodb
}

type PlayData struct {
	Playlist       Playlist `json:"playlist"`
	CurrentTrackID string   `json:"currentTrackId"`
}
type PlaylistEntry struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int64    `json:"created"`
	CreatedTime time.Time
}

type Playlist []PlaylistEntry

func NewRadioPlaylistClient(sender string, m *mongodb) *RadioPlaylist {
	p := RadioPlaylist{
		RadioName: sender,
		mongo:     m,
	}
	p.BaseURL, _ = url.Parse("https://onlineradiobox.com")
	p.Client = NewClient(nil)
	p.Client.SetBaseURL(p.BaseURL.String())
	newpath := filepath.Join(".", "data")
	os.MkdirAll(newpath, os.ModePerm)
	return &p
}
func (p *RadioPlaylist) GetFileName() string {
	return filepath.Join("data", fmt.Sprintf("%s.json", p.RadioName))

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
	p.pd, err = p.Client.GetPlaylist(fmt.Sprintf("/json/de/%s/playlist/0?tz=-1000&rnd=0.36527829491055281", p.RadioName))
	if err != nil {
		log.Printf("Cannot get Playlist %v", err)
	}
	return &p.pd.Playlist, err
}

func (p *RadioPlaylist) SavePlaylistToFile(data *Playlist) error {
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("cannot marshal playlist %v", err)
		return err
	}
	ioutil.WriteFile(p.GetFileName(), b, 0600)
	return nil
}

func (p *RadioPlaylist) MergePlaylistsAndSave() {
	p.mongo.EnsureIndex(p.RadioName)
	old := p.mongo.GetLastEntry(p.RadioName)
	log.Printf("lastentry %v", old)
	actualPlaylist, err := p.GetRadioData()
	if err != nil {
		log.Printf("cannot get Playlist %v", err)
	}


	log.Printf("Playlist %v", actualPlaylist)
	log.Printf("Playlist len %d", len(*actualPlaylist))
	oldPlaylist := p.ReadPlaylistFromFile()
	log.Printf("Old list %v", oldPlaylist)
	log.Printf("Old list len %d", len(*oldPlaylist))

	for _, v := range *oldPlaylist {
		v.CreatedTime = time.Unix(v.Created, 0)
		if v.Created > old.Created {
			p.mongo.Insert(p.RadioName, v)
		}
		found := false
		for _, vv := range *actualPlaylist {
			if v.Created == vv.Created {
				found = true
			}
		}
		if !found {
			*actualPlaylist = append(*actualPlaylist, v)
			log.Printf("Append %v to playlist", v)
		}
	}

	for _, vv := range *actualPlaylist {
		vv.CreatedTime = time.Unix(vv.Created, 0)
		if vv.Created > old.Created {
			p.mongo.Insert(p.RadioName, vv)
		}
	}


	sort.Slice(*actualPlaylist, func(i, j int) bool {
		return (*actualPlaylist)[i].Created < (*actualPlaylist)[j].Created
	})
	p.SavePlaylistToFile(actualPlaylist)
}
