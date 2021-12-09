package dao

import "homework7And8/model"

func InsertSecret(secret model.Secret) error {
	return nil
}

func DeleteSecret(id int) error {
	return nil
}

func UpdateSecretByContent(id int, content string) error {
	return nil
}

func UpdateSecret(secret model.Secret) error {
	return nil
}

// GetComments 获取一级评论
func GetComments(id int) ([]model.Comment, error) {
	return nil, nil
}

// GetAllComments 获取所有评论
func GetAllComments() ([]model.Comment, error) {
	return nil, nil
}

func GetSecretBrief(id int) (model.Secret, error) {
	return model.Secret{}, nil
}

func GetSecretDetails(id int) (model.SecretDetails, error) {
	return model.SecretDetails{}, nil
}
