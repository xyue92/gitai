package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yourusername/gitai/internal/config"
)

// CommitSelector provides interactive selection for commit parameters
type CommitSelector struct {
	Config *config.Config
}

// NewCommitSelector creates a new CommitSelector
func NewCommitSelector(cfg *config.Config) *CommitSelector {
	return &CommitSelector{
		Config: cfg,
	}
}

// SelectType prompts user to select a commit type
func (cs *CommitSelector) SelectType() (string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Emoji }} {{ .Name | cyan }} - {{ .Desc }}",
		Inactive: "  {{ .Emoji }} {{ .Name }} - {{ .Desc }}",
		Selected: "{{ .Emoji }} {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select commit type",
		Items:     cs.Config.Types,
		Templates: templates,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return cs.Config.Types[idx].Name, nil
}

// SelectScope prompts user to select or enter a scope
func (cs *CommitSelector) SelectScope() (string, error) {
	// If config has predefined scopes, offer them
	if len(cs.Config.Scopes) > 0 {
		// Add "No scope" and "Custom..." options
		items := append([]string{"(No scope)", "(Custom...)"}, cs.Config.Scopes...)

		prompt := promptui.Select{
			Label: "Select scope",
			Items: items,
			Size:  10,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return "", err
		}

		if result == "(No scope)" {
			return "", nil
		}

		if result == "(Custom...)" {
			return cs.promptCustomScope()
		}

		return result, nil
	}

	// No predefined scopes, ask directly
	return cs.promptCustomScope()
}

// promptCustomScope prompts for a custom scope
func (cs *CommitSelector) promptCustomScope() (string, error) {
	prompt := promptui.Prompt{
		Label:   "Enter scope (or leave empty)",
		Default: "",
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// ConfirmAction asks user what to do with the generated message
func (cs *CommitSelector) ConfirmAction(message string) (Action, error) {
	actions := []ActionItem{
		{Name: "use", Display: "✅ Use this message", Action: ActionUse},
		{Name: "regenerate", Display: "🔄 Regenerate", Action: ActionRegenerate},
		{Name: "edit", Display: "✏️  Edit manually", Action: ActionEdit},
		{Name: "cancel", Display: "❌ Cancel", Action: ActionCancel},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Display | cyan }}",
		Inactive: "  {{ .Display }}",
		Selected: "{{ .Display | green }}",
	}

	prompt := promptui.Select{
		Label:     "What do you want to do?",
		Items:     actions,
		Templates: templates,
		Size:      4,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return ActionCancel, err
	}

	return actions[idx].Action, nil
}

// EditMessage allows user to manually edit the message
func (cs *CommitSelector) EditMessage(original string) (string, error) {
	validate := func(input string) error {
		if len(input) == 0 {
			return fmt.Errorf("commit message cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Edit commit message",
		Default:  original,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// ConfirmActionAfterEdit asks user what to do after editing the message
func (cs *CommitSelector) ConfirmActionAfterEdit(message string) (Action, error) {
	actions := []ActionItem{
		{Name: "use", Display: "✅ Use this message", Action: ActionUse},
		{Name: "regenerate", Display: "🔄 Regenerate based on this input", Action: ActionRegenerateFromEdit},
		{Name: "edit", Display: "✏️  Edit again", Action: ActionEdit},
		{Name: "cancel", Display: "❌ Cancel", Action: ActionCancel},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Display | cyan }}",
		Inactive: "  {{ .Display }}",
		Selected: "{{ .Display | green }}",
	}

	prompt := promptui.Select{
		Label:     "What do you want to do?",
		Items:     actions,
		Templates: templates,
		Size:      4,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return ActionCancel, err
	}

	return actions[idx].Action, nil
}

// Confirm asks for yes/no confirmation
func (cs *CommitSelector) Confirm(message string) (bool, error) {
	prompt := promptui.Select{
		Label: message,
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == "Yes", nil
}

// PromptTicket asks user to input ticket/issue number
func (cs *CommitSelector) PromptTicket(prefix string) (string, error) {
	label := "Enter ticket/issue number"
	if prefix != "" {
		label = fmt.Sprintf("Enter ticket number (e.g., %s-123)", prefix)
	}

	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return fmt.Errorf("ticket number cannot be empty")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// Action represents what to do with generated message
type Action int

const (
	ActionUse Action = iota
	ActionRegenerate
	ActionEdit
	ActionCancel
	ActionRegenerateFromEdit
)

// ActionItem represents a selectable action
type ActionItem struct {
	Name    string
	Display string
	Action  Action
}
