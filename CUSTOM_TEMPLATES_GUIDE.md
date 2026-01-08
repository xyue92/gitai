# Custom Company Commit Standards Guide

## Quick Start

### Method 1: Use Pre-made Templates

We provide 4 commonly used company commit standard templates that you can use directly:

```bash
# View all templates
ls examples/company-templates/

# Copy the template you need to your project root
cp examples/company-templates/jira-integration.yaml .gitcommit.yaml

# Start using
git add .
gitai commit
```

### Method 2: Paste Your Company Standards

1. Create a configuration file:
```bash
gitai config --init
```

2. Edit `.gitcommit.yaml` and paste your company's commit standards in the `custom_prompt` section:

```yaml
custom_prompt: |
  [Paste your company's commit standards document here]

  For example:
  Commit format requirements:
  - Must include Jira ticket: [PROJ-123]
  - Must include reviewer: @name
  - Must explain business impact
  - Must describe testing status
```

3. Use:
```bash
git add .
gitai commit
```

The AI will strictly follow the standards you pasted to generate commit messages!

---

## Pre-made Template Documentation

### 1. Jira Integration Template (Recommended for Enterprises)

**Use case**: Teams using Jira for requirement management

**File**: `examples/company-templates/jira-integration.yaml`

**Features**:
- ✅ Mandatory Jira ticket number
- ✅ Includes Reviewer field
- ✅ Distinguishes business impact from technical details
- ✅ Complete footer information

**Generated Example**:
```
feat(auth): [AUTH-456] add OAuth2 login support

Implemented OAuth2 authentication flow for enterprise SSO.
This allows users to login using their company credentials.

Business Impact:
- Reduces password management overhead
- Improves security compliance
- Enables enterprise customer onboarding

Technical Details:
- Added OAuth2 library integration
- Implemented callback endpoint
- Updated user session management

Jira: AUTH-456
Reviewer: @tech-lead
```

**How to use**:
```bash
cp examples/company-templates/jira-integration.yaml .gitcommit.yaml
# Modify the example ticket number prefix to match your company's
vim .gitcommit.yaml
```

---

### 2. Chinese Enterprise Template (Recommended for Domestic Teams)

**Use case**: Chinese companies requiring Chinese commits with PRD documentation management

**File**: `examples/company-templates/chinese-enterprise.yaml`

**Features**:
- ✅ Complete Chinese descriptions
- ✅ Includes impact scope
- ✅ Includes testing status
- ✅ Links to PRD requirement documents
- ✅ Complies with Chinese enterprise development process

**Generated Example**:
```
feat(user-center): add user login feature

Implemented JWT-based user login authentication, users can login to the system via phone number and verification code.

Changes:
- Added login API and verification code sending API
- Implemented JWT token generation and verification logic
- Added login state management middleware

Reason for Changes:
- Meets product V2.0 user login requirements
- Improves system security, replacing original simple password login

Business Value:
- Improves user login experience, reduces password forgetting rate
- Enhances system security, meets Security Level Protection 2.0 requirements

Impact Scope: User center module, API gateway, frontend login page
Testing Status: Completed unit tests, integration tests, UAT tests
Related Requirements: PRD-2024-001-User-Login-Redesign
```

**How to use**:
```bash
cp examples/company-templates/chinese-enterprise.yaml .gitcommit.yaml
# Modify scopes to match your project's module names
vim .gitcommit.yaml
```

---

### 3. Google Style Template

**Use case**: Teams pursuing concise professionalism

**File**: `examples/company-templates/google-style.yaml`

**Features**:
- ✅ 50-character short title
- ✅ Detailed body explanation
- ✅ Bug number references
- ✅ Test descriptions

**Generated Example**:
```
Add user authentication module

This adds JWT-based authentication to secure API endpoints.
Users can now login and receive tokens for authenticated requests.

Bug: 12345
Test: Added unit tests for auth flow
```

---

### 4. Angular Convention Template

**Use case**: Angular projects or teams following strict Conventional Commits

**File**: `examples/company-templates/angular-style.yaml`

**Features**:
- ✅ Strict conventional commits
- ✅ Breaking changes tracking
- ✅ Issue references
- ✅ No emoji (professional style)

**Generated Example**:
```
feat(parser): add ability to parse arrays

The parser can now handle array syntax in configuration files.
This enables users to define multiple values for a single key.

Closes #456
BREAKING CHANGE: Array syntax changes the configuration format
```

---

## Customize Your Own Standards

### Step 1: Create Configuration File

```bash
gitai config --init
```

### Step 2: Edit custom_prompt

Open `.gitcommit.yaml` and find the `custom_prompt` section:

```yaml
custom_prompt: |
  # Paste your company's commit standards here
```

### Step 3: Paste Company Standards

Paste your company's commit standards document directly. For example:

#### Example 1: Mandatory Work Order Number

```yaml
custom_prompt: |
  Commit Standards:
  - Format: <type>(<module>): [work-order-number] <description>
  - Work order format: WO-YYYYMMDD-XXX
  - Work order number is mandatory

  Example:
  feat(payment): [WO-20250106-001] add Alipay payment

  Implemented Alipay QR code payment feature with automatic order amount calculation.

  Work Order: WO-20250106-001
  Tester: @qa-engineer
```

#### Example 2: Mandatory Reviewers

```yaml
custom_prompt: |
  Commit Standards:
  - All commits must include Code Reviewer
  - All commits must include QA Tester
  - Format: Reviewed-by: @username, Tested-by: @username

  Example:
  feat(api): add payment endpoint

  Added new payment processing endpoint for Stripe integration.

  Reviewed-by: @john-smith
  Tested-by: @jane-doe
```

#### Example 3: Mandatory Impact Scope

```yaml
custom_prompt: |
  Commit Requirements:
  - Must specify "Impact Scope" (frontend/backend/database/all)
  - Must specify "Requires Release" (yes/no)
  - Must specify "Rollback Plan"

  Example:
  feat(orders): add order cancellation feature

  Users can cancel unpaid orders within 30 minutes.

  Impact Scope: Backend API + Frontend order page
  Requires Release: Yes
  Rollback Plan: Rollback to previous stable version
```

---

## Advanced Tips

### Tip 1: Different Standards for Multiple Projects

If you have multiple projects with different standards:

```bash
# Project A - Using Jira
cd ~/projects/project-a
cp ~/gitai/examples/company-templates/jira-integration.yaml .gitcommit.yaml

# Project B - Using Chinese standards
cd ~/projects/project-b
cp ~/gitai/examples/company-templates/chinese-enterprise.yaml .gitcommit.yaml

# Project C - Custom standards
cd ~/projects/project-c
gitai config --init
vim .gitcommit.yaml  # Paste company standards
```

GitAI will automatically use each project's `.gitcommit.yaml`!

### Tip 2: Team Shared Configuration

Commit the configuration file to your git repository for the whole team to share:

```bash
# 1. Create team configuration
gitai config --init

# 2. Edit to team standards
vim .gitcommit.yaml

# 3. Commit to repository
git add .gitcommit.yaml
git commit -m "chore: add team commit message standards"
git push

# 4. Team members automatically use it after pulling
git pull  # Other members run this
gitai commit  # Automatically uses team standards
```

### Tip 3: Use Environment Variables to Distinguish Environments

```yaml
custom_prompt: |
  # Development environment can be simpler
  # Production environment must be detailed

  {% if env == "production" %}
  Must include:
  - Complete test report
  - Launch checklist
  - Rollback plan
  {% else %}
  Can simplify format
  {% endif %}
```

---

## Testing Your Configuration

After configuration is complete, test if it meets expectations:

```bash
# 1. Make some changes
echo "test" > test.txt
git add test.txt

# 2. Test with dry-run mode
gitai commit --dry-run

# 3. Check if generated commit meets company standards
# If not, adjust custom_prompt and try again
```

---

## Common Questions

### Q: Will AI strictly follow my standards?

**A**: Yes! AI will strictly generate commits according to your `custom_prompt`. We explicitly instruct in the prompt:
> "IMPORTANT: Follow the above guidelines strictly when generating the commit message."

### Q: Can I mix Chinese and English?

**A**: Yes! But we recommend specifying the main language:
```yaml
language: "zh"  # or "en"
custom_prompt: |
  Mainly Chinese, technical terms can be in English
  Example: feat(API): add JWT authentication
```

### Q: Can I force specific fields?

**A**: Absolutely! Clearly state in `custom_prompt`:
```yaml
custom_prompt: |
  Must include the following fields:
  Ticket: XXX-123
  Reviewer: @name
  Testing: description

  If any field is missing, the commit is invalid!
```

AI will generate commits containing these fields.

### Q: What if company standards are very long?

**A**: No problem, `custom_prompt` supports multi-line long text:
```yaml
custom_prompt: |
  [Paste your complete multi-page standards document]
  ...
  ...
  [All requirements pasted here]
```

### Q: Can I reference external files?

**A**: Not currently supported, but you can copy and paste. We recommend writing standards directly in the config file because:
- Easier version control
- Simpler team sharing
- No dependency on external files

---

## Real-World Cases

### Case 1: A Fintech Company

**Requirements**:
- Must include JIRA ticket
- Must include security reviewer
- Must state if customer data is involved

**Configuration**:
```yaml
custom_prompt: |
  Fintech Company Commit Standards:

  Format: <type>(<module>): [JIRA-XXX] <description>

  Must include:
  Security Review: @security-lead
  Customer Data: Yes/No
  Compliance: Checked/Waived

  Example:
  feat(payment): [PAY-789] add encryption for card data

  Implemented AES-256 encryption for credit card storage.

  Security Review: @security-lead
  Customer Data: Yes
  Compliance: Checked - Meets PCI-DSS requirements
```

### Case 2: A Large Internet Company

**Requirements**:
- Chinese commits
- Must link to PRD
- Must describe gradual rollout plan

**Configuration**:
```yaml
language: "zh"
custom_prompt: |
  Large Company Commit Standards:

  Format: <type>(<business-domain>): <requirement-number> <description>

  Required fields:
  - Related PRD: PRD-YYYYMMDD-XXX
  - Rollout Plan: description
  - Monitoring Metrics: list

  Example:
  feat(recommendation): PRD-20250106-001 add personalized recommendation algorithm

  Implemented collaborative filtering-based personalized recommendation.

  Related PRD: PRD-20250106-001
  Rollout Plan: 10% -> 30% -> 100%, monitor 24 hours each stage
  Monitoring Metrics: CTR, conversion rate, page dwell time
```

---

## Summary

Using GitAI's custom template feature, you can:

✅ **Fully automate** company commit standards
✅ **Zero learning cost**: Just paste company documents
✅ **Team collaboration**: Commit config file to repository for sharing
✅ **Multi-project support**: Independent configuration per project
✅ **Internationalization**: Supports Chinese, English, and other languages

Start using:
```bash
# Choose a template
cp examples/company-templates/jira-integration.yaml .gitcommit.yaml

# Or write your own
gitai config --init

# Start enjoying automated commits
gitai commit
```

🎉 Say goodbye to manual commits, let AI help you strictly follow company standards!
