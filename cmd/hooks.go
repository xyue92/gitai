package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xyue92/gitai/internal/hooks"
)

var (
	hooksForce   bool
	hooksRestore bool
	hooksAll     bool
)

var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "Manage Git hooks for GitAI automation",
	Long: `Install, uninstall, or check status of Git hooks that enable GitAI automation.

Git hooks allow GitAI to automatically generate commit messages when you run 'git commit',
eliminating the need to manually run 'gitai commit'.

Examples:
  gitai hooks install    # Install prepare-commit-msg hook
  gitai hooks status     # Check hook installation status
  gitai hooks uninstall  # Remove GitAI hooks`,
}

var hooksInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install GitAI git hooks",
	Long: `Install Git hooks to enable automatic commit message generation.

After installation, GitAI will automatically generate commit messages when you run:
  git commit

The hook will be skipped when:
  - Using 'git commit -m "message"' (message already provided)
  - During merge commits
  - During rebase/cherry-pick operations
  - Using 'git commit --no-verify' (explicitly skip hooks)

If existing hooks are found, they will be backed up automatically.`,
	Example: `  # Install the prepare-commit-msg hook
  gitai hooks install

  # Force reinstall (overwrite existing hook)
  gitai hooks install --force

  # Install all available hooks
  gitai hooks install --all`,
	RunE: runHooksInstall,
}

var hooksUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall GitAI git hooks",
	Long: `Remove GitAI git hooks from the repository.

If you previously had hooks that were backed up during installation,
you can restore them using the --restore flag.`,
	Example: `  # Remove GitAI hooks
  gitai hooks uninstall

  # Remove hooks and restore backups
  gitai hooks uninstall --restore

  # Remove all GitAI hooks
  gitai hooks uninstall --all`,
	RunE: runHooksUninstall,
}

var hooksStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check GitAI hooks installation status",
	Long:  `Display the current status of GitAI hooks in the repository.`,
	Example: `  # Check hook status
  gitai hooks status`,
	RunE: runHooksStatus,
}

func init() {
	rootCmd.AddCommand(hooksCmd)
	hooksCmd.AddCommand(hooksInstallCmd)
	hooksCmd.AddCommand(hooksUninstallCmd)
	hooksCmd.AddCommand(hooksStatusCmd)

	// Install flags
	hooksInstallCmd.Flags().BoolVarP(&hooksForce, "force", "f", false,
		"Force installation, overwrite existing hooks")
	hooksInstallCmd.Flags().BoolVarP(&hooksAll, "all", "a", false,
		"Install all available hooks (prepare-commit-msg, commit-msg, pre-commit)")

	// Uninstall flags
	hooksUninstallCmd.Flags().BoolVarP(&hooksRestore, "restore", "r", false,
		"Restore backed-up hooks after uninstalling")
	hooksUninstallCmd.Flags().BoolVarP(&hooksAll, "all", "a", false,
		"Uninstall all GitAI hooks")
}

func runHooksInstall(cmd *cobra.Command, args []string) error {
	// Check if in git repository
	if !hooks.IsGitRepository() {
		return fmt.Errorf("not a git repository\nRun this command from within a git repository")
	}

	// Create hook manager
	hm, err := hooks.NewHookManager()
	if err != nil {
		return err
	}

	// Determine which hooks to install
	hookTypes := []string{hooks.PrepareCommitMsg}
	if hooksAll {
		hookTypes = []string{
			hooks.PrepareCommitMsg,
			hooks.CommitMsg,
			hooks.PreCommit,
		}
	}

	fmt.Println("🔧 Installing GitAI hooks...")
	fmt.Println()

	// Install hooks
	installedCount := 0
	for _, hookType := range hookTypes {
		fmt.Printf("  Installing %s hook... ", hookType)

		if err := hm.Install(hookType, hooksForce); err != nil {
			fmt.Printf("❌ Failed\n")
			fmt.Printf("     Error: %v\n", err)
			continue
		}

		fmt.Printf("✅ Done\n")
		installedCount++
	}

	fmt.Println()
	if installedCount > 0 {
		fmt.Println("✨ GitAI hooks installed successfully!")
		fmt.Println()
		fmt.Println("📋 What's next?")
		fmt.Println("  1. Make some changes: git add <files>")
		fmt.Println("  2. Commit normally: git commit")
		fmt.Println("  3. GitAI will automatically generate the commit message!")
		fmt.Println()
		fmt.Println("💡 Tips:")
		fmt.Println("  - Use 'git commit -m \"message\"' to skip GitAI")
		fmt.Println("  - Use 'git commit --no-verify' to bypass all hooks")
		fmt.Println("  - Run 'gitai hooks status' to check installation")
		fmt.Println("  - Run 'gitai hooks uninstall' to remove hooks")
	} else {
		fmt.Println("⚠️  No hooks were installed")
	}

	return nil
}

func runHooksUninstall(cmd *cobra.Command, args []string) error {
	// Check if in git repository
	if !hooks.IsGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	// Create hook manager
	hm, err := hooks.NewHookManager()
	if err != nil {
		return err
	}

	// Determine which hooks to uninstall
	hookTypes := []string{hooks.PrepareCommitMsg}
	if hooksAll {
		hookTypes = []string{
			hooks.PrepareCommitMsg,
			hooks.CommitMsg,
			hooks.PreCommit,
		}
	}

	fmt.Println("🔧 Uninstalling GitAI hooks...")
	fmt.Println()

	// Uninstall hooks
	uninstalledCount := 0
	for _, hookType := range hookTypes {
		fmt.Printf("  Removing %s hook... ", hookType)

		if err := hm.Uninstall(hookType, hooksRestore); err != nil {
			// Ignore "file not found" errors
			if !os.IsNotExist(err) {
				fmt.Printf("❌ Failed\n")
				fmt.Printf("     Error: %v\n", err)
				continue
			}
			fmt.Printf("⚠️  Not found\n")
			continue
		}

		fmt.Printf("✅ Done\n")
		uninstalledCount++
	}

	fmt.Println()
	if uninstalledCount > 0 {
		fmt.Println("✨ GitAI hooks uninstalled successfully!")
		if hooksRestore {
			fmt.Println("   Previous hooks have been restored.")
		}
		fmt.Println()
		fmt.Println("📋 Back to manual mode:")
		fmt.Println("  - Use 'gitai commit' to generate commit messages")
		fmt.Println("  - Or run 'gitai hooks install' to reinstall automation")
	} else {
		fmt.Println("ℹ️  No hooks were uninstalled")
	}

	return nil
}

func runHooksStatus(cmd *cobra.Command, args []string) error {
	// Check if in git repository
	if !hooks.IsGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	// Create hook manager
	hm, err := hooks.NewHookManager()
	if err != nil {
		return err
	}

	// Get status
	status, err := hm.Status()
	if err != nil {
		return err
	}

	fmt.Println("📊 GitAI Hooks Status")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	anyInstalled := false
	for _, hookType := range []string{hooks.PrepareCommitMsg, hooks.CommitMsg, hooks.PreCommit} {
		hs := status[hookType]

		fmt.Printf("🔗 %s\n", hookType)

		if hs.Installed {
			if hs.IsGitAI {
				fmt.Printf("   Status:  ✅ Installed (GitAI)\n")
				anyInstalled = true
			} else {
				fmt.Printf("   Status:  ⚠️  Installed (Not GitAI)\n")
			}
		} else {
			fmt.Printf("   Status:  ❌ Not installed\n")
		}

		if hs.HasBackup {
			fmt.Printf("   Backup:  ✅ Available\n")
		}

		fmt.Println()
	}

	if anyInstalled {
		fmt.Println("✨ GitAI hooks are active!")
		fmt.Println()
		fmt.Println("💡 Usage:")
		fmt.Println("  - Just run 'git commit' and GitAI will generate messages")
		fmt.Println("  - Use 'git commit -m \"message\"' to skip automation")
		fmt.Println("  - Use 'git commit --no-verify' to bypass all hooks")
	} else {
		fmt.Println("ℹ️  GitAI hooks are not installed")
		fmt.Println()
		fmt.Println("📋 To enable automation:")
		fmt.Println("  gitai hooks install")
	}

	return nil
}
