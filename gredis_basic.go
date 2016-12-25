package gredis

import (
	"strconv"
)

// List of basic commands
var (
	AuthCommand     = []byte("AUTH")
	SelectCommand   = []byte("SELECT")
	EchoCommand     = []byte("ECHO")
	PingCommand     = []byte("PING")
	ShutdownCommand = []byte("SHUTDOWN")
	// TODO: use custom name "COMMANDS" instead of "COMMAND" due to incompatibility with redis-cli
	CommandCommand = []byte("COMMANDS")
	KeysCommand    = []byte("KEYS")
	ExistsCommand  = []byte("EXISTS")
	ExpireCommand  = []byte("EXPIRE")
)

// Auth requests for authentication in a password-protected GRedis server. GRedis can be instructed to
// require a password before allowing clients to execute commands.
//
// Returns true if success, otherwise false.
func (client *Client) Auth(password string) (bool, error) {
	_, err := client.Do(AuthCommand, []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

// Select the DB with having the specified zero-based numeric index.
//
// Returns true if success, otherwise false.
func (client *Client) Select(db int) (bool, error) {
	_, err := client.Do(SelectCommand, []byte(strconv.Itoa(db)))
	if err != nil {
		return false, err
	}

	return true, nil
}

// Echo returns a copy of the argument as a bulk if success, otherwise nil.
func (client *Client) Echo(message string) ([]byte, error) {
	msg, err := client.Do(EchoCommand, []byte(message))
	if err != nil {
		return nil, err
	}

	return msg.BulkString(), nil
}

// Ping returns `PONG` if success, otherwise empty string.
func (client *Client) Ping() (string, error) {
	msg, err := client.Do(PingCommand)
	if err != nil {
		return "", err
	}

	return msg.String(), nil
}

// PingMsg returns a copy of the argument as a bulk if success, otherwise nil.
func (client *Client) PingMsg(message string) ([]byte, error) {
	msg, err := client.Do(EchoCommand, []byte(message))
	if err != nil {
		return nil, err
	}

	return msg.BulkString(), nil
}

// Shutdown behavior is the following:
//  Send command to the server
//  Close connection to the server.
func (client *Client) Shutdown() error {
	err := client.Send(ShutdownCommand)
	if err != nil {
		return err
	}

	client.Close()
	return nil
}

// Command returns Bulk Array of all supported commands.
func (client *Client) Command() ([][]byte, error) {
	msg, err := client.Do(CommandCommand)
	if err != nil {
		return nil, err
	}

	arr := msg.Array()
	res := make([][]byte, 0, len(arr))
	for _, cmd := range arr {
		res = append(res, cmd.BulkString())
	}

	return res, nil
}

// Keys returns Bulk Array of all keys matching **regexp** pattern.
func (client *Client) Keys(pattern string) ([][]byte, error) {
	msg, err := client.Do(KeysCommand, []byte(pattern))
	if err != nil {
		return nil, err
	}

	arr := msg.Array()
	res := make([][]byte, 0, len(arr))
	for _, cmd := range arr {
		res = append(res, cmd.BulkString())
	}

	return res, nil
}

// Exists returns if keys exist with count of such keys.
//
// It is possible to specify multiple keys instead of a single one. In such a case, it returns the total
// number of keys existing.
//
// The user should be aware that if the same existing key is mentioned in the arguments multiple times,
// it will be counted multiple times. So if `somekey` exists, `Exists("somekey", "somekey")` will return 2.
func (client *Client) Exists(key string, keys ...string) (int, error) {
	msg, err := client.Do(ExistsCommand, toBulkArray(keys, key)...)
	if err != nil {
		return 0, err
	}

	return msg.Int(), nil
}

// Expire sets a timeout on key. After the timeout has expired, the key will automatically be deleted.
//  1 if the timeout was set.
//  0 if key does not exist or the timeout could not be set.
func (client *Client) Expire(key string, seconds int) (int, error) {
	msg, err := client.Do(ExpireCommand, []byte(key), []byte(strconv.Itoa(seconds)))
	if err != nil {
		return 0, err
	}

	return msg.Int(), nil
}
