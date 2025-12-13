package models

import (
	"encoding/json"
	"fmt"
)

// BlobTypeData тип данных бинарные данные
const BlobTypeData = "BLOB"

// Binary модель бинарныъ данных
type Binary struct {
	Name     string `json:"name"`
	Bytes    []byte `json:"bytes"`
	Comment  string `json:"-"`
	Metadata string `json:"-"`
}

// String интерфейс вывода в консоль
func (b *Binary) String() string {
	return fmt.Sprintf("Name: %s\n Note: %s\n", b.Name)
}

// ToVaultObject преобраование модели Секрета в модель объекта хранилища
func (b *Binary) ToVaultObject() VaultObject {
	var vaultObject VaultObject
	payload, _ := json.Marshal(b)
	vaultObject.Payload = payload
	vaultObject.Comment = b.Comment
	vaultObject.Metadata = b.Metadata
	vaultObject.DataType = BlobTypeData
	return vaultObject
}
