package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello: " + request
	return nil
}

func main() {
	fmt.Println("hello world!")

	for i := 0; i < 10; i++ {
		fmt.Printf("%d,", i)
	}

	rpc.RegisterName("HelloService", new(HelloService))
	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("Accept error: ", err)
	}

	rpc.ServeConn(conn)
}
