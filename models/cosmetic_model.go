package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cosmetic struct {
	Cos_brand       string               `json:"cos_brand,omitempty"`
	Cos_name        string               `json:"cos_name,omitempty"`
	Cos_desc        string               `json:"cos_desc,omitempty"`
	Cos_cate        string               `json:"cos_cate,omitempty"`
	Cos_form        string               `json:"cos_form,omitempty"`
	Cos_img         string               `json:"cos_img,omitempty"`
	Cos_claims      string               `json:"cos_claims,omitempty"`
	Cos_istryon     bool                 `json:"cos_is-try-on,omitempty"`
	Cos_color_img   []interface{}        `json:"cos_color-img,omitempty"`
	Cos_tryon_name  []interface{}        `json:"cos_try-on-name,omitempty"`
	Cos_tryon_color []interface{}        `json:"cos_try-on-color,omitempty"`
	Ing_id          []primitive.ObjectID `json:"cos_ing_id,omitempty"`
}
