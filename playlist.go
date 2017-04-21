package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const (
	VideoTitle    = ".pl-video-title-link"
	PlaylistTitle = ".pl-header-title"
)

type Playlist struct {
	title  string
	videos []*Video
}

// Example "https://www.youtube.com/playlist?list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
func getPlaylist(url string) (playlist Playlist, err error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return playlist, err
	}

	playlist.title = doc.Find(PlaylistTitle).Text()
	fmt.Printf("playlist.title = %+v\n", playlist.title)
	doc.Find(VideoTitle).Each(func(i int, titleLink *goquery.Selection) {
		videoTitle := titleLink.Text()
		link, exists := titleLink.Attr("href")
		if exists {
			id, err := parseVideoIDFromURL(link)
			if err == nil {
				video := &Video{title: videoTitle, id: id}
				playlist.videos = append(playlist.videos, video)
			} else {
				fmt.Printf("Error parsing: %s", link)
			}
		}
	})
	return playlist, nil
}
