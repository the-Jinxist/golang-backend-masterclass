package gapi

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/pb"
	"backend_masterclass/util"
	val "backend_masterclass/val"
	"backend_masterclass/worker"
	"context"
	"time"

	"github.com/hibiken/asynq"
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

	createUserTxParams := db.CreateUserTxParams{
		CreateUserParams: arg,
		AfterCreate: func(user db.Users) error {
			//We're supposed to use a DB transaction to do this two requets
			//We will be using Redis to send the verification email here
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			//Here, we're seeing how to send a golang task to a particular queue and not the default one. When you change the queue, you also have
			//to tell the task processor to look for tasks in the queue too
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.CriticalQueue),
			}

			err := s.taskDistrubutor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
			return err
		},
	}

	result, err := s.store.CreateUserTx(ctx, createUserTxParams)

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
		User: convertUser(result.User),
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
