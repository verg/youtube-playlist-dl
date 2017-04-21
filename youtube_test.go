package main

import (
	"fmt"
	"testing"
)

func TestParseVideoIDFromURL(t *testing.T) {
	url := "/watch?v=oFE3tp5esLw&index=6&list=PL6MuV0DF6AurABItm5OzSdVrEgJ_DxWVD"
	expected := "oFE3tp5esLw"
	got, err := parseVideoIDFromURL(url)
	if err != nil {
		t.Errorf(fmt.Sprintf("Got error %s", err))
	}
	if got != expected {
		t.Errorf("Expected: %s, got: %s", expected, got)
	}

}
