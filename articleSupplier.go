package main

import (
	beanstalk "github.com/JalfResi/gobeanstalk"
	"gopkg.in/yaml.v2"
)

type ArticleURLSupplier interface {
	GetArticleURL() *ArticleURL
	Done(fu *ArticleURL)
}

type ArticleSupplier struct {
	bsConn *beanstalk.Conn
	minTTR int
}

func NewArticleSupplier(bs *beanstalk.Conn, minTTR int, srcTube string) *ArticleSupplier {
	fs := &ArticleSupplier{
		bsConn: bs,
		minTTR: minTTR,
	}
	fs.SetSrcTube(srcTube)

	return fs
}

func (as *ArticleSupplier) SetSrcTube(srcTube string) {
	// Watch our source tube or bail
	_, err := as.bsConn.Watch(srcTube)
	if err != nil {
		logError.Fatalf("Could not watch tube %s: %v\n", srcTube, err)
	}
}

func (as *ArticleSupplier) Done(au *ArticleURL) {
	as.bsConn.Delete(au.job.ID)
}

func (as *ArticleSupplier) Retry(au *ArticleURL) {
	as.bsConn.Release(au.job.ID, 1, 0)
}

func (as *ArticleSupplier) GetArticleURL() *ArticleURL {
	for {
		job, err := as.bsConn.Reserve()
		if err != nil {
			logError.Fatal(err)
		}

		// First we check the TTR of the job
		// if it is lower than *timeout, we
		// put it back in the tube with an
		// increased TTR. This ensures that
		// all jobs we deal with can be dealt
		// within a sensible timeframe, otherwise
		// the job will keep failing and will be
		// automatically reclaimed by beanstalk
		stats := as.getJobTTR(job)
		if stats.TTR < as.minTTR {
			as.increaseJobTTR(job, stats, as.minTTR)
			logError.Printf("Increased job %d TTR to %d from %d\n", job.ID, as.minTTR, stats.TTR)
			continue
		}

		au, err := NewArticleUrl(job, stats)
		if err != nil {
			as.bsConn.Bury(job.ID, 1)
			logError.Printf("Bad Article URL format; burying: %s\n", err)
			continue
		}
		logInfo.Printf("Article URL: %s\n", au)
		return au
	}
}

func (as *ArticleSupplier) getJobTTR(job *beanstalk.Job) *StatsJob {
	rawJobStats, err := as.bsConn.StatsJob(job.ID)
	if err != nil {
		logError.Fatalf("Job %d StatsJob failed: %s\n", job.ID, err)
	}

	statsJob := StatsJob{}
	err = yaml.Unmarshal(rawJobStats, &statsJob)
	if err != nil {
		logError.Fatalf("Job %d yaml: %s\n", job.ID, err)
	}

	return &statsJob
}

// NOTE:
// Uses globals! Naughty!
func (as *ArticleSupplier) increaseJobTTR(job *beanstalk.Job, stats *StatsJob, newTTR int) {
	as.bsConn.Use(config.srcTube)
	as.bsConn.PutUnique(job.Body, stats.Pri, 1, newTTR) // We can set the delay to 1 because the delay is already up and will be reset when we crawl the feed
	as.bsConn.Delete(job.ID)
}
