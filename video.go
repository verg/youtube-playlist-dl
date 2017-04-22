package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kennygrant/sanitize"
)

const videoInfoURL = "http://youtube.com/get_video_info?video_id="

type Video struct {
	id    string
	title string
}

func NewVideo(id string) *Video {
	return &Video{id: id}
}

func (video *Video) Download(stream VideoStream, path string) error {
	filename := video.title + "." + stream.extention()
	path = sanitize.Path(filepath.Join(path, filename))
	fmt.Printf("Downloading: %s\n", video.title)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(stream.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errStr := fmt.Sprintf("Got response code: %d for %s", resp.StatusCode, video.title)
		return errors.New(errStr)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (video *Video) GetVideoStreams() (streams Streams, err error) {
	infoUrl := videoInfoURL + video.id
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
			fmt.Printf("Error decoding stream: %d, %s\n", i, err)
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
		fmt.Printf("Error decoding stream: %s\n", err)
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
