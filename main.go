package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

const (
	LOG_PREFIX_INFO  = "INFO "
	LOG_PREFIX_ERROR = "ERROR "
)

var (
	logInfo  *log.Logger = log.New(ioutil.Discard, LOG_PREFIX_INFO, log.Ldate|log.Ltime)
	logError *log.Logger = log.New(os.Stderr, LOG_PREFIX_ERROR, log.Ldate|log.Ltime)
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if config.configTest {
		flag.PrintDefaults()

		fmt.Println("\nCurrent configuration")
		fmt.Printf("verbose: %+v\n", config.verbose)
		fmt.Printf("debug: %+v\n", config.debug)
		fmt.Printf("config-test: %+v\n", config.configTest)
		fmt.Printf("src-tube: %+v\n", config.srcTube)
		fmt.Printf("beanstalkd: %+v\n", config.beanstalkdHost)
		fmt.Printf("memcachedb: %+v\n", config.memcachedbHost)
		fmt.Printf("max-fetch-retries: %+v\n", config.maxRetryAttempts)
		fmt.Printf("timeout: %+v\n", config.timeout)
		fmt.Printf("workers: %+v\n", config.initialWorkerCount)
		os.Exit(0)
	}

	if config.verbose {
		logInfo = log.New(os.Stdout, LOG_PREFIX_INFO, log.Ldate|log.Ltime)
	}

	if config.debug {
		logInfo.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		logError.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	quit := make(chan bool)
	stack := &Stack{}

	// Hook up workers here...
	config.Lock()
	c := config.initialWorkerCount
	config.Unlock()

	for _, worker := range stack.Inc(c) {
		go worker.DoWork(newWorkerConfig(config))
	}

	log.Printf("Running %d workers\n", stack.Len())

	/*
	   // Send kill signal over this slice of chans
	   for _, worker := range stack.Dec(3) {
	       worker.DieGracefully()
	   }
	   log.Printf("Running %d workers\n", stack.Len())
	*/

	<-quit // This will eventually be replaced with the command server
}
