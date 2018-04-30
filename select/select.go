package main

import (
	"time"
	"fmt"
)

func sample(m chan string) {
	for i := 0; i < 19; i++ {
		m <- fmt.Sprintf("I'm sample1 num:%d", i)
		time.Sleep(3 * time.Second)
	}
}

func sample2(m chan int) {
	for i := 0; i < 10; i++ {
		m <- i
		time.Sleep(60 * time.Second)
	}
}

func consumer(ch1 chan string, ch2 chan int) {
	for {
		select {
		case str, ch1Check := <-ch1:
			if !ch1Check {
				fmt.Println("ch1 failed!")
			} else {
				fmt.Println(str)
			}
		case p, ch2Check := <-ch2:
			if !ch2Check {
				fmt.Println("ch1 failed!")
			} else {
				fmt.Println(p)
			}
		}
	}
}
func main() {
	ch1 := make(chan string, 3)
	ch2 := make(chan int, 5)
	for i := 0; i < 10; i++ {
		go sample(ch1)
		go sample2(ch2)
	}

	go consumer(ch1, ch2)
	time.Sleep(5 * time.Minute)
}
