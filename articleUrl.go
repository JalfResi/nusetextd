package main

import (
	"fmt"
	"net/url"

	beanstalk "github.com/JalfResi/gobeanstalk"
)

// ArticleURL struct
type ArticleURL struct {
	url   *url.URL
	Hash  string
	job   *beanstalk.Job
	stats *StatsJob
}

// NewArticleURL ArticleURL constructor
func NewArticleURL(job *beanstalk.Job, stats *StatsJob) (*ArticleURL, error) {
	u, parseErr := url.Parse(string(job.Body))
	if parseErr != nil {
		return nil, parseErr
	}

	return &ArticleURL{
		url:   u,
		Hash:  fmt.Sprintf("%32s_article", generateHash(u.String())),
		job:   job,
		stats: stats,
	}, nil
}

func (a *ArticleURL) String() string {
	return a.url.String()
}
