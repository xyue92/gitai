# GitAI - AI-Powered Git Commit Message Generator

GitAI is a CLI tool that uses local Ollama models to generate intelligent, context-aware Git commit messages following the Conventional Commits format.

## Features

- **🪝 Git Hooks Automation**: One-time setup, then just `git commit` - GitAI automatically generates messages
- **Privacy First**: Uses local Ollama models - your code never leaves your machine
- **Interactive**: Select commit type, scope, and review generated messages
- **Context-Aware**: Analyzes project structure, recent commits, and README to generate better messages
- **Intelligent Diff Analysis**: Smart analysis of code changes with function/class extraction and complexity evaluation
- **Multilingual**: Support for 10 languages with auto-detection (en, zh, ja, ko, de, fr, es, pt, ru, it)
- **Commit Statistics**: Comprehensive analysis of commit history with visual insights and recommendations
- **Configurable**: Customize commit types, scopes, and templates per project
- **Smart**: Understands git diff and generates meaningful commit messages
- **Self-Updating**: Built-in update command to keep GitAI up to date

## Prerequisites

1. **Ollama** - For running local AI models (required)
2. **Git** - Already installed on most systems (required)
3. **Go 1.21+** - Only needed if building from source (optional)

### Install Ollama

**macOS:**
```bash
# Option 1: Using Homebrew (Recommended)
brew install ollama

# Option 2: Download from website
# Visit https://ollama.com/download and download the macOS app
```

**Linux:**
```bash
curl -fsSL https://ollama.com/install.sh | sh
```

**Windows:**
```bash
# Download installer from https://ollama.com/download
```

**Start Ollama:**
```bash
# macOS: Ollama runs automatically after installation
# You can also start it manually:
ollama serve

# Linux: Start Ollama service
ollama serve
```

### Pull an AI Model

```bash
# Recommended for code (fast and accurate)
ollama pull qwen2.5-coder:7b

# Alternative models
ollama pull mistral:7b
ollama pull codellama:7b
```

## Installation

### Quick Install (Recommended)

**macOS and Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.ps1 | iex
```

### Homebrew (macOS/Linux)

```bash
# Add the tap (first time only)
brew tap xyue92/tap

# Install
brew install gitai

# Update via Homebrew
brew upgrade gitai
```

### Manual Installation

Download the pre-built binary for your platform from [GitHub Releases](https://github.com/xyue92/gitai/releases/latest):

**macOS:**
```bash
# For Apple Silicon (M1/M2/M3)
wget https://github.com/xyue92/gitai/releases/latest/download/gitai-darwin-arm64
chmod +x gitai-darwin-arm64
sudo mv gitai-darwin-arm64 /usr/local/bin/gitai

# For Intel
wget https://github.com/xyue92/gitai/releases/latest/download/gitai-darwin-amd64
chmod +x gitai-darwin-amd64
sudo mv gitai-darwin-amd64 /usr/local/bin/gitai
```

**Linux:**
```bash
# For AMD64
wget https://github.com/xyue92/gitai/releases/latest/download/gitai-linux-amd64
chmod +x gitai-linux-amd64
sudo mv gitai-linux-amd64 /usr/local/bin/gitai

# For ARM64
wget https://github.com/xyue92/gitai/releases/latest/download/gitai-linux-arm64
chmod +x gitai-linux-arm64
sudo mv gitai-linux-arm64 /usr/local/bin/gitai
```

**Windows:**
1. Download [gitai-windows-amd64.exe](https://github.com/xyue92/gitai/releases/latest/download/gitai-windows-amd64.exe)
2. Rename to `gitai.exe`
3. Move to a directory in your PATH (e.g., `C:\Program Files\GitAI\`)

### Build from Source

If you prefer to build from source or want to contribute:

```bash
# Clone the repository
git clone https://github.com/xyue92/gitai.git
cd gitai

# Install dependencies
go mod download

# Build
go build -o gitai

# Move to PATH (optional)
sudo mv gitai /usr/local/bin/
```

### Install with Go

```bash
go install github.com/xyue92/gitai@latest
```

## Quick Start

### 🚀 Recommended: Automated Mode (with Git Hooks)

**One-time setup, then forget about it!**

1. **Navigate to your git repository & install hooks**:
```bash
cd your-project
gitai hooks install
```

2. **Use Git normally - GitAI works automatically**:
```bash
git add .
git commit        # GitAI automatically generates the message!
```

That's it! From now on, every `git commit` will auto-generate messages.

### 🔧 Alternative: Manual Mode

If you prefer manual control:

1. **Navigate to your git repository**:
```bash
cd your-project
```

2. **Stage your changes**:
```bash
git add .
```

3. **Run GitAI**:
```bash
gitai commit
```

4. **Follow the interactive prompts**:
   - Select commit type (feat, fix, docs, etc.)
   - Enter scope (optional)
   - Review AI-generated message
   - Choose to use, regenerate, edit, or cancel

## Updating

### Self-Update (Recommended for v0.2.0+)

GitAI includes a built-in self-update command that works regardless of how you installed it:

```bash
# Update to the latest version
gitai update

# Check for updates without installing
gitai update --check

# Force update even if already on latest version
gitai update --force
```

The update command will:
- Check GitHub releases for the latest version
- Download the appropriate binary for your platform
- Verify the binary using SHA256 checksums
- Replace your current binary with the new version

> **Note for v0.1.0 users:** The `update` command was introduced in v0.2.0. If you're on v0.1.0, use one of the alternative methods below to upgrade first.

### Alternative Update Methods

**Via Install Script (works for all versions):**
```bash
# macOS/Linux
curl -sSL https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.sh | bash

# Windows (PowerShell)
iwr -useb https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.ps1 | iex
```

**Via Homebrew:**
```bash
brew upgrade gitai
```

## Usage

### Basic Commands

#### Generate and Commit
```bash
gitai commit
```

#### Generate Only (No Commit)
```bash
gitai generate
```

#### Dry Run
```bash
gitai commit --dry-run
```

#### Update GitAI
```bash
gitai update
```

#### Manage Git Hooks (Automation)
```bash
# Install hooks for automatic message generation
gitai hooks install

# Check hook status
gitai hooks status

# Uninstall hooks
gitai hooks uninstall
```

#### View Commit Statistics
```bash
# Show stats for last 100 commits (default)
gitai stats

# Analyze more commits
gitai stats --limit 500

# Export to JSON
gitai stats --export stats.json
```

### Command Flags

#### Commit Command
```bash
gitai commit [flags]

Flags:
  -d, --dry-run          Show message without committing
  -t, --type string      Commit type (skip selection)
  -s, --scope string     Commit scope (skip selection)
  -l, --language string  Message language (en/zh)
  -m, --model string     Ollama model to use
```

#### Examples
```bash
# Skip interactive type selection
gitai commit --type feat --scope api

# Use different model
gitai commit --model mistral:7b

# Generate Chinese commit message
gitai commit --language zh

# Just see what would be generated
gitai commit --dry-run
```

#### Stats Command
```bash
gitai stats [flags]

Flags:
  -n, --limit int       Number of commits to analyze (default 100)
  -e, --export string   Export statistics to JSON file
```

### Configuration Commands

#### Initialize Config File
```bash
gitai config --init
```

This creates `.gitcommit.yaml` in your current directory.

#### Show Current Config
```bash
gitai config --show
```

## Configuration

GitAI looks for configuration in this order:
1. `.gitcommit.yaml` (current directory)
2. `~/.gitcommit.yaml` (home directory)
3. Default configuration

### Example Configuration

Create `.gitcommit.yaml` in your project root:

```yaml
# Ollama model to use
model: "qwen2.5-coder:7b"

# Language for commit messages
language: "en"

# Commit message template
template: "{type}{scope}: {emoji} {message}"

# Commit types
types:
  - name: "feat"
    desc: "A new feature"
    emoji: "✨"

  - name: "fix"
    desc: "A bug fix"
    emoji: "🐛"

  - name: "docs"
    desc: "Documentation changes"
    emoji: "📝"

# Project-specific scopes
scopes:
  - "api"
  - "ui"
  - "auth"
  - "db"
```

See [.gitcommit.example.yaml](.gitcommit.example.yaml) for a complete example.

## How It Works

1. **Analyzes Context**: Reads git diff, recent commits, README, and project structure
2. **Builds Smart Prompt**: Creates a detailed prompt with context for the AI model
3. **Generates Message**: Sends prompt to local Ollama model
4. **Interactive Review**: Shows generated message and allows editing
5. **Commits**: Executes `git commit` with the final message

## Supported Models

GitAI works with any Ollama model, but code-specialized models work best:

| Model | Size | Speed | Quality | Recommended |
|-------|------|-------|---------|-------------|
| qwen2.5-coder:7b | 4.7GB | Fast | Excellent | ⭐ Yes |
| codellama:7b | 3.8GB | Fast | Very Good | ⭐ Yes |
| mistral:7b | 4.1GB | Fast | Good | Yes |
| deepseek-coder:6.7b | 3.8GB | Fast | Very Good | Yes |
| llama3:8b | 4.7GB | Medium | Good | OK |

## Conventional Commits

GitAI follows the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Default Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code formatting
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding/updating tests
- `chore`: Build/tool changes
- `ci`: CI configuration
- `build`: Build system changes

## Examples

### Example Output

```
📝 Git Commit AI Assistant

Changed files (3):
  ✓ src/api/auth.go (+45, -12)
  ✓ internal/middleware/jwt.go (+23, -5)
  ✓ README.md (+10, -0)

? Select commit type:
  ✨ feat - A new feature
▸ 🐛 fix - A bug fix
  📝 docs - Documentation changes

? Select scope: api

🤖 Generating commit message...

Generated message:
┌─────────────────────────────────────────────┐
│ fix(api): resolve JWT token expiration bug │
└─────────────────────────────────────────────┘

? What do you want to do?
▸ ✅ Use this message
  🔄 Regenerate
  ✏️  Edit manually
  ❌ Cancel

✨ Commit created successfully!

Commit message:
fix(api): resolve JWT token expiration bug

Files changed:
  src/api/auth.go
  internal/middleware/jwt.go
  README.md

View commit: git show HEAD
```

## Troubleshooting

### Error: Cannot connect to Ollama

**Solution**: Make sure Ollama is running:
```bash
ollama serve
```

### Error: Model not found

**Solution**: Pull the model first:
```bash
ollama pull qwen2.5-coder:7b
```

### Error: No staged changes found

**Solution**: Stage your changes:
```bash
git add <files>
```

### Error: Not a git repository

**Solution**: Initialize git:
```bash
git init
```

### Slow generation

**Solutions**:
- Use a smaller model (`mistral:7b` instead of larger models)
- Reduce `max_diff_length` in config
- Ensure Ollama has enough RAM allocated

## Development

### Project Structure

```
gitai/
├── cmd/              # CLI commands
│   ├── commit.go     # Commit command
│   ├── generate.go   # Generate command
│   └── config.go     # Config command
├── internal/
│   ├── ai/           # AI/Ollama integration
│   ├── git/          # Git operations
│   ├── config/       # Configuration management
│   └── ui/           # User interface
└── main.go
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build for current platform
go build -o gitai

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o gitai-linux-amd64
GOOS=darwin GOARCH=arm64 go build -o gitai-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o gitai-windows-amd64.exe
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Documentation

📚 **Comprehensive Guides:**
- [Git Hooks Automation Guide](docs/GIT_HOOKS_GUIDE.md) - Set up automatic commit message generation ⭐ NEW
- [Intelligent Diff Analysis Guide](INTELLIGENT_DIFF_GUIDE.md) - Deep dive into smart diff analysis
- [Multilingual Support Guide](MULTILINGUAL_GUIDE.md) - Using GitAI in multiple languages
- [Commit Statistics Guide](docs/STATS_GUIDE.md) - Analyzing commit patterns and quality
- [Multilingual Quick Start](docs/QUICKSTART_MULTILINGUAL.md) - Get started with multilingual commits
- [Multilingual Features](MULTILINGUAL_FEATURES.md) - Technical implementation details

## Acknowledgments

- [Ollama](https://ollama.com/) - Local AI model runtime
- [Conventional Commits](https://www.conventionalcommits.org/) - Commit message specification
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompts
