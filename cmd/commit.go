package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/xyue92/gitai/internal/ai"
	"github.com/xyue92/gitai/internal/config"
	"github.com/xyue92/gitai/internal/git"
	"github.com/xyue92/gitai/internal/ui"
)

var (
	dryRun          bool
	typeFlag        string
	scopeFlag       string
	langFlag        string
	modelFlag       string
	ticketFlag      string
	subjectLenFlag  string
	streamFlag      bool
	promptScopeFlag bool
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
	commitCmd.Flags().StringVarP(&ticketFlag, "ticket", "k", "", "Ticket/issue number (e.g., JIRA-123)")
	commitCmd.Flags().StringVarP(&subjectLenFlag, "subject-length", "n", "", "Subject length (short/normal)")
	commitCmd.Flags().BoolVarP(&streamFlag, "stream", "S", false, "Enable streaming output")
	commitCmd.Flags().BoolVarP(&promptScopeFlag, "prompt-scope", "p", false, "Prompt for scope selection")
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
	if subjectLenFlag != "" {
		cfg.SubjectLength = subjectLenFlag
	}
	if promptScopeFlag {
		cfg.PromptScope = true
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

	// Analyze file types to detect mixed commits
	files := make([]string, len(fileChanges))
	for i, fc := range fileChanges {
		files[i] = fc.File
	}
	typeHints := git.AnalyzeFileTypes(files)

	// Create selector for interactive prompts
	selector := ui.NewCommitSelector(cfg)

	// Warn if multiple commit types detected
	if len(typeHints) > 1 {
		display.ShowWarning("⚠️  Multiple commit types detected in staged files:")
		for _, hint := range typeHints {
			display.ShowInfo(fmt.Sprintf("  • %s: %d file(s)", hint.Type, len(hint.Files)))
			for _, file := range hint.Files {
				fmt.Printf("    - %s\n", file)
			}
		}
		fmt.Println()
		display.ShowInfo("💡 Best practice: Split into separate commits for better history")
		display.ShowInfo("   You can unstage files with: git restore --staged <file>")
		fmt.Println()

		// Ask user if they want to continue
		shouldContinue, err := selector.Confirm("Do you want to continue with mixed commit types?")
		if err != nil || !shouldContinue {
			return fmt.Errorf("commit cancelled - please split your changes into separate commits")
		}
		fmt.Println()
	}

	// Select commit type
	commitType := typeFlag
	if commitType == "" {
		commitType, err = selector.SelectType()
		if err != nil {
			return fmt.Errorf("type selection cancelled")
		}
	}

	// Select scope - only prompt if explicitly requested via flag or config
	scope := scopeFlag
	if scopeFlag == "" && cfg.PromptScope {
		scope, err = selector.SelectScope()
		if err != nil {
			return fmt.Errorf("scope selection cancelled")
		}
	}

	// Handle ticket number
	ticket := ticketFlag
	if ticket == "" && cfg.RequireTicket {
		// Try to extract from branch name first
		ctx, _ := git.GetProjectContext()
		autoTicket := git.ExtractTicketFromBranch(ctx.BranchName, cfg.TicketPattern)

		if autoTicket != "" {
			// Found ticket in branch name, ask for confirmation
			display.ShowInfo(fmt.Sprintf("Found ticket number in branch: %s", autoTicket))
			useAuto, err := selector.Confirm("Use this ticket number?")
			if err == nil && useAuto {
				ticket = autoTicket
			}
		}

		// If still no ticket, prompt user
		if ticket == "" {
			ticket, err = selector.PromptTicket(cfg.TicketPrefix)
			if err != nil {
				return fmt.Errorf("ticket number required but not provided")
			}
		}
	} else if ticket == "" && cfg.TicketPrefix != "" {
		// Optional ticket with prefix - try to extract from branch
		ctx, _ := git.GetProjectContext()
		ticket = git.ExtractTicketFromBranch(ctx.BranchName, cfg.TicketPattern)
	}

	// Format ticket number if needed
	if ticket != "" {
		ticket = git.FormatTicketNumber(ticket, cfg.TicketPrefix)
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
		TicketNumber:   ticket,
		SubjectLength:  cfg.SubjectLength,
	}

	prompt := promptBuilder.Build()

	// Generate commit message
	client := ai.NewOllamaClient(cfg.Model)

	var finalMessage string
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		startTime := time.Now()
		var message string
		var err error

		if streamFlag {
			// Use streaming mode
			fmt.Print("\n")
			message, err = client.GenerateStream(prompt, func(chunk string) {
				fmt.Print(chunk)
			})
			fmt.Print("\n\n")
		} else {
			// Use non-streaming mode
			message, err = client.Generate(prompt)
		}

		elapsed := time.Since(startTime)

		if err != nil {
			return fmt.Errorf("failed to generate commit message: %w", err)
		}

		// Clean up the message
		message = cleanCommitMessage(message)

		// Display time taken
		display.ShowInfo(fmt.Sprintf("[time elapsed: %.2fs]", elapsed.Seconds()))
		fmt.Println()

		// Display generated message (only in non-streaming mode, as streaming already showed it)
		if !streamFlag {
			display.ShowCommitMessage(message)
		} else {
			// In streaming mode, show the final cleaned message in a box
			display.ShowCommitMessage(message)
		}

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
			// Handle edit mode with post-edit options
			editedMessage := message
			for {
				edited, err := selector.EditMessage(editedMessage)
				if err != nil {
					return fmt.Errorf("edit cancelled")
				}
				editedMessage = edited

				// Display the edited message
				display.ShowCommitMessage(editedMessage)

				// Ask what to do after editing
				postEditAction, err := selector.ConfirmActionAfterEdit(editedMessage)
				if err != nil {
					return fmt.Errorf("action selection cancelled")
				}

				switch postEditAction {
				case ui.ActionUse:
					finalMessage = editedMessage
					goto commit
				case ui.ActionRegenerateFromEdit:
					// Use the edited message as a prompt to regenerate
					display.ShowGenerating()
					promptBuilder.CustomPrompt = fmt.Sprintf("Generate a commit message based on this user input: %s\n\nAlso consider the following diff:", editedMessage)
					prompt = promptBuilder.Build()
					// Break out of edit loop and regenerate
					goto regenerate
				case ui.ActionEdit:
					// Continue editing loop
					continue
				case ui.ActionCancel:
					return fmt.Errorf("commit cancelled by user")
				}
			}
		case ui.ActionCancel:
			return fmt.Errorf("commit cancelled by user")
		}

	regenerate:
		// Continue the outer loop to regenerate
		display.ShowGenerating()
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
