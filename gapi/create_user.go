package gapi

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/pb"
	"backend_masterclass/util"
	"context"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	arg := db.CreateUserParams{
		HashedPassword: hashedPassword,
		Username:       request.GetUsername(),
		FullName:       request.GetFullname(),
		Email:          request.GetEmail(),
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				{
					return nil, status.Errorf(codes.AlreadyExists, "username already exists")
				}
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}
