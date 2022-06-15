package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"pbgo"
)

func main() {
	server := grpc.NewServer()
	pbgo.RegisterHelloStudentServiceServer(server, new(pbgo.HelloStudentServiceServerImp))
	listen, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
	}
	server.Serve(listen)
}
