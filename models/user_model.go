package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Email     string             `json:"email,omitempty" validate:"required"`
	Password  string             `json:"password,omitempty" validate:"required"`
	Firstname string             `json:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty"`
	Admin     string             `json:"admin,omitempty"`
}
