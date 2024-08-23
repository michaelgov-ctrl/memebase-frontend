package models

import "go.mongodb.org/mongo-driver/mongo"

type AuthModel struct {
	DB *mongo.Client
}
