package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Reward struct {
	ID          primitive.ObjectID `bson:"_id"`
	R_name      string             `json:"r_name,omitempty"`
	R_brand     string             `json:"r_brand,omitempty"`
	R_desc      string             `json:"r_desc,omitempty"`
	R_point     string             `json:"r_point,omitempty"`
	R_img       string             `json:"r_img,omitempty"`
	R_frequency int                `json:"r_frequency,omitempty"`
}
