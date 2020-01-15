package main

import (
	"context"
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

	// Setup logging File for the Server
	lf, fileError := os.OpenFile("auth.log", os.O_RDWR, os.ModePerm)
	if fileError != nil {
		log.Warn("Error in opening logging file: auth.log")
	}

	// The server File
	s = Server{
		logFile: lf,
	}

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
	mw := io.MultiWriter(os.Stdout, s.logFile)
	if fileError != nil {
		mw = io.MultiWriter(os.Stdout)
	}
	log.SetOutput(mw)

	// Setup the database
	dbClient, dbErr := getDatabaseConnection("mongodb://localhost:27017")
	if dbErr != nil {
		os.Exit(1)
	}
	s.dbClient = dbClient

	// Setting up the server
	log.Infoln("Starting Authentication Server")
	l, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		log.Fatalln("error in listening to the port 8080")
	}

	gs := grpc.NewServer()

	RegisterAuthenticationServiceServer(gs, s)

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
	_ = s.dbClient.Disconnect(context.Background())
}
