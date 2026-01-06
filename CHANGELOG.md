# Changelog

All notable changes to GitAI will be documented in this file.

## [Unreleased]

### Added
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
