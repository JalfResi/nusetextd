package main

import (
	"encoding/json"

	"github.com/bitly/go-nsq"
)

// NsqReportRecorder stores an TextRazorResults in NSQ
type NsqReportRecorder struct {
	w           *nsq.Producer
	TopicsTopic string
}

// NewNsqReportRecorder is a ReportRecorder constructor
func NewNsqReportRecorder(hostname string) (*NsqReportRecorder, error) {
	config := nsq.NewConfig()
	w, err := nsq.NewProducer(hostname, config)
	if err != nil {
		return nil, err
	}

	return &NsqReportRecorder{
		w:           w,
		TopicsTopic: "article_topics",
	}, nil
}

// StoreTopics If there is an error executing any of the inserts, all pervious inserts
// for this TextRazorResult is reolledback, ensuring we dont have a partial
// TextRazorResult written to the database.
func (rr *NsqReportRecorder) StoreTopics(r *TextRazorResult) error {

	b, err := json.Marshal(
		struct {
			URL    string
			Topics []TextRazorTopic
		}{
			r.URL,
			r.Response.Topics,
		})
	if err != nil {
		return err
	}

	err = rr.w.Publish(rr.TopicsTopic, b)
	if err != nil {
		return err
	}

	return nil
}
