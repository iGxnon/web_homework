package service

import (
	"tree-hollow/dao"
	"tree-hollow/model"
	"tree-hollow/utils"
)

func AddComment(comment model.Comment) error {
	if comment.CommentType == model.TypeSecret {
		return dao.UpdateSecretCommentsCnt(comment.ParentId, 1)
	}
	return dao.InsertComment(comment)
}

func GetChildCommentsById(parentId int, commentType model.CommentType) ([]model.Comment, error) {
	return dao.SelectChildCommentsById(parentId, commentType)
}

func GetCommentDetails(commentId int) (commentDetails model.CommentDetails, err error) {
	details, err := dao.SelectCommentDetails(commentId)
	if !details.IsOpen {
		details.SnitchName = "***"
	}
	return details, err
}

func GetCommentBrief(commentId int) (comment model.Comment, err error) {
	brief, err := dao.SelectCommentBrief(commentId)
	if !brief.IsOpen {
		brief.SnitchName = "***"
	}
	return brief, err
}

func CheckCommentIdMatchName(id int, name string) (bool, error) {
	brief, err := dao.SelectCommentBrief(id)
	if err != nil {
		return false, err
	}
	return brief.SnitchName == name, nil
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
