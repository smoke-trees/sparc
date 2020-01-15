package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	UserName string

	FirstName  string
	MiddleName string
	LastName   string

	level AuthLevel

	Password string
}

// getDatabaseConnection returns the connection to the database required to perform cru
func getDatabaseConnection(cs string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cs)

	conn, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Error(log.Fields{"err": err, "msg": "Error in connecting to data base", "connection string": cs})
		return nil, err
	}

	return conn, nil
}

// CRUD operations on the database

func addUser(details *UserDetails, password string, client *mongo.Client, context context.Context) error {
	usersCollection := client.Database("UserDatabase").Collection("users")

	_, err := usersCollection.InsertOne(context,
		User{
			UserName:   details.Username,
			FirstName:  details.FirstName,
			MiddleName: details.MiddleName,
			LastName:   details.LastName,
			level:      details.LevelOfAuth,
			Password:   password,
		})

	if err != nil {
		log.Error(log.Fields{"err": err, "msg": "Error in creating a new User"})
		return err
	}

	return nil
}

func getUser(username string, client *mongo.Client, c context.Context) (UserDetails, error) {
	var user User

	userCollection := client.Database("UserDatabase").Collection("users")

	query := bson.D{{
		"username", username,
	}}

	res := userCollection.FindOne(c, query)
	err := res.Decode(user)

	if err != nil {
		log.Error(log.Fields{"err": err, "msg": "Error in decoding user data"})
		return UserDetails{}, err
	}

	return UserDetails{
		Username:    user.UserName,
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName,
		LastName:    user.LastName,
		LevelOfAuth: user.level,
	}, nil

}

func getUserPassword(username string, client *mongo.Client, c context.Context) (string, error) {
	var user User

	userCollection := client.Database("UserDatabase").Collection("users")

	query := bson.D{{
		"username", username,
	}}

	res := userCollection.FindOne(c, query)
	err := res.Decode(user)

	if err != nil {
		log.Error(log.Fields{"err": err, "msg": "Error in decoding user data"})
		return "", err
	}

	return user.Password, nil
}

func editUserDetails(username string, details UserDetails, client *mongo.Client, context context.Context) {
	panic("Implement me!") //TODO
}

func editUserPassword(username string, newPassword string, client *mongo.Client, context context.Context) {
	panic("Implement me!")
}

func deleteUser(username string, client *mongo.Client, context context.Context) error {

	userCollection := client.Database("UserDatabase").Collection("users")

	query := bson.D{{
		"username", username,
	}}

	_, err := userCollection.DeleteOne(context, query)

	if err != nil {
		log.Error(log.Fields{"err": err, "msg": "Error in deleting user data"})
		return err
	}

	return nil
}
