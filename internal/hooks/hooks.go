package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// Hook types
	PrepareCommitMsg = "prepare-commit-msg"
	CommitMsg        = "commit-msg"
	PreCommit        = "pre-commit"
)

// HookManager manages git hooks installation and removal
type HookManager struct {
	repoPath  string
	hooksPath string
}

// NewHookManager creates a new hook manager
func NewHookManager() (*HookManager, error) {
	// Get git repository root
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("not a git repository")
	}

	repoPath := strings.TrimSpace(string(output))
	hooksPath := filepath.Join(repoPath, ".git", "hooks")

	return &HookManager{
		repoPath:  repoPath,
		hooksPath: hooksPath,
	}, nil
}

// Install installs GitAI hooks
func (hm *HookManager) Install(hookType string, force bool) error {
	hookPath := filepath.Join(hm.hooksPath, hookType)

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		if !force {
			// Backup existing hook
			if err := hm.backupHook(hookType); err != nil {
				return fmt.Errorf("failed to backup existing hook: %w", err)
			}
		}
	}

	// Create hooks directory if it doesn't exist
	if err := os.MkdirAll(hm.hooksPath, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Get hook template
	content := hm.getHookTemplate(hookType)
	if content == "" {
		return fmt.Errorf("unsupported hook type: %s", hookType)
	}

	// Write hook file
	if err := os.WriteFile(hookPath, []byte(content), 0755); err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	return nil
}

// Uninstall removes GitAI hooks
func (hm *HookManager) Uninstall(hookType string, restore bool) error {
	hookPath := filepath.Join(hm.hooksPath, hookType)

	// Remove hook file
	if err := os.Remove(hookPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove hook: %w", err)
	}

	// Restore backup if requested
	if restore {
		if err := hm.restoreHook(hookType); err != nil {
			// Don't fail if no backup exists
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to restore backup: %w", err)
			}
		}
	}

	return nil
}

// Status checks the status of GitAI hooks
func (hm *HookManager) Status() (map[string]HookStatus, error) {
	status := make(map[string]HookStatus)
	hookTypes := []string{PrepareCommitMsg, CommitMsg, PreCommit}

	for _, hookType := range hookTypes {
		hookPath := filepath.Join(hm.hooksPath, hookType)
		backupPath := hookPath + ".gitai-backup"

		hs := HookStatus{
			Type:      hookType,
			Installed: false,
			HasBackup: false,
			IsGitAI:   false,
		}

		// Check if hook exists
		if content, err := os.ReadFile(hookPath); err == nil {
			hs.Installed = true
			// Check if it's a GitAI hook
			if strings.Contains(string(content), "GitAI") {
				hs.IsGitAI = true
			}
		}

		// Check if backup exists
		if _, err := os.Stat(backupPath); err == nil {
			hs.HasBackup = true
		}

		status[hookType] = hs
	}

	return status, nil
}

// HookStatus represents the status of a hook
type HookStatus struct {
	Type      string
	Installed bool
	HasBackup bool
	IsGitAI   bool
}

// backupHook creates a backup of existing hook
func (hm *HookManager) backupHook(hookType string) error {
	hookPath := filepath.Join(hm.hooksPath, hookType)
	backupPath := hookPath + ".gitai-backup"

	// Read existing hook
	content, err := os.ReadFile(hookPath)
	if err != nil {
		return err
	}

	// Don't backup if it's already a GitAI hook
	if strings.Contains(string(content), "GitAI") {
		return nil
	}

	// Write backup
	return os.WriteFile(backupPath, content, 0755)
}

// restoreHook restores a backed-up hook
func (hm *HookManager) restoreHook(hookType string) error {
	hookPath := filepath.Join(hm.hooksPath, hookType)
	backupPath := hookPath + ".gitai-backup"

	// Read backup
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}

	// Restore hook
	if err := os.WriteFile(hookPath, content, 0755); err != nil {
		return err
	}

	// Remove backup
	return os.Remove(backupPath)
}

// getHookTemplate returns the hook script template
func (hm *HookManager) getHookTemplate(hookType string) string {
	switch hookType {
	case PrepareCommitMsg:
		return hm.getPrepareCommitMsgTemplate()
	case CommitMsg:
		return hm.getCommitMsgTemplate()
	case PreCommit:
		return hm.getPreCommitTemplate()
	default:
		return ""
	}
}

// getPrepareCommitMsgTemplate returns prepare-commit-msg hook template
func (hm *HookManager) getPrepareCommitMsgTemplate() string {
	return `#!/bin/sh
# GitAI - Auto-generated Git Hook
# This hook automatically generates commit messages using GitAI

COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2
SHA1=$3

# Skip if committing with -m, --amend, merge, or squash
if [ "$COMMIT_SOURCE" = "message" ] || [ "$COMMIT_SOURCE" = "merge" ] || [ "$COMMIT_SOURCE" = "squash" ]; then
    exit 0
fi

# Skip if template already has content (except comments)
if [ -s "$COMMIT_MSG_FILE" ]; then
    # Check if file has non-comment content
    if grep -qv '^#' "$COMMIT_MSG_FILE" 2>/dev/null; then
        exit 0
    fi
fi

# Check if gitai is available
if ! command -v gitai >/dev/null 2>&1; then
    echo "GitAI hook installed but gitai command not found in PATH" >&2
    echo "Run 'gitai hooks uninstall' to remove this hook" >&2
    exit 0
fi

# Check if there are staged changes
if ! git diff --cached --quiet 2>/dev/null; then
    # Generate commit message with GitAI
    echo "🤖 Generating commit message with GitAI..." >&2

    # Run gitai generate and capture output
    if GENERATED_MSG=$(gitai generate --quiet 2>&1); then
        # Write generated message to commit message file
        echo "$GENERATED_MSG" > "$COMMIT_MSG_FILE"
        echo "✅ Commit message generated. Edit if needed." >&2
    else
        echo "⚠️  GitAI generation failed. Write message manually." >&2
        echo "   Error: $GENERATED_MSG" >&2
    fi
fi

exit 0
`
}

// getCommitMsgTemplate returns commit-msg hook template
func (hm *HookManager) getCommitMsgTemplate() string {
	return `#!/bin/sh
# GitAI - Commit Message Validator
# This hook validates commit messages

COMMIT_MSG_FILE=$1

# Read the commit message
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Skip merge commits
if echo "$COMMIT_MSG" | grep -q '^Merge'; then
    exit 0
fi

# Basic validation (optional - can be customized)
# Check if message follows conventional commits format
if ! echo "$COMMIT_MSG" | grep -qE '^(feat|fix|docs|style|refactor|perf|test|chore|ci|build)(\(.+\))?: .+'; then
    echo "⚠️  Commit message doesn't follow Conventional Commits format" >&2
    echo "   Expected: <type>(<scope>): <subject>" >&2
    echo "   Example: feat(auth): add user login" >&2
    echo "" >&2
    echo "   Use 'git commit --no-verify' to skip this check" >&2
    # Don't fail, just warn
    # exit 1
fi

exit 0
`
}

// getPreCommitTemplate returns pre-commit hook template
func (hm *HookManager) getPreCommitTemplate() string {
	return `#!/bin/sh
# GitAI - Pre-commit Hook
# This hook runs before commit to check for issues

# Check if there are staged changes
if git diff --cached --quiet; then
    echo "⚠️  No staged changes to commit" >&2
    exit 1
fi

# Additional checks can be added here:
# - Linting
# - Formatting
# - Tests
# etc.

exit 0
`
}

// GetGitAIPath returns the path to gitai binary
func GetGitAIPath() (string, error) {
	// Try to find gitai in PATH
	path, err := exec.LookPath("gitai")
	if err != nil {
		return "", fmt.Errorf("gitai not found in PATH")
	}
	return path, nil
}

// IsGitRepository checks if current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}
