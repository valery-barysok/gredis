package gredis

// List of key value commands
var (
	SetCommand = []byte("SET")
	GetCommand = []byte("GET")
	DelCommand = []byte("DEL")
)

// Set key to hold the string value. If key already holds a value, it is overwritten, regardless of its type.
// Any previous time to live associated with the key is discarded on successful `Set` operation.
//
// Returns true if success, otherwise false.
func (client *Client) Set(key string, value string) (bool, error) {
	_, err := client.Do(SetCommand, []byte(key), []byte(value))
	if err != nil {
		return false, err
	}

	return true, nil
}

// Get the value of key. If the key does not exist the special value nil is returned. An error is returned
// if the value stored at key is not a string, because `GET` only handles string values.
func (client *Client) Get(key string) ([]byte, error) {
	msg, err := client.Do(GetCommand, []byte(key))
	if err != nil {
		return nil, err
	}

	return msg.BulkString(), nil
}

// Del removes the specified keys. A key is ignored if it does not exist.
//
// Returns the number of keys that were removed.
func (client *Client) Del(key string, keys ...string) (int, error) {
	msg, err := client.Do(DelCommand, toBulkArray(keys, key)...)
	if err != nil {
		return 0, err
	}

	return msg.Int(), nil
}
