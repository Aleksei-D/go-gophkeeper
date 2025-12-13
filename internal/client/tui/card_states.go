package tui

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
	"go-gophkeeper/internal/models"
	errors2 "go-gophkeeper/internal/utils/errors"
	"strconv"
	"time"
)

const (
	getCardMessage    = "Найти данные карты"
	addCardMessage    = "Добавить новую карту"
	editCardMessage   = "Редактировать данные карты"
	deleteCardMessage = "Удалить данные карты"
)

type cardState struct{}

func (c *cardState) Action(_ context.Context, _ *Terminal) (State, error) {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{getCardMessage, addCardMessage, backButtonMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case getCardMessage:
		return &getCardState{}, nil
	case addCardMessage:
		return &addCardState{}, nil
	case backButtonMessage:
		return &mainMenuState{}, nil
	case exitButtonMessage:
		return nil, nil
	}
	return nil, nil
}

type getCardState struct{}

func (c *getCardState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	resourcePrompt := promptui.Prompt{
		Label:     "Введите номер карты",
		Templates: enterTemplates,
		Validate:  checkLuhnAlgorithm,
	}

	cardNumber, err := resourcePrompt.Run()
	if err != nil {
		return &cardState{}, fmt.Errorf("Prompt failed %v\n", err)
	}

	vaultObject, err := terminal.vaultRepo.Get(ctx, cardNumber, terminal.User.Login, models.CardTypeData)
	if err != nil {
		return &cardState{}, err
	}

	card, err := vaultObject.ToCard()
	if err != nil {
		return &cardState{}, err
	}
	return &manageCardState{card: card}, nil
}

type manageCardState struct {
	card *models.Card
}

func (m *manageCardState) Action(_ context.Context, _ *Terminal) (State, error) {
	fmt.Print(m.card)
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{editCardMessage, deleteCardMessage, backButtonMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case editCardMessage:
		return &editCardState{card: m.card}, nil
	case backButtonMessage:
		return &mainMenuState{}, nil
	case deleteCardMessage:
		return &confirmDeleteState{confirmState: &deleteCardState{card: m.card}, cancelState: m}, nil
	default:
		return nil, nil
	}
}

type addCardState struct{}

func (a *addCardState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	prevState := &cardState{}
	card, err := getCardFromInput(ctx, true, terminal)
	if err != nil {
		return prevState, fmt.Errorf("Prompt failed %v\n", err)
	}

	err = terminal.vaultRepo.Add(ctx, card.ToVaultObject())
	if err != nil {
		return prevState, err
	}
	fmt.Println("Секрет успешно добавлены")
	return prevState, nil
}

type editCardState struct {
	card *models.Card
}

func (e *editCardState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	prevState := &manageCardState{card: e.card}
	card, err := getCardFromInput(ctx, true, terminal)
	if err != nil {
		return prevState, fmt.Errorf("Prompt failed %v\n", err)
	}

	card.Number = e.card.Number
	err = terminal.vaultRepo.Add(ctx, card.ToVaultObject())
	if err != nil {
		return prevState, err
	}
	fmt.Println("Секрет успешн обновлен")
	return prevState, nil
}

func checkLuhnAlgorithm(cardNumber string) error {
	sum := 0
	double := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(cardNumber[i]))
		if err != nil {
			return errors2.ErrCardNumber
		}

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	if sum%10 != 0 {
		return errors2.ErrCardNumber
	}
	return nil
}

func isExpirationMonthValid(expirationMonth string) error {
	digit, err := strconv.Atoi(expirationMonth)
	if err != nil {
		return errors2.ErrExpirationMonth
	}
	if digit < 1 || digit > 12 {
		return errors2.ErrExpirationMonth
	}
	return nil
}

func isExpirationYearValid(expirationYear string) error {
	digit, err := strconv.Atoi(expirationYear)
	if err != nil {
		return errors2.ErrExpirationYear
	}
	now := time.Now()
	if digit < 10 || digit > now.Year() {
		return errors2.ErrExpirationYear
	}
	return nil
}

func isCVVValid(cvv string) error {
	if len(cvv) != 3 {
		return errors2.ErrCVV
	}
	_, err := strconv.Atoi(cvv)
	if err != nil {
		return errors2.ErrCVV
	}
	return nil
}

func getCardFromInput(ctx context.Context, isNew bool, terminal *Terminal) (models.Card, error) {
	var cardNumber string
	var card models.Card
	if isNew {
		cardNumberPrompt := promptui.Prompt{
			Label:     "Введите номер карты",
			Templates: enterTemplates,
			Validate:  newExistValidator(ctx, terminal, models.CardTypeData),
		}
		cardNumberInput, err := cardNumberPrompt.Run()
		if err != nil {
			return card, fmt.Errorf("Prompt failed %v\n", err)
		}
		cardNumber = cardNumberInput
	}

	cardHolderPrompt := promptui.Prompt{
		Label:     "Введите имя владельца карты",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}
	cardHolder, err := cardHolderPrompt.Run()
	if err != nil {
		return card, fmt.Errorf("Prompt failed %v\n", err)
	}

	expirationMonthPrompt := promptui.Prompt{
		Label:     "Введите пароль",
		Templates: enterTemplates,
		Validate:  isExpirationMonthValid,
	}
	expirationMonth, err := expirationMonthPrompt.Run()
	if err != nil {
		return card, fmt.Errorf("Prompt failed %v\n", err)
	}

	expirationYearPrompt := promptui.Prompt{
		Label:     "Введите пароль",
		Templates: enterTemplates,
		Validate:  isExpirationYearValid,
	}
	expirationYear, err := expirationYearPrompt.Run()
	if err != nil {
		return card, fmt.Errorf("Prompt failed %v\n", err)
	}

	CVVPrompt := promptui.Prompt{
		Label:     "Введите пароль",
		Templates: enterTemplates,
		Validate:  isCVVValid,
	}
	CVV, err := CVVPrompt.Run()
	if err != nil {
		return card, fmt.Errorf("Prompt failed %v\n", err)
	}

	commentPrompt := promptui.Prompt{
		Label:     "Введите Комментарий",
		Templates: enterTemplates,
	}
	comment, err := commentPrompt.Run()
	if err != nil {
		return card, fmt.Errorf("Prompt failed %v\n", err)
	}

	expirationMonthDigit, err := strconv.Atoi(expirationMonth)
	if err != nil {
		return card, fmt.Errorf("Prompt failed %v\n", err)
	}

	expirationYearDigit, err := strconv.Atoi(expirationYear)
	if err != nil {
		return card, fmt.Errorf("Prompt failed %v\n", err)
	}

	CVVDigit, err := strconv.Atoi(CVV)
	if err != nil {
		return card, fmt.Errorf("Prompt failed %v\n", err)
	}

	card.Number = cardNumber
	card.CardHolder = cardHolder
	card.ExpirationMonth = int32(expirationMonthDigit)
	card.ExpirationYear = int32(expirationYearDigit)
	card.CVV = int32(CVVDigit)
	card.Comment = comment

	return card, nil
}

type deleteCardState struct {
	terminal *Terminal
	card     *models.Card
}

func (d *deleteCardState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	vaultObject := d.card.ToVaultObject()
	err := terminal.vaultRepo.Add(ctx, *vaultObject.ToDelete())
	if err != nil {
		return &cardState{}, err
	}
	return &cardState{}, nil
}
