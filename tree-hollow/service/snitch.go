package service

import (
	"golang.org/x/crypto/bcrypt"
	"tree-hollow/dao"
	"tree-hollow/model"
)

func CheckPassword(username, password string) (bool, error) {
	hasPwd, err := dao.SelectSnitchPasswordFromName(username)
	if err != nil {
		return false, err
	}
	// 验证加盐加密
	err = bcrypt.CompareHashAndPassword([]byte(hasPwd), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func RegisterSnitch(snitch model.Snitch) error {
	// 加盐加密
	var hashPwd, err = bcrypt.GenerateFromPassword([]byte(snitch.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = dao.InsertSnitch(snitch.Name, string(hashPwd))
	if err != nil {
		return err
	}
	return nil
}
