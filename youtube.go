package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/kennygrant/sanitize"
)

const MaxConncurentDownloads = 3

const (
	quality_usage = `(optional) quality of video. e.g. "medium". "max" or "min" will select the highest or lowest availible quality.`
)

type Preferences struct {
	quality string
}

var wg sync.WaitGroup

func main() {
	preferences := GetPreferencesFromFlags()
	url := flag.Arg(0)
	playlist, err := getPlaylist(url)
	if err != nil {
		fmt.Println("Error getting playlist")
		os.Exit(1)
	}

	dir := playlist.title
	makePlaylistDir(dir)

	semaphore := make(chan struct{}, MaxConncurentDownloads)
	for _, video := range playlist.videos {
		streams, err := video.GetVideoStreams()
		if err != nil {
			fmt.Printf("Error getting streams for %s\n", video.title)
			continue
		}
		stream, err := streams.ChooseStream(preferences)
		if err != nil {
			fmt.Printf("Error choosing stream for %s\n", video.title)
			continue
		}

		// Download MaxConncurentDownloads at a time
		wg.Add(1)
		go func(video *Video, stream VideoStream) {
			semaphore <- struct{}{}
			err = video.Download(stream, dir)
			if err != nil {
				fmt.Printf("Error Downloading %s, %s", video.title, err)
			}
			<-semaphore
			wg.Done()
		}(video, stream)
	}
	wg.Wait()
}

func GetPreferencesFromFlags() Preferences {
	quality := flag.String("q", NoPreferedQualityString, quality_usage)
	flag.Parse()
	checkUsage()
	return Preferences{quality: *quality}
}

func checkUsage() {
	if flag.NArg() < 1 {
		fmt.Printf("Usage: %s url\n", os.Args[0])
		os.Exit(1)
	}
}

func makePlaylistDir(dirName string) {
	path := sanitize.Path(dirName)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Printf("Error creating dir %s\n", dirName)
		os.Exit(1)
	}
}
