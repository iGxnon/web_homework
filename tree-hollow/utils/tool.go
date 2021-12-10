package utils

import (
	"tree-hollow/model"
	"tree-hollow/service"
)

// BFS 对评论图进行广搜
func BFS(root model.Comment) ([]model.Comment, error) {
	ans := make([]model.Comment, 0)
	queue := NewQueue()
	queue.Offer(root)
	for !queue.IsEmpty() {
		poll := queue.Poll().(model.Comment)
		ans = append(ans, poll)
		// 树中不会出现双向箭头，所以不需要维护closeList
		comments, err := service.GetChildCommentsById(poll.Id, poll.CommentType)
		if err != nil {
			return nil, err
		}
		for _, comment := range comments {
			queue.Offer(comment)
		}
	}
	return ans, nil
}
