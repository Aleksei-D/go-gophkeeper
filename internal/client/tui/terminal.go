package tui

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
	"go-gophkeeper/internal/client/api"
	"go-gophkeeper/internal/domain"
	"go-gophkeeper/internal/models"
)

const (
	enterExistUserMessage = "Вход для зарегистрированного пользователя"
	enterNewUserMessage   = "Регистрация нового пользователя"
	backButtonMessage     = "Назад"
	exitButtonMessage     = "Выйти"
	confirmMessage        = "Подтвердить"
	negativeMessage       = "Отменить"
)

// Terminal взаимодействует с юзером , окно вввода и выбора
type Terminal struct {
	User         *models.User
	currentState State
	authClient   *api.AuthClient
	vaultRepo    domain.VaultRepository
}

// NewTerminal возвращает новый  Terminal
func NewTerminal(authClient *api.AuthClient, vaultRepo domain.VaultRepository) *Terminal {
	return &Terminal{
		authClient: authClient,
		vaultRepo:  vaultRepo,
	}
}

// Login авторизация юзера в терминале
func (t *Terminal) Login(ctx context.Context) error {
	t.currentState = &welcomeState{}
	nextState, err := t.currentState.Action(ctx, t)
	if err != nil {
		return err
	}
	t.currentState = nextState
	return nil
}

// Run запуск терминала
func (t *Terminal) Run(ctx context.Context) error {
	for t.currentState != nil {
		nextState, err := t.currentState.Action(ctx, t)
		if err != nil {
			return err
		}
		t.currentState = nextState
	}
	return nil
}

// State интерфейс стейтов в терминале
type State interface {
	Action(ctx context.Context, terminal *Terminal) (State, error)
}

type welcomeState struct{}

func (w *welcomeState) Action(_ context.Context, _ *Terminal) (State, error) {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{enterExistUserMessage, enterNewUserMessage, "Выйти"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case enterExistUserMessage:
		return &authState{}, nil
	case enterNewUserMessage:
		return &registerState{}, nil
	default:
		return nil, nil
	}
}

var enterTemplates = &promptui.PromptTemplates{
	Prompt:  "{{ . }} ",
	Valid:   "{{ . | green }} ",
	Invalid: "{{ . | red }} ",
	Success: "{{ . | bold }} ",
}

func noEmptyValidate(input string) error {
	if input == "" {
		return fmt.Errorf("недопустим пустой ввод")
	}
	return nil
}
