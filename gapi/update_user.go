package gapi

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/pb"
	"backend_masterclass/util"
	val "backend_masterclass/val"
	"context"
	"database/sql"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	violations := ValidateUpdateUserRequest(request)
	if len(violations) > 0 {
		//We're finding an internal way of showing errors with fields by using the `errdetails.BadRequest` struct
		err := invalidArgumentError(violations)
		return nil, err

	}

	arg := db.UpdateUserParams{
		Username: request.GetUsername(),
		FullName: sql.NullString{
			String: request.GetFullname(),
			Valid:  len(request.GetFullname()) > 0,
		},
		Email: sql.NullString{
			String: request.GetEmail(),
			Valid:  len(request.GetEmail()) > 0,
		},
	}

	if len(request.GetPassword()) > 0 {
		hashedPassword, err := util.HashPassword(request.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err.Error())
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := s.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "cannot find user with the provided username: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err.Error())
	}

	response := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}

func ValidateUpdateUserRequest(request *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if len(request.GetEmail()) > 0 {
		if emailErr := val.ValidateEmail(request.GetEmail()); emailErr != nil {
			violations = append(violations, FieldViolation("email", emailErr))
		}
	}

	if usernameErr := val.ValidateUserName(request.GetUsername()); usernameErr != nil {
		violations = append(violations, FieldViolation("username", usernameErr))
	}

	if len(request.GetFullname()) > 0 {
		if fullnameErr := val.ValidateFullName(request.GetFullname()); fullnameErr != nil {
			violations = append(violations, FieldViolation("fullname", fullnameErr))
		}
	}

	if len(request.GetPassword()) > 0 {
		if passwordErr := val.ValidatePassword(request.GetPassword()); passwordErr != nil {
			violations = append(violations, FieldViolation("password", passwordErr))
		}
	}

	return violations

}
