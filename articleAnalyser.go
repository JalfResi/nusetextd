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
	Analyse(u *ArticleURL) *TextRazorResult
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

func (a *Analyser) Analyse(u *ArticleURL) (*TextRazorResult, error) {

	if config.RequestLimitMet() {
		return nil, requestLimitMet
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

	return tr.Analysis(c)
}
