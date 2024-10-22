package callback_data

import (
	"errors"
	"strings"
)

var (
	errWrongPrefix       = errors.New("data string does not start with prefix")
	errWrongKeyValuePair = errors.New("length of split kv pair is not 2")
)

type CallbackData struct {
	Prefix string
	Data   map[string]string
	Sep    string
	KVSep  string
}

func (c *CallbackData) GetFrom(s string) error {
	c.Data = make(map[string]string)

	if !strings.HasPrefix(s, c.Prefix) {
		return errWrongPrefix
	}
	data := s[len(c.Prefix):]
	kvpairs := strings.Split(data, c.Sep)
	for _, pair := range kvpairs {
		kv := strings.Split(pair, c.KVSep)
		if len(kv) != 2 {
			return errWrongKeyValuePair
		}
		key, value := kv[0], kv[1]
		c.Data[key] = value
	}
	return nil
}

func (c *CallbackData) String() string {
	s := c.Prefix
	data := make([]string, 0, 5)
	for k, v := range c.Data {
		kv := strings.Join([]string{k, v}, c.KVSep)
		data = append(data, kv)
	}
	s += strings.Join(data, c.Sep)
	return s
}
