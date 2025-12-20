package tui

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
	"go-gophkeeper/internal/models"
)

const (
	getSecretMessage    = "Найти Секрет"
	addSecretMessage    = "Добавить новый Секрет"
	editSecretMessage   = "Редактировать Секрет"
	deleteSecretMessage = "Удалить секрет"
)

type secretsState struct{}

func (c *secretsState) Action(_ context.Context, _ *Terminal) (State, error) {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{getSecretMessage, addSecretMessage, backButtonMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case getSecretMessage:
		return &getSecretState{}, nil
	case addSecretMessage:
		return &addSecretState{}, nil
	case backButtonMessage:
		return &mainMenuState{}, nil
	default:
		return nil, nil
	}
}

type getSecretState struct{}

func (c *getSecretState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	namePrompt := promptui.Prompt{
		Label:     "Введите название секрета",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}

	secretName, err := namePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	vaultObject, err := terminal.vaultRepo.Get(ctx, secretName, terminal.User.Login, models.SecretTypeData)
	if err != nil {
		return nil, err
	}

	secret, err := vaultObject.ToSecret()
	if err != nil {
		return nil, err
	}
	return &manageSecretState{Secret: secret}, nil
}

type manageSecretState struct {
	Secret *models.Secret
}

func (m *manageSecretState) Action(_ context.Context, _ *Terminal) (State, error) {
	fmt.Print(m.Secret)
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{editSecretMessage, deleteSecretMessage, backButtonMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case editSecretMessage:
		return &editSecretState{Secret: m.Secret}, nil
	case deleteSecretMessage:
		return &confirmDeleteState{confirmState: &deleteSecretState{secret: m.Secret}, cancelState: m}, nil
	case backButtonMessage:
		return &mainMenuState{}, nil
	default:
		return nil, nil
	}
}

type addSecretState struct{}

func (a *addSecretState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	prevState := &secretsState{}

	secret, err := getSecretFromInput(ctx, true, terminal)
	if err != nil {
		return nil, err
	}

	err = terminal.vaultRepo.Add(ctx, secret.ToVaultObject())
	if err != nil {
		return prevState, err
	}

	fmt.Println("Секрет успешно добавлены")
	return prevState, nil
}

type editSecretState struct {
	Secret *models.Secret
}

func (e *editSecretState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	prevState := &manageSecretState{Secret: e.Secret}
	secret, err := getSecretFromInput(ctx, false, terminal)
	if err != nil {
		return prevState, err
	}

	secret.Name = e.Secret.Name
	err = terminal.vaultRepo.Add(ctx, secret.ToVaultObject())
	if err != nil {
		return prevState, err
	}
	fmt.Println("Секрет успешн обновлен")
	return prevState, nil
}

func getSecretFromInput(ctx context.Context, isNew bool, terminal *Terminal) (models.Secret, error) {
	var secretName string
	var secret models.Secret
	if isNew {
		namePrompt := promptui.Prompt{
			Label:     "Введите название Секрета",
			Templates: enterTemplates,
			Validate:  newExistValidator(ctx, terminal, models.SecretTypeData),
		}
		name, err := namePrompt.Run()
		if err != nil {
			return secret, fmt.Errorf("Prompt failed %v\n", err)
		}
		secretName = name
	}

	loginPrompt := promptui.Prompt{
		Label:     "Введите логин",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}
	login, err := loginPrompt.Run()
	if err != nil {
		return secret, fmt.Errorf("Prompt failed %v\n", err)
	}

	passwordPrompt := promptui.Prompt{
		Label:     "Введите пароль",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return secret, fmt.Errorf("Prompt failed %v\n", err)
	}

	commentPrompt := promptui.Prompt{
		Label:     "Введите Комментарий",
		Templates: enterTemplates,
	}
	comment, err := commentPrompt.Run()
	if err != nil {
		return secret, fmt.Errorf("Prompt failed %v\n", err)
	}

	secret.Name = secretName
	secret.Login = login
	secret.Password = password
	secret.Comment = comment
	return secret, nil
}

type deleteSecretState struct {
	terminal *Terminal
	secret   *models.Secret
}

func (d *deleteSecretState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	vaultObject := d.secret.ToVaultObject()
	err := terminal.vaultRepo.Add(ctx, *vaultObject.ToDelete())
	if err != nil {
		return &secretsState{}, err
	}
	return &secretsState{}, nil
}
