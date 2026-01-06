package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gitai/internal/ai"
	"github.com/yourusername/gitai/internal/config"
	"github.com/yourusername/gitai/internal/git"
	"github.com/yourusername/gitai/internal/ui"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate commit message without committing",
	Long:  "Generate a commit message for staged changes without actually committing",
	RunE:  runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&typeFlag, "type", "t", "", "Commit type (skip selection)")
	generateCmd.Flags().StringVarP(&scopeFlag, "scope", "s", "", "Commit scope (skip selection)")
	generateCmd.Flags().StringVarP(&langFlag, "language", "l", "", "Message language (en/zh)")
	generateCmd.Flags().StringVarP(&modelFlag, "model", "m", "", "Ollama model to use")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	display := ui.NewDisplay()
	display.ShowHeader()

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
	}

	prompt := promptBuilder.Build()

	// Generate commit message
	client := ai.NewOllamaClient(cfg.Model)
	message, err := client.Generate(prompt)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Clean up the message
	message = cleanCommitMessage(message)

	// Display generated message
	fmt.Println()
	display.ShowCommitMessage(message)

	display.ShowInfo("Copy this message and use it with: git commit -m \"<message>\"")

	return nil
}
