package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

	pb "github.com/MrCroxx/etclock/proto"
	"google.golang.org/grpc"
)

const (
	version = "1.0"
)

var (
	// 172.18.105.5:2301,172.18.105.5:2302172.18.105.5:2303
	endpoints = flag.String("endpoints", "", "etcd endpoints")
	port      = flag.String("port", ":5000", "grpc port")
)

func main() {

	flag.Parse()

	gs := grpc.NewServer()

	ls, err := NewLockerServer(LockerServerConfig{
		Endpoints: strings.Split(*endpoints, ","),
	})
	if err != nil {
		panic(err)
	}
	pb.RegisterLockerServer(gs, ls)

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Serving grpc on %s\n", *port)

	panic(gs.Serve(lis))

}
