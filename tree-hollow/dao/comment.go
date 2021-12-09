package dao

import "tree-hollow/model"

func InsertComment(comment model.Comment) error {
	sqlStr := "INSERT INTO comment(parent_id, comment_type, content, snitch_name, is_open, comment_time, update_time) " + "values(?, ?, ?, ?);"
	_, err := dB.Exec(sqlStr, comment.ParentId, comment.CommentType, comment.Content, comment.SnitchName, comment.IsOpen, comment.CommentTime, comment.UpdateTime)
	return err
}

func UpdateCommentByContent(id int, content string) error {
	return nil
}

func UpdateComment(comment model.Comment) error {
	return nil
}

func DeleteComment(id int) error {
	return nil
}

// GetChildCommentsById 获取一级Comment
func GetChildCommentsById(commentId int) ([]model.Comment, error) {
	return nil, nil
}

func GetCommentDetails(commentId int) (model.CommentDetails, error) {
	return model.CommentDetails{}, nil
}

func GetCommentBrief(commentId int) (model.Comment, error) {
	return model.Comment{}, nil
}
