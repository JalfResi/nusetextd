package main

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// ErrRequestLimitMet error
var ErrRequestLimitMet = errors.New("Request limit met")

// ArticleAnalyser interface
type ArticleAnalyser interface {
	Analyse(u *ArticleURL) *TextRazorResult
}

// Analyser struct
type Analyser struct {
	apiKey            string
	downloadUserAgent string
}

// NewAnalyser Analyser constructor
func NewAnalyser(key string) *Analyser {
	return &Analyser{
		apiKey:            key,
		downloadUserAgent: fmt.Sprintf("NuseAgent Article Downloader v1.0 (%s)", url.QueryEscape("http://nuseagent.com/")),
	}
}

// Analyse method
func (a *Analyser) Analyse(u *ArticleURL) (*TextRazorResult, error) {

	if config.RequestLimitMet() {
		return nil, ErrRequestLimitMet
	}

	c := NewTimeoutClient(time.Duration(config.timeout) * time.Second)

	tr := NewTextRazorRequest(a.apiKey)
	tr.DownloadUserAgent = a.downloadUserAgent
	tr.URL = u.String()
	tr.CleanupMode = ModeCleanHTML
	tr.CleanupReturnCleaned = false
	tr.CleanupReturnRaw = false
	tr.SetExtractors(
		ExtractorTopics,
		ExtractorEntities,
		ExtractorWords,
		ExtractorPhrases,
		ExtractorDependancyTrees,
		ExtractorRelations,
		ExtractorEntailments,
		ExtractorSenses,
	)

	data, err := tr.Fetch(c)
	if err != nil {
		return nil, err
	}
	result, err := tr.Analysis(data)
	config.IncRequestCount()

	return result, err
}
