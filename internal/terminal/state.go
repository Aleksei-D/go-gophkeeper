package client

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type State interface {
	Action()
}

type AuthState struct {
	nextAction chan<- State
}

func (a *AuthState) Action() {
	credValidate := func(input string) error {
		if input == "" {
			return fmt.Errorf("login not must be an empty")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	loginPrompt := promptui.Prompt{
		Label:     "Enter Login",
		Templates: templates,
		Validate:  credValidate,
	}

	login, err := loginPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You answered %s\n", login)

	passwordPrompt := promptui.Prompt{
		Label:     "Enter password",
		Templates: templates,
		Validate:  credValidate,
	}

	password, err := passwordPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You answered %s\n", password)

	// user := models.User{Login: login, Password: password}
	ac := AuthState{}
	a.nextAction <- &ac
}
