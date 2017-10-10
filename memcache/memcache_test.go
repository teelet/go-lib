package memcache

import (
	"fmt"
	"testing"
	"time"
)

func Test_test(t *testing.T) {
	config := &Config{
		NodeName: "default",
		Addr:     "10.16.133.61:11211|10.16.133.61:11211",
	}
	mc, _ := GetMc(config)
	res1, _ := mc.Get("foo")
	fmt.Println(string(res1))

	mc.Set("foo1", []byte("foo1111"))

	mc.Setex("foo2", []byte("foo222"), 3)

	res2, _ := mc.MGet([]string{"foo", "foo1", "foo2"})
	for k, v := range res2 {
		fmt.Println(k, string(v))
	}

	time.Sleep(5 * time.Second)
	fmt.Println("==============")
	res3, _ := mc.MGet([]string{"foo", "foo1", "foo2"})
	for k, v := range res3 {
		fmt.Println(k, string(v))
	}

	mc.Del("foo1")
	fmt.Println("==============")
	res4, _ := mc.MGet([]string{"foo", "foo1", "foo2"})
	for k, v := range res4 {
		fmt.Println(k, string(v))
	}
}
