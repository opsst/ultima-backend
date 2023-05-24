package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id            primitive.ObjectID `json:"id,omitempty"`
	Email         string             `json:"email,omitempty"`
	Password      string             `json:"password,omitempty"`
	Firstname     string             `json:"firstname,omitempty"`
	Lastname      string             `json:"lastname,omitempty"`
	Admin         string             `json:"admin,omitempty"`
	Point         []interface{}      `json:"point,omitempty"`
	Firebasetoken string             `json:"firebasetoken,omitempty"`
	Fb_login      string             `json:"fb_login,omitempty"`
	Google_login  string             `json:"google_login,omitempty"`
}
