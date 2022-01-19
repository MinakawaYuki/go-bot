package model

import "github.com/jinzhu/gorm"

type PartyPic struct {
	gorm.Model
	PartyId  int    `json:"party_id"`
	PicUrl   string `json:"pic_url"`
	IsDelete int    `json:"is_delete" gorm:"column:is_delete"`
}
