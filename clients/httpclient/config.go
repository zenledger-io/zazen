package httpclient

import "time"

var (
	DefaultMaxIdleConns          = 50
	DefaultIdleConnTimeout       = 120 * time.Second
	DefaultKeepAlive             = 60 * time.Second
	DefaultTimeout               = 10 * time.Second
	DefaultDialTimeout           = 5 * time.Second
	DefaultTLSHandshakeTimeout   = 5 * time.Second
	DefaultExpectContinueTimeout = 1 * time.Second
	DefaultResponseHeaderTimeout = 5 * time.Second
)

// config

type Config struct {
	MaxIdleConns          int
	IdleConnTimeout       time.Duration
	KeepAlive             time.Duration
	Timeout               time.Duration
	DialTimeout           time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	ResponseHeaderTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		MaxIdleConns:          DefaultMaxIdleConns,
		IdleConnTimeout:       DefaultIdleConnTimeout,
		KeepAlive:             DefaultKeepAlive,
		Timeout:               DefaultTimeout,
		DialTimeout:           DefaultDialTimeout,
		TLSHandshakeTimeout:   DefaultTLSHandshakeTimeout,
		ExpectContinueTimeout: DefaultExpectContinueTimeout,
		ResponseHeaderTimeout: DefaultResponseHeaderTimeout,
	}
}
