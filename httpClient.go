package main

// Based on code found in the following gist: https://gist.github.com/dmichael/5710968

import (
	"net"
	"net/http"
	"time"
)

// Config struct
type Config struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

// TimeoutDialer function
func TimeoutDialer(config *Config) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(config.ReadWriteTimeout))
		return conn, nil
	}
}

// NewTimeoutClient timeout HTTP client constructor
func NewTimeoutClient(args ...interface{}) *http.Client {
	// Default configuration
	config := &Config{
		ConnectTimeout:   1 * time.Second,
		ReadWriteTimeout: 1 * time.Second,
	}

	// merge the default with user input if there is one
	if len(args) == 1 {
		timeout := args[0].(time.Duration)
		config.ConnectTimeout = timeout
		config.ReadWriteTimeout = timeout
	}

	if len(args) == 2 {
		config.ConnectTimeout = args[0].(time.Duration)
		config.ReadWriteTimeout = args[1].(time.Duration)
	}

	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(config),
		},
	}
}
