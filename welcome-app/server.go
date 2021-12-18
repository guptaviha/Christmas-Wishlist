package main

import (
	"fmt"
	"log"
	"net"
	"welcome-app/feed"
	"welcome-app/login"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen to login: %v", err)
	}

	l := login.Server{}
	f := feed.Server{}
	Server := grpc.NewServer()
	// fServer := grpc.NewServer()

	login.RegisterAuthServiceServer(Server, &l)
	feed.RegisterFeedServiceServer(Server, &f)

	if err := Server.Serve(lis); err != nil {
		log.Fatalf("failed to serve for login: %s", err)
	}
}
