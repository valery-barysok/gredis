package gredis

import (
	. "github.com/onsi/gomega"
	"testing"

	"github.com/valery-barysok/gredisd/app"
	"github.com/valery-barysok/gredisd/gredisd-app"
)

func TestValidDial(t *testing.T) {
	RegisterTestingT(t)

	gApp := gredisd.NewApp(&app.Options{})
	go gApp.Run()
	defer gApp.Shutdown()

	opts, err := NewOptions("gredis://localhost")
	if Expect(err).ToNot(HaveOccurred()) {
		_, err := Dial(opts)
		Expect(err).ToNot(HaveOccurred())
	}
}

func TestValidDialWithAuthAndDatabase(t *testing.T) {
	RegisterTestingT(t)

	gredisApp := gredisd.NewApp(&app.Options{
		Auth: "password",
	})
	go gredisApp.Run()
	defer gredisApp.Shutdown()

	opts, err := NewOptions("gredis://:password@localhost/12?password=ignored&db=1")
	if Expect(err).ToNot(HaveOccurred()) {
		_, err := Dial(opts)
		Expect(err).ToNot(HaveOccurred())
	}
}

func TestIntegrationForAllCommandsAtOnce(t *testing.T) {
	RegisterTestingT(t)

	gApp := gredisd.NewApp(&app.Options{
		Auth:          "password",
		TraceProtocol: true,
	})
	go gApp.Run()
	defer gApp.Shutdown()

	opts, err := NewOptions("gredis://:password@localhost:16379/12?password=ignored&db=1")
	if Expect(err).ToNot(HaveOccurred()) {
		client, err := Dial(opts)
		Expect(err).ToNot(HaveOccurred())

		pingRes, err := client.Ping()
		Expect(err).ToNot(HaveOccurred())
		Expect(pingRes).To(Equal("PONG"))

		msg := "test echo cmd"
		bulkRes, err := client.PingMsg(msg)
		Expect(err).ToNot(HaveOccurred())
		Expect(bulkRes).To(BeEquivalentTo(msg))

		bulkRes, err = client.Echo(msg)
		Expect(err).ToNot(HaveOccurred())
		Expect(bulkRes).To(BeEquivalentTo(msg))

		list, err := client.Command()
		Expect(err).ToNot(HaveOccurred())
		Expect(list).ToNot(Equal(nil))

		// Valid regexp
		list, err = client.Keys(".*")
		Expect(err).ToNot(HaveOccurred())
		if Expect(list).ToNot(Equal(nil)) {
			Expect(len(list)).To(Equal(0))
		}

		// Invalid regexp
		list, err = client.Keys(")")
		Expect(err).To(HaveOccurred())
		if Expect(list).ToNot(Equal(nil)) {
			Expect(len(list)).To(Equal(0))
		}

		key := "key"
		exists, err := client.Exists(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		keyValue := "key_value"
		success, err := client.Set(key, keyValue)
		Expect(err).ToNot(HaveOccurred())
		Expect(success).To(Equal(true))

		exists, err = client.Exists(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(1))

		exists, err = client.Exists(key, key, key)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(3))

		success, err = client.Select(opts.DB + 1)
		Expect(err).ToNot(HaveOccurred())
		Expect(success).To(Equal(true))

		exists, err = client.Exists(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		success, err = client.Select(opts.DB)
		Expect(err).ToNot(HaveOccurred())
		Expect(success).To(Equal(true))

		exists, err = client.Exists(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).ToNot(Equal(0))

		value, err := client.Get(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(keyValue))

		keys, err := client.Keys(".*")
		Expect(err).ToNot(HaveOccurred())
		if Expect(keys).To(HaveLen(1)) {
			Expect(keys[0]).To(BeEquivalentTo(key))
		}

		cnt, err := client.Del(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(1))

		exists, err = client.Exists(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		listKey := "list_key"
		listKeyValue1 := "list_key_value1"
		listKeyValue2 := "list_key_value2"
		listKeyValue3 := "list_key_value3"
		listKeyValue4 := "list_key_value4"
		listKeyValue5 := "list_key_value5"
		listKeyValue6 := "list_key_value6"
		listKeyValue7 := "list_key_value7"
		listKeyValue8 := "list_key_value8"

		cnt, err = client.LPush(listKey, listKeyValue1, listKeyValue2, listKeyValue3)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(3))

		cnt, err = client.RPush(listKey, listKeyValue4, listKeyValue5, listKeyValue6)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(6))

		value, err = client.LPop(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(listKeyValue3))

		value, err = client.RPop(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(listKeyValue6))

		l, err := client.LLen(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(l).To(Equal(4))

		cnt, err = client.LInsert(listKey, true, listKeyValue2, listKeyValue7)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(5))

		cnt, err = client.LInsert(listKey, false, listKeyValue2, listKeyValue8)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(6))

		value, err = client.LIndex(listKey, 0)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(listKeyValue7))

		value, err = client.LIndex(listKey, 2)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(listKeyValue8))

		values, err := client.LRange(listKey, 0, -1)
		Expect(err).ToNot(HaveOccurred())
		if Expect(values).To(HaveLen(6)) {
			Expect(values).To(Equal([][]byte{
				[]byte(listKeyValue7),
				[]byte(listKeyValue2),
				[]byte(listKeyValue8),
				[]byte(listKeyValue1),
				[]byte(listKeyValue4),
				[]byte(listKeyValue5),
			}))
		}

		keys, err = client.Keys(".*")
		Expect(err).ToNot(HaveOccurred())
		if Expect(keys).To(HaveLen(1)) {
			Expect(keys[0]).To(BeEquivalentTo(listKey))
		}

		cnt, err = client.Del(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(1))

		exists, err = client.Exists(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		dictKey := "dict_key"
		dictKeyField := "dict_key_field"
		dictKeyField2 := "dict_key_field2"
		dictKeyFieldValue := "dict_key_field_value"
		dictKeyFieldValue2 := "dict_key_field_value2"

		exists, err = client.HExists(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		inserted, err := client.HSet(dictKey, dictKeyField, dictKeyFieldValue)
		Expect(err).ToNot(HaveOccurred())
		Expect(inserted).To(Equal(1))

		updated, err := client.HSet(dictKey, dictKeyField, dictKeyFieldValue)
		Expect(err).ToNot(HaveOccurred())
		Expect(updated).To(Equal(0))

		inserted, err = client.HSet(dictKey, dictKeyField2, dictKeyFieldValue2)
		Expect(err).ToNot(HaveOccurred())
		Expect(inserted).To(Equal(1))

		updated, err = client.HSet(dictKey, dictKeyField2, dictKeyFieldValue2)
		Expect(err).ToNot(HaveOccurred())
		Expect(updated).To(Equal(0))

		exists, err = client.HExists(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(1))

		value, err = client.HGet(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(dictKeyFieldValue))

		deleted, err := client.HDel(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(deleted).To(Equal(1))

		exists, err = client.HExists(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		keys, err = client.Keys(".*")
		Expect(err).ToNot(HaveOccurred())
		if Expect(keys).To(HaveLen(1)) {
			Expect(keys[0]).To(BeEquivalentTo(dictKey))
		}

		cnt, err = client.Del(dictKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(1))

		exists, err = client.Exists(dictKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		exists, err = client.HExists(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		client.Close()
	}
}
