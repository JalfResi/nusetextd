package main

import (
	"os"
	"testing"
)

func TestAnalysis(t *testing.T) {
	apiKey := "014dd0eee816fa4938f2364251273bc93c8ac0d04410ca8187676b88"

	var f *os.File
	f, err := os.Open("./error.json")
	if err != nil {
		t.Fatalf("Failed to open error.json: %v", err)
	}
	defer f.Close()

	tr := NewTextRazorRequest(apiKey)
	_, err = tr.Analysis(f)
	if err != nil {
		t.Fatal(err)
	}

}
