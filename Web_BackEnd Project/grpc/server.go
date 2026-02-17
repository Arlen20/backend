package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "web_backend_project/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedQuizServiceServer
	pb.UnimplementedTransactionServiceServer
	pb.UnimplementedUserServiceServer
	pb.UnimplementedNotificationServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// Quiz Service Implementation
func (s *Server) CreateQuiz(ctx context.Context, req *pb.CreateQuizRequest) (*pb.QuizResponse, error) {
	// TODO: Implement database operations
	quiz := &pb.Quiz{
		Id:          "1", // Generate UUID
		Title:       req.Title,
		Description: req.Description,
		Questions:   req.Questions,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	return &pb.QuizResponse{Quiz: quiz}, nil
}

func (s *Server) GetQuiz(ctx context.Context, req *pb.GetQuizRequest) (*pb.QuizResponse, error) {
	// TODO: Implement database operations
	return &pb.QuizResponse{}, nil
}

func (s *Server) UpdateQuiz(ctx context.Context, req *pb.UpdateQuizRequest) (*pb.QuizResponse, error) {
	// TODO: Implement database operations
	return &pb.QuizResponse{}, nil
}

func (s *Server) DeleteQuiz(ctx context.Context, req *pb.DeleteQuizRequest) (*pb.DeleteQuizResponse, error) {
	// TODO: Implement database operations
	return &pb.DeleteQuizResponse{Success: true}, nil
}

func (s *Server) ListQuizzes(ctx context.Context, req *pb.ListQuizzesRequest) (*pb.ListQuizzesResponse, error) {
	// TODO: Implement database operations
	return &pb.ListQuizzesResponse{}, nil
}

// Transaction Service Implementation
func (s *Server) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error) {
	// TODO: Implement database operations
	transaction := &pb.Transaction{
		Id:          "1", // Generate UUID
		UserId:      req.UserId,
		Amount:      req.Amount,
		Type:        req.Type,
		Description: req.Description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	return &pb.TransactionResponse{Transaction: transaction}, nil
}

func (s *Server) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.TransactionResponse, error) {
	// TODO: Implement database operations
	return &pb.TransactionResponse{}, nil
}

func (s *Server) UpdateTransaction(ctx context.Context, req *pb.UpdateTransactionRequest) (*pb.TransactionResponse, error) {
	// TODO: Implement database operations
	return &pb.TransactionResponse{}, nil
}

func (s *Server) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.DeleteTransactionResponse, error) {
	// TODO: Implement database operations
	return &pb.DeleteTransactionResponse{Success: true}, nil
}

func (s *Server) ListTransactions(ctx context.Context, req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
	// TODO: Implement database operations
	return &pb.ListTransactionsResponse{}, nil
}

// User Service Implementation
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	// TODO: Implement database operations
	user := &pb.User{
		Id:        "1", // Generate UUID
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	return &pb.UserResponse{User: user}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	// TODO: Implement database operations
	return &pb.UserResponse{}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	// TODO: Implement database operations
	return &pb.UserResponse{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// TODO: Implement database operations
	return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// TODO: Implement database operations
	return &pb.ListUsersResponse{}, nil
}

func (s *Server) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	// TODO: Implement authentication logic
	return &pb.AuthenticateUserResponse{}, nil
}

// Notification Service Implementation
func (s *Server) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	// TODO: Implement email sending logic
	return &pb.SendEmailResponse{
		Success: true,
		Message: "Email sent successfully",
	}, nil
}

func (s *Server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	// TODO: Implement notification sending logic
	return &pb.SendNotificationResponse{
		Success:        true,
		NotificationId: "1", // Generate UUID
	}, nil
}

func (s *Server) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	// TODO: Implement database operations
	return &pb.GetNotificationsResponse{}, nil
}

func (s *Server) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.MarkNotificationAsReadResponse, error) {
	// TODO: Implement database operations
	return &pb.MarkNotificationAsReadResponse{Success: true}, nil
}

func (s *Server) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.DeleteNotificationResponse, error) {
	// TODO: Implement database operations
	return &pb.DeleteNotificationResponse{Success: true}, nil
}

func StartGRPCServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterQuizServiceServer(server, NewServer())
	pb.RegisterTransactionServiceServer(server, NewServer())
	pb.RegisterUserServiceServer(server, NewServer())
	pb.RegisterNotificationServiceServer(server, NewServer())

	log.Printf("Starting gRPC server on port %d", port)
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
