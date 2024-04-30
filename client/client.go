package main

import (
	"context"
	"log"
	"time"

	pb "go_assignment4/user"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	u, err := c.AddUser(ctx, &pb.User{Name: "Miras", Email: "orazkhanm@gmail.com"})
	if err != nil {
		log.Fatalf("Failed to add the user: %v", err)
	}
	log.Printf("User is added: %v", u)

	u, err = c.GetUser(ctx, &pb.User{Id: u.Id})
	if err != nil {
		log.Fatalf("Failed to get the user: %v", err)
	}
	log.Printf("Success: %v", u)
}
