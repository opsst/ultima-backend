package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fragrance struct {
	ID      primitive.ObjectID   `bson:"_id"`
	P_name  string               `json:"p_name,omitempty"`
	P_brand string               `json:"p_brand,omitempty"`
	P_desc  string               `json:"p_desc,omitempty"`
	P_cate  string               `json:"p_cate,omitempty"`
	P_img   string               `json:"p_img,omitempty"`
	L_link  []interface{}        `json:"l_link,omitempty"`
	Ing_id  []primitive.ObjectID `json:"ing_id,omitempty"`
}
