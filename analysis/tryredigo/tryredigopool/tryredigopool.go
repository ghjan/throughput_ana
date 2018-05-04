package main

import (
	"errors"
	"flag"
	"time"

	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	MYKEY = "mykey"
)

/*
配置场景
再来看下主要参数

MaxIdle
表示连接池空闲连接列表的长度限制
空闲列表是一个栈式的结构，先进后出
MaxActive
表示连接池中最大连接数限制
主要考虑到服务端支持的连接数上限，以及应用之间”瓜分”连接数
IdleTimeout
空闲连接的超时设置，一旦超时，将会从空闲列表中摘除
该超时时间时间应该小于服务端的连接超时设置

区分两种使用场景：

高频调用的场景，需要尽量压榨redis的性能：
调高MaxIdle的大小，该数目小于maxActive，由于作为一个缓冲区一样的存在，扩大缓冲区自然没有问题
调高MaxActive，考虑到服务端的支持上限，尽量调高
IdleTimeout由于是高频使用场景，设置短一点也无所谓，需要注意的一点是MaxIdle设置的长了，队列中的过期连接可能会增多，这个时候IdleTimeout也要相应变化
低频调用的场景，调用量远未达到redis的负载，稳定性为重：
MaxIdle可以设置的小一些
IdleTimeout相应地设置小一些
MaxActive随意，够用就好，容易检测到异常
 */

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

var (
	pool          *redis.Pool
	redisServer   = flag.String("redisServer", ":6379", "")
	redisPassword = flag.String("redisPassword", "", "")
)

func main() {
	flag.Parse()
	pool = newPool(*redisServer, *redisPassword)
	// 从池里获取连接
	c := pool.Get()
	// 用完后将连接放回连接池
	defer c.Close()

	err := test1(c)
	if err != nil {
		fmt.Println("test1 error", err)
	}
	err = test2(c)
	if err != nil {
		fmt.Println("test2 error", err)
	}
	err = test3(c)
	if err != nil {
		fmt.Println("test3 error", err)
	}

	err = test4(c)
	if err != nil {
		fmt.Println("test4 error", err)
	}

	err = test5(c)
	if err != nil {
		fmt.Println("test5 error", err)
	}

	err = test6(c, MYKEY)
	if err != nil {
		fmt.Println("test6 error", err)
	}

	err = test7(c)
	if err != nil {
		fmt.Println("test7 error", err)
	}
	err = testPipeline(c)
	if err != nil {
		fmt.Println("testPipeline error", err)
	}

}

//test1 读写 这里写入的值永远不会过期
func test1(c redis.Conn) (err error) {
	_, err = c.Do("SET", MYKEY, "superWang")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	username, err := redis.String(c.Do("GET", MYKEY))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}
	return err
}

//test2 设置过期呢，可以使用SET的附加参数
func test2(c redis.Conn) (err error) {
	_, err = c.Do("SET", MYKEY, "superWang", "EX", "5")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	username, err := redis.String(c.Do("GET", MYKEY))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}

	time.Sleep(8 * time.Second)

	username, err = redis.String(c.Do("GET", MYKEY))
	if err != nil {
		fmt.Println("redis get failed:", err)
		return nil
	} else {
		fmt.Printf("Get mykey: %v \n", username)
		return errors.New("redis should not existes because of time out")
	}
}

//test3 检测值是否存在  EXISTS key
func test3(c redis.Conn) (err error) {
	_, err = c.Do("SET", MYKEY, "superWang")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	is_key_exit, err := redis.Bool(c.Do("EXISTS", "mykey1"))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("exists or not: %v \n", is_key_exit)
	}
	return err
}

//test4 删除 DEL key [key …]
func test4(c redis.Conn) (err error) {
	_, err = c.Do("SET", MYKEY, "superWang")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	username, err := redis.String(c.Do("GET", MYKEY))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}

	_, err = c.Do("DEL", MYKEY)
	if err != nil {
		fmt.Println("redis delelte failed:", err)
	}

	username, err = redis.String(c.Do("GET", MYKEY))
	if err != nil {
		fmt.Println("redis get failed:", err)
		return nil
	} else {
		fmt.Printf("Get mykey: %v \n", username)
		return errors.New("redis should get failed because of deleted just now")
	}
}

//test5  读写json到redis
func test5(c redis.Conn) (err error) {
	key := "profile"
	imap := map[string]string{"username": "666", "phonenumber": "888"}
	value, _ := json.Marshal(imap)

	n, err := c.Do("SETNX", key, value)
	if err != nil {
		fmt.Println(err)
	}
	if n == int64(1) {
		fmt.Println("success")
	}

	var imapGet map[string]string

	valueGet, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		fmt.Println(err)
	}

	errShal := json.Unmarshal(valueGet, &imapGet)
	if errShal != nil {
		fmt.Println(err)
	}
	fmt.Println(imapGet["username"])
	fmt.Println(imapGet["phonenumber"])
	return err
}

//test6 设置过期时间
func test6(c redis.Conn, key string) (err error) {
	// 设置过期时间为24小时
	n, err := c.Do("EXPIRE", key, 24*3600)
	if n == int64(1) {
		fmt.Println("success")
	}
	return err
}

//test7 列表操作
func test7(c redis.Conn) (err error) {
	_, err = c.Do("lpush", "runoobkey", "redis")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	_, err = c.Do("lpush", "runoobkey", "mongodb")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	_, err = c.Do("lpush", "runoobkey", "mysql")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	values, _ := redis.Values(c.Do("lrange", "runoobkey", "0", "100"))

	for _, v := range values {
		fmt.Println(string(v.([]byte)))
	}
	return err
}

func testPipeline(c redis.Conn) (err error) {
	c.Send("SET", "foo", "bar")
	c.Send("GET", "foo")
	c.Flush()
	c.Receive()           // reply from SET
	v, err := c.Receive() // reply from GET
	if err == nil {
		fmt.Printf("testPipeline, v:%v\n", v)
	}
	return err

}
