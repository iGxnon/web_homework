package service

import (
	"tree-hollow/dao"
	"tree-hollow/model"
)

func AddComment(comment model.Comment) error {
	return dao.InsertComment(comment)
}

func GetChildCommentsById(parentId int, commentType model.CommentType) ([]model.Comment, error) {
	return dao.GetChildCommentsById(parentId, commentType)
}

func DeleteComment(id int) error {
	return dao.DeleteComment(id)
}

func UpdateComment(comment model.Comment) error {
	return dao.UpdateComment(comment)
}

func UpdateCommentByContent(id int, content string) error {
	return dao.UpdateCommentByContent(id, content)
}

// GetAllChildComment 获取该节点以下所有评论
func GetAllChildComment(parentId int, commentType model.CommentType) ([]model.Comment, error) {
	return nil, nil
}
