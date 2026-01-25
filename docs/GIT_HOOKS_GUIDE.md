# 🪝 GitAI Git Hooks - 自动化集成指南

## 概述

GitAI Git Hooks 让你的提交流程完全自动化。安装后，你只需运行标准的 `git commit` 命令，GitAI 就会自动生成提交消息 —— 无需手动运行 `gitai commit`。

## 快速开始

### 一键安装

```bash
# 在你的 git 仓库中运行
gitai hooks install
```

就这么简单！从此以后：

```bash
# 传统方式（仍然可用）
git add .
gitai commit

# 新方式（自动化！）
git add .
git commit          # GitAI 自动生成提交消息
```

## 工作原理

GitAI 使用 Git 的 `prepare-commit-msg` hook，在你运行 `git commit` 时自动介入：

```
┌──────────────┐
│ git add .    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ git commit   │
└──────┬───────┘
       │
       ▼
┌──────────────────────┐
│ prepare-commit-msg   │  ← GitAI hook 在这里运行
│ (GitAI 自动生成消息)  │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ 打开编辑器           │  ← 显示生成的消息，可编辑
│ (预填充 AI 消息)      │
└──────┬───────────────┘
       │
       ▼
┌──────────────┐
│ 提交完成      │
└──────────────┘
```

## 命令详解

### 安装 Hooks

```bash
# 基本安装（只安装 prepare-commit-msg）
gitai hooks install

# 强制安装（覆盖已有的 hooks）
gitai hooks install --force

# 安装所有可用的 hooks
gitai hooks install --all
```

**安装后会发生什么：**
1. 在 `.git/hooks/` 目录创建 hook 脚本
2. 如果已有同名 hook，会自动备份为 `.gitai-backup`
3. 新的 hook 会调用 `gitai generate --quiet` 生成消息

### 查看状态

```bash
gitai hooks status
```

**输出示例：**
```
📊 GitAI Hooks Status
============================================================

🔗 prepare-commit-msg
   Status:  ✅ Installed (GitAI)
   Backup:  ✅ Available

🔗 commit-msg
   Status:  ❌ Not installed

🔗 pre-commit
   Status:  ⚠️  Installed (Not GitAI)

✨ GitAI hooks are active!

💡 Usage:
  - Just run 'git commit' and GitAI will generate messages
  - Use 'git commit -m "message"' to skip automation
  - Use 'git commit --no-verify' to bypass all hooks
```

### 卸载 Hooks

```bash
# 基本卸载
gitai hooks uninstall

# 卸载并恢复之前的 hooks
gitai hooks uninstall --restore

# 卸载所有 GitAI hooks
gitai hooks uninstall --all
```

## Hook 类型

GitAI 支持三种类型的 Git Hooks：

### 1. prepare-commit-msg（推荐）

**用途**: 自动生成提交消息

**何时运行**: 在提交编辑器打开前

**行为**:
- 自动调用 `gitai generate --quiet`
- 将生成的消息预填充到编辑器
- 你可以在提交前编辑消息

**何时跳过**:
- 使用 `git commit -m "message"` 时
- 合并提交（merge commits）
- 变基/樱桃拣选操作
- 使用 `--no-verify` 标志时

### 2. commit-msg（可选）

**用途**: 验证提交消息格式

**何时运行**: 提交消息写入后，实际提交前

**行为**:
- 检查消息是否符合 Conventional Commits 格式
- 如果格式不正确，显示警告（不阻止提交）

### 3. pre-commit（可选）

**用途**: 提交前检查

**何时运行**: 实际提交前，最早执行

**行为**:
- 检查是否有暂存的更改
- 可扩展以添加 linting、formatting、测试等

## 实际使用场景

### 场景 1: 日常开发工作流

```bash
# 1. 修改代码
vim src/auth.go

# 2. 暂存更改
git add src/auth.go

# 3. 提交（GitAI 自动生成消息）
git commit

# 编辑器打开，显示：
# feat(auth): implement JWT token validation
#
# - Add JWT token validation middleware
# - Support RS256 and HS256 algorithms
# - Add token expiration checking
# - Include comprehensive error handling

# 4. 保存并关闭编辑器 → 提交完成！
```

### 场景 2: 快速提交（跳过 GitAI）

```bash
# 直接指定消息（GitAI 不会运行）
git commit -m "fix: typo in README"
```

### 场景 3: 紧急提交（绕过所有 hooks）

```bash
# 完全跳过所有 hooks
git commit --no-verify
```

### 场景 4: 修改 AI 生成的消息

```bash
git commit

# 编辑器打开后：
# 1. 查看 GitAI 生成的消息
# 2. 根据需要修改
# 3. 保存提交
```

## 高级配置

### 自定义 Hook 行为

如果你想自定义 hook 行为，可以直接编辑 `.git/hooks/prepare-commit-msg`：

```bash
#!/bin/sh
# GitAI - Auto-generated Git Hook

COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2

# 添加自定义逻辑...

# 例如：只在特定分支使用 GitAI
BRANCH=$(git branch --show-current)
if [ "$BRANCH" = "main" ] || [ "$BRANCH" = "develop" ]; then
    # 在 main/develop 分支使用 GitAI
    if GENERATED_MSG=$(gitai generate --quiet --type feat 2>&1); then
        echo "$GENERATED_MSG" > "$COMMIT_MSG_FILE"
    fi
fi

exit 0
```

### 与现有 Hooks 集成

如果你已经有自己的 hooks：

1. **方法 1：安装时自动备份**
   ```bash
   gitai hooks install  # 自动备份到 .gitai-backup
   ```

2. **方法 2：手动合并**
   ```bash
   # 查看你的原始 hook
   cat .git/hooks/prepare-commit-msg.gitai-backup

   # 编辑 GitAI hook，合并原始逻辑
   vim .git/hooks/prepare-commit-msg
   ```

3. **方法 3：链式调用**
   ```bash
   #!/bin/sh
   # prepare-commit-msg

   # 运行原始 hook
   .git/hooks/prepare-commit-msg.gitai-backup "$@"

   # 然后运行 GitAI
   if GENERATED_MSG=$(gitai generate --quiet 2>&1); then
       echo "$GENERATED_MSG" > "$1"
   fi
   ```

## 团队协作

### 推荐设置

为了让整个团队受益，建议：

1. **在项目 README 中说明**:
   ```markdown
   ## 开发设置

   我们使用 GitAI 自动生成提交消息：

   ```bash
   # 安装 GitAI
   brew install xyue92/tap/gitai

   # 安装 hooks
   gitai hooks install
   ```

2. **添加到入职文档**

3. **可选：添加到 .gitattributes**（确保一致性）

### 不强制使用

Git hooks 是本地的，不会提交到仓库。这意味着：
- ✅ 团队成员可以选择是否使用
- ✅ 不会影响未安装 GitAI 的人
- ✅ 灵活且无侵入性

## CI/CD 集成

虽然 hooks 是本地的，但你可以在 CI 中验证提交消息：

```yaml
# .github/workflows/commit-lint.yml
name: Commit Message Validation

on: [pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Install GitAI
        run: |
          curl -sSL https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.sh | bash

      - name: Validate Commits
        run: |
          # 验证 PR 中的所有提交
          for commit in $(git rev-list origin/main..HEAD); do
            msg=$(git log -1 --pretty=%B $commit)
            # 检查格式...
          done
```

## 故障排查

### 问题：Hook 不工作

**检查列表：**

1. **Hook 已安装？**
   ```bash
   gitai hooks status
   ```

2. **GitAI 在 PATH 中？**
   ```bash
   which gitai
   # 应该显示路径，如 /usr/local/bin/gitai
   ```

3. **Hook 有执行权限？**
   ```bash
   ls -la .git/hooks/prepare-commit-msg
   # 应该是 -rwxr-xr-x (可执行)
   ```

4. **手动测试**
   ```bash
   gitai generate --quiet
   # 应该输出生成的消息
   ```

### 问题：每次都提示输入类型

**原因**: 在 quiet 模式下，没有指定类型时会使用默认值 "feat"

**解决方案**:
```bash
# 在 .gitcommit.yaml 中设置默认类型
# 或修改 hook 脚本指定类型：
gitai generate --quiet --type feat
```

### 问题：Ollama 连接失败

**症状**: Hook 运行时报错 "failed to connect to Ollama"

**解决方案**:
```bash
# 确保 Ollama 正在运行
ollama serve

# 检查 Ollama 状态
curl http://localhost:11434/api/version
```

### 问题：想临时禁用 hook

```bash
# 方法 1：使用 --no-verify
git commit --no-verify

# 方法 2：直接指定消息
git commit -m "your message"

# 方法 3：临时卸载
gitai hooks uninstall
# ... 做你的事情
gitai hooks install
```

### 问题：Hook 与其他工具冲突

如果你使用其他工具（如 Husky、pre-commit framework）：

**方案 1：手动合并**
```bash
# 编辑你的 hook，添加 GitAI 调用
vim .git/hooks/prepare-commit-msg
```

**方案 2：使用 GitAI 但不用 hook**
```bash
# 卸载 hook
gitai hooks uninstall

# 继续手动使用
gitai commit
```

## 最佳实践

### ✅ 推荐做法

1. **在个人项目中先尝试**
   - 熟悉工作流程
   - 了解如何调整生成的消息

2. **逐步推广到团队**
   - 先给愿意尝试的成员
   - 收集反馈并改进配置

3. **配置团队特定模板**
   ```yaml
   # .gitcommit.yaml
   custom_prompt: |
     Our team standards:
     - Always include ticket number
     - Use imperative mood
     - Max 50 chars subject
   ```

4. **保留编辑的自由**
   - Hook 只是生成草稿
   - 鼓励团队审查和修改

5. **定期更新 GitAI**
   ```bash
   gitai update
   ```

### ❌ 避免的做法

1. ��不验证就直接提交**
   - 始终审查 AI 生成的消息

2. **在所有场景都强制使用**
   - 给团队选择权
   - 允许 `--no-verify`

3. **忽略配置调优**
   - 根据项目调整 templates
   - 利用 custom_prompt

## 性能优化

### 减少延迟

```yaml
# .gitcommit.yaml
ai_optimization:
  # 使用更小的模型
  model: "qwen2.5-coder:3b"  # 更快

  # 减少上下文
  max_diff_length: 1000

  # 简化提交（无 body）
  detailed_commit: false
```

### 缓存策略

未来版本可能支持：
```yaml
cache:
  enabled: true
  duration: 3600  # 1 hour
```

## 卸载与清理

### 完全移除

```bash
# 1. 卸载 hooks
gitai hooks uninstall --restore

# 2. 删除配置（可选）
rm .gitcommit.yaml

# 3. 卸载 GitAI（可选）
brew uninstall gitai
# 或
rm /usr/local/bin/gitai
```

### 只保留手动模式

```bash
# 卸载 hooks，但保留 GitAI
gitai hooks uninstall

# 继续使用 gitai commit
```

## 常见问题

**Q: Hooks 会被提交到仓库吗？**
A: 不会。`.git/hooks/` 目录不会被 Git 跟踪，hooks 只存在于本地。

**Q: 团队成员必须安装吗？**
A: 不必须。Hooks 是可选的，未安装的成员仍可正常工作。

**Q: 如何在 Windows 上使用？**
A: Windows Git Bash 完全支持 Git Hooks，安装方式相同。

**Q: 可以自定义生成逻辑吗？**
A: 可以。通过编辑 `.gitcommit.yaml` 或直接修改 hook 脚本。

**Q: 性能影响如何？**
A: 通常增加 1-3 秒延迟（取决于模型大小）。可以通过使用更小的模型优化。

**Q: 支持 Mono-repo 吗？**
A: 支持。每个子项目可以有自己的 `.gitcommit.yaml` 配置。

## 总结

GitAI Git Hooks 提供：

✅ **零摩擦自动化** - 无需改变工作流程
✅ **完全可选** - 团队成员可自由选择
✅ **智能生成** - AI 驱动的高质量消息
✅ **灵活配置** - 适应各种团队需求
✅ **安全备份** - 保护现有 hooks
✅ **易于管理** - 简单的安装/卸载

开始使用：
```bash
gitai hooks install
git add .
git commit  # 就这么简单！
```

---

**相关文档：**
- [GitAI 主文档](../README.md)
- [配置指南](../CUSTOM_TEMPLATES_GUIDE.md)
- [提交统计](STATS_GUIDE.md)
