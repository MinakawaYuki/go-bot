package model

import "github.com/jinzhu/gorm"

type Party struct {
	gorm.Model
	Likes     int    `json:"likes" gorm:"column:likes"`
	Dislikes  int    `json:"dislikes" gorm:"column:dislikes"`
	IsDelete  int    `json:"is_delete" gorm:"column:is_delete"`
	Tags      string `json:"tags" gorm:"column:tags"`
	IsMagna   int    `json:"is_magna" gorm:"column:is_magna"`
	IsOldWang int    `json:"is_oldwang" gorm:"column:is_oldwang"`
	Property  int    `json:"property" gorm:"column:property"`

	PartyPic []PartyPic `json:"party_pic" gorm:"foreignKey:PartyId"`
}
