package service

import (
	"homework7And8/dao"
	"homework7And8/model"
)

func AddComment(comment model.Comment) error {
	return dao.InsertComment(comment)
}
