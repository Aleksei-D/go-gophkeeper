package tui

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
)

const (
	selectSecretsMessage = "Секреты"
	selectCardsMessage   = "Банковские карты"
	selectBinaryMessage  = "Файлы"
	selectNoteMessage    = "Записки"
)

type mainMenuState struct{}

func (a *mainMenuState) Action(_ context.Context, _ *Terminal) (State, error) {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{selectSecretsMessage, selectCardsMessage, selectBinaryMessage, selectNoteMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case selectSecretsMessage:
		return &secretsState{}, nil
	case selectCardsMessage:
		return &cardState{}, nil
	case selectNoteMessage:
		return &noteState{}, nil
	case selectBinaryMessage:
		return &binaryState{}, nil

	case exitButtonMessage:
		return nil, nil
	}
	return nil, nil
}

type confirmDeleteState struct {
	confirmState State
	cancelState  State
}

func (c *confirmDeleteState) Action(_ context.Context, _ *Terminal) (State, error) {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{confirmMessage, negativeMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case confirmMessage:
		return c.confirmState, nil
	default:
		return c.cancelState, nil
	}
}
