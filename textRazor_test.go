package main

import (
	"io/ioutil"
	"testing"
)

func TestAnalysis(t *testing.T) {
	apiKey := "014dd0eee816fa4938f2364251273bc93c8ac0d04410ca8187676b88"

	data, err := ioutil.ReadFile("./error.json")
	if err != nil {
		t.Fatal(err)
	}

	tr := NewTextRazorRequest(apiKey)
	_, err = tr.Analysis(data)
	if err != nil {
		t.Fatal(err)
	}

}
