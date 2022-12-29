package gapi

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/pb"
	"backend_masterclass/util"
	val "backend_masterclass/val"
	"context"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := ValidateCreateUserRequest(request)
	if len(violations) > 0 {
		//We're finding an internal way of showing errors with fields by using the `errdetails.BadRequest` struct
		err := invalidArgumentError(violations)
		return nil, err

	}

	hashedPassword, err := util.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err.Error())
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
					return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err.Error())
				}
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err.Error())
	}

	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}

func ValidateCreateUserRequest(request *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if emailErr := val.ValidateEmail(request.GetEmail()); emailErr != nil {
		violations = append(violations, FieldViolation("email", emailErr))
	}

	if usernameErr := val.ValidateUserName(request.GetUsername()); usernameErr != nil {
		violations = append(violations, FieldViolation("username", usernameErr))
	}

	if passwordErr := val.ValidatePassword(request.GetPassword()); passwordErr != nil {
		violations = append(violations, FieldViolation("password", passwordErr))
	}

	if fullnameErr := val.ValidateFullName(request.GetFullname()); fullnameErr != nil {
		violations = append(violations, FieldViolation("full_name", fullnameErr))
	}

	return violations

}
