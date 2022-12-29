package gapi

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/pb"
	"backend_masterclass/util"
	val "backend_masterclass/val"
	"context"
	"database/sql"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	violations := ValidateLoginUserRequest(request)
	if len(violations) > 0 {
		//We're finding an internal way of showing errors with fields by using the `errdetails.BadRequest` struct
		err := invalidArgumentError(violations)
		return nil, err

	}

	user, err := s.store.GetUser(ctx, request.GetEmail())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err.Error())
		}

		return nil, status.Errorf(codes.Internal, "error while finding user %s", err.Error())
	}

	err = util.CheckPassword(request.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while checking password: %s", err.Error())
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating token: %s", err.Error())
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating refresh token: %s", err.Error())
	}

	metadata := s.extractMetadata(ctx)
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
		ClientIp:     metadata.ClientIP,
		ID:           refreshPayload.ID,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating refresh token: %s", err.Error())
	}

	response := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return response, nil
}

func ValidateLoginUserRequest(request *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if userNameErr := val.ValidateUserName(request.GetEmail()); userNameErr != nil {
		violations = append(violations, FieldViolation("email/userbame", userNameErr))
	}

	if passwordErr := val.ValidatePassword(request.GetPassword()); passwordErr != nil {
		violations = append(violations, FieldViolation("password", passwordErr))
	}

	return violations

}
