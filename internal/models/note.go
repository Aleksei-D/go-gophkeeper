package models

import (
	"encoding/json"
	"fmt"
)

// NoteTypeData тип данных заметки
const NoteTypeData = "NOTE"

// Note модель заметки
type Note struct {
	Name     string `json:"name"`
	Note     string `json:"note"`
	Comment  string `json:"-"`
	Metadata string `json:"-"`
}

// String интерфейс вывода в консоль
func (n *Note) String() string {
	return fmt.Sprintf("Name: %s\n Note: %s\n", n.Name, n.Note)
}

// ToVaultObject преобраование модели Секрета в модель объекта хранилища
func (n *Note) ToVaultObject() VaultObject {
	var vaultObject VaultObject
	payload, _ := json.Marshal(n)
	vaultObject.Payload = payload
	vaultObject.Comment = n.Comment
	vaultObject.Metadata = n.Metadata
	vaultObject.DataType = NoteTypeData
	return vaultObject
}
