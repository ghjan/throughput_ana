package main

import (
	"errors"
	"fmt"

	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	MYKEY = "mykey"
)

func main() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()

	err = test1(c)
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
