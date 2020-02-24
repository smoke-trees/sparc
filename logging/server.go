package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

type Server struct {
	logFile        *os.File
	logFileError   error
	databaseClient SMDataServiceClient
}

func (s *Server) DataLog(c context.Context, dr *DataLogRequest) (*DataLogResponse, error) {
	dlr, err := s.databaseClient.DataLog(c, dr)
	if err != nil {
		log.Error(err)
	}
	return dlr, err
}

func (s *Server) CloseFile() {
	if s.logFileError == nil {
		err := s.logFile.Close()
		if err != nil {
			log.Warn("Error in closing the logging File")
		}
	}
}
