package main

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Connects to TextRazor and outputs the returned result
// in an AnalysisReport

var requestLimitMet = errors.New("Request limit met")

type ArticleAnalyser interface {
	Analyse(u *ArticleURL) *TopicsReport
}

type Analyser struct {
	apiKey            string
	downloadUserAgent string
}

func NewAnalyser(key string) *Analyser {
	return &Analyser{
		apiKey:            key,
		downloadUserAgent: fmt.Sprintf("NuseAgent Article Downloader v1.0 (%s)", url.QueryEscape("http://nuseagent.com/")),
	}
}

func (a *Analyser) Analyse(u *ArticleURL) (*TopicsReport, error) {

	if config.RequestLimitMet() {
		return nil, requestLimitMet
	}

	c := NewTimeoutClient(time.Duration(config.timeout) * time.Second)

	tr := NewTextRazorRequest(a.apiKey)
	tr.DownloadUserAgent = a.downloadUserAgent
	tr.Url = u.String()
	tr.CleanupMode = MODE_CLEANHTML
	tr.CleanupReturnCleaned = false
	tr.CleanupReturnRaw = false
	tr.SetExtractors(EXTRACTOR_TOPICS)

	r, err := tr.Analysis(c)
	if err != nil {
		return nil, err
	}

	config.IncRequestCount()
	return NewTopicsReport(tr, r), nil
}
