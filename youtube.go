package main

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	videoInfoURL = "http://youtube.com/get_video_info?video_id="
)

type YouTubeVideo struct {
	id    string
	title string
}

func NewYouTubeVideo(id string) *YouTubeVideo {
	return &YouTubeVideo{id: id}
}

func parseVideoIDFromURL(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	if len(params["v"]) == 0 {
		return "", errors.New("Invalid url")
	}

	return params["v"][0], nil
}

func main() {

	// playlist := "https://www.youtube.com/playlist?list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
	// videos, _ := getPlaylist(playlist)
	url := "https://www.youtube.com/watch?v=HFFpKw1ecR4"
	id, _ := parseVideoIDFromURL(url)
	fmt.Printf("id = %+v\n", id)
	// v := NewYouTubeVideo(id)
}
