package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetStagedDiff returns the diff of staged changes
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	diff := string(output)
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no staged changes found\nStage your changes first:\n  $ git add <files>")
	}

	return diff, nil
}

// GetChangedFiles returns list of files with staged changes
func GetChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}

	return files, nil
}

// GetDiffStats returns statistics about the changes
func GetDiffStats() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--stat")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff stats: %w", err)
	}

	return string(output), nil
}

// GetChangedFilesWithStats returns files with their change statistics
func GetChangedFilesWithStats() ([]FileChange, error) {
	cmd := exec.Command("git", "diff", "--cached", "--numstat")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	var changes []FileChange
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		changes = append(changes, FileChange{
			File:      parts[2],
			Additions: parts[0],
			Deletions: parts[1],
		})
	}

	return changes, nil
}

// FileChange represents statistics for a single file
type FileChange struct {
	File      string
	Additions string
	Deletions string
}

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// CommitTypeHint represents a suggested commit type based on file analysis
type CommitTypeHint struct {
	Type  string
	Files []string
}

// AnalyzeFileTypes analyzes changed files and suggests commit types
// Returns a list of detected commit types with their associated files
func AnalyzeFileTypes(files []string) []CommitTypeHint {
	typeMap := make(map[string][]string)

	for _, file := range files {
		commitType := detectFileType(file)
		typeMap[commitType] = append(typeMap[commitType], file)
	}

	// Convert map to slice and sort by priority
	var hints []CommitTypeHint
	priorityOrder := []string{"feat", "fix", "docs", "test", "style", "refactor", "perf", "build", "ci", "chore"}

	for _, typeStr := range priorityOrder {
		if files, exists := typeMap[typeStr]; exists {
			hints = append(hints, CommitTypeHint{
				Type:  typeStr,
				Files: files,
			})
		}
	}

	return hints
}

// detectFileType determines the likely commit type based on file path
func detectFileType(file string) string {
	file = strings.ToLower(file)

	// Documentation files
	if strings.HasSuffix(file, ".md") ||
	   strings.HasSuffix(file, ".mdx") ||
	   strings.HasSuffix(file, ".rst") ||
	   strings.Contains(file, "readme") ||
	   strings.Contains(file, "docs/") ||
	   strings.Contains(file, "documentation/") ||
	   strings.HasSuffix(file, ".txt") && strings.Contains(file, "doc") {
		return "docs"
	}

	// Test files
	if strings.Contains(file, "_test.") ||
	   strings.Contains(file, ".test.") ||
	   strings.Contains(file, "/test/") ||
	   strings.Contains(file, "/tests/") ||
	   strings.Contains(file, "__tests__/") ||
	   strings.Contains(file, ".spec.") {
		return "test"
	}

	// CI/CD files
	if strings.Contains(file, ".github/workflows/") ||
	   strings.Contains(file, ".gitlab-ci") ||
	   strings.Contains(file, "jenkinsfile") ||
	   strings.Contains(file, ".circleci/") ||
	   strings.Contains(file, ".travis.yml") {
		return "ci"
	}

	// Build/config files
	if strings.HasSuffix(file, "package.json") ||
	   strings.HasSuffix(file, "package-lock.json") ||
	   strings.HasSuffix(file, "go.mod") ||
	   strings.HasSuffix(file, "go.sum") ||
	   strings.HasSuffix(file, "cargo.toml") ||
	   strings.HasSuffix(file, "cargo.lock") ||
	   strings.HasSuffix(file, "pom.xml") ||
	   strings.HasSuffix(file, "build.gradle") ||
	   strings.HasSuffix(file, "dockerfile") ||
	   strings.HasSuffix(file, "makefile") ||
	   strings.HasSuffix(file, ".yaml") && strings.Contains(file, "config") ||
	   strings.HasSuffix(file, ".yml") && strings.Contains(file, "config") ||
	   strings.HasSuffix(file, ".toml") ||
	   strings.HasSuffix(file, ".json") && strings.Contains(file, "config") {
		return "build"
	}

	// Style files
	if strings.HasSuffix(file, ".css") ||
	   strings.HasSuffix(file, ".scss") ||
	   strings.HasSuffix(file, ".sass") ||
	   strings.HasSuffix(file, ".less") {
		return "style"
	}

	// Default to feat for code files, chore for others
	codeExtensions := []string{".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".c", ".cpp", ".rs", ".rb", ".php", ".swift", ".kt"}
	for _, ext := range codeExtensions {
		if strings.HasSuffix(file, ext) {
			return "feat"
		}
	}

	return "chore"
}
