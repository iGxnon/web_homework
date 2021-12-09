package service

import (
	"tree-hollow/dao"
	"tree-hollow/model"
)

func GetSecretBrief(id int) (model.Secret, error) {
	return dao.GetSecretBrief(id)
}

func GetSecretsFromSnitchName(name string) ([]model.Secret, error) {
	return dao.GetSecretsFromSnitchName(name)
}

func AddSecret(secret model.Secret) error {
	return dao.InsertSecret(secret)
}

func DeleteSecretFromId(id int) error {
	return dao.DeleteSecret(id)
}

func UpdateSecret(secret model.Secret) error {
	return dao.UpdateSecret(secret)
}
