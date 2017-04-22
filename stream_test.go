package main

import (
	"fmt"
	"testing"
)

const (
	urlString = "https://example.com/"
	format    = `video/mp4; codecs="avc1.42001E, mp4a.40.2"`
)

var small VideoStream = NewVideoStream(urlString, QualitySmallString, format)
var medium VideoStream = NewVideoStream(urlString, QualityMediumString, format)
var large VideoStream = NewVideoStream(urlString, Quality720PString, format)

func prefsFor(quality string) Preferences {
	return Preferences{quality: quality}
}

func TestChooseMinStream(t *testing.T) {
	streams := genStreams()
	expected := small
	got, err := streams.ChooseStream(prefsFor(MinQualityString))
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	if got != expected {
		t.Errorf("Expected: %v, Got: %v", expected, got)
	}
}

func TestChooseMaxStream(t *testing.T) {
	streams := genStreams()
	expected := large
	got, err := streams.ChooseStream(prefsFor(MaxQualityString))
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	if got != expected {
		t.Errorf("Expected: %v, Got: %v", expected, got)
	}
}

func TestChooseStreamByName(t *testing.T) {
	streams := genStreams()
	expected := medium
	got, err := streams.ChooseStream(prefsFor(QualityMediumString))
	if err != nil {
		t.Errorf(fmt.Sprint(err))
	}
	if got != expected {
		t.Errorf("Expected: %v, Got: %v", expected, got)
	}
}

func TestExtension(t *testing.T) {
	expected := "mp4"
	if got := medium.extention(); got != expected {
		t.Errorf("Expected: %v, Got: %v", expected, got)
	}
}

func genStreams() Streams {
	return Streams([]VideoStream{small, medium, large})
}
