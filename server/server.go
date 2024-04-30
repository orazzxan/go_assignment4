package main

import (
	"context"
	pb "go_assignment4/user"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	mu     sync.Mutex
	users  map[int32]*pb.User
	nextID int32
}

func (s *userServiceServer) AddUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.Id = s.nextID
	s.nextID++
	s.users[user.Id] = user
	return user, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	return user, nil
}

func (s *userServiceServer) ListUsers(req *pb.Empty, stream pb.UserService_ListUsersServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.users {
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userServiceServer{users: make(map[int32]*pb.User)})
	log.Println("Server listening at", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
