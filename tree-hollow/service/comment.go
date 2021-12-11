package service

import (
	"tree-hollow/dao"
	"tree-hollow/model"
	"tree-hollow/utils"
)

func AddComment(comment model.Comment) error {
	return dao.InsertComment(comment)
}

func GetChildCommentsById(parentId int, commentType model.CommentType) ([]model.Comment, error) {
	return dao.SelectChildCommentsById(parentId, commentType)
}

func GetCommentDetails(commentId int) (commentDetails model.CommentDetails, err error) {
	return dao.SelectCommentDetails(commentId)
}

func GetCommentBrief(commentId int) (comment model.Comment, err error) {
	return dao.SelectCommentBrief(commentId)
}

func CheckCommentIdMatchName(id int, name string) (bool, error) {
	return dao.CheckCommentIdMatchName(id, name)
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

// GetAllChildComment 获取该评论以下所有评论
func GetAllChildComment(parentId int) ([]model.Comment, error) {
	root := model.Comment{
		Id: parentId,
	}
	ans := make([]model.Comment, 0)
	queue := utils.NewQueue()
	queue.Offer(root)
	for !queue.IsEmpty() {
		poll := queue.Poll().(model.Comment)
		ans = append(ans, poll)
		// 树中不会出现双向箭头，所以不需要维护closeList
		comments, err := GetChildCommentsById(poll.Id, model.TypeComment)
		if err != nil {
			return nil, err
		}
		for _, comment := range comments {
			queue.Offer(comment)
		}
	}
	// 除掉自身
	return ans[1:], nil
}
