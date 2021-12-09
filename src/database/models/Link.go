package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Link struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `max:"100"`
	URL  string             `required:"true"`
}
