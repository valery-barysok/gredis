package gredis

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultHost         = "localhost"
	defaultPort         = "16379"
	defaultTimeout      = time.Minute
	defaultReadTimeout  = 2 * time.Second
	defaultWriteTimeout = 2 * time.Second
)

var errInvalidURLFormat = errors.New("invalid URL format")

// Options provides setting for client connection to GRedis server
type Options struct {
	Host          string
	Port          string
	DB            int
	Password      string
	Timeout       time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	TraceProtocol bool
}

// NewOptions supported URLs are in any of these formats:
//	gredis://HOST[:PORT][?db=DATABASE[&password=PASSWORD]]
//	gredis://HOST[:PORT][?password=PASSWORD[&db=DATABASE]]
//	gredis://[:PASSWORD@]HOST[:PORT][/DATABASE]
//	gredis://[:PASSWORD@]HOST[:PORT][?db=DATABASE]
//	gredis://HOST[:PORT]/DATABASE[?password=PASSWORD]
func NewOptions(rawURL string) (*Options, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errInvalidURLFormat
	}

	if u.Scheme != "gredis" {
		return nil, fmt.Errorf("invalid gredis URL scheme: %s", u.Scheme)
	}

	opts := Options{
		Timeout:      defaultTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	opts.Host, opts.Port, err = net.SplitHostPort(u.Host)
	if err != nil {
		opts.Host = u.Host
		opts.Port = defaultPort
	}
	if opts.Host == "" {
		opts.Host = defaultHost
	}

	opts.Password = u.Query().Get("password")
	if u.User != nil {
		opts.Password, _ = u.User.Password()
	}

	db := u.Query().Get("db")
	if db != "" {
		opts.DB, err = strconv.Atoi(db)
		if err != nil {
			return nil, fmt.Errorf("invalid database: %s", db)
		}
	}

	if len(u.Path) > 1 {
		db := string(u.Path[1:])
		opts.DB, err = strconv.Atoi(db)
		if err != nil {
			return nil, fmt.Errorf("invalid database: %s", db)
		}
	}

	return &opts, nil
}
