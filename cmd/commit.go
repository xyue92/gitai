package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gitai/internal/ai"
	"github.com/yourusername/gitai/internal/config"
	"github.com/yourusername/gitai/internal/git"
	"github.com/yourusername/gitai/internal/ui"
)

var (
	dryRun   bool
	typeFlag string
	scopeFlag string
	langFlag  string
	modelFlag string
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate and commit with AI",
	Long:  "Analyze git diff and generate commit message using local Ollama AI",
	RunE:  runCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Show message without committing")
	commitCmd.Flags().StringVarP(&typeFlag, "type", "t", "", "Commit type (skip selection)")
	commitCmd.Flags().StringVarP(&scopeFlag, "scope", "s", "", "Commit scope (skip selection)")
	commitCmd.Flags().StringVarP(&langFlag, "language", "l", "", "Message language (en/zh)")
	commitCmd.Flags().StringVarP(&modelFlag, "model", "m", "", "Ollama model to use")
}

func runCommit(cmd *cobra.Command, args []string) error {
	display := ui.NewDisplay()

	// Show header
	if dryRun {
		display.ShowDryRun()
	} else {
		display.ShowHeader()
	}

	// Check if in git repository
	if !git.IsGitRepository() {
		return fmt.Errorf("not a git repository\nInitialize git first:\n  $ git init")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Override config with flags
	if modelFlag != "" {
		cfg.Model = modelFlag
	}
	if langFlag != "" {
		cfg.Language = langFlag
	}

	// Get staged changes
	diff, err := git.GetStagedDiff()
	if err != nil {
		return err
	}

	// Get changed files with stats
	fileChanges, err := git.GetChangedFilesWithStats()
	if err != nil {
		return fmt.Errorf("failed to get file changes: %w", err)
	}

	// Display changed files
	display.ShowChangedFiles(fileChanges)

	// Create selector for interactive prompts
	selector := ui.NewCommitSelector(cfg)

	// Select commit type
	commitType := typeFlag
	if commitType == "" {
		commitType, err = selector.SelectType()
		if err != nil {
			return fmt.Errorf("type selection cancelled")
		}
	}

	// Select scope
	scope := scopeFlag
	if scopeFlag == "" {
		scope, err = selector.SelectScope()
		if err != nil {
			return fmt.Errorf("scope selection cancelled")
		}
	}

	// Get project context
	display.ShowGenerating()
	ctx, err := git.GetProjectContext()
	if err != nil {
		// Context gathering is best-effort, continue without it
		ctx = git.ProjectContext{}
	}

	// Build prompt
	promptBuilder := &ai.PromptBuilder{
		CommitType: commitType,
		Scope:      scope,
		Diff:       diff,
		Context: ai.ProjectContext{
			ProjectName:   ctx.ProjectName,
			RecentCommits: ctx.RecentCommits,
			BranchName:    ctx.BranchName,
			ChangedFiles:  ctx.ChangedFiles,
			ReadmeSnippet: ctx.ReadmeSnippet,
			DiffStats:     ctx.DiffStats,
		},
		Language:       cfg.Language,
		DetailedCommit: cfg.DetailedCommit,
		CustomPrompt:   cfg.CustomPrompt,
	}

	prompt := promptBuilder.Build()

	// Generate commit message
	client := ai.NewOllamaClient(cfg.Model)

	var finalMessage string
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		message, err := client.Generate(prompt)
		if err != nil {
			return fmt.Errorf("failed to generate commit message: %w", err)
		}

		// Clean up the message
		message = cleanCommitMessage(message)

		// Display generated message
		display.ShowCommitMessage(message)

		// Ask user what to do
		action, err := selector.ConfirmAction(message)
		if err != nil {
			return fmt.Errorf("action selection cancelled")
		}

		switch action {
		case ui.ActionUse:
			finalMessage = message
			goto commit
		case ui.ActionRegenerate:
			display.ShowGenerating()
			continue
		case ui.ActionEdit:
			edited, err := selector.EditMessage(message)
			if err != nil {
				return fmt.Errorf("edit cancelled")
			}
			finalMessage = edited
			goto commit
		case ui.ActionCancel:
			return fmt.Errorf("commit cancelled by user")
		}
	}

	return fmt.Errorf("max retries reached")

commit:
	// Dry run mode - don't commit
	if dryRun {
		fmt.Println()
		display.ShowInfo("Would commit with message:")
		display.ShowCommitMessage(finalMessage)
		return nil
	}

	// Perform the commit
	if err := git.CommitWithMessage(finalMessage); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	// Show success
	fmt.Println()
	files := make([]string, len(fileChanges))
	for i, fc := range fileChanges {
		files[i] = fc.File
	}
	display.ShowCommitSuccess(finalMessage, files)

	return nil
}

// cleanCommitMessage removes common prefixes and cleans up the message
func cleanCommitMessage(message string) string {
	message = strings.TrimSpace(message)

	// Remove common unwanted prefixes
	unwantedPrefixes := []string{
		"commit message:",
		"here is the commit message:",
		"suggested commit message:",
		"generated message:",
	}

	lower := strings.ToLower(message)
	for _, prefix := range unwantedPrefixes {
		if strings.HasPrefix(lower, prefix) {
			message = strings.TrimSpace(message[len(prefix):])
			break
		}
	}

	// Remove surrounding quotes if present
	message = strings.Trim(message, "\"'`")

	return message
}
