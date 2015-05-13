package main

import (
	"sort"
)

// Stores the ArticleURL along with the TextRazor output
// Ready to be stored in a database

type TopicsReport struct {
	topics sort.StringSlice
}

func NewTopicsReport(t *TextRazorResult) *TopicsReport {
	ar := &TopicsReport{}
	for _, topic := range t.Response.Topics {
		if topic.Score == 1.0 {
			ar.topics = append(ar.topics, topic.Label)
		}
	}
	ar.topics.Sort()
	return ar
}
