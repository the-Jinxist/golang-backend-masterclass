package token

import "time"

//The idea is to declare a token maker interface so any implementation method we use, Paseto or JWT will implement
//this interface

type Maker interface {

	//CreateToken creates and signs a token for a specific username and valid duration
	CreateToken(username string, duration time.Duration) (string, error)

	//VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
