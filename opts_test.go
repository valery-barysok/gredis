package gredis

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestOptionsWithValidUrl(t *testing.T) {
	RegisterTestingT(t)

	successCases := []struct {
		url  string
		opts *Options
	}{
		{
			"gredis://HOST",
			&Options{
				Host:         "HOST",
				Port:         "16379",
				DB:           0,
				Password:     "",
				Timeout:      defaultTimeout,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
		},
		{
			"gredis://HOST:PORT",
			&Options{
				Host:         "HOST",
				Port:         "PORT",
				DB:           0,
				Password:     "",
				Timeout:      defaultTimeout,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
		},
		{
			"gredis://HOST:PORT?db=1",
			&Options{
				Host:         "HOST",
				Port:         "PORT",
				DB:           1,
				Password:     "",
				Timeout:      defaultTimeout,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
		},
		{
			"gredis://HOST:PORT?db=1&password=PASSWORD",
			&Options{
				Host:         "HOST",
				Port:         "PORT",
				DB:           1,
				Password:     "PASSWORD",
				Timeout:      defaultTimeout,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
		},
		{
			"gredis://HOST:PORT/2?db=1&password=PASSWORD",
			&Options{
				Host:         "HOST",
				Port:         "PORT",
				DB:           2,
				Password:     "PASSWORD",
				Timeout:      defaultTimeout,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
		},
		{
			"gredis://HOST?password=PASSWORD",
			&Options{
				Host:         "HOST",
				Port:         "16379",
				DB:           0,
				Password:     "PASSWORD",
				Timeout:      defaultTimeout,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
		},
		{
			"gredis://:PASSWORD1@HOST:PORT/2?db=1&password=PASSWORD2",
			&Options{
				Host:         "HOST",
				Port:         "PORT",
				DB:           2,
				Password:     "PASSWORD1",
				Timeout:      defaultTimeout,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
		},
	}

	for _, c := range successCases {
		opts, err := NewOptions(c.url)
		if Expect(err).ToNot(HaveOccurred(), c.url) {
			Expect(opts).To(Equal(c.opts), c.url)
		}
	}
}

func TestOptionsWithInvalidUrl(t *testing.T) {
	RegisterTestingT(t)

	failureCases := []struct {
		url string
		err string
	}{
		{
			"gredis://foo bar",
			"invalid URL format",
		},
		{
			"localhost",
			"invalid gredis URL scheme: ",
		},
		{
			"http://github.com/valery-barysok",
			"invalid gredis URL scheme: http",
		},
		{
			"gredis://localhost:6379/abc123",
			"invalid database: abc123",
		},
		{
			"gredis://HOST:PORT?db=DATABASE",
			"invalid database: DATABASE",
		},
	}

	for _, c := range failureCases {
		_, err := NewOptions(c.url)
		if Expect(err).To(HaveOccurred(), c.url) {
			Expect(err.Error()).To(Equal(c.err), c.url)
		}
	}
}
