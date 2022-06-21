package main

import (
	"geerpc"
	"log"
	"net"
	"reflect"
	"sync"
	"time"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addr chan string) {
	var foo Foo
	regErr := geerpc.Registry(&foo)
	if regErr != nil {
		log.Fatalf("type %s registry fail....", reflect.TypeOf(foo).Name())
		return
	}
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	// 开启server
	addr := make(chan string)
	go startServer(addr)

	// 创建client
	client, _ := geerpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i + 1}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d\n", args.Num1, args.Num2, reply)
			time.Sleep(time.Second)
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Minute)

}
