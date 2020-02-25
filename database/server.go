package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/harshithvarma/spark/database/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
)

var collection *mongo.Collection

type Server struct{}

type DataFormat struct {
	Id             primitive.ObjectID   `bson:"_id,omitempty"`
	MeterId        string               `bson:"meter_id"`
	CustomerId     string               `bson:"customer_id"`
	LastUpdated    *timestamp.Timestamp `bson:"time_stamp"`
	EnergyConsumed float32              `bson:"energy_consumed"`
}

func (*Server) DataLog(ctx context.Context, req *proto.DataLogRequest) (*proto.DataLogResponse, error) {
	fmt.Println("Data is being logged..")
	data := req.GetData()

	dataLog := DataFormat{
		MeterId:        data.GetMeterId(),
		CustomerId:     data.GetCustomerId(),
		LastUpdated:    data.GetLastUpdated(),
		EnergyConsumed: data.GetEnergyConsumed(),
	}

	res, err := collection.InsertOne(context.Background(), dataLog)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error in logging : %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot convert to OID"))
	}

	return &proto.DataLogResponse{
		LogId:                oid.Hex(),
		Status:               1,
		Response:             "Data logged successfully!",
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}, nil
}

func (*Server) ReadData(ctx context.Context, req *proto.ReadDataRequest) (*proto.ReadDataResponse, error) {
	fmt.Println("Reading the data..")
	logid := req.GetLogId()
	oid, err := primitive.ObjectIDFromHex(logid)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("cannot parse the ID"))
	}

	data := &DataFormat{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find the data with specified ID: %v", err))
	}
	return &proto.ReadDataResponse{
		Data: &proto.SMData{
			Id:             data.Id.Hex(),
			MeterId:        data.MeterId,
			CustomerId:     data.CustomerId,
			LastUpdated:    data.LastUpdated,
			EnergyConsumed: data.EnergyConsumed,
		},
	}, nil
}

func (*Server) UpdateData(ctx context.Context, req *proto.UpdateDataRequest) (*proto.UpdateDataResponse, error) {
	fmt.Println("Updating data..")
	data := req.GetData()
	oid, err := primitive.ObjectIDFromHex(data.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"))
	}

	tempData := &DataFormat{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(tempData); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find the data with specified ID: %v", err))
	}

	tempData.MeterId = data.GetMeterId()
	tempData.CustomerId = data.GetCustomerId()
	tempData.LastUpdated = data.GetLastUpdated()
	tempData.EnergyConsumed = data.GetEnergyConsumed()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, tempData)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update the info in mongoDB"))
	}

	return &proto.UpdateDataResponse{
		Status:               1,
		Response:             "Data updated successfully",
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}, nil

}

func (*Server) DeleteData(ctx context.Context, req *proto.DeleteDataRequest) (*proto.DeleteDataResponse, error) {
	fmt.Println("Deleting..")
	logid := req.GetLogId()
	oid, err := primitive.ObjectIDFromHex(logid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot parse the ID"))
	}

	filter := bson.M{"_id": oid}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete the data from MongoDB : %v", err))
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find the data in MongoDB:%v", err))
	}

	return &proto.DeleteDataResponse{
		LogId:                logid,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}, nil
}

func (*Server) DisplayAllData(req *proto.DisplayAllDataRequest, stream proto.SMDataService_DisplayAllDataServer) error {
	fmt.Println("Displaying the data logged in the database..")

	cursor, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error : %v", err))
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		data := &DataFormat{}
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB : %v", err))
		}
		stream.Send(&proto.DisplayAllDataResponse{
			Data: &proto.SMData{
				Id:             data.Id.Hex(),
				MeterId:        data.MeterId,
				CustomerId:     data.CustomerId,
				LastUpdated:    data.LastUpdated,
				EnergyConsumed: data.EnergyConsumed,
			},
		})

		if err := cursor.Err(); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Internal error : %v", err))
		}

		return nil
	}
	return nil
}

func main() {

	// Read port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	//if we crash the code, we can point out the source of the error
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Connecting to mongoDB")

	//connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database service started")
	collection = client.Database("SMDatabase").Collection("Customer meter data")

	lis, error := net.Listen("tcp", "0.0.0.0:"+port)
	if error != nil {
		log.Fatalf("Failed to listen : %V", error)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	proto.RegisterSMDataServiceServer(s, &Server{})
	reflection.Register(s)

	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %V")
		}
	}()

	//wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//block until the signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing the MongoDB connection")
	client.Disconnect(context.TODO())
	fmt.Println("Exiting")
}
