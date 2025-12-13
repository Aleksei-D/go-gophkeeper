package tui

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
	"go-gophkeeper/internal/models"
)

const (
	getNoteMessage    = "Найти Заметки"
	addNoteMessage    = "Добавить новую Заметку"
	editNoteMessage   = "Редактировать Заметку"
	deleteNoteMessage = "Удалить Заметку"
)

type noteState struct{}

func (c *noteState) Action(_ context.Context, _ *Terminal) (State, error) {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{getNoteMessage, addNoteMessage, backButtonMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case getNoteMessage:
		return &getNoteState{}, nil
	case addNoteMessage:
		return &addNoteState{}, nil
	case backButtonMessage:
		return &mainMenuState{}, nil
	default:
		return nil, nil
	}
}

type getNoteState struct{}

func (c *getNoteState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	namePrompt := promptui.Prompt{
		Label:     "Введите название Заметки",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}

	noteName, err := namePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	vaultObject, err := terminal.vaultRepo.Get(ctx, noteName, terminal.User.Login, models.NoteTypeData)
	if err != nil {
		return nil, err
	}

	note, err := vaultObject.ToNote()
	if err != nil {
		return nil, err
	}
	return &manageNoteState{Note: note}, nil
}

type manageNoteState struct {
	Note *models.Note
}

func (m *manageNoteState) Action(_ context.Context, _ *Terminal) (State, error) {
	fmt.Print(m.Note)
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{editNoteMessage, deleteNoteMessage, backButtonMessage, exitButtonMessage},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed %v\n", err)
	}

	switch result {
	case editNoteMessage:
		return &editNoteState{Note: m.Note}, nil
	case deleteNoteMessage:
		return &confirmDeleteState{confirmState: &deleteNoteState{note: m.Note}, cancelState: m}, nil
	case backButtonMessage:
		return &mainMenuState{}, nil
	default:
		return nil, nil
	}
}

type addNoteState struct{}

func (a *addNoteState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	prevState := &noteState{}

	note, err := getNoteFromInput(ctx, true, terminal)
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

type editNoteState struct {
	Note *models.Note
}

func (e *editNoteState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	prevState := &manageNoteState{Note: e.Note}
	note, err := getNoteFromInput(ctx, false, terminal)
	if err != nil {
		return prevState, err
	}

	note.Name = e.Note.Name
	err = terminal.vaultRepo.Add(ctx, note.ToVaultObject())
	if err != nil {
		return prevState, err
	}
	fmt.Println("Секрет успешн обновлен")
	return prevState, nil
}

func getNoteFromInput(ctx context.Context, isNew bool, terminal *Terminal) (models.Note, error) {
	var noteName string
	var note models.Note
	if isNew {
		namePrompt := promptui.Prompt{
			Label:     "Введите название Заметки",
			Templates: enterTemplates,
			Validate:  newExistValidator(ctx, terminal, models.NoteTypeData),
		}
		name, err := namePrompt.Run()
		if err != nil {
			return note, fmt.Errorf("Prompt failed %v\n", err)
		}
		noteName = name
	}

	notePrompt := promptui.Prompt{
		Label:     "Введите текст замеетки",
		Templates: enterTemplates,
		Validate:  noEmptyValidate,
	}
	noteContent, err := notePrompt.Run()
	if err != nil {
		return note, fmt.Errorf("Prompt failed %v\n", err)
	}

	commentPrompt := promptui.Prompt{
		Label:     "Введите Комментарий",
		Templates: enterTemplates,
	}
	comment, err := commentPrompt.Run()
	if err != nil {
		return note, fmt.Errorf("Prompt failed %v\n", err)
	}

	note.Name = noteName
	note.Note = noteContent
	note.Comment = comment
	return note, nil
}

type deleteNoteState struct {
	terminal *Terminal
	note     *models.Note
}

func (d *deleteNoteState) Action(ctx context.Context, terminal *Terminal) (State, error) {
	vaultObject := d.note.ToVaultObject()
	err := terminal.vaultRepo.Add(ctx, *vaultObject.ToDelete())
	if err != nil {
		return &noteState{}, err
	}
	return &noteState{}, nil
}
