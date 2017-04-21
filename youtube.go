package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	quality_usage = `(optional) quality of video. e.g. "medium". "max" or "min" will select the highest or lowest availible quality.`
)

func main() {
	playlistURL := "https://www.youtube.com/playlist?list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
	playlist, _ := getPlaylist(playlistURL)
	fmt.Printf("playlist = %+v\n", playlist)
	// videos, _ := getPlaylist(playlist)
	quality := flag.String("q", NoPreferedQualityString, quality_usage)
	flag.Parse()
	checkUsage()
	url := flag.Arg(0)

	id, _ := parseVideoIDFromURL(url)
	v := NewVideo(id)
	streams, _ := v.GetVideoStreams()
	for _, s := range streams {
		fmt.Printf("s = %+v\n", s)
		fmt.Println("")
	}
	stream, _ := streams.ChooseStream(*quality)
	fmt.Printf("stream = %+v\n", stream)
}

func checkUsage() {
	if flag.NArg() < 1 {
		fmt.Printf("Usage: %s url\n", os.Args[0])
		os.Exit(1)
	}
}
