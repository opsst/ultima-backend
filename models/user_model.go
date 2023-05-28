package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id             primitive.ObjectID `json:"id,omitempty"`
	Email          string             `json:"email,omitempty"`
	Password       string             `json:"password,omitempty"`
	Profile_img    string             `json:"profile_img,omitempty"`
	Firstname      string             `json:"firstname,omitempty"`
	Lastname       string             `json:"lastname,omitempty"`
	Admin          string             `json:"admin,omitempty"`
	Used_Point_URL []interface{}      `json:"used_point_url,omitempty"`
	Point          int                `json:"point,omitempty"`
	Firebasetoken  string             `json:"firebasetoken,omitempty"`
	Fb_login       string             `json:"fb_login,omitempty"`
	Google_login   string             `json:"google_login,omitempty"`
}
