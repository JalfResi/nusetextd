package main

import (
	"fmt"
)

// Stores an AnalysisReport in a MySQL table

type ReportRecorderer interface {
	StoreReport(r *TopicsReport)
}

type ReportRecorder struct{}

func NewReportRecorder(mysqlHost, mysqlUsername, mysqlPassword string) *ReportRecorder {
	return &ReportRecorder{}
}

func (rr *ReportRecorder) StoreTopicsReport(r *TopicsReport) {
	for _, topic := range r.topics {
		fmt.Println(topic)
	}
}
