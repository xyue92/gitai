package ai

import (
	"fmt"
	"strings"

	"github.com/xyue92/gitai/internal/i18n"
)

// ProjectContext holds information about the project for context-aware commit messages
type ProjectContext struct {
	ProjectName   string
	RecentCommits []string
	BranchName    string
	ChangedFiles  []string
	ReadmeSnippet string
	DiffStats     string
}

// PromptBuilder constructs prompts for AI commit message generation
type PromptBuilder struct {
	CommitType       string
	Scope            string
	Diff             string
	Context          ProjectContext
	Language         string
	Languages        []string // Multiple languages for multilingual commits
	DetailedCommit   bool     // If true, generate multi-line commit with body
	CustomPrompt     string   // Custom company/team commit guidelines
	TicketNumber     string   // Ticket/issue number (e.g., JIRA-123)
	SubjectLength    string   // Subject length: "short" (36 chars) or "normal" (72 chars)
	RegenerateCount  int      // Number of times regenerated (adds variation hints)
}

// Build constructs the complete prompt for Ollama
func (pb *PromptBuilder) Build() string {
	var prompt strings.Builder

	// Header
	prompt.WriteString("You are a Git commit message generator expert.\n\n")

	// Project context
	if pb.Context.ProjectName != "" || pb.Context.BranchName != "" || len(pb.Context.RecentCommits) > 0 {
		prompt.WriteString("PROJECT CONTEXT:\n")

		if pb.Context.ProjectName != "" {
			prompt.WriteString(fmt.Sprintf("- Project: %s\n", pb.Context.ProjectName))
		}

		if pb.Context.BranchName != "" {
			prompt.WriteString(fmt.Sprintf("- Branch: %s\n", pb.Context.BranchName))
		}

		if len(pb.Context.RecentCommits) > 0 {
			prompt.WriteString("- Recent commits style:\n")
			for _, commit := range pb.Context.RecentCommits {
				prompt.WriteString(fmt.Sprintf("  * %s\n", commit))
			}
		}

		if pb.Context.ReadmeSnippet != "" {
			prompt.WriteString(fmt.Sprintf("- Project description: %s\n", pb.Context.ReadmeSnippet))
		}

		prompt.WriteString("\n")
	}

	// Custom company/team guidelines (if provided)
	if pb.CustomPrompt != "" {
		prompt.WriteString("COMPANY/TEAM COMMIT GUIDELINES:\n")
		prompt.WriteString(pb.CustomPrompt)
		prompt.WriteString("\n\n")
		prompt.WriteString("IMPORTANT: Follow the above guidelines strictly when generating the commit message.\n\n")
	}

	// Task description
	prompt.WriteString("TASK:\n")
	prompt.WriteString(fmt.Sprintf("Generate a %s commit message for the following changes.\n", pb.CommitType))

	if pb.Scope != "" {
		prompt.WriteString(fmt.Sprintf("Scope: %s\n", pb.Scope))
	}

	if pb.TicketNumber != "" {
		prompt.WriteString(fmt.Sprintf("Ticket/Issue Number: %s\n", pb.TicketNumber))
		prompt.WriteString(fmt.Sprintf("IMPORTANT: Include the ticket number [%s] in the commit message.\n", pb.TicketNumber))
	}

	// Determine language(s) to use
	language := pb.Language
	if language == "" {
		language = "en"
	}

	// Check if multilingual mode
	isMultilingual := len(pb.Languages) > 1
	if isMultilingual {
		prompt.WriteString("\nLANGUAGES:\n")
		prompt.WriteString(i18n.GetMultilingualInstructions(pb.Languages))
		prompt.WriteString("\n")
	} else {
		// Single language mode
		effectiveLang := language
		if len(pb.Languages) == 1 {
			effectiveLang = pb.Languages[0]
		}

		// Get language-specific template
		langTemplate := i18n.GetLanguageTemplate(effectiveLang, pb.CommitType, pb.Scope)
		prompt.WriteString(langTemplate.LanguageInstruction)
		prompt.WriteString("\n\n")
	}

	// Changed files
	if len(pb.Context.ChangedFiles) > 0 {
		prompt.WriteString("CHANGED FILES:\n")
		for _, file := range pb.Context.ChangedFiles {
			prompt.WriteString(fmt.Sprintf("- %s\n", file))
		}
		prompt.WriteString("\n")
	}

	// Diff stats
	if pb.Context.DiffStats != "" {
		prompt.WriteString("CHANGES SUMMARY:\n")
		prompt.WriteString(pb.Context.DiffStats)
		prompt.WriteString("\n\n")
	}

	// Actual diff (truncated if too long)
	prompt.WriteString("CHANGES:\n")
	diff := pb.Diff
	// Limit diff to ~2000 characters to avoid token limits
	if len(diff) > 2000 {
		diff = diff[:2000] + "\n... (truncated)"
	}
	prompt.WriteString(diff)
	prompt.WriteString("\n\n")

	// Requirements - different based on detailed mode
	prompt.WriteString("REQUIREMENTS:\n")
	prompt.WriteString("1. Follow Conventional Commits format\n")

	// Set subject length based on configuration
	maxLength := 72
	if pb.SubjectLength == "short" {
		maxLength = 36
	}
	prompt.WriteString(fmt.Sprintf("2. Subject line: concise summary (max %d characters)\n", maxLength))

	if isMultilingual {
		// Multilingual mode requirements
		prompt.WriteString("3. Primary language for subject line: " + pb.Languages[0] + "\n")
		prompt.WriteString("4. Provide translations for all configured languages\n")
		if pb.DetailedCommit {
			prompt.WriteString("5. Body: explain WHAT changed and WHY (2-4 bullet points in primary language)\n")
			prompt.WriteString("6. Include translation section with subject line in each language\n")
			prompt.WriteString("7. Start subject line with lowercase letter after the type\n")
			prompt.WriteString("8. Separate subject and body with a blank line\n\n")
		} else {
			prompt.WriteString("5. Focus on WHAT changed and WHY (concise)\n")
			prompt.WriteString("6. Include translation section with subject line in each language\n")
			prompt.WriteString("7. Start with lowercase letter after the type\n\n")
		}
	} else {
		// Single language mode
		if pb.DetailedCommit {
			// Detailed mode: include body with explanations
			prompt.WriteString("3. Body: explain WHAT changed and WHY (2-4 bullet points)\n")
			prompt.WriteString("4. Focus on the motivation and impact, not implementation details\n")
			prompt.WriteString("5. Start subject line with lowercase letter after the type\n")
			prompt.WriteString("6. Separate subject and body with a blank line\n\n")
		} else {
			// Concise mode: subject line only
			prompt.WriteString("3. Focus on WHAT changed and WHY (concise)\n")
			prompt.WriteString("4. Start with lowercase letter after the type\n")
			prompt.WriteString("5. Generate ONLY the subject line, no body or explanation\n\n")
		}
	}

	// Output format
	prompt.WriteString("OUTPUT FORMAT:\n")

	// Get effective language for examples
	effectiveLang := language
	if len(pb.Languages) > 0 {
		effectiveLang = pb.Languages[0]
	}

	// Get language-specific template
	langTemplate := i18n.GetLanguageTemplate(effectiveLang, pb.CommitType, pb.Scope)

	// Build format string based on ticket presence
	var formatStr string
	if pb.TicketNumber != "" {
		if pb.Scope != "" {
			formatStr = fmt.Sprintf("%s(%s): [%s] <subject line>", pb.CommitType, pb.Scope, pb.TicketNumber)
		} else {
			formatStr = fmt.Sprintf("%s: [%s] <subject line>", pb.CommitType, pb.TicketNumber)
		}
	} else {
		if pb.Scope != "" {
			formatStr = fmt.Sprintf("%s(%s): <subject line>", pb.CommitType, pb.Scope)
		} else {
			formatStr = fmt.Sprintf("%s: <subject line>", pb.CommitType)
		}
	}

	if isMultilingual {
		// Multilingual format
		prompt.WriteString(formatStr + "\n\n")
		if pb.DetailedCommit {
			prompt.WriteString("<body with bullet points in primary language>\n\n")
			prompt.WriteString("Translations:\n")
			for i := 1; i < len(pb.Languages); i++ {
				langName := pb.Languages[i]
				if lang, ok := i18n.GetLanguage(pb.Languages[i]); ok {
					langName = lang.NativeName
				}
				prompt.WriteString(fmt.Sprintf("- [%s] <translated subject line>\n", langName))
			}
			prompt.WriteString("\n")
		} else {
			prompt.WriteString("\nTranslations:\n")
			for i := 1; i < len(pb.Languages); i++ {
				langName := pb.Languages[i]
				if lang, ok := i18n.GetLanguage(pb.Languages[i]); ok {
					langName = lang.NativeName
				}
				prompt.WriteString(fmt.Sprintf("- [%s] <translated subject line>\n", langName))
			}
			prompt.WriteString("\n")
		}

		prompt.WriteString("Example:\n")
		prompt.WriteString(langTemplate.ExampleSubject + "\n\n")
		if pb.DetailedCommit && len(langTemplate.ExampleBody) > 0 {
			for _, line := range langTemplate.ExampleBody {
				prompt.WriteString("- " + line + "\n")
			}
			prompt.WriteString("\n")
		}

		// Add translation examples for other languages
		if len(pb.Languages) > 1 {
			prompt.WriteString("Translations:\n")
			for i := 1; i < len(pb.Languages); i++ {
				otherTemplate := i18n.GetLanguageTemplate(pb.Languages[i], pb.CommitType, pb.Scope)
				langName := pb.Languages[i]
				if lang, ok := i18n.GetLanguage(pb.Languages[i]); ok {
					langName = lang.NativeName
				}
				prompt.WriteString(fmt.Sprintf("- [%s] %s\n", langName, otherTemplate.ExampleSubject))
			}
			prompt.WriteString("\n")
		}

		prompt.WriteString("Generate the multilingual commit message now:\n")
	} else if pb.DetailedCommit {
		// Single language, detailed format with body
		prompt.WriteString(formatStr + "\n\n<body with bullet points>\n\n")
		prompt.WriteString("Example:\n")
		prompt.WriteString(langTemplate.ExampleSubject + "\n\n")
		if len(langTemplate.ExampleBody) > 0 {
			for _, line := range langTemplate.ExampleBody {
				prompt.WriteString("- " + line + "\n")
			}
		}
		prompt.WriteString("\nGenerate the commit message now (subject + body with details):\n")
	} else {
		// Single language, concise format - subject only
		prompt.WriteString(formatStr + "\n\n")
		prompt.WriteString("Example:\n")
		prompt.WriteString(langTemplate.ExampleSubject + "\n\n")
		prompt.WriteString("Generate the commit message now (ONLY the subject line):\n")
	}

	// Add variation hint if this is a regeneration
	if pb.RegenerateCount > 0 {
		variationHints := []string{
			"Try a different perspective or emphasis in the subject line.",
			"Consider alternative wording or focus on different aspects.",
			"Rephrase with a fresh approach while maintaining accuracy.",
			"Use different verbs or structure to convey the same changes.",
			"Focus on a different aspect of the changes for variety.",
		}
		hintIndex := pb.RegenerateCount % len(variationHints)
		prompt.WriteString(fmt.Sprintf("\nNOTE: This is regeneration attempt #%d. %s\n", pb.RegenerateCount, variationHints[hintIndex]))
	}

	return prompt.String()
}

// NewPromptBuilder creates a new PromptBuilder with default values
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		Language: "en",
		Context:  ProjectContext{},
	}
}
