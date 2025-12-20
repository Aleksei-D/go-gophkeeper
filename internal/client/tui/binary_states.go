package tui

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
	"go-gophkeeper/internal/models"
	errors2 "go-gophkeeper/internal/utils/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

const (
	getBinaryMessage    = "Найти файл"
	addBinaryMessage    = "Добавить новый Файл"
	deleteBinaryMessage = "Удалить Файл"
)

type binaryState struct{}

func (c *binaryState) Action(_ context.Context, _ *Terminal) (State, error) {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{getNoteMessage, addNoteMessage, backButtonMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case getBinaryMessage:
		return &getBinaryState{}, nil
	case addBinaryMessage:
		return &addBinaryState{}, nil
	case backButtonMessage:
		return &mainMenuState{}, nil
	default:
		return nil, nil
	}
}

type getBinaryState struct{}

func (c *getBinaryState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	namePrompt := promptui.Prompt{
		Label:     "Введите название файла",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}

	binaryName, err := namePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	vaultObject, err := terminal.vaultRepo.Get(ctx, binaryName, terminal.User.Login, models.BlobTypeData)
	if err != nil {
		return nil, err
	}

	binary, err := vaultObject.ToBinary()
	if err != nil {
		return nil, err
	}
	return &manageBinaryState{Binary: binary}, nil
}

type manageBinaryState struct {
	Binary *models.Binary
}

func (m *manageBinaryState) Action(_ context.Context, _ *Terminal) (State, error) {
	fmt.Print(m.Binary)
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{deleteBinaryMessage, backButtonMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case deleteNoteMessage:
		return &confirmDeleteState{confirmState: &deleteBinaryState{binary: m.Binary}, cancelState: m}, nil
	case backButtonMessage:
		return &mainMenuState{}, nil
	default:
		return nil, nil
	}
}

type addBinaryState struct{}

func (a *addBinaryState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	prevState := &binaryState{}

	note, err := getBinaryFromInput(ctx, true, terminal)
	if err != nil {
		return nil, err
	}

	err = terminal.vaultRepo.Add(ctx, note.ToVaultObject())
	if err != nil {
		return prevState, err
	}

	fmt.Println("Секрет успешно добавлены")
	return prevState, nil
}

func newExistValidator(ctx context.Context, terminal *Terminal, dateType string) func(string) error {
	return func(input string) error {
		_, err := terminal.vaultRepo.Get(ctx, input, terminal.User.Login, dateType)
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.NotFound:
					return errors2.ErrNoContent
				default:
					return err
				}
			}
		}
		return nil
	}
}

func getBinaryFromInput(ctx context.Context, isNew bool, terminal *Terminal) (models.Binary, error) {
	var binaryName string
	var binary models.Binary
	if isNew {
		namePrompt := promptui.Prompt{
			Label:     "Введите путь имя",
			Templates: enterTemplates,
			Validate:  newExistValidator(ctx, terminal, models.BlobTypeData),
		}
		name, err := namePrompt.Run()
		if err != nil {
			return binary, fmt.Errorf("Prompt failed %v\n", err)
		}
		binaryName = name
	}

	filePathPrompt := promptui.Prompt{
		Label:     "Введите текст путь до файла",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}
	filePath, err := filePathPrompt.Run()
	if err != nil {
		return binary, fmt.Errorf("Prompt failed %v\n", err)
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return binary, err
	}

	commentPrompt := promptui.Prompt{
		Label:     "Введите Комментарий",
		Templates: enterTemplates,
	}
	comment, err := commentPrompt.Run()
	if err != nil {
		return binary, fmt.Errorf("Prompt failed %v\n", err)
	}

	binary.Name = binaryName
	binary.Bytes = buf
	binary.Comment = comment
	return binary, nil
}

type deleteBinaryState struct {
	terminal *Terminal
	binary   *models.Binary
}

func (d *deleteBinaryState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	vaultObject := d.binary.ToVaultObject()
	err := terminal.vaultRepo.Add(ctx, *vaultObject.ToDelete())
	if err != nil {
		return &noteState{}, err
	}
	return &noteState{}, nil
}
