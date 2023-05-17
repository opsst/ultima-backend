package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cosmetic struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	P_name   string             `json:"p_name,omitempty" validate:"required"`
	P_brand  string             `json:"p_brand,omitempty" validate:"required"`
	P_desc   string             `json:"p_desc,omitempty"`
	P_cate   string             `json:"p_cate,omitempty"`
	P_img    string             `json:"p_img,omitempty"`
	Ing_id   []interface{}      `json:"ing_id,omitempty"`
	CreateAt time.Time          `json:"createAt,omitempty"`
}
