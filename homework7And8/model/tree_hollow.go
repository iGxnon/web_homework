package model

import (
	"homework7And8/dao"
	"time"
)

type TreeHollow struct {
	Id         int `json:"id"`
	CommentCnt int
	Content    string
	UserName   string
	IsOpen     bool
	PostTime   time.Time
	UpdateTime time.Time
}

// TreeHollowEntity
// remain for nukkit as an entity
type TreeHollowEntity struct {
	Prefix             string   `json:"prefix"`
	ModelJson          string   `json:"model_json"`
	ModelSerializedIMG string   `json:"model_serialized_img"`
	Loc                Location `json:"loc"`
}

// TreeHollowEntityWithDetails
// the details of nukkit tree hollow
type TreeHollowEntityWithDetails struct {
	FormTitle string             // default as [Prefix]
	FormText  [dao.StrCnt]string // text in form
	Animation string             // animation of model entity
	Particle  string             // particle of model entity
	TreeHollowEntity
}

type Location struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Z         float64 `json:"z"`
	Yaw       int     `json:"yaw"`
	LevelName string  `json:"level_name"`
}
