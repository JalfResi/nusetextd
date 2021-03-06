package main

// Config should provide a observer/listener interface
// as the following changes to a config property will
// require executing additional functions
//
// NOTE:
// It could be possible that the setter functions trigger the necessary
// calls, but there may be an issue with scope so I will have to explore
// this further.
//

import (
	"flag"
	"sync"

	"github.com/JalfResi/flagenv"
)

// NusefeedConfig struct
// This must be mutex locked as multiple connections could be
// modifying the config options
type NusefeedConfig struct {
	sync.Mutex
	verbose             bool
	configTest          bool
	debug               bool
	srcTube             string
	destTube            string
	beanstalkdHost      string
	memcachedbHost      string
	maxRetryAttempts    uint64
	timeout             int
	initialWorkerCount  int
	totalRequestLimit   int
	currentRequestCount int
	textRazorAPIKey     string
}

// IncRequestCount method
func (c *NusefeedConfig) IncRequestCount() {
	c.Lock()
	defer c.Unlock()
	c.currentRequestCount++
}

// RequestCount method
func (c *NusefeedConfig) RequestCount() int {
	c.Lock()
	defer c.Unlock()
	return c.currentRequestCount
}

// RequestLimitMet method
func (c *NusefeedConfig) RequestLimitMet() bool {
	c.Lock()
	defer c.Unlock()
	return (c.currentRequestCount >= c.totalRequestLimit)
}

var config = &NusefeedConfig{}

func init() {
	config.Lock()
	defer config.Unlock()

	flag.BoolVar(&config.verbose, "verbose", false, "Display verbose information messages")
	flag.BoolVar(&config.debug, "debug", false, "Display debug messages")
	flag.BoolVar(&config.configTest, "test", false, "Display config options")
	flag.StringVar(&config.srcTube, "src-tube", "articles", "The source tube")
	flag.StringVar(&config.beanstalkdHost, "beanstalk", "127.0.0.1:11300", "The beanstalk host")
	flag.StringVar(&config.memcachedbHost, "memcache", "127.0.0.1:11211", "The memcache host")
	flag.Uint64Var(&config.maxRetryAttempts, "max-fetch-retries", 3, "The maximum number of attempts to fetch a feed url")
	flag.IntVar(&config.timeout, "timeout", 30, "The http connection timeout")
	flag.IntVar(&config.initialWorkerCount, "workers", 2, "The initial worker count")
	flag.IntVar(&config.totalRequestLimit, "requests", 500, "The maximum TextRazor requests in a 24hr period")
	flag.StringVar(&config.textRazorAPIKey, "key", "", "The TextRazor API key")

	flagenv.Prefix = "NUSETEXT_"
	flagenv.Parse()
	flag.Parse()
}

func newWorkerConfig(c *NusefeedConfig) *WorkerConfig {
	c.Lock()
	defer c.Unlock()

	return &WorkerConfig{
		srcTube:          c.srcTube,
		destTube:         c.destTube,
		beanstalkdHost:   c.beanstalkdHost,
		memcachedbHost:   c.memcachedbHost,
		maxRetryAttempts: c.maxRetryAttempts,
		timeout:          c.timeout,
	}
}
