package model

import "time"

type Secret struct {
	Id         int       `json:"id"`
	CommentCnt int       `json:"comment_cnt"`
	Content    string    `json:"content"`
	SnitchName string    `json:"snitch_name"`
	IsOpen     bool      `json:"is_open"`
	PostTime   time.Time `json:"post_time"`
	UpdateTime time.Time `json:"update_time"`
}

type SecretDetails struct {
	Comments []Comment
	Secret
}
