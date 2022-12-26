package gapi

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(dbUser db.Users) *pb.User {
	return &pb.User{
		Username:          dbUser.Username,
		Fullname:          dbUser.FullName,
		Email:             dbUser.Email,
		PasswordChangedAt: timestamppb.New(dbUser.PasswordChangedAt),
		CreatedAt:         timestamppb.New(dbUser.CreatedAt),
	}
}
