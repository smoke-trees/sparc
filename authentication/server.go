package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/pascaldekloe/jwt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type Server struct {
	logFile      *os.File
	logFileError error
	dbClient     *mongo.Client
}

// NewUser creates a new user
// Status codes: 0 - Success
//				 1 - Not Success
func (s Server) NewUser(c context.Context, u *NewUserRequest) (*NewUserResponse, error) {
	err := addUser(u.User, u.Password, s.dbClient, c)
	if err != nil {
		return &NewUserResponse{
			Status:  1,
			Message: fmt.Sprintf("Error in Creating New User: %v", err),
		}, err
	}
	return &NewUserResponse{
		Status:  0,
		Message: fmt.Sprintf("Success"),
	}, nil
}

func (s Server) LoginUser(c context.Context, l *LoginRequest) (*LoginResponse, error) {
	var claims jwt.Claims

	password, err := getUserPassword(l.Username, s.dbClient, c)

	if err != nil {
		return &LoginResponse{
			Status:  1,
			Message: "User Not Found",
		}, nil
	}

	if password != l.Password {
		return &LoginResponse{
			Status:  2,
			Message: "Password doesn't match",
		}, nil
	}

	user, _ := getUser(l.Username, s.dbClient, context.Background())

	claims = jwt.Claims{
		Registered: jwt.Registered{
			Issuer: "smoketrees",
		},
		Set: map[string]interface{}{
			"username":   user.Username,
			"firstname":  user.FirstName,
			"lastname":   user.LastName,
			"middlename": user.MiddleName,
		},
	}

	sign, err := claims.HMACSign(jwt.HS256, []byte("smonke"))
	if err != nil {
		log.Error("Error in Signing jwt: ", err)
		return &LoginResponse{
			Status: 5,
		}, err
	}

	return &LoginResponse{
		Status:  0,
		Message: "Success",
		Token:   string(sign),
	}, nil
}

func (s Server) VerifyUser(c context.Context, v *VerifyRequest) (*VerifyResponse, error) {
	user, err := getUser(v.Username, s.dbClient, c)
	if err != nil {
		return &VerifyResponse{
			Status:  1,
			Message: "User Not Found",
			Granted: false,
		}, err
	}
	if user.LevelOfAuth < (v.AuthLevelRequested) {
		return &VerifyResponse{
			Status:  2,
			Message: "User not has level Clearance",
			Granted: false,
		}, errors.New("does not have clearance")
	}
	return &VerifyResponse{
		Status:  0,
		Message: "Success",
		Granted: true,
	}, nil
}

func (s *Server) CloseFile() {
	if s.logFileError == nil {
		err := s.logFile.Close()
		if err != nil {
			log.Warn("Error in closing the logging File")
		}
	}
}
