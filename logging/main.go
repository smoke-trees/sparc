package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var s Server

func main() {
	lf, fileError := os.OpenFile("auth.log", os.O_RDWR, os.ModePerm)
	if fileError != nil {
		log.Warn("Error in opening logging file: auth.log")
	}
	s.logFile = lf

	// setup logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             false,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	})
	mw := io.MultiWriter(os.Stdout, lf)
	if fileError != nil {
		mw = io.MultiWriter(os.Stdout)
	}
	log.SetOutput(mw)

	// Setting up the server
	log.Infoln("Starting Logging Server")

	dbc, dbcErr := grpc.Dial("0.0.0.0:8081", grpc.WithInsecure())
	if dbcErr != nil {
		log.Fatalf("Error in Connecting to Database Service")
	}
	s.databaseClient = NewSMDataServiceClient(dbc)

	l, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		log.Fatalln("error in listening to the port 8080")
	}

	gs := grpc.NewServer()
	RegisterLoggingServiceServer(gs, &s)

	// Start the server
	go func() {
		if err := gs.Serve(l); err != nil {
			log.Fatalln("error in listening to port 8080")
		}
	}()

	signalChannel := make(chan os.Signal)

	signal.Notify(signalChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-signalChannel
	log.Info("Gracefully Shutting down the server")

	// Shutdown Routine
	s.CloseFile()
}
