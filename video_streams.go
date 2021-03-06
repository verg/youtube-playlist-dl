package main

import (
	"errors"
	"fmt"
	"strings"
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

func (streams Streams) ChooseStream(preferences Preferences) (stream VideoStream, err error) {
	if len(streams) == 0 {
		return stream, errors.New("Empty streams struct")
	}
	quality := preferences.quality
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

func (stream VideoStream) extention() string {
	slash := strings.Index(stream.format, `/`)
	ext := stream.format[slash+1:]
	if semiColon := strings.Index(ext, `;`); semiColon >= 0 {
		ext = ext[:semiColon]
	}
	return ext
}
