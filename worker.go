package main

import (
	beanstalk "github.com/JalfResi/gobeanstalk"
)

type WorkerConfig struct {
	srcTube          string
	destTube         string
	beanstalkdHost   string
	memcachedbHost   string
	maxRetryAttempts uint64
	timeout          int
	mysqlHost        string
	mysqlUsername    string
	mysqlPassword    string
}

type Worker chan struct{}

// DoWork does the following:
// - Pulls a URL out of the srcTube
// - Makes a GET/POST request to textrazor
// - Stores the results in MySQL
// - Deletes the job from Beanstalk
// - Creates a new job in destTube with the URL
//
// So we will need:
// - An ArticleSupplier to read article urls from the queue
// - An ArticleURL to represent an article url
// - An ReportRecorder to store the returned TextRazor report
// - An ArticleAnalyser to contact TextRazor and return a TextRazor report
// -
//
func (w Worker) DoWork(c *WorkerConfig) {

	// The following is a worker

	// Connect to beanstalkd
	bs, err := beanstalk.Dial(c.beanstalkdHost)
	if err != nil {
		logError.Fatalf("Beanstalk connect failed: %s\n", err)
	}

	as := NewArticleSupplier(bs, c.timeout, c.srcTube)
	aa := NewAnalyser(config.textRazorAPIKey)
	rr := NewReportRecorder(c.mysqlHost, c.mysqlUsername, c.mysqlPassword)

	for {
		article := as.GetArticleURL()
		report, err := aa.Analyse(article)
		if err != nil {
			if err == requestLimitMet {
				logError.Printf("%s: %d\n", err, config.totalRequestLimit)
				w.DieGracefully()
				// should possibly wait until the next day
				// reset the config.currentRequestCount and
				// start up the number of workers to continue
				// for the next day?
				//
				// What do we do about articleUrls that may be building up in
				// beanstalkd? Should we bin off old (i.e. yesterday)
				// article urls to deal with the backlog? Maybe we can
				// bury them? (but then how do we deal with the bury list?)
				return
			}

			if err.Error() == "Unauthorized" {
				as.Retry(article)
				logInfo.Printf("Got '%s' from TextRazor. Retrying\n", err)
				continue
			}

			logError.Printf("%+v %T\n", err, err)
			// Possibly bury continuinly failing jobs?
			as.Done(article)
			continue
		}
		as.Done(article)
		err = rr.StoreTopicsReport(report)
		if err != nil {
			logError.Println(err)
		}
	}
}

func (w Worker) DieGracefully() {
	close(w)
}
