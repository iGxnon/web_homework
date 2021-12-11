package dao

import (
	"database/sql"
	"log"
	"time"
	"tree-hollow/model"
)

func InsertComment(comment model.Comment) error {
	sqlStr := "INSERT INTO comment(parent_id, comment_type, content, snitch_name, is_open, comment_time, update_time) values(?, ?, ?, ?, ?, ?, ?);"
	_, err := dB.Exec(sqlStr, comment.ParentId, comment.CommentType, comment.Content, comment.SnitchName, comment.IsOpen, comment.CommentTime, comment.UpdateTime)
	if err != nil {
		return err
	}
	if comment.CommentType == model.TypeSecret {
		return UpdateSecretCommentsCnt(comment.ParentId, 1)
	}
	return nil
}

func UpdateCommentByContent(id int, content string) error {
	sqlStr := "UPDATE comment SET content=?,update_time=? WHERE id=?;"
	_, err := dB.Exec(sqlStr, content, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateComment(comment model.Comment) error {
	sqlStr := "UPDATE comment SET parent_id=?,comment_type=?,content=?,snitch_name=?,is_open=?,comment_time=?,update_time=? WHERE id=?;"
	_, err := dB.Exec(sqlStr, comment.ParentId, comment.CommentType, comment.Content, comment.SnitchName, comment.IsOpen, comment.CommentTime, comment.UpdateTime, comment.Id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteComment(id int) error {
	sqlStr := "DELETE FROM comment WHERE id=?;"
	_, err := dB.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	return nil
}

// SelectChildCommentsById 获取一级Comment
func SelectChildCommentsById(parentId int, commentType model.CommentType) ([]model.Comment, error) {
	sqlStr := "SELECT * FROM comment WHERE parent_id = ? AND comment_type = ? ;"
	stmt, err := dB.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Fatalln("未正常释放SQL")
		}
	}(stmt)

	rows, err := stmt.Query(parentId, commentType)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalln("未正常释放SQL")
		}
	}(rows)

	comments := make([]model.Comment, 0)

	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(&comment.Id, &comment.ParentId, &comment.CommentType, &comment.Content, &comment.SnitchName, &comment.IsOpen, &comment.CommentTime, &comment.UpdateTime)
		if err != nil {
			return nil, err
		}
		if !comment.IsOpen {
			comment.SnitchName = "***"
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func SelectCommentDetails(commentId int) (commentDetails model.CommentDetails, err error) {
	sqlStr := "SELECT * FROM comment WHERE id=?;"
	row := dB.QueryRow(sqlStr, commentId)
	err = row.Scan(&commentDetails.Id, &commentDetails.ParentId, &commentDetails.CommentType, &commentDetails.Content, &commentDetails.SnitchName, &commentDetails.IsOpen, &commentDetails.CommentTime, &commentDetails.UpdateTime)
	if err != nil {
		return model.CommentDetails{}, err
	}
	comments, err := SelectChildCommentsById(commentDetails.Id, model.TypeComment)
	commentDetails.ChildComment = comments
	if !commentDetails.IsOpen {
		commentDetails.SnitchName = "***"
	}
	return
}

func SelectCommentBrief(commentId int) (comment model.Comment, err error) {
	sqlStr := "SELECT * FROM comment WHERE id=?;"
	row := dB.QueryRow(sqlStr, commentId)
	err = row.Scan(&comment.Id, &comment.ParentId, &comment.CommentType, &comment.Content, &comment.SnitchName, &comment.IsOpen, &comment.CommentTime, &comment.UpdateTime)
	if err != nil {
		return model.Comment{}, err
	}
	if !comment.IsOpen {
		comment.SnitchName = "***"
	}
	return
}
