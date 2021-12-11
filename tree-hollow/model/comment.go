package model

import "time"

type CommentType int

const (
	TypeSecret  CommentType = iota // 代表一条 Secret 的评论
	TypeComment                    // 代表一条 Comment 的评论
)

type Comment struct {
	Id          int         `json:"id"`
	ParentId    int         `json:"parent_id"`
	CommentType CommentType `json:"comment_type"`
	Content     string      `json:"content"`
	SnitchName  string      `json:"snitch_name"`
	IsOpen      bool        `json:"is_open"`
	CommentTime time.Time   `json:"comment_time"`
	UpdateTime  time.Time   `json:"update_time"`
}

type CommentDetails struct {
	Comment
	ChildComment []Comment
}
