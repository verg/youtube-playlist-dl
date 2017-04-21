package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	videoInfoURL = "http://youtube.com/get_video_info?video_id="
)

type Video struct {
	id    string
	title string
}

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

func NewVideo(id string) *Video {
	return &Video{id: id}
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

func (y *Video) GetVideoStreams() (streams Streams, err error) {
	infoUrl := videoInfoURL + y.id
	resp, err := http.Get(infoUrl)
	if err != nil {
		return streams, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return streams, errors.New(fmt.Sprintf("Got response code: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return streams, err
	}

	return parseVideoInfo(body)
}

func parseVideoInfo(body []byte) (streams Streams, err error) {
	parsed, err := url.ParseQuery(string(body))
	if err != nil {
		return streams, err
	}

	if status := parsed["status"][0]; status != "ok" {
		errString := "Error requesting video info."
		reason, reason_exists := parsed["reason"]
		if reason_exists {
			errString += fmt.Sprintf("Reason: %s", reason[0])
		}
		return streams, errors.New(errString)
	}

	streamsData := strings.Split(parsed["url_encoded_fmt_stream_map"][0], ",")
	for i, streamData := range streamsData {
		query, err := url.ParseQuery(streamData)
		if err != nil {
			fmt.Printf("Error decoding stream: %d, %s", i, err)
			continue
		}

		stream, err := streamFromQueryData(query)
		if err == nil {
			streams = append(streams, stream)
		}
	}

	return streams, nil
}

func streamFromQueryData(streamData map[string][]string) (stream VideoStream, err error) {
	err = ensureFields(streamData, "quality", "type", "url")
	if err != nil {
		fmt.Printf("Error decoding stream: %s", err)
		return stream, err
	}
	quality := streamData["quality"][0]
	format := streamData["type"][0]
	url := streamData["url"][0]
	return NewVideoStream(url, quality, format), nil
}

func ensureFields(data map[string][]string, fields ...string) error {
	for _, field := range fields {
		if _, exists := data[field]; !exists {
			return errors.New(fmt.Sprintf("Missing field: %s", field))
		}
	}
	return nil
}

func (streams Streams) ChooseStream(quality string) (stream VideoStream, err error) {
	fmt.Printf("quality = %+v\n", quality)
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
	// playlist := "https://www.youtube.com/playlist?list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
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
