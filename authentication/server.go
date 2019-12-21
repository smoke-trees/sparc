package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

type Server struct {
	logFile      *os.File
	logFileError error
}

func (s Server) NewUser(context.Context, *NewUserRequest) (*NewUserResponse, error) {
	panic("implement me")
}

func (s Server) LoginUser(context.Context, *LoginRequest) (*LoginResponse, error) {
	panic("implement me")
}

func (s Server) VerifyUser(context.Context, *VerifyRequest) (*VerifyResponse, error) {
	panic("implement me")
}

func (s *Server) CloseFile() {
	if s.logFileError == nil {
		err := s.logFile.Close()
		if err != nil {
			log.Warn("Error in closing the logging File")
		}
	}
}
