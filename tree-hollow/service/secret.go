package service

import (
	"time"
	"tree-hollow/dao"
	"tree-hollow/model"
)

func GetSecretBrief(id int) (model.Secret, error) {
	brief, err := dao.SelectSecretBrief(id)
	if !brief.IsOpen {
		brief.SnitchName = "***"
	}
	return brief, err
}

func GetSecretDetails(id int) (model.SecretDetails, error) {
	details, err := dao.SelectSecretDetails(id)
	if !details.IsOpen {
		details.SnitchName = "***"
	}
	return details, err
}

func GetSecretsFromSnitchName(name string) ([]model.Secret, error) {
	return dao.SelectSecretsFromSnitchName(name)
}

func GetCommentCnt(id int) (int, error) {
	brief, err := GetSecretBrief(id)
	if err != nil {
		return 0, err
	}
	return brief.CommentCnt, nil
}

func GetPostTime(id int) (time.Time, error) {
	brief, err := GetSecretBrief(id)
	if err != nil {
		return time.Time{}, err
	}
	return brief.PostTime, nil
}

func AddSecret(secret model.Secret) error {
	return dao.InsertSecret(secret)
}

func CheckSecretIdMatchName(id int, name string) (bool, error) {
	brief, err := dao.SelectSecretBrief(id)
	if err != nil {
		return false, err
	}
	return brief.SnitchName == name, nil
}

func DeleteSecretFromId(id int) error {
	return dao.DeleteSecret(id)
}

func UpdateSecret(secret model.Secret) error {
	return dao.UpdateSecret(secret)
}
