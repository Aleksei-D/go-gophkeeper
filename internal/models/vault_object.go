package models

import (
	"encoding/json"
	"go-gophkeeper/internal/pb/vault"
	"time"
)

// VaultObject модель объекта Хранилища
type VaultObject struct {
	Login     string
	Name      string
	Metadata  string
	Comment   string
	Payload   []byte
	DataType  string
	IsDeleted bool
	UpdateAt  time.Time
}

// ToDelete сброс данных для удаления
func (v *VaultObject) ToDelete() *VaultObject {
	v.Name = ""
	v.Metadata = ""
	v.Comment = ""
	v.Payload = nil
	v.DataType = ""
	v.UpdateAt = time.Now()
	v.IsDeleted = true
	return v
}

// ToCard преобразования модели в б. Карту
func (v *VaultObject) ToCard() (*Card, error) {
	var card *Card
	err := json.Unmarshal(v.Payload, &card)
	if err != nil {
		return nil, err
	}

	card.Comment = v.Comment
	card.Metadata = v.Metadata
	return card, err
}

// ToSecret преобразования модели в Секрет
func (v *VaultObject) ToSecret() (*Secret, error) {
	var secret Secret
	err := json.Unmarshal(v.Payload, &secret)
	if err != nil {
		return nil, err
	}

	secret.Comment = v.Comment
	secret.Metadata = v.Metadata
	return &secret, err
}

// VaultObjects список моделей VaultObject
type VaultObjects []VaultObject

// ToProto преобразование VaultObject в прото
func (v *VaultObject) ToProto() *vault.VaultMessage {
	return &vault.VaultMessage{
		Login:     v.Login,
		Name:      v.Name,
		Metadata:  v.Metadata,
		Comment:   v.Comment,
		Payload:   v.Payload,
		DataType:  v.DataType,
		IsDeleted: v.IsDeleted,
		UploadAt:  v.UpdateAt.String(),
	}
}

// ToNote преобразования модели в Заметку
func (v *VaultObject) ToNote() (*Note, error) {
	var note Note
	err := json.Unmarshal(v.Payload, &note)
	if err != nil {
		return nil, err
	}

	note.Name = v.Name
	note.Comment = v.Comment
	note.Metadata = v.Metadata
	return &note, err
}

// ToBinary преобразования модели в blob
func (v *VaultObject) ToBinary() (*Binary, error) {
	var binary Binary
	err := json.Unmarshal(v.Payload, &binary)
	if err != nil {
		return nil, err
	}

	binary.Name = v.Name
	binary.Comment = v.Comment
	binary.Metadata = v.Metadata
	return &binary, err
}

// ToProto преобразование VaultObjects в прото
func (v *VaultObjects) ToProto() *vault.VaultListMessage {
	messageProto := make([]*vault.VaultMessage, len(*v))
	for i, obj := range *v {
		messageProto[i] = obj.ToProto()
	}
	return &vault.VaultListMessage{VaultObjects: messageProto}
}

// NewVaultObjectListFromProto из Прото преобразование в VaultObjects
func NewVaultObjectListFromProto(protoMessage *vault.VaultListMessage) (VaultObjects, error) {
	var vaultObjects VaultObjects
	for _, message := range protoMessage.VaultObjects {
		vaultObject, err := NewVaultObjectFromProto(message)
		if err != nil {
			return vaultObjects, err
		}
		vaultObjects = append(vaultObjects, vaultObject)
	}
	return vaultObjects, nil
}

// NewVaultObjectFromProto из Прото преобразование в VaultObject
func NewVaultObjectFromProto(protoMessage *vault.VaultMessage) (VaultObject, error) {
	var vaultObject VaultObject
	date, err := time.Parse("2006-01-02T15:04:05-07:00", protoMessage.UploadAt)
	if err != nil {
		return vaultObject, err
	}
	vaultObject.Login = protoMessage.Login
	vaultObject.Name = protoMessage.Name
	vaultObject.Metadata = protoMessage.Metadata
	vaultObject.Comment = protoMessage.Comment
	vaultObject.Payload = protoMessage.Payload
	vaultObject.DataType = protoMessage.DataType
	vaultObject.IsDeleted = protoMessage.IsDeleted
	vaultObject.UpdateAt = date
	return vaultObject, nil
}
