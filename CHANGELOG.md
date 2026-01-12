# Changelog

All notable changes to GitAI will be documented in this file.

## [Unreleased]

## [0.2.0] - 2026-01-12

### Added
- **Self-Update Command**: New `gitai update` command for automatic updates
  - `gitai update` - Update to the latest version
  - `gitai update --check` - Check for updates without installing
  - `gitai update --force` - Force update even if already on latest version
- Automatic platform detection (macOS, Linux, Windows) and architecture (amd64, arm64)
- SHA256 checksum verification for secure updates
- Atomic binary replacement with automatic backup and rollback mechanism
- Unified `checksums.txt` file generation in release workflow

### Changed
- Enhanced GitHub Actions release workflow to generate unified checksums file
- Updated release notes template to include self-update instructions
- Improved documentation with dedicated "Updating" section in README
- Added migration notes for v0.1.0 users

### Migration Guide for v0.1.0 Users
The `update` command was introduced in v0.2.0. To upgrade from v0.1.0, use one of these methods:

1. **Install Script** (Recommended):
   ```bash
   # macOS/Linux
   curl -sSL https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.sh | bash
   ```

2. **Homebrew**:
   ```bash
   brew upgrade gitai
   ```

3. **Manual Download**: Download from [releases](https://github.com/xyue92/gitai/releases/latest)

Once on v0.2.0+, use `gitai update` for all future updates.

## [0.1.0] - 2026-01-XX

### Added
- **Ticket/Issue number integration**: Complete support for Jira, GitHub Issues, and custom ticketing systems
  - New `--ticket` / `-k` flag for command line
  - Auto-extraction from branch names (e.g., `feature/PROJ-123-description`)
  - Interactive prompt when `require_ticket: true`
  - Smart formatting with `ticket_prefix` (e.g., `123` → `PROJ-123`)
  - Configurable patterns via `ticket_pattern` regex
  - Support for multiple formats: JIRA-123, #456, GH-789, custom formats
  - AI automatically includes ticket number in commit message
  - Full documentation in `TICKET_INTEGRATION_GUIDE.md`

- **Custom company templates**: Full support for company/team commit standards
  - New `custom_prompt` field in config for company-specific guidelines
  - AI strictly follows your company's commit message requirements
  - Perfect for teams with mandatory formats, ticket numbers, or specific sections
  - Example templates provided for popular formats:
    - Google Style
    - Jira Integration (with ticket numbers and reviewers)
    - Chinese Enterprise (完整中文支持)
    - Angular Convention
  - Templates located in `examples/company-templates/`

- **Detailed commit messages**: New `detailed_commit` configuration option
  - When `true` (default): Generates multi-line commits with subject + body (2-4 bullet points explaining changes)
  - When `false`: Generates concise single-line commits (subject only)
  - Addresses issue where all commits were too short
  - Provides better context about WHAT changed and WHY

### Changed
- Updated prompt generation to support custom company guidelines
- Custom prompt is inserted before task description for maximum AI compliance
- Enhanced example config with detailed instructions on using custom_prompt
- Updated prompt generation to support both detailed and concise modes
- Modified default behavior to generate detailed commits by default
- Improved prompt instructions to clearly separate subject line and body requirements

### Example

**Before** (concise only):
```
feat(api): add user authentication
```

**After** (detailed mode - default):
```
feat(api): add user authentication endpoint

- Implement JWT-based authentication
- Add login and logout endpoints
- Include token validation middleware
```

## [1.0.0] - 2025-01-06

### Added
- Initial release of GitAI
- Local Ollama integration for AI-powered commit messages
- Interactive commit type and scope selection
- Context-aware message generation (README, recent commits, branch)
- Full YAML configuration support
- Multi-language support (English/Chinese)
- Comprehensive error handling
- Dry-run and generate-only modes
- Complete documentation (README, Quick Start, Examples, Installation)

### Features
- CLI commands: `commit`, `generate`, `config`
- Support for Conventional Commits format
- Customizable commit types and scopes
- Project context collection
- Beautiful terminal UI with colors
- 1,581 lines of Go code
- 13 source files + tests
- 46+ KB of documentation
