package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserDAO struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Name     string
	Email    string
	Password string
}
