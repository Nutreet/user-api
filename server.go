package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nutreet/common"
	proto "github.com/nutreet/common/gen/user"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Start() {
	service := NewLogger(NewUserService())
	userServiceServer := NewUserServiceServer(service)

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down gRPC server...")
		server.GracefulStop()

		log.Println("gRPC server stopped")
	}()

	log.Println("Starting GRPC user server")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", common.Constants.USER_GRPC_PORT))
	if err != nil {
		panic(err)
	}

	log.Printf("GRPC user server up and running on port %s", common.Constants.USER_GRPC_PORT)
	proto.RegisterUserServiceServer(server, userServiceServer)
	err = server.Serve(ln)

	if err != nil {
		panic(err)
	}
}

type UserServiceServer struct {
	service UserService
	proto.UnimplementedUserServiceServer
}

func NewUserServiceServer(service UserService) *UserServiceServer {
	return &UserServiceServer{service: service}
}

func (s *UserServiceServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	validationError := ValidateRegisterRequest(req)
	if validationError != nil {
		return nil, status.Error(codes.InvalidArgument, validationError.Error())
	}

	user, err := s.service.Register(ctx, req)

	if err != nil {
		if _, ok := err.(*UserAlreadyExistsError); ok {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.RegisterResponse{
		Uid: user.UID,
	}, nil
}

func (s *UserServiceServer) GetAuthenticatedUser(ctx context.Context, req *proto.GetAuthenticatedUserRequest) (*proto.GetAuthenticatedUserResponse, error) {
	_, err := s.service.GetAutenticatedUser(ctx, req.Token)

	if err != nil {
		if _, ok := err.(*UserAlreadyExistsError); ok {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.GetAuthenticatedUserResponse{
		Uid: "uid",
	}, nil
}
