package main

import (
	"codec"
	"encoding/json"
	"fmt"
	"geerpc"
	"log"
	"net"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}
func main() {

	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)

	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewGobCodec(conn)

	for i := 0; i < 5; i++ {
		h := codec.NewHeader("Foo.Sum", uint64(i), "")
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))

		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}
