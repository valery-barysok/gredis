# gredis

Go client for [GRedis](https://github.com/valery-barysok/gredisd)

[![License][License-Image]][License-Url] [![ReportCard][ReportCard-Image]][ReportCard-Url] [![Build Status][Travis-Image]][Travis-Url]


## GRedis API

##### NewOptions(rawURL string) (*Options, error)

  NewOptions supported URLs are in any of these formats:

  - gredis://HOST[:PORT][?db=DATABASE[&password=PASSWORD]]
  - gredis://HOST[:PORT][?password=PASSWORD[&db=DATABASE]]
  - gredis://[:PASSWORD@]HOST[:PORT][/DATABASE]
  - gredis://[:PASSWORD@]HOST[:PORT][?db=DATABASE]
  - gredis://HOST[:PORT]/DATABASE[?password=PASSWORD]
  
##### Dial(opts *Options) (*Client, error)

  Dial establish connection to GRedis server with specified options

## Client Low Level API

##### Send(cmd []byte, args ...[]byte) error

  Sends command to GRedis server

##### Receive() (*resp.Item, error)

  Receives reply from GRedis server

##### Do(cmd []byte, args ...[]byte) (*resp.Item, error)

  Sends command to GRedis server and receives reply from GRedis server

## Client High Level API

### Basic Commands

##### [**Auth(password string) (bool, error)**](https://github.com/valery-barysok/gredisd#auth-password)

  Request for authentication in a password-protected GRedis server. GRedis can be instructed to
  require a password before allowing clients to execute commands.

  Returns true if success, otherwise false.

##### [**Select(db int) (bool, error)**](https://github.com/valery-barysok/gredisd#select-index)

  Select the DB with having the specified zero-based numeric index.
  
  Returns true if success, otherwise false.

##### [**Echo(message string) ([]byte, error)**](https://github.com/valery-barysok/gredisd#echo-message)

  Returns a copy of the argument as a bulk if success, otherwise nil.

##### [**Ping() (string, error)**](https://github.com/valery-barysok/gredisd#ping-message)

  Returns `PONG` if success, otherwise empty string.

##### [**PingMsg(message string) ([]byte, error)**](https://github.com/valery-barysok/gredisd#ping-message)

  Returns a copy of the argument as a bulk if success, otherwise nil.

##### [**Shutdown() error**](https://github.com/valery-barysok/gredisd#shutdown)

  The command behavior is the following:
     Send command to the server
     Close connection to the server.

##### [**Command() ([][]byte, error)**](https://github.com/valery-barysok/gredisd#commands)

  Returns Bulk Array of all supported commands.

##### [**Keys(pattern string) ([][]byte, error)**](https://github.com/valery-barysok/gredisd#keys-pattern)

  Returns Bulk Array of all keys matching **regexp** pattern.

##### [**Exists(key string, keys ...string) (int, error)**](https://github.com/valery-barysok/gredisd#exists-key-key-)

  Returns if keys exist with count of such keys.

  It is possible to specify multiple keys instead of a single one. In such a case, it returns the total
  number of keys existing.

  The user should be aware that if the same existing key is mentioned in the arguments multiple times,
  it will be counted multiple times. So if `somekey` exists, `Exists("somekey", "somekey")` will return 2.
  
##### [**Expire(key string, seconds int) (int, error)**](https://github.com/valery-barysok/gredisd#expire-key-seconds)

  Expire sets a timeout on key. After the timeout has expired, the key will automatically be deleted.

  - 1 if the timeout was set.
  - 0 if key does not exist or the timeout could not be set.

### Key Value Commands

##### [**Set(key string, value string) (bool, error)**](https://github.com/valery-barysok/gredisd#set-key-value-ex-seconds-px-milliseconds-nxxx)

  Set key to hold the string value. If key already holds a value, it is overwritten, regardless of its type.
  Any previous time to live associated with the key is discarded on successful `Set` operation.

  Returns true if success, otherwise false.

##### [**Get(key string) ([]byte, error)**](https://github.com/valery-barysok/gredisd#get-key)

  Get the value of key. If the key does not exist the special value nil is returned. An error is returned
  if the value stored at key is not a string, because `GET` only handles string values.

##### [**Del(key string, keys ...string) (int, error)**](https://github.com/valery-barysok/gredisd#del-key-key-)

  Removes the specified keys. A key is ignored if it does not exist.

  Returns the number of keys that were removed.

### Key Value List Commands

##### [**LPush(key string, value string, values ...string) (int, error)**](https://github.com/valery-barysok/gredisd#lpush-key-value-value-)

  Insert all the specified values at the head of the list stored at key. If key does not exist, it is
  created as empty list before performing the push operations. When key holds a value that is not a list,
  an error is returned.

  It is possible to push multiple elements using a single command call just specifying multiple arguments
  at the end of the command. Elements are inserted one after the other to the head of the list, from the
  leftmost element to the rightmost element. So for instance the command `LPush("mylist", "a", "b", "c")`
  will result into a list containing `c` as first element, `b` as second element and `a` as third element.
  
  Returns the length of the list after the push operations.

##### [**RPush(key string, value string, values ...string) (int, error)**](https://github.com/valery-barysok/gredisd#rpush-key-value-value-)

  Insert all the specified values at the tail of the list stored at key. If key does not exist, it is
  created as empty list before performing the push operation. When key holds a value that is not a list,
  an error is returned.

  It is possible to push multiple elements using a single command call just specifying multiple arguments
  at the end of the command. Elements are inserted one after the other to the tail of the list, from the
  leftmost element to the rightmost element. So for instance the command `RPush("mylist", "a", "b", "c")`
  will result into a list containing `a` as first element, `b` as second element and `c` as third element.
  
  Returns the length of the list after the push operations.

##### [**LPop(key string) ([]byte, error)**](https://github.com/valery-barysok/gredisd#lpop-key)

  Removes and returns the first element of the list stored at key.

##### [**RPop(key string) ([]byte, error)**](https://github.com/valery-barysok/gredisd#rpop-key)

  Removes and returns the last element of the list stored at key.

##### [**LLen(key string) (int, error)**](https://github.com/valery-barysok/gredisd#llen-key)

  Returns the length of the list stored at key. If key does not exist, it is interpreted as an empty list
  and `0` is returned. An error is returned when the value stored at key is not a list.

##### [**LInsert(key string, before bool, pivot string, value string) (int, error)**](https://github.com/valery-barysok/gredisd#linsert-key-beforeafter-pivot-value)

  Inserts value in the list stored at key either before or after the reference value pivot.

  When key does not exist, it is considered an empty list and no operation is performed.

  An error is returned when key exists but does not hold a list value.

##### [**LIndex(key string, index int) ([]byte, error)**](https://github.com/valery-barysok/gredisd#lindex-key-index)

  Returns the element at index index in the list stored at key. The index is zero-based, so 0 means the
  first element, 1 the second element and so on. Negative indices can be used to designate elements
  starting at the tail of the list. Here, -1 means the last element, -2 means the penultimate and so
  forth.

  When the value at key is not a list, an error is returned.

##### [**LRange(key string, start int, stop int) ([][]byte, error)**](https://github.com/valery-barysok/gredisd#lrange-key-start-stop)

  Returns the specified elements of the list stored at key. The offsets start and stop are zero-based
  indexes, with 0 being the first element of the list (the head of the list), 1 being the next element
  and so on.

  These offsets can also be negative numbers indicating offsets starting at the end of the list.
  For example, -1 is the last element of the list, -2 the penultimate, and so on.

### Key Value Dict Commands

##### [**HSet(key string, field string, value string) (int, error)**](https://github.com/valery-barysok/gredisd#hset-key-field-value)

  Sets field in the hash stored at key to value. If key does not exist, a new key holding a hash is
  created. If field already exists in the hash, it is overwritten.
  
  - 1 if field is a new field in the hash and value was set.
  - 0 if field already exists in the hash and the value was updated.

##### [**HGet(key string, field string) ([]byte, error)**](https://github.com/valery-barysok/gredisd#hget-key-field)

  Returns the value associated with field in the hash stored at key.

##### [**HDel(key string, field string, fields ...string) (int, error)**](https://github.com/valery-barysok/gredisd#hdel-key-field-field-)

  Removes the specified fields from the hash stored at key. Specified fields that do not exist
  within this hash are ignored. If key does not exist, it is treated as an empty hash and this
  command returns 0.

##### [**HLen(key string) (int, error)**](https://github.com/valery-barysok/gredisd#hlen-key)

  Returns the number of fields contained in the hash stored at key or 0 when key does not exist.

##### [**HExists(key string, field string) (int, error)**](https://github.com/valery-barysok/gredisd#hexists-key-field)

  Returns if field is an existing field in the hash stored at key.
  
  - 1 if the hash contains field.
  - 0 if the hash does not contain field, or key does not exist.

[License-Url]: http://opensource.org/licenses/Apache-2.0
[License-Image]: https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square
[ReportCard-Url]: http://goreportcard.com/report/valery-barysok/gredis
[ReportCard-Image]: http://goreportcard.com/badge/github.com/valery-barysok/gredis?style=flat-square
[Travis-Image]: https://img.shields.io/travis/valery-barysok/gredis/master.svg?style=flat-square
[Travis-Url]: https://travis-ci.org/valery-barysok/gredis