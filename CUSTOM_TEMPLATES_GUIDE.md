# 自定义公司Commit规范使用指南

## 快速开始

### 方式1：使用现成模板

我们提供了4个常用的公司commit规范模板，直接复制使用：

```bash
# 查看所有模板
ls examples/company-templates/

# 复制你需要的模板到项目根目录
cp examples/company-templates/jira-integration.yaml .gitcommit.yaml

# 开始使用
git add .
gitai commit
```

### 方式2：粘贴公司规范

1. 创建配置文件：
```bash
gitai config --init
```

2. 编辑 `.gitcommit.yaml`，在 `custom_prompt` 部分粘贴你公司的commit规范：

```yaml
custom_prompt: |
  [直接粘贴你们公司的commit规范文档]

  比如：
  提交格式要求：
  - 必须包含Jira ticket: [PROJ-123]
  - 必须包含reviewer: @姓名
  - 必须说明业务影响
  - 必须说明测试情况
```

3. 使用：
```bash
git add .
gitai commit
```

AI会严格按照你粘贴的规范生成commit消息！

---

## 现成模板说明

### 1. Jira集成模板 (推荐企业使用)

**适用场景**：使用Jira管理需求的团队

**文件**：`examples/company-templates/jira-integration.yaml`

**特点**：
- ✅ 强制包含Jira ticket编号
- ✅ 包含Reviewer字段
- ✅ 区分业务影响和技术细节
- ✅ 完整的footer信息

**生成示例**：
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

**使用方法**：
```bash
cp examples/company-templates/jira-integration.yaml .gitcommit.yaml
# 修改模板中的示例ticket编号为你们公司的前缀
vim .gitcommit.yaml
```

---

### 2. 中国企业模板 (推荐国内团队)

**适用场景**：国内企业，要求中文commit，有PRD文档管理

**文件**：`examples/company-templates/chinese-enterprise.yaml`

**特点**：
- ✅ 完全中文描述
- ✅ 包含影响范围
- ✅ 包含测试情况
- ✅ 关联PRD需求文档
- ✅ 符合国内企业开发流程

**生成示例**：
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

**使用方法**：
```bash
cp examples/company-templates/chinese-enterprise.yaml .gitcommit.yaml
# 修改scope为你们项目的模块名
vim .gitcommit.yaml
```

---

### 3. Google风格模板

**适用场景**：追求简洁专业的团队

**文件**：`examples/company-templates/google-style.yaml`

**特点**：
- ✅ 50字符简短标题
- ✅ 详细的body解释
- ✅ Bug编号引用
- ✅ Test说明

**生成示例**：
```
Add user authentication module

This adds JWT-based authentication to secure API endpoints.
Users can now login and receive tokens for authenticated requests.

Bug: 12345
Test: Added unit tests for auth flow
```

---

### 4. Angular规范模板

**适用场景**：Angular项目或遵循严格Conventional Commits的团队

**文件**：`examples/company-templates/angular-style.yaml`

**特点**：
- ✅ 严格的conventional commits
- ✅ Breaking changes追踪
- ✅ Issue引用
- ✅ 无emoji（专业风格）

**生成示例**：
```
feat(parser): add ability to parse arrays

The parser can now handle array syntax in configuration files.
This enables users to define multiple values for a single key.

Closes #456
BREAKING CHANGE: Array syntax changes the configuration format
```

---

## 自定义你自己的规范

### 步骤1：创建配置文件

```bash
gitai config --init
```

### 步骤2：编辑 custom_prompt

打开 `.gitcommit.yaml`，找到 `custom_prompt` 部分：

```yaml
custom_prompt: |
  # 在这里粘贴你公司的commit规范
```

### 步骤3：粘贴公司规范

把你们公司的commit规范文档直接粘贴进去。例如：

#### 示例1：强制包含工单号

```yaml
custom_prompt: |
  提交规范：
  - 格式：<类型>(<模块>): [工单号] <描述>
  - 工单号格式：WO-YYYYMMDD-XXX
  - 必须包含工单号

  示例：
  feat(支付): [WO-20250106-001] 新增支付宝支付

  实现支付宝扫码支付功能，支持订单金额自动计算。

  工单号: WO-20250106-001
  测试人: @测试工程师
```

#### 示例2：强制包含审核人

```yaml
custom_prompt: |
  Commit规范：
  - 所有commit必须包含Code Reviewer
  - 所有commit必须包含QA Tester
  - 格式：Reviewed-by: @username, Tested-by: @username

  示例：
  feat(api): add payment endpoint

  Added new payment processing endpoint for Stripe integration.

  Reviewed-by: @john-smith
  Tested-by: @jane-doe
```

#### 示例3：强制包含影响范围

```yaml
custom_prompt: |
  提交要求：
  - 必须说明"影响范围"（前端/后端/数据库/全部）
  - 必须说明"是否需要发版"（是/否）
  - 必须说明"回滚方案"

  示例：
  feat(订单): 新增订单取消功能

  用户可以在30分钟内取消未支付订单。

  影响范围: 后端API + 前端订单页面
  是否需要发版: 是
  回滚方案: 回滚到上一个稳定版本即可
```

---

## 高级技巧

### 技巧1：多项目不同规范

如果你有多个项目，每个项目用不同规范：

```bash
# 项目A - 使用Jira
cd ~/projects/project-a
cp ~/gitai/examples/company-templates/jira-integration.yaml .gitcommit.yaml

# 项目B - 使用中文规范
cd ~/projects/project-b
cp ~/gitai/examples/company-templates/chinese-enterprise.yaml .gitcommit.yaml

# 项目C - 自定义规范
cd ~/projects/project-c
gitai config --init
vim .gitcommit.yaml  # 粘贴公司规范
```

GitAI会在每个项目目录自动使用该项目的 `.gitcommit.yaml`！

### 技巧2：团队共享配置

把配置文件提交到git仓库，整个团队共享：

```bash
# 1. 创建团队配置
gitai config --init

# 2. 编辑为团队规范
vim .gitcommit.yaml

# 3. 提交到仓库
git add .gitcommit.yaml
git commit -m "chore: add team commit message standards"
git push

# 4. 团队成员拉取后自动使用
git pull  # 其他成员执行
gitai commit  # 自动使用团队规范
```

### 技巧3：使用环境变量区分环境

```yaml
custom_prompt: |
  # 开发环境可以简单一些
  # 生产环境必须详细

  {% if env == "production" %}
  必须包含：
  - 完整测试报告
  - 上线检查清单
  - 回滚预案
  {% else %}
  可以简化格式
  {% endif %}
```

---

## 测试你的配置

配置完成后，测试是否符合预期：

```bash
# 1. 做一些改动
echo "test" > test.txt
git add test.txt

# 2. 使用dry-run模式测试
gitai commit --dry-run

# 3. 检查生成的commit是否符合公司规范
# 如果不符合，调整 custom_prompt 再试
```

---

## 常见问题

### Q: AI会严格遵守我的规范吗？

**A**: 会！AI会严格按照你的 `custom_prompt` 生成commit。我们在prompt中明确指示：
> "IMPORTANT: Follow the above guidelines strictly when generating the commit message."

### Q: 可以用中英文混合吗？

**A**: 可以！但建议指定主要语言：
```yaml
language: "zh"  # 或 "en"
custom_prompt: |
  中文为主，专业术语可以用英文
  例如：feat(API): 新增JWT认证
```

### Q: 能强制包含特定字段吗？

**A**: 完全可以！在 `custom_prompt` 中明确要求：
```yaml
custom_prompt: |
  必须包含以下字段：
  Ticket: XXX-123
  Reviewer: @name
  Testing: description

  如果缺少任何字段，commit无效！
```

AI会生成包含这些字段的commit。

### Q: 公司规范很长怎么办？

**A**: 没问题，`custom_prompt` 支持多行长文本：
```yaml
custom_prompt: |
  [粘贴你们完整的几页规范文档]
  ...
  ...
  [所有要求都粘贴进来]
```

### Q: 可以引用外部文件吗？

**A**: 目前不支持，但你可以复制粘贴。我们建议把规范直接写在配置文件中，这样：
- 版本控制更方便
- 团队共享更简单
- 不依赖外部文件

---

## 真实案例

### 案例1：某金融科技公司

**需求**：
- 必须包含JIRA ticket
- 必须包含安全审核人
- 必须说明是否涉及客户数据

**配置**：
```yaml
custom_prompt: |
  金融科技公司提交规范：

  格式：<type>(<module>): [JIRA-XXX] <description>

  必须包含：
  Security Review: @security-lead
  Customer Data: Yes/No
  Compliance: Checked/Waived

  示例：
  feat(payment): [PAY-789] add encryption for card data

  Implemented AES-256 encryption for credit card storage.

  Security Review: @security-lead
  Customer Data: Yes
  Compliance: Checked - Meets PCI-DSS requirements
```

### 案例2：某互联网大厂

**需求**：
- 中文commit
- 必须关联PRD
- 必须说明灰度方案

**配置**：
```yaml
language: "zh"
custom_prompt: |
  大厂提交规范：

  格式：<类型>(<业务域>): <需求编号> <描述>

  必填项：
  - 关联PRD: PRD-YYYYMMDD-XXX
  - 灰度方案: 描述
  - 监控指标: 列表

  示例：
  feat(推荐系统): PRD-20250106-001 新增个性化推荐算法

  实现基于协同过滤的个性化推荐功能。

  关联PRD: PRD-20250106-001
  灰度方案: 10% -> 30% -> 100%，每阶段观察24小时
  监控指标: CTR、转化率、页面停留时长
```

---

## 总结

使用GitAI的自定义模板功能，你可以：

✅ **完全自动化**公司commit规范
✅ **零学习成本**：直接粘贴公司文档
✅ **团队协作**：配置文件提交到仓库共享
✅ **多项目支持**：每个项目独立配置
✅ **国际化**：支持中英文及其他语言

开始使用：
```bash
# 选择一个模板
cp examples/company-templates/jira-integration.yaml .gitcommit.yaml

# 或者自己写
gitai config --init

# 开始享受自动化commit
gitai commit
```

🎉 从此告别手写commit，AI帮你严格遵守公司规范！
