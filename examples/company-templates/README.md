# Company Commit Message Templates

This directory contains example configurations for different company/team commit message conventions. Choose one that matches your organization's standards or use them as inspiration to create your own.

## Available Templates

### 1. Google Style ([google-style.yaml](google-style.yaml))
**Best for**: Teams following Google's engineering practices

**Features**:
- 50-character summary line
- Detailed explanation in body
- Bug/issue number references
- Clean, professional format

**Example Output**:
```
Add user authentication module

This adds JWT-based authentication to secure API endpoints.
Users can now login and receive tokens for authenticated requests.

Bug: 12345
Test: Added unit tests for auth flow
```

---

### 2. Jira Integration ([jira-integration.yaml](jira-integration.yaml))
**Best for**: Teams using Jira/Atlassian tools

**Features**:
- Mandatory Jira ticket numbers
- Reviewer mentions
- Business impact section
- Technical details section

**Example Output**:
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

---

### 3. Chinese Enterprise ([chinese-enterprise.yaml](chinese-enterprise.yaml))
**适用于**: 中国企业开发团队

**特点**:
- 完全中文描述
- 包含影响范围说明
- 测试情况说明
- 关联需求文档

**示例输出**:
```
feat(用户中心): 新增用户登录功能

实现了基于JWT的用户登录认证功能，用户可以通过手机号和验证码登录系统。

改动内容：
- 新增登录接口和验证码发送接口
- 实现JWT token生成和验证逻辑
- 添加登录状态管理中间件

改动原因：
- 满足产品V2.0版本用户登录需求
- 提升系统安全性，替代原有简单密码登录

业务价值：
- 提升用户登录体验，降低密码遗忘率
- 增强系统安全性，符合等保2.0要求

影响范围: 用户中心模块、API网关、前端登录页面
测试情况: 已完成单元测试、集成测试、UAT测试
关联需求: PRD-2024-001-用户登录改造
```

---

### 4. Angular Convention ([angular-style.yaml](angular-style.yaml))
**Best for**: Angular projects or teams using Angular's commit convention

**Features**:
- Strict conventional commits format
- No emojis (professional)
- Breaking change tracking
- Issue references

**Example Output**:
```
feat(parser): add ability to parse arrays

The parser can now handle array syntax in configuration files.
This enables users to define multiple values for a single key.

Previous behavior required separate keys for each value, which
was verbose and error-prone.

Closes #456
BREAKING CHANGE: Array syntax changes the configuration format
```

---

## How to Use

### Option 1: Copy to Your Project

```bash
# Copy the template you want to your project
cp examples/company-templates/jira-integration.yaml .gitcommit.yaml

# Edit to match your team's specifics
vim .gitcommit.yaml
```

### Option 2: Create Custom Template

1. Start with `.gitcommit.example.yaml`
2. Modify the `custom_prompt` section with your company's guidelines
3. Save as `.gitcommit.yaml` in your project root

```yaml
custom_prompt: |
  [Paste your company's commit message requirements here]

  Format: <type>(<scope>): [TICKET-123] <description>

  Required sections:
  - What changed
  - Why it changed
  - Business impact

  Example:
  feat(api): [PROJ-456] add payment processing

  Integrated Stripe payment gateway for subscription billing.

  Business Impact: Enables recurring revenue model
  Testing: Unit + integration tests added
  Ticket: PROJ-456
```

### Option 3: Multiple Projects

If you work on multiple projects with different standards:

```bash
# Project A (uses Jira)
cd ~/projects/project-a
cp ~/gitai/examples/company-templates/jira-integration.yaml .gitcommit.yaml

# Project B (uses Chinese)
cd ~/projects/project-b
cp ~/gitai/examples/company-templates/chinese-enterprise.yaml .gitcommit.yaml

# Project C (Angular style)
cd ~/projects/project-c
cp ~/gitai/examples/company-templates/angular-style.yaml .gitcommit.yaml
```

---

## Customization Guide

### Key Fields to Customize

1. **custom_prompt** - Your company's specific requirements
2. **scopes** - Your project's modules/components
3. **types** - If your company uses custom types
4. **language** - "en" or "zh" or other
5. **detailed_commit** - true (detailed) or false (concise)

### Example: Adding Ticket Numbers

```yaml
custom_prompt: |
  MANDATORY: Include ticket number in format [PROJ-123]

  Format: <type>(<scope>): [TICKET] <description>

  Example:
  feat(auth): [AUTH-789] add SSO support
```

### Example: Required Reviewers

```yaml
custom_prompt: |
  All commits must include:

  Reviewer: @username
  Approved-By: @tech-lead

  Example footer:
  Reviewer: @john-doe
  Approved-By: @tech-lead
```

### Example: Testing Requirements

```yaml
custom_prompt: |
  Must include testing section:

  Testing:
  - Unit tests: [Coverage %]
  - Integration tests: [Pass/Fail]
  - Manual testing: [Scenarios tested]
```

---

## Tips for Writing Good Custom Prompts

1. **Be Specific**: The AI follows your instructions literally
2. **Provide Examples**: Include example commits in your format
3. **List Requirements**: Use bullet points for clarity
4. **Include Format**: Show exact structure expected
5. **Explain Sections**: Describe what each section should contain

### Good Custom Prompt Example

```yaml
custom_prompt: |
  Follow company commit standards:

  REQUIRED FORMAT:
  <type>(<module>): [TICKET] <summary>

  <detailed explanation>

  Impact: <who is affected>
  Testing: <what was tested>
  Ticket: PROJ-XXX

  REQUIREMENTS:
  - Summary max 50 characters
  - Always include ticket number
  - Explain business impact
  - Describe test coverage

  EXAMPLE:
  feat(payments): [PAY-123] add refund processing

  Implemented automatic refund processing for canceled orders.
  Customers now receive refunds within 24 hours.

  Impact: All customers using credit card payments
  Testing: Unit tests + manual testing with test cards
  Ticket: PAY-123
```

---

## Testing Your Template

After creating your template:

```bash
# Test in dry-run mode
cd your-project
git add some-file.js
gitai commit --dry-run

# Check if the output matches your company's format
# If not, adjust the custom_prompt and try again
```

---

## Contributing

Have a template for a popular company format? Submit a PR!

Common formats we'd like to add:
- GitHub Flow style
- GitLab style
- Microsoft/Azure DevOps style
- Linux Kernel style
- Other enterprise formats

---

## FAQ

**Q: Can I use multiple templates?**
A: Use one `.gitcommit.yaml` per project. Switch between projects to use different templates.

**Q: Does this work with all Ollama models?**
A: Yes, but code-focused models (qwen2.5-coder, codellama) work best.

**Q: Can I mix English and Chinese?**
A: Yes, but specify the primary language in the config for best results.

**Q: How strict is the AI at following my template?**
A: Very strict! The AI is instructed to follow your custom_prompt exactly.

**Q: Can I include company-specific terminology?**
A: Absolutely! That's the main purpose. Add your domain language to the custom_prompt.

---

## Support

If you need help creating a custom template for your company:
1. Check the examples above
2. Read the main [README.md](../../README.md)
3. Open an issue with your requirements
