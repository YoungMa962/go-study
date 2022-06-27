package main

import (
	"context"
	"geerpc/client"
	"geerpc/server"
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
	regErr := server.Registry(&foo)
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
	server.Accept(l)
}

func main() {
	// 开启server
	addr := make(chan string)
	go startServer(addr)

	// 创建client
	clt, err := client.Dial("tcp", <-addr)
	if err != nil {
		log.Fatalf("connection error %s", err.Error())
		return
	}
	defer func() { _ = clt.Close() }()
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i + 1}
			var reply int
			// 调用 超时处理 通过context 交给用户控制
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := clt.Call(ctx, "Foo.Sum", args, &reply); err != nil {
				log.Printf("call Foo.Sum error:%s", err.Error())
				return
			}
			log.Printf("%d + %d = %d\n", args.Num1, args.Num2, reply)
			time.Sleep(time.Second)
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Minute)
}
