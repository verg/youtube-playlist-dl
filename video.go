package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Video struct {
	id    string
	title string
}

func NewVideo(id string) *Video {
	return &Video{id: id}
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
