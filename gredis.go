package gredis

import (
	"github.com/valery-barysok/resp"
	"net"
	"os"
	"sync"
	"time"
)

var defaultProtocol *resp.Protocol
var defaultTraceProtocol *resp.Protocol

func init() {
	defaultProtocol = resp.NewProtocol()
	defaultTraceProtocol = resp.NewProtocolWithLogging(os.Stdout)
}

// Client is connection to GRedis server
type Client struct {
	opts *Options

	mu   sync.Mutex
	conn net.Conn
	r    *resp.Reader
	w    *resp.Writer
}

// Dial establish connection to GRedis server with specified options
func Dial(opts *Options) (*Client, error) {
	dialer := net.Dialer{
		Timeout: opts.Timeout,
	}

	conn, err := dialer.Dial("tcp", net.JoinHostPort(opts.Host, opts.Port))
	if err != nil {
		return nil, err
	}

	var protocol *resp.Protocol
	if opts.TraceProtocol {
		protocol = defaultTraceProtocol
	} else {
		protocol = defaultProtocol
	}

	client := &Client{
		opts: opts,
		conn: conn,
		r:    resp.NewReader(conn, protocol),
		w:    resp.NewWriter(conn, protocol),
	}

	if opts.Password != "" {
		if _, err := client.Auth(opts.Password); err != nil {
			conn.Close()
			return nil, err
		}
	}

	if opts.DB != 0 {
		if _, err := client.Select(opts.DB); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return client, nil
}

// Close flushes all pending writes and disconnect from GRedis server
func (client *Client) Close() {
	client.Flush()
	client.conn.Close()
}

// Send sends command to GRedis server
func (client *Client) Send(cmd []byte, args ...[]byte) error {
	err := client.w.WriteCmd(cmd, args...)
	if err != nil {
		return err
	}

	return client.Flush()
}

// Flush flushes all pending writes to GRedis server
func (client *Client) Flush() error {
	if client.opts.WriteTimeout != 0 {
		client.conn.SetWriteDeadline(time.Now().Add(client.opts.ReadTimeout))
		defer client.conn.SetWriteDeadline(time.Time{})
	}

	return client.w.Flush()
}

// Receive receives reply from GRedis server
func (client *Client) Receive() (*resp.Message, error) {
	if client.opts.ReadTimeout != 0 {
		client.conn.SetReadDeadline(time.Now().Add(client.opts.ReadTimeout))
		defer client.conn.SetReadDeadline(time.Time{})
	}

	msg, err := client.r.Read()
	if err != nil {
		return nil, err
	}

	if msg.IsError() {
		return nil, msg.Err()
	}

	return msg, nil
}

// Do sends command to GRedis server and receives reply from GRedis server
func (client *Client) Do(cmd []byte, args ...[]byte) (*resp.Message, error) {
	err := client.Send(cmd, args...)
	if err != nil {
		return nil, err
	}

	return client.Receive()
}

func toBulkArray(args []string, keys ...string) [][]byte {
	res := make([][]byte, 0, len(args)+len(keys))

	for _, value := range keys {
		res = append(res, []byte(value))
	}

	for _, value := range args {
		res = append(res, []byte(value))
	}

	return res
}
