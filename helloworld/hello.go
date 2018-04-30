package main

import (
	"fmt"
	"time"
)

/**
FIFO, first In first Out
 {"","",""}
 */
var message = make(chan string, 3)

func sample(m chan string) {
	m <- "hello goroutine!1"
	m <- "hello goroutine!2"
	m <- "hello goroutine!3"
	m <- "hello goroutine!4"
}

func sample2(m chan string) {
	time.Sleep(2 * time.Second)
	str := <-m
	str += " I'm goroutine!"
	m <- str
	close(m)
}

func consumer(m chan string) {
	for str := range m {
		fmt.Println(str)
	}
	//
	//for {
	//	fmt.Println(<-m)
	//}
}
func main() {

	go sample(message)
	go sample2(message)
	time.Sleep(2 * time.Second)
	go consumer(message)
	time.Sleep(1 * time.Second)

	fmt.Println("hello world!")
}
