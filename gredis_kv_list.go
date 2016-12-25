package gredis

import "strconv"

// List of key value list commands
var (
	LPushCommand   = []byte("LPUSH")
	RPushCommand   = []byte("RPUSH")
	LPopCommand    = []byte("LPOP")
	RPopCommand    = []byte("RPOP")
	LLenCommand    = []byte("LLEN")
	LInsertCommand = []byte("LINSERT")
	LIndexCommand  = []byte("LINDEX")
	LRangeCommand  = []byte("LRANGE")
)

var (
	insertBefore = []byte("BEFORE")
	insertAfter  = []byte("AFTER")
)

// LPush inserts all the specified values at the head of the list stored at key. If key does not exist, it is
// created as empty list before performing the push operations. When key holds a value that is not a list,
// an error is returned.
//
// It is possible to push multiple elements using a single command call just specifying multiple arguments
// at the end of the command. Elements are inserted one after the other to the head of the list, from the
// leftmost element to the rightmost element. So for instance the command `LPush("mylist", "a", "b", "c")`
// will result into a list containing `c` as first element, `b` as second element and `a` as third element.
//
// Returns the length of the list after the push operations.
func (client *Client) LPush(key string, value string, values ...string) (int, error) {
	msg, err := client.Do(LPushCommand, toBulkArray(values, key, value)...)
	if err != nil {
		return 0, err
	}

	return msg.Int(), nil
}

// RPush inserts all the specified values at the tail of the list stored at key. If key does not exist, it is
//created as empty list before performing the push operation. When key holds a value that is not a list,
//an error is returned.
//
//It is possible to push multiple elements using a single command call just specifying multiple arguments
//at the end of the command. Elements are inserted one after the other to the tail of the list, from the
//leftmost element to the rightmost element. So for instance the command `RPush("mylist", "a", "b", "c")`
//will result into a list containing `a` as first element, `b` as second element and `c` as third element.
//
//Returns the length of the list after the push operations.
func (client *Client) RPush(key string, value string, values ...string) (int, error) {
	msg, err := client.Do(RPushCommand, toBulkArray(values, key, value)...)
	if err != nil {
		return 0, err
	}

	return msg.Int(), nil
}

// LPop removes and returns the first element of the list stored at key.
func (client *Client) LPop(key string) ([]byte, error) {
	msg, err := client.Do(LPopCommand, []byte(key))
	if err != nil {
		return nil, err
	}

	return msg.BulkString(), nil
}

// RPop removes and returns the last element of the list stored at key.
func (client *Client) RPop(key string) ([]byte, error) {
	msg, err := client.Do(RPopCommand, []byte(key))
	if err != nil {
		return nil, err
	}

	return msg.BulkString(), nil
}

// LLen returns the length of the list stored at key. If key does not exist, it is interpreted as an empty list
// and `0` is returned. An error is returned when the value stored at key is not a list.
func (client *Client) LLen(key string) (int, error) {
	msg, err := client.Do(LLenCommand, []byte(key))
	if err != nil {
		return 0, err
	}

	return msg.Int(), nil
}

// LInsert inserts value in the list stored at key either before or after the reference value pivot.
//
// When key does not exist, it is considered an empty list and no operation is performed.
//
// An error is returned when key exists but does not hold a list value.
func (client *Client) LInsert(key string, before bool, pivot string, value string) (int, error) {
	place := insertBefore
	if !before {
		place = insertAfter
	}

	msg, err := client.Do(LInsertCommand, []byte(key), place, []byte(pivot), []byte(value))
	if err != nil {
		return 0, err
	}

	return msg.Int(), nil
}

// LIndex returns the element at index index in the list stored at key. The index is zero-based, so 0 means the
// first element, 1 the second element and so on. Negative indices can be used to designate elements
// starting at the tail of the list. Here, -1 means the last element, -2 means the penultimate and so
// forth.
//
// When the value at key is not a list, an error is returned.
func (client *Client) LIndex(key string, index int) ([]byte, error) {
	msg, err := client.Do(LIndexCommand, []byte(key), []byte(strconv.Itoa(index)))
	if err != nil {
		return nil, err
	}

	return msg.BulkString(), nil
}

// LRange returns the specified elements of the list stored at key. The offsets start and stop are zero-based
// indexes, with 0 being the first element of the list (the head of the list), 1 being the next element
// and so on.
//
// These offsets can also be negative numbers indicating offsets starting at the end of the list.
// For example, -1 is the last element of the list, -2 the penultimate, and so on.
func (client *Client) LRange(key string, start int, stop int) ([][]byte, error) {
	msg, err := client.Do(LRangeCommand, []byte(key), []byte(strconv.Itoa(start)), []byte(strconv.Itoa(stop)))
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
