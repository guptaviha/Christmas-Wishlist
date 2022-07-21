package main

import (
	"fmt"
	"log"
	"net"
	"welcome-app/feed"
	"welcome-app/login"
	"welcome-app/user"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("Failed to listen to login: %v", err)
	}

	l := login.Server{}
	f := feed.Server{}
	u := user.Server{}
	Server := grpc.NewServer()

	login.RegisterAuthServiceServer(Server, &l)
	feed.RegisterFeedServiceServer(Server, &f)
	user.RegisterUserServiceServer(Server, &u)

	if err := Server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve for login: %s", err)
	}
}
