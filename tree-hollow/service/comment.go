package service

import (
	"tree-hollow/dao"
	"tree-hollow/model"
)

func AddComment(comment model.Comment) error {
	return dao.InsertComment(comment)
}
