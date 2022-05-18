package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// chanel 进行线程间通信
	chanel := make(chan int, 10)
	// 通过go关键字开启一个协程
	go producer(chanel)
	go consumer(chanel)
	time.Sleep(30 * time.Second)
}

func producer(chanel chan int) {
	for {
		random := makeRandomInt()
		fmt.Printf("producer[producing random int : %d]\n", random)
		chanel <- random
		time.Sleep(time.Second)
	}
}

func consumer(chanel chan int) {
	//for {
	//	value := <-chanel
	//	fmt.Printf("consumer[consuming random int : %d]\n", value)
	//	time.Sleep(2 * time.Second)
	//}
	for value := range chanel {
		fmt.Printf("consumer[consuming random int : %d]\n", value)
		time.Sleep(2 * time.Second)
	}

}
func makeRandomInt() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(100)
}
