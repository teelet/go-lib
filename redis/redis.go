package redis

import (
	"fmt"
	"sync"
	"time"

	rds "github.com/garyburd/redigo/redis"
)

var redisIns sync.Map

//Rao rao
type Rao struct {
	conn *rds.Conn
	pool *rds.Pool
}

//Config config
type Config struct {
	NodeName        string `json:"nodeName"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Password        string `json:"password"`
	Database        int    `json:"database"`
	MaxOpen         int    `json:"maxOpen"`
	MaxIdle         int    `json:"maxIdle"`
	ConnIdleTimeout int    `json:"connIdleTimeout"`
}

//GetRedis get redis
func GetRedis(config *Config) (*Rao, error) {
	rao := new(Rao)
	conn, err := NewRedis(config)
	if err != nil {
		return nil, err
	}
	rao.conn = conn
	rao.pool = nil
	return rao, nil
}

//NewRedis new redis
func NewRedis(config *Config) (*rds.Conn, error) {
	conn, err := rds.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return nil, err
	}
	if config.Password != "" {
		if _, err := conn.Do("AUTH", config.Password); err != nil {
			conn.Close()
			return nil, err
		}
	}
	if _, err := conn.Do("SELECT", config.Database); err != nil {
		conn.Close()
		return nil, err
	}
	return &conn, nil
}

//GetPool get pool
func GetPool(config *Config) (*Rao, error) {
	rao := new(Rao)
	nodeName := config.NodeName + "_pool"
	if pool, ok := redisIns.Load(nodeName); ok {
		rao.pool = pool.(*rds.Pool)
		rao.conn = nil
		return rao, nil
	}
	pool := NewPool(config)
	redisIns.Store(nodeName, pool)
	rao.pool = pool
	rao.conn = nil
	return rao, nil
}

//NewPool new pool
func NewPool(config *Config) *rds.Pool {
	dialFunc := func() (rds.Conn, error) {
		conn, err := rds.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
		if err != nil {
			return nil, err
		}
		if config.Password != "" {
			if _, err := conn.Do("AUTH", config.Password); err != nil {
				conn.Close()
				return nil, err
			}
		}
		if _, err := conn.Do("SELECT", config.Database); err != nil {
			conn.Close()
			return nil, err
		}
		return conn, nil
	}
	pool := &rds.Pool{
		MaxActive:   config.MaxOpen,
		MaxIdle:     config.MaxIdle,
		IdleTimeout: time.Duration(config.ConnIdleTimeout) * time.Second,
		Dial:        dialFunc,
	}
	return pool
}

//Close close
func (rao *Rao) Close() error {
	if rao.conn != nil {
		if err := (*(rao.conn)).Close(); err != nil {
			return err
		}
	}
	return nil
}

//Get get
func (rao *Rao) Get(key string) (interface{}, error) {
	var conn rds.Conn
	var value interface{}
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("GET", key); err != nil {
		return nil, err
	}
	switch value.(type) {
	case nil:
		return nil, nil
	default:
		return string(value.([]byte)), nil
	}
}

//Set set
func (rao *Rao) Set(key string, val interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("SET", key, val); err != nil {
		return err
	}
	return nil
}

//Del del
func (rao *Rao) Del(key string) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("DEL", key); err != nil {
		return err
	}
	return nil
}

//Exists exists
func (rao *Rao) Exists(key string) (bool, error) {
	var conn rds.Conn
	var value interface{}
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("EXISTS", key); err != nil {
		return false, err
	}
	if value.(int64) == 1 {
		return true, nil
	}
	return false, nil
}

//Expire expire
func (rao *Rao) Expire(key string, expire interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("EXPIRE", key, expire); err != nil {
		return err
	}
	return nil
}

//TTL ttl
func (rao *Rao) TTL(key string) (int64, error) {
	var conn rds.Conn
	var value interface{}
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("TTL", key); err != nil {
		return -2, err
	}
	return value.(int64), nil
}

//Setex setex
func (rao *Rao) Setex(key string, expire int64, val interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("SETEX", key, expire, val); err != nil {
		return err
	}
	return nil
}

//MSet mset
func (rao *Rao) MSet(args []interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("MSET", args...); err != nil {
		return err
	}
	return nil
}

//MGet mget
func (rao *Rao) MGet(keys []interface{}) (map[string]interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	var tmp2 []interface{}
	var value = make(map[string]interface{}, len(keys))
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if tmp, err = conn.Do("MGET", keys...); err != nil {
		return nil, err
	}
	tmp2 = tmp.([]interface{})
	for k, v := range keys {
		kk := v.(string)
		if tmp2[k] != nil {
			value[kk] = string(tmp2[k].([]byte))
		} else {
			value[kk] = nil
		}
	}
	return value, nil
}

//Incr incr
func (rao *Rao) Incr(key string) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("INCR", key); err != nil {
		return err
	}
	return nil
}

//IncrBy incrby
func (rao *Rao) IncrBy(key string, num int64) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("INCRBY", key, num); err != nil {
		return err
	}
	return nil
}

//Decr decr
func (rao *Rao) Decr(key string) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("DECR", key); err != nil {
		return err
	}
	return nil
}

//DecrBy decrby
func (rao *Rao) DecrBy(key string, num int64) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("DECRBY", key, num); err != nil {
		return err
	}
	return nil
}

//HSet hset
func (rao *Rao) HSet(key string, field string, val interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("HSET", key, field, val); err != nil {
		return err
	}
	return nil
}

//HGet hget
func (rao *Rao) HGet(key string, field string) (interface{}, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("HGET", key, field); err != nil {
		return nil, err
	}
	return rtti(value), nil
}

//HMSet hmset
func (rao *Rao) HMSet(key string, args []interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	args = append([]interface{}{key}, args...)
	if _, err = conn.Do("HMSET", args...); err != nil {
		return err
	}
	return nil
}

//HMGet hmget
func (rao *Rao) HMGet(key string, fields []interface{}) (interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	var tmp2 []interface{}
	var value = make(map[string]interface{}, len(fields))
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	args := append([]interface{}{key}, fields...)
	if tmp, err = conn.Do("HMGET", args...); err != nil {
		return nil, err
	}
	tmp2 = tmp.([]interface{})
	for k, v := range fields {
		kk := v.(string)
		if tmp2[k] != nil {
			value[kk] = rtti(tmp2[k])
		} else {
			value[kk] = nil
		}
	}
	return value, nil
}

//HDel hdel
func (rao *Rao) HDel(key string, fields []interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	args := append([]interface{}{key}, fields...)
	if _, err = conn.Do("HDEL", args...); err != nil {
		return err
	}
	return nil
}

//HGetAll hgetall
func (rao *Rao) HGetAll(key string) (interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	var value = make(map[string]interface{})
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if tmp, err = conn.Do("HGETALL", key); err != nil {
		return nil, err
	}
	tmp2 := tmp.([]interface{})
	for k, v := range tmp2 {
		if k%2 != 0 {
			continue
		}
		kk := string(v.([]byte))
		value[kk] = rtti(tmp2[k+1])
	}
	return value, nil
}

//HLen hlen
func (rao *Rao) HLen(key string) (int64, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("HLEN", key); err != nil {
		return -2, err
	}
	return value.(int64), nil
}

//LPush lpush
func (rao *Rao) LPush(key string, args []interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	args = append([]interface{}{key}, args...)
	if _, err = conn.Do("LPUSH", args...); err != nil {
		return err
	}
	return nil
}

//LPop lpop
func (rao *Rao) LPop(key string) (interface{}, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("LPOP", key); err != nil {
		return nil, err
	}
	return rtti(value), nil
}

//RPush rpush
func (rao *Rao) RPush(key string, args []interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	args = append([]interface{}{key}, args...)
	if _, err = conn.Do("RPUSH", args...); err != nil {
		return err
	}
	return nil
}

//RPop rpop
func (rao *Rao) RPop(key string) (interface{}, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("RPOP", key); err != nil {
		return nil, err
	}
	return rtti(value), nil
}

//LLen llen
func (rao *Rao) LLen(key string) (int64, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("LLEN", key); err != nil {
		return -2, err
	}
	return value.(int64), nil
}

//LRange lrange
func (rao *Rao) LRange(key string, start, stop int64) (interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	var value []interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if tmp, err = conn.Do("LRANGE", key, start, stop); err != nil {
		return nil, err
	}
	for _, v := range tmp.([]interface{}) {
		value = append(value, rtti(v))
	}
	return value, nil
}

//LGet lget
func (rao *Rao) LGet(key string, index int) (interface{}, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("LGET", key, index); err != nil {
		return nil, err
	}
	return rtti(value), nil
}

//RGet rget
func (rao *Rao) RGet(key string, index int) (interface{}, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("RGET", key, index); err != nil {
		return nil, err
	}
	return rtti(value), nil
}

//SAdd sadd
func (rao *Rao) SAdd(key string, args []interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	args = append([]interface{}{key}, args...)
	if _, err = conn.Do("SADD", args...); err != nil {
		return err
	}
	return nil
}

//SMembers smembers
func (rao *Rao) SMembers(key string) ([]interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	var value []interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if tmp, err = conn.Do("SMEMBERS", key); err != nil {
		return nil, err
	}
	for _, v := range tmp.([]interface{}) {
		value = append(value, rtti(v))
	}
	return value, nil
}

//SCard scard
func (rao *Rao) SCard(key string) (int64, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("SCARD", key); err != nil {
		return -2, err
	}
	return value.(int64), nil
}

//SRem srem
func (rao *Rao) SRem(key string, args []interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	args = append([]interface{}{key}, args...)
	if _, err = conn.Do("SREM", args...); err != nil {
		return err
	}
	return nil
}

//SIsMember sismember
func (rao *Rao) SIsMember(key string, val interface{}) (bool, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("SISMEMBER", key, val); err != nil {
		return false, err
	}
	if value.(int64) == 1 {
		return true, nil
	}

	return false, nil
}

//ZAdd zadd
func (rao *Rao) ZAdd(key string, score, val interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if _, err = conn.Do("ZADD", key, score, val); err != nil {
		return err
	}
	return nil
}

//ZCard zcard
func (rao *Rao) ZCard(key string) (int64, error) {
	var conn rds.Conn
	var err error
	var value interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if value, err = conn.Do("ZCARD", key); err != nil {
		return -2, err
	}
	return value.(int64), nil
}

//ZRem zrem
func (rao *Rao) ZRem(key string, args []interface{}) error {
	var conn rds.Conn
	var err error
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	args = append([]interface{}{key}, args...)
	if _, err = conn.Do("ZREM", args...); err != nil {
		return err
	}
	return nil
}

//ZRange zragne
func (rao *Rao) ZRange(key string, start, stop int64, withScores bool) ([]interface{}, map[string]interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if withScores == true {
		var value = make(map[string]interface{})
		var keys []interface{}
		if tmp, err = conn.Do("ZRANGE", key, start, stop, "WITHSCORES"); err != nil {
			return nil, nil, err
		}
		tmp2 := tmp.([]interface{})
		for k, v := range tmp2 {
			if k%2 != 0 {
				continue
			}
			kk := string(v.([]byte))
			value[kk] = rtti(tmp2[k+1])
			keys = append(keys, kk)
		}
		return keys, value, nil
	} else if withScores == false {
		var value []interface{}
		if tmp, err = conn.Do("ZRANGE", key, start, stop); err != nil {
			return nil, nil, err
		}
		for _, v := range tmp.([]interface{}) {
			value = append(value, rtti(v))
		}
		return value, nil, nil
	}
	return nil, nil, nil
}

//ZRevRange zrevragne
func (rao *Rao) ZRevRange(key string, start, stop int64, withScores bool) ([]interface{}, map[string]interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if withScores == true {
		var value = make(map[string]interface{})
		var keys []interface{}
		if tmp, err = conn.Do("ZREVRANGE", key, start, stop, "WITHSCORES"); err != nil {
			return nil, nil, err
		}
		tmp2 := tmp.([]interface{})
		for k, v := range tmp2 {
			if k%2 != 0 {
				continue
			}
			kk := string(v.([]byte))
			value[kk] = rtti(tmp2[k+1])
			keys = append(keys, kk)
		}
		return keys, value, nil
	} else if withScores == false {

		var value []interface{}
		if tmp, err = conn.Do("ZREVRANGE", key, start, stop); err != nil {
			return nil, nil, err
		}
		for _, v := range tmp.([]interface{}) {
			value = append(value, rtti(v))
		}
		return value, nil, nil
	}
	return nil, nil, nil
}

//ZRangeByScore zragnebyscore
func (rao *Rao) ZRangeByScore(key string, min, max int64, withScores bool) ([]interface{}, map[string]interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	if min > max {
		min, max = max, min
	}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if withScores == true {
		var value = make(map[string]interface{})
		var keys []interface{}
		if tmp, err = conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES"); err != nil {
			return nil, nil, err
		}
		tmp2 := tmp.([]interface{})
		for k, v := range tmp2 {
			if k%2 != 0 {
				continue
			}
			kk := string(v.([]byte))
			value[kk] = rtti(tmp2[k+1])
			keys = append(keys, kk)
		}
		return keys, value, nil
	} else if withScores == false {
		var value []interface{}
		if tmp, err = conn.Do("ZRANGEBYSCORE", key, min, max); err != nil {
			return nil, nil, err
		}
		for _, v := range tmp.([]interface{}) {
			value = append(value, rtti(v))
		}
		return value, nil, nil
	}
	return nil, nil, nil
}

//ZRevRangeByScore zrevragnebyscore
func (rao *Rao) ZRevRangeByScore(key string, min, max int64, withScores bool) ([]interface{}, map[string]interface{}, error) {
	var conn rds.Conn
	var err error
	var tmp interface{}
	if min > max {
		min, max = max, min
	}
	if rao.conn != nil {
		conn = *(rao.conn)
	} else if rao.pool != nil {
		conn = rao.pool.Get()
		defer conn.Close()
	}
	if withScores == true {
		var value = make(map[string]interface{})
		var keys []interface{}
		if tmp, err = conn.Do("ZREVRANGEBYSCORE", key, max, min, "WITHSCORES"); err != nil {
			return nil, nil, err
		}
		tmp2 := tmp.([]interface{})
		for k, v := range tmp2 {
			if k%2 != 0 {
				continue
			}
			kk := string(v.([]byte))
			value[kk] = rtti(tmp2[k+1])
			keys = append(keys, kk)
		}
		return keys, value, nil
	} else if withScores == false {
		var value []interface{}
		if tmp, err = conn.Do("ZREVRANGEBYSCORE", key, max, min); err != nil {
			return nil, nil, err
		}
		for _, v := range tmp.([]interface{}) {
			value = append(value, rtti(v))
		}
		return value, nil, nil
	}
	return nil, nil, nil
}

func rtti(val interface{}) interface{} {
	var value interface{}
	switch val.(type) {
	case nil:
		value = nil
	case bool:
		value = bool(val.(bool))
	case byte:
		value = byte(val.(byte))
	case int8:
		value = int8(val.(int8))
	case int16:
		value = int16(val.(int16))
	case int32:
		value = int32(val.(int32))
	case int:
		value = int(val.(int))
	case int64:
		value = int64(val.(int64))
	case float32:
		value = float32(val.(float32))
	case float64:
		value = float64(val.(float64))
	case []byte:
		value = string(val.([]byte))
	default:
		value = string(val.([]byte))
	}
	return value
}
