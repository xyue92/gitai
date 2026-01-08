# Ticket/Issue Number Integration Guide

GitAI supports multiple ways to provide and manage ticket numbers (Jira, GitHub Issues, etc.), allowing your commits to automatically include work order information.

## 🚀 Quick Start

### Method 1: Command Line Parameter (Simplest)

```bash
# Directly specify ticket number
gitai commit --ticket PROJ-123

# Short form
gitai commit -k JIRA-456
```

Generated commit will automatically include:
```
feat(api): [PROJ-123] add user authentication

- Implement JWT-based authentication
- Add login and logout endpoints
```

---

### Method 2: Auto-extract from Branch Name (Recommended)

If your branch naming includes ticket numbers, GitAI will automatically recognize them:

```bash
# Branch name examples
git checkout -b feature/PROJ-123-add-login
git checkout -b bugfix/JIRA-456-fix-auth
git checkout -b hotfix/GH-789

# GitAI will automatically extract PROJ-123, JIRA-456, GH-789
gitai commit
```

**Supported branch naming formats**:
- `feature/PROJ-123-description`
- `bugfix/JIRA-456-fix-something`
- `PROJ-789` (using ticket number directly as branch name)
- `#123` (GitHub Issues)
- `GH-456` (GitHub format)

GitAI will ask for confirmation:
```
📝 Git Commit AI Assistant

Found ticket number in branch: PROJ-123
? Use this ticket number?
▸ Yes
  No
```

---

### Method 3: Configure as Required Field

Configure in `.gitcommit.yaml`:

```yaml
# Require ticket number
require_ticket: true
ticket_prefix: "PROJ"                    # Default prefix
ticket_pattern: "[A-Z]+-\\d+"            # Extraction rule
```

**Effect**:
- If branch name contains ticket, auto-extract
- If not, prompt for input:

```
? Enter ticket number (e.g., PROJ-123): ▌
```

Enter `123`, will be automatically formatted as `PROJ-123`

---

## 📝 Configuration Details

### Basic Configuration

```yaml
# .gitcommit.yaml

# Whether to require ticket number
require_ticket: false              # true=required, false=optional

# Default ticket prefix (automatically added when user inputs only numbers)
ticket_prefix: "PROJ"              # e.g., input "123" → "PROJ-123"

# Regular expression to extract ticket from branch name
ticket_pattern: "[A-Z]+-\\d+"      # Matches ABC-123, PROJ-456
```

### Jira Project Configuration

```yaml
require_ticket: true
ticket_prefix: "JIRA"
ticket_pattern: "[A-Z]+-\\d+"

custom_prompt: |
  Format requirement: <type>(<scope>): [JIRA-XXX] <description>

  Must include:
  Jira: JIRA-XXX
  Reviewer: @username
```

Usage:
```bash
# Method 1: Auto-extract from branch name
git checkout -b feature/JIRA-123-new-feature
gitai commit

# Method 2: Manually specify
gitai commit --ticket JIRA-123

# Method 3: Input only numbers, auto-add prefix
? Enter ticket number (e.g., JIRA-123): 123
# Automatically becomes JIRA-123
```

---

### GitHub Issues Configuration

```yaml
require_ticket: true
ticket_prefix: "GH"
ticket_pattern: "#\\d+|GH-\\d+"    # Matches #123 or GH-123
```

Usage:
```bash
# Branch name example
git checkout -b fix/#123-bug-fix
gitai commit
# Auto-extracts #123
```

---

### Chinese Enterprise Configuration

```yaml
require_ticket: true
ticket_prefix: "WO"                # Work Order
ticket_pattern: "WO-\\d{8}-\\d+"   # WO-20250106-001

custom_prompt: |
  Commit format: <type>(<module>): [work-order-number] <description>

  Work order format: WO-YYYYMMDD-XXX

  Must include:
  Work Order: WO-YYYYMMDD-XXX
  Tester: @name
```

---

## 🎯 Use Cases

### Scenario 1: Jira Workflow

**Team Standard**: All commits must link to Jira tickets

**Configuration**:
```yaml
require_ticket: true
ticket_prefix: "PROJ"
ticket_pattern: "PROJ-\\d+"
```

**Daily Usage**:
```bash
# 1. Create branch from Jira
git checkout -b feature/PROJ-1234-add-payment

# 2. Develop code
vim src/payment.js

# 3. Auto-extract ticket on commit
git add .
gitai commit

# ✅ Generates: feat(payment): [PROJ-1234] add payment gateway
```

---

### Scenario 2: GitHub Flow

**Team Standard**: Issue-driven development

**Configuration**:
```yaml
require_ticket: false              # Optional
ticket_prefix: "GH"
ticket_pattern: "#\\d+|GH-\\d+"
```

**Usage**:
```bash
# 1. Create branch from Issue
git checkout -b fix/#456-memory-leak

# 2. GitAI auto-extracts #456
gitai commit
```

---

### Scenario 3: Multiple Projects with Different Standards

**Project A (Jira)**:
```yaml
# ~/projects/project-a/.gitcommit.yaml
require_ticket: true
ticket_prefix: "PROJA"
```

**Project B (GitHub)**:
```yaml
# ~/projects/project-b/.gitcommit.yaml
require_ticket: true
ticket_prefix: "GH"
ticket_pattern: "#\\d+"
```

GitAI will automatically use the current project's configuration!

---

## 🔧 Advanced Features

### Custom Extraction Rules

Supports custom regular expressions:

```yaml
# Match complex formats
ticket_pattern: "(PROJ|TASK|BUG)-\\d+"    # PROJ-123 or TASK-456

# Match multiple formats
ticket_pattern: "[A-Z]{2,10}-\\d+|#\\d+"  # ABC-123 or #456
```

### Smart Formatting

If user input is incomplete, auto-complete:

```yaml
ticket_prefix: "JIRA"
```

User input:
- `123` → Automatically formatted as `JIRA-123`
- `JIRA-123` → Remains unchanged
- `PROJECT-456` → Remains unchanged

### Branch Name Patterns

Supports the following branch naming patterns:

```
✅ feature/PROJ-123-description
✅ bugfix/JIRA-456-fix-bug
✅ hotfix/PROJ-789
✅ PROJ-123
✅ fix/#123
✅ feature/GH-456-new-feature
✅ 123-add-feature (requires ticket_prefix configuration)
```

---

## 📋 Best Practices

### 1. Unified Branch Naming Convention

```bash
# Recommended format
<type>/<ticket>-<description>

# Examples
feature/PROJ-123-add-login
bugfix/PROJ-456-fix-crash
hotfix/PROJ-789-security-patch
```

### 2. Configure Git Branch Template

In `~/.gitconfig` or project `.git/config`:

```ini
[alias]
    # Auto-prompt for ticket when creating branch
    nb = "!f() { \
        read -p 'Ticket number: ' ticket; \
        read -p 'Description: ' desc; \
        git checkout -b \"feature/$ticket-$desc\"; \
    }; f"
```

Usage:
```bash
git nb
Ticket number: PROJ-123
Description: add-payment
# Creates branch: feature/PROJ-123-add-payment
```

### 3. Team Configuration Template

Create a shared team configuration template:

```bash
# .gitcommit.team.yaml (commit to repository)
require_ticket: true
ticket_prefix: "PROJ"
ticket_pattern: "PROJ-\\d+"

custom_prompt: |
  Team commit standards:
  - Must include Jira ticket: [PROJ-XXX]
  - Must include Reviewer: @username
```

After team members pull:
```bash
git pull
cp .gitcommit.team.yaml .gitcommit.yaml
```

---

## ❓ Frequently Asked Questions

### Q: Can I require tickets?

**A**: Yes! Set `require_ticket: true`

```yaml
require_ticket: true
```

If user doesn't provide ticket, GitAI will error:
```
❌ Error: ticket number required but not provided
```

### Q: What if branch name doesn't contain ticket?

**A**: GitAI will prompt for input:

```
? Enter ticket number (e.g., PROJ-123): ▌
```

### Q: Can I skip tickets?

**A**: If `require_ticket: false`, you can skip:

```bash
# Can use normally without providing ticket
gitai commit
```

### Q: Which ticket systems are supported?

**A**: All systems are supported, just configure the correct pattern:

- ✅ Jira
- ✅ GitHub Issues
- ✅ GitLab Issues
- ✅ Azure DevOps
- ✅ Custom ticketing systems

### Q: Where does the ticket number appear?

**A**: In the commit subject line:

```
feat(api): [PROJ-123] add new endpoint

Detailed description...

Jira: PROJ-123
```

### Q: Can I customize ticket format?

**A**: Yes! Specify in `custom_prompt`:

```yaml
custom_prompt: |
  Format requirements:
  - Ticket number must be at the front: [PROJ-123]
  - Or in footer: Ticket: PROJ-123
```

---

## 🎨 Example Commit Output

### With Jira Ticket

```
feat(auth): [JIRA-456] add OAuth2 login support

Implemented OAuth2 authentication flow for enterprise SSO.
This allows users to login using their company credentials.

Business Impact:
- Enables enterprise customer onboarding
- Improves security compliance

Technical Details:
- Added OAuth2 library integration
- Implemented callback endpoint

Jira: JIRA-456
Reviewer: @tech-lead
```

### With GitHub Issue

```
fix(api): [#123] resolve memory leak in connection pool

Fixed connection pool not releasing connections properly.
This was causing server crashes under high load.

- Implement proper connection cleanup
- Add connection timeout handling
- Update connection pool configuration

Closes #123
```

### With Work Order Number (Chinese)

```
feat(payment): [WO-20250106-001] add Alipay payment feature

Implemented Alipay QR code payment interface integration.

Changes:
- Added Alipay SDK integration
- Implemented payment callback handling
- Added payment status synchronization

Business Value:
- Supports more payment methods
- Improves user experience

Work Order: WO-20250106-001
Tester: @qa-engineer
```

---

## 🚀 Summary

GitAI provides 4 ways to manage ticket numbers:

1. **Command Line Parameter** - `--ticket PROJ-123`
2. **Auto-extract from Branch Name** - Extract from `feature/PROJ-123-xxx`
3. **Interactive Input** - Prompt user for input
4. **Configure Default Value** - Use `ticket_prefix` for auto-completion

Choose the method that suits your team to automate commit standards!
