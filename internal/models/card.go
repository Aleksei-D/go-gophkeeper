package models

import (
	"encoding/json"
	"fmt"
)

// CardTypeData тип данных б.карта
const CardTypeData = "CARD"

// Card модель Б. карты
type Card struct {
	Number          string `json:"number"`
	CardHolder      string `json:"cardholder"`
	ExpirationMonth int32  `json:"expiration_month"`
	ExpirationYear  int32  `json:"expiration_year"`
	CVV             int32  `json:"cvv"`
	Comment         string `json:"-"`
	Metadata        string `json:"-"`
}

// ToVaultObject преобраование модели карты в модель объекта хранилища
func (c *Card) ToVaultObject() VaultObject {
	var vaultObject VaultObject
	payload, _ := json.Marshal(c)
	vaultObject.Payload = payload
	vaultObject.Comment = c.Comment
	vaultObject.Metadata = c.Metadata
	vaultObject.DataType = CardTypeData
	return vaultObject
}

// String интерфейс вывода в консоль
func (c *Card) String() string {
	return fmt.Sprintf(
		"Card: %s\n Holder: %s\n Expiration: %d/%d\ncvv: %d\nComment: %s\n",
		c.CardHolder, c.CardHolder, c.ExpirationMonth, c.ExpirationYear, c.CVV, c.Comment)
}
