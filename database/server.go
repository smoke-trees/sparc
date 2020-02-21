package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/harshithvarma/spark/database/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*Server) DataLog(req *proto.DataLogRequest) (*proto.DataLogResponse, error) {
	fmt.Print("Data is being logged..")
	data := req.GetData()

	log := DataFormat{
		MeterId:        data.GetMeterId(),
		CustomerId:     data.GetCustomerId(),
		LastUpdated:    data.GetLastUpdated(),
		EnergyConsumed: data.GetEnergyConsumed(),
	}

	res, err := collection.InsertOne(context.Background(), log)
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

func (*Server) ReadData(req *proto.ReadDataRequest) *proto.ReadDataResponse {
	fmt.Print("Reading the data..")
	logid := req.GetLogId()
	oid, err := primitive.ObjectIDFromHex(logid)
	if err != nil {
		return nil
	}

	data := &DataFormat{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil
	}
	return &proto.ReadDataResponse{
		Data:                 &proto.SMData{
			Id:                   data.Id.Hex(),
			MeterId:              data.MeterId,
			CustomerId:           data.CustomerId,
			LastUpdated:          data.LastUpdated,
			EnergyConsumed:       data.EnergyConsumed,
		},
	}
}


