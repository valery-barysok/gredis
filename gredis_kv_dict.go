package gredis

// List of key value dict commands
var (
	HSetCommand    = []byte("HSET")
	HGetCommand    = []byte("HGET")
	HDelCommand    = []byte("HDEL")
	HLenCommand    = []byte("HLEN")
	HExistsCommand = []byte("HEXISTS")
)

// HSet sets field in the hash stored at key to value. If key does not exist, a new key holding a hash is
// created. If field already exists in the hash, it is overwritten.
//  1 if field is a new field in the hash and value was set.
//  0 if field already exists in the hash and the value was updated.
func (client *Client) HSet(key string, field string, value string) (int, error) {
	item, err := client.Do(HSetCommand, []byte(key), []byte(field), []byte(value))
	if err != nil {
		return 0, err
	}

	return item.Int(), nil
}

// HGet returns the value associated with field in the hash stored at key.
func (client *Client) HGet(key string, field string) ([]byte, error) {
	item, err := client.Do(HGetCommand, []byte(key), []byte(field))
	if err != nil {
		return nil, err
	}

	if item.IsNil() {
		return nil, nil
	}

	return item.BulkString(), nil
}

// HDel removes the specified fields from the hash stored at key. Specified fields that do not exist
// within this hash are ignored. If key does not exist, it is treated as an empty hash and this
// command returns 0.
func (client *Client) HDel(key string, field string, fields ...string) (int, error) {
	item, err := client.Do(HDelCommand, toBulkArray(fields, key, field)...)
	if err != nil {
		return 0, err
	}

	return item.Int(), nil
}

// HLen returns the number of fields contained in the hash stored at key or 0 when key does not exist.
func (client *Client) HLen(key string) (int, error) {
	item, err := client.Do(HLenCommand, []byte(key))
	if err != nil {
		return 0, err
	}

	return item.Int(), nil
}

// HExists returns if field is an existing field in the hash stored at key.
//  1 if the hash contains field.
//  0 if the hash does not contain field, or key does not exist.
func (client *Client) HExists(key string, field string) (int, error) {
	item, err := client.Do(HExistsCommand, []byte(key), []byte(field))
	if err != nil {
		return 0, err
	}

	return item.Int(), nil
}
