package redis

import (
	"fmt"
	"os"
	"testing"
)

func Test_test(t *testing.T) {

	config := &Config{
		NodeName: "default",
		Host:     "127.0.0.1",
		Port:     6379,
		Database: 0,
		Password: "",
	}

	rao, _ := GetRedis(config)
	defer rao.Close()

	//pipeline

	var argArray = []map[string][]interface{}{
		{
			"SETEX": {"test_1", 60, 1},
		},
		{
			"SETEX": {"test_2", 60, 2},
		},
	}

	err1 := rao.Pipeline(argArray)
	fmt.Println(err1)

	os.Exit(0)

	//set
	err := rao.Set("111_222", 123)

	rao1, _ := GetPool(config)
	//set
	err = rao1.Set("222_111", 123)

	//del
	rao.Del("111_222")

	//exists
	val1, err := rao.Exists("111_222")
	fmt.Println(val1, err)

	//expire
	rao.Expire("111_222", 100)

	//ttl
	val2, err := rao.TTL("111_2244")
	fmt.Println(val2, err)

	//setex
	err = rao.Setex("111_333", 100, "setex")

	//mset
	args := []interface{}{"m1", 111, "m2", 222, "m3", 333}
	rao.MSet(args)

	//mget
	val3, _ := rao.MGet([]interface{}{"m1", "m2", "m3", "m4"})
	fmt.Println(val3)

	//incr
	_, err = rao.Incr("111_444")
	//fmt.Println(err)

	//incrby
	_, err = rao.IncrBy("111_444", 100)
	//fmt.Println(err)

	//incr
	_, err = rao.Decr("111_555")
	fmt.Println(err)

	//decrby
	_, err = rao.DecrBy("111_555", 200)
	fmt.Println(err)

	//hset
	err = rao.HSet("my_h", "m3", 333)
	fmt.Println(err)

	//hget
	val4, _ := rao.HGet("my_h", "m3")
	fmt.Println(val4)

	//hmset
	err = rao.HMSet("my_h", []interface{}{"key1", 1, "key2", 2})
	fmt.Println(err)

	//hdel
	//err = rao.HDel("my_h", []interface{}{"key1", "key2", "key100"})

	//hget
	val5, _ := rao.HMGet("my_h", []interface{}{"key1", "key2", "key100"})
	fmt.Println(val5)

	//hgetall
	val6, _ := rao.HGetAll("my_h")
	fmt.Println(val6)

	//hlen
	val7, _ := rao.HLen("my_h")
	fmt.Println(val7)

	//lpush
	args = []interface{}{"l1", "l2", "l3", "l4", 100}
	err = rao.LPush("my_l", args)
	fmt.Println(err)

	//lpop
	val8, _ := rao.LPop("my_l")
	fmt.Println(val8)

	//rpush
	args = []interface{}{"l5", "l6", 0.1}
	err = rao.RPush("my_l", args)
	fmt.Println(err)

	//rpop
	val9, _ := rao.RPop("my_l")
	fmt.Println(val9)

	//llen
	val10, _ := rao.LLen("my_l")
	fmt.Println(val10)

	//lrange
	val11, _ := rao.LRange("my_l", 0, 5)
	fmt.Println(val11)

	//sadd
	args = []interface{}{"s1", "s2", "s3"}
	err = rao.SAdd("my_s", args)
	fmt.Println(err)

	//smembers
	val12, _ := rao.SMembers("my_s")
	fmt.Println(val12)

	//scard
	val13, _ := rao.SCard("my_s")
	fmt.Println(val13)

	//srem
	err = rao.SRem("my_s", []interface{}{"s2", "s1"})
	fmt.Println(err)

	//sismember
	val14, _ := rao.SIsMember("my_s", "s3")
	fmt.Println(val14)

	//zadd
	err = rao.ZAdd("my_z", 100, "z1")
	err = rao.ZAdd("my_z", 101, "z2")

	//zrem
	//err = rao.ZRem("my_z", []interface{}{"z1", "z3"})

	//zcard
	val15, _ := rao.ZCard("my_z")
	fmt.Println(val15)

	//zrange
	val16, val19, _ := rao.ZRange("my_z", 0, -1, true)
	fmt.Println(val16, val19)

	//zrevrange
	val17, val18, _ := rao.ZRevRange("my_z", 0, -1, true)
	fmt.Println(val17, val18)

	//zrangebyscore
	val20, val21, _ := rao.ZRangeByScore("my_z", 0, 222, true)
	fmt.Println(val20, val21)

	//zrevrangebyscore
	val22, val23, _ := rao.ZRevRangeByScore("my_z", 0, 222, true)
	fmt.Println(val22, val23)

}
