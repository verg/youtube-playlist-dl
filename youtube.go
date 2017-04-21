package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	videoInfoURL = "http://youtube.com/get_video_info?video_id="
)

type Streams []VideoStream

type VideoStream struct {
	quality Quality
	url     string
	format  string
}

func NewVideoStream(url, quality, format string) VideoStream {
	return VideoStream{
		quality: qualityStringToIotaMap[quality],
		url:     url,
		format:  format,
	}
}

type Quality uint8

var qualityStringToIotaMap = map[string]Quality{
	Quality720PString:   Quality720p,
	QualityMediumString: QualityMedium,
	QualitySmallString:  QualitySmall,
}

const (
	QualitySmall Quality = iota
	QualityMedium
	Quality720p

	NoPreferedQualityString = ""
	MaxQualityString        = "max"
	MinQualityString        = "min"
	Quality720PString       = "hd720"
	QualityMediumString     = "medium"
	QualitySmallString      = "small"
)

func (streams Streams) ChooseStream(quality string) (stream VideoStream, err error) {
	if len(streams) == 0 {
		return stream, errors.New("Empty streams struct")
	}
	switch quality {
	case NoPreferedQualityString:
		return streams[0], nil // choose arbitrarily
	case MaxQualityString:
		return streams.findMaxStream(), nil
	case MinQualityString:
		return streams.findMinStream(), nil
	default: // Search by name
		return streams.findByQuality(quality)
	}
}

func (streams Streams) findByQuality(qualityString string) (VideoStream, error) {
	quality, exists := qualityStringToIotaMap[qualityString]
	if !exists {
		errString := fmt.Sprintf("%s quality isn't a defined quality", qualityString)
		return VideoStream{}, errors.New(errString)
	}
	for _, stream := range streams {
		if stream.quality == quality {
			return stream, nil
		}
	}
	errString := fmt.Sprintf("No Matching Stream for %s quality", qualityString)
	return VideoStream{}, errors.New(errString)
}

func (streams Streams) findMinStream() VideoStream {
	var min VideoStream
	for i, stream := range streams {
		if i == 0 || stream.quality < min.quality {
			min = stream
		}
	}
	return min
}

func (streams Streams) findMaxStream() VideoStream {
	var max VideoStream
	for i, stream := range streams {
		if i == 0 || stream.quality > max.quality {
			max = stream
		}
	}
	return max
}

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
