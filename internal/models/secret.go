package models

import (
	"encoding/json"
	"fmt"
)

// SecretTypeData тип данных Секрет
const SecretTypeData = "SECRET"

// Secret модель Секрета
type Secret struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Comment  string `json:"-"`
	Metadata string `json:"-"`
}

// String интерфейс вывода в консоль
func (s *Secret) String() string {
	return fmt.Sprintf("Name: %s\n Login: %s\n Password: %s\n Comment: %s\n", s.Name, s.Login, s.Password, s.Comment)
}

// ToVaultObject преобраование модели Секрета в модель объекта хранилища
func (s *Secret) ToVaultObject() VaultObject {
	var vaultObject VaultObject
	payload, _ := json.Marshal(s)
	vaultObject.Payload = payload
	vaultObject.Comment = s.Comment
	vaultObject.Metadata = s.Metadata
	vaultObject.DataType = SecretTypeData
	return vaultObject
}
