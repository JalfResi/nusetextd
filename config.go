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

// This must be mutex locked as multiple connections could be
// modifying the config options over the command connection

type NusefeedConfig struct {
	sync.Mutex
	verbose            bool
	configTest         bool
	debug              bool
	srcTube            string
	destTube           string
	beanstalkdHost     string
	memcachedbHost     string
	maxRetryAttempts   uint64
	timeout            int
	initialWorkerCount int
	textRazorAPIKey    string
}

var config *NusefeedConfig = &NusefeedConfig{}

func init() {
	config.Lock()
	defer config.Unlock()

	flag.BoolVar(&config.verbose, "verbose", false, "Display verbose information messages")
	flag.BoolVar(&config.debug, "debug", false, "Display debug messages")
	flag.BoolVar(&config.configTest, "test", false, "Display config options")
	flag.StringVar(&config.srcTube, "src-tube", "default", "The source tube")
	flag.StringVar(&config.destTube, "dest-tube", "articles", "The destination tube")
	flag.StringVar(&config.beanstalkdHost, "beanstalk", "127.0.0.1:11300", "The beanstalk host")
	flag.StringVar(&config.memcachedbHost, "memcache", "127.0.0.1:11211", "The memcache host")
	flag.Uint64Var(&config.maxRetryAttempts, "max-fetch-retries", 3, "The maximum number of attempts to fetch a feed url")
	flag.IntVar(&config.timeout, "timeout", 30, "The http connection timeout")
	flag.IntVar(&config.initialWorkerCount, "workers", 100, "The initial worker count")
	flag.StringVar(&config.textRazorAPIKey, "key", "014dd0eee816fa4938f2364251273bc93c8ac0d04410ca8187676b88", "The TextRazor API key")

	flagenv.UseUpperCaseFlagNames = true
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
