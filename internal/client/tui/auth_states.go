package tui

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
	"go-gophkeeper/internal/models"
)

type authState struct{}

func (a *authState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	user, err := getUserFromInput()
	if err != nil {
		return nil, err
	}

	user, err = terminal.authClient.UserLogin(ctx, user)
	if err != nil {
		return nil, err
	}

	terminal.User = user
	return nil, nil
}

type registerState struct{}

func (a *registerState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	user, err := getUserFromInput()
	if err != nil {
		return nil, err
	}

	user, err = terminal.authClient.UserRegister(ctx, user)
	if err != nil {
		return nil, err
	}

	terminal.User = user
	return nil, nil
}

func getUserFromInput() (*models.User, error) {
	loginPrompt := promptui.Prompt{
		Label:     "Введите Login",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}

	login, err := loginPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	passwordPrompt := promptui.Prompt{
		Label:     "Введите password",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}

	password, err := passwordPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}
	return &models.User{Login: login, Password: password}, nil
}
