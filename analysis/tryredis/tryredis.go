package main

import (
	"fmt"
	"log"

	"github.com/mediocregopher/radix.v2/redis"
)

const (
	REDIS_ADDRESS = "localhost:6379"
)

var redisCli redis.Client

func init() {
	redisCli, err := redis.Dial("tcp", REDIS_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	// Importantly, use defer to ensure the connection is always properly
	// closed before exiting the main() function.
	defer redisCli.Close()
}

func main() {
	hmset()
	hmget()

	hget()
}
func hget() {
	// Issue a HGET command to retrieve the title for a specific album, and use
	// the Str() helper method to convert the reply to a string.
	title, err := redisCli.Cmd("HGET", "album:1", "title").Str()
	if err != nil {
		log.Fatal(err)
	}
	// Similarly, get the artist and convert it to a string.
	artist, err := redisCli.Cmd("HGET", "album:1", "artist").Str()
	if err != nil {
		log.Fatal(err)
	}
	// And the price as a float64...
	price, err := redisCli.Cmd("HGET", "album:1", "price").Float64()
	if err != nil {
		log.Fatal(err)
	}
	// And the number of likes as an integer.
	likes, err := redisCli.Cmd("HGET", "album:1", "likes").Int()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s by %s: £%.2f [%d likes]\n", title, artist, price, likes)
}
func hmget() {
	resp := redisCli.Cmd("HMGET", "album:1", "title", "artist", "price", "likes")
	fmt.Println(resp.String())
}
func hmset() {
	// Establish a connection to the Redis server listening on port 6379 of the
	// local machine. 6379 is the default port, so unless you've already
	// changed the Redis configuration file this should work.
	// Send our command across the connection. The first parameter to Cmd()
	// is always the name of the Redis command (in this example HMSET),
	// optionally followed by any necessary arguments (in this example the
	// key, followed by the various hash fields and values).
	//HMSET key field1 value1 field2 value2
	//127.0.0.1:6379> HMSET album:1 title "Electric Ladyland" artist "Jimi Hendrix" price 4.95 likes 8
	resp := redisCli.Cmd("HMSET", "album:1", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)
	// Check the Err field of the *Resp object for any errors.
	if resp.Err != nil {
		log.Fatal(resp.Err)
	}
	fmt.Println(resp.String())
	fmt.Println("Electric Ladyland added!")
}
