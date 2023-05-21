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
	Cos_istryon     bool                 `json:"cos_is-try-on,omitempty"`
	Cos_color_img   []interface{}        `json:"cos_color-img,omitempty"`
	Cos_tryon_name  []interface{}        `json:"cos_try-on-name,omitempty"`
	Cos_tryon_color []interface{}        `json:"cos_try-on-color,omitempty"`
	Ing_id          []primitive.ObjectID `json:"cos_ing_id,omitempty"`
}
