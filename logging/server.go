package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Server struct {
	logFile      *os.File
	logFileError error
}

func (s *Server) CloseFile() {
	if s.logFileError == nil {
		err := s.logFile.Close()
		if err != nil {
			log.Warn("Error in closing the logging File")
		}
	}
}
