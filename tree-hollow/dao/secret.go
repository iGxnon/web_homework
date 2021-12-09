package dao

import (
	"database/sql"
	"log"
	"tree-hollow/model"
)

func InsertSecret(secret model.Secret) error {
	sqlStr := "INSERT INTO secret(comment_cnt, content, snitch_name, is_open, post_time, update_time) values(?,?,?,?,?,?);"
	_, err := dB.Exec(sqlStr, secret.CommentCnt, secret.Content, secret.SnitchName, secret.IsOpen, secret.PostTime, secret.UpdateTime)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSecret(id int) error {
	sqlStr := "DELETE FROM secret WHERE id = ?;"
	_, err := dB.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSecret(secret model.Secret) error {
	sqlStr := "UPDATE secret SET comment_cnt=?,content=?,snitch_name=?,is_open=?,post_time=?,update_time=? WHERE id = ?;"
	_, err := dB.Exec(sqlStr, secret.CommentCnt, secret.Content, secret.SnitchName, secret.IsOpen, secret.PostTime, secret.UpdateTime, secret.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSecretCommentsCnt(id, delta int) error {
	var cnt int
	sqlStr := "SELECT comment_cnt FROM secret WHERE id = ?;"
	row := dB.QueryRow(sqlStr, id)
	err := row.Scan(&cnt)
	if err != nil {
		return err
	}
	cnt += delta

	sqlStr = "UPDATE secret SET comment_cnt = ? WHERE id = ?;"
	_, err = dB.Exec(sqlStr, cnt, id)
	if err != nil {
		return err
	}
	return nil
}

func GetSecretBrief(id int) (secretDetails model.Secret, err error) {
	sqlStr := "SELECT * FROM secret WHERE id = ? ;"
	row := dB.QueryRow(sqlStr, id)
	err = row.Scan(&secretDetails.Id, &secretDetails.CommentCnt, &secretDetails.Content, &secretDetails.SnitchName, &secretDetails.IsOpen, &secretDetails.PostTime, &secretDetails.UpdateTime)
	if err != nil {
		return model.Secret{}, err
	}
	return
}

func GetSecretDetails(id int) (secretDetails model.SecretDetails, err error) {
	sqlStr := "SELECT * FROM secret WHERE id = ? ;"
	row := dB.QueryRow(sqlStr, id)
	err = row.Scan(&secretDetails.Id, &secretDetails.CommentCnt, &secretDetails.Content, &secretDetails.SnitchName, &secretDetails.IsOpen, &secretDetails.PostTime, &secretDetails.UpdateTime)
	if err != nil {
		return model.SecretDetails{}, err
	}
	comments, err := GetChildCommentsById(id, model.TypeSecret)
	if err != nil {
		return model.SecretDetails{}, err
	}
	secretDetails.Comments = comments
	return
}

func GetSecretsFromSnitchName(name string) ([]model.Secret, error) {
	sqlStr := "SELECT * FROM secret WHERE snitch_name = ? ;"
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

	rows, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalln("未正常释放SQL")
		}
	}(rows)

	secrets := make([]model.Secret, 0)
	for rows.Next() {
		var secret model.Secret
		err := rows.Scan(&secret.Id, &secret.CommentCnt, &secret.Content, &secret.SnitchName, &secret.IsOpen, &secret.PostTime, &secret.UpdateTime)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}
	return secrets, nil
}
