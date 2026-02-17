package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"web_backend_project/internal/domain"
	pb "web_backend_project/web_backend_project/proto"
)

// Server представляет gRPC сервер
type Server struct {
	userUseCase domain.UserUseCase
	grpcServer  *grpc.Server
	pb.UnimplementedUserServiceServer
}

// NewGRPCServer создает новый gRPC сервер
func NewGRPCServer(userUseCase domain.UserUseCase) *Server {
	return &Server{
		userUseCase: userUseCase,
		grpcServer:  grpc.NewServer(),
	}
}

// Start запускает gRPC сервер на указанном порту
func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	pb.RegisterUserServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	log.Printf("gRPC server listening on :%s", port)
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop останавливает gRPC сервер
func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// GetUsers обрабатывает запрос на получение списка пользователей
func (s *Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := s.userUseCase.GetUsers(ctx, int(req.Page), int(req.Limit), req.Filter, req.SortBy, req.SortOrder)
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:        user.ID.Hex(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
		})
	}

	return &pb.GetUsersResponse{
		Users: pbUsers,
	}, nil
}

// GetUser обрабатывает запрос на получение информации о пользователе
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userUseCase.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        user.ID.Hex(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
		},
	}, nil
}

// CreateUser обрабатывает запрос на создание пользователя
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &domain.User{
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Username:  req.User.Username,
		Email:     req.User.Email,
	}

	id, err := s.userUseCase.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &pb.CreateUserResponse{
		Id: id.Hex(),
	}, nil
}

// UpdateUser обрабатывает запрос на обновление пользователя
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.User.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user := &domain.User{
		ID:        id,
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Username:  req.User.Username,
		Email:     req.User.Email,
	}

	if err := s.userUseCase.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &pb.UpdateUserResponse{
		Success: true,
	}, nil
}

// DeleteUser обрабатывает запрос на удаление пользователя
func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	if err := s.userUseCase.DeleteUser(ctx, id); err != nil {
		return nil, fmt.Errorf("error deleting user: %w", err)
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}
