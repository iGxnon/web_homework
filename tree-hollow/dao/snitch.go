package dao

import "log"

func InsertSnitch(username, hashPwd string) error {
	sqlStr := "INSERT INTO snitch(name, password)  values(?, ?);"
	exec, err := dB.Exec(sqlStr, username, hashPwd)
	if err != nil {
		return err
	}
	rowsAffected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	log.Print(username+" registered, rows affect ", rowsAffected)
	return nil
}

func SelectSnitchPasswordFromName(name string) (string, error) {
	var pwd string
	sqlStr := "SELECT password FROM snitch where name = ? "
	row := dB.QueryRow(sqlStr, name)
	err := row.Scan(&pwd)
	if err != nil {
		return "", err
	}
	return pwd, nil
}
