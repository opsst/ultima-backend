package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cosmetic struct {
	Id              primitive.ObjectID   `bson:"_id"`
	Cos_brand       string               `json:"cos_brand,omitempty"`
	Cos_name        string               `json:"cos_name,omitempty"`
	Cos_desc        string               `json:"cos_desc,omitempty"`
	Cos_cate        string               `json:"cos_cate,omitempty"`
	Cos_img         []interface{}        `json:"cos_img,omitempty"`
	Cos_istryon     bool                 `json:"cos_istryon,omitempty"`
	Cos_color_img   []interface{}        `json:"cos_color_img,omitempty"`
	Cos_tryon_name  []interface{}        `json:"cos_tryon_name,omitempty"`
	Cos_tryon_color []interface{}        `json:"cos_tryon_color,omitempty"`
	L_link          []interface{}        `json:"l_link,omitempty"`
	Ing_id          []primitive.ObjectID `json:"ing_id,omitempty"`
}
