package main

import (
	"fmt"
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

	// Read environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Print(port)

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

	l, err := net.Listen("tcp", "0.0.0.0:"+port)

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
