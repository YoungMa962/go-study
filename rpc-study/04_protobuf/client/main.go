package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"pbgo"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:9999", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pbgo.NewHelloStudentServiceClient(conn)

	reply, err := client.Hello(context.Background(), &pbgo.Student{Name: "YangMa", Male: false})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Student %s say Hello", reply.GetName())

	channel, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 开启线程 不断发送
	go func() {
		index := 0
		for {
			if index == 10 {
				channel.CloseSend()
				break
			}
			if err := channel.Send(&pbgo.Student{Name: "Student" + strconv.Itoa(index)}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
			index++
		}
	}()
	// 不断响应 请求
	for {
		reply, err := channel.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetName())
	}
}
