package memcache

import (
	"fmt"
	"strings"
	"sync"

	mc "github.com/bradfitz/gomemcache/memcache"
)

func test() {
	fmt.Println("hello memcache")

	m := mc.New()
	m.Set(&mc.Item{Key: "foo", Value: []byte("my value")})

	it, _ := m.Get("foo")
	fmt.Println(string(it.Value))
}

var mcIns sync.Map

//Config config
type Config struct {
	NodeName string `json:"nodeName"`
	Addr     string `json:"addr"` // like "127.0.0.1:11211|127.0.0.1:11212|127.0.0.1:11213"
}

//Mao mao
type Mao struct {
	conn *mc.Client
}

//GetMc get mc
func GetMc(config *Config) (*Mao, error) {
	mao := new(Mao)

	if conn, ok := mcIns.Load(config.NodeName); ok {
		mao.conn = conn.(*mc.Client)
		return mao, nil
	}
	conn, err := NewMc(config)
	if err != nil {
		return nil, err
	}
	mcIns.Store(config.NodeName, conn)
	mao.conn = conn
	return mao, nil
}

//NewMc new mc
func NewMc(config *Config) (*mc.Client, error) {
	nodes := strings.Split(config.Addr, "|")
	if len(nodes) <= 0 {
		return nil, nil
	}
	m := mc.New(nodes...)
	if _, err := m.Get("foo"); err != nil {
		return nil, err
	}
	return m, nil
}

//Get get
func (mao *Mao) Get(key string) ([]byte, error) {
	item, err := mao.conn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

//Setex setex
func (mao *Mao) Setex(key string, value []byte, expire int32) error {
	item := &mc.Item{
		Key:        key,
		Value:      value,
		Expiration: expire,
	}
	return mao.conn.Set(item)
}

//Set set
func (mao *Mao) Set(key string, value []byte) error {
	item := &mc.Item{
		Key:   key,
		Value: value,
	}
	return mao.conn.Set(item)
}

//Del del
func (mao *Mao) Del(key string) error {
	return mao.conn.Delete(key)
}

//MGet mget
func (mao *Mao) MGet(keys []string) (map[string][]byte, error) {
	items, err := mao.conn.GetMulti(keys)
	if err != nil {
		return nil, err
	}
	res := make(map[string][]byte)
	for key, item := range items {
		res[key] = item.Value
	}
	return res, nil
}
