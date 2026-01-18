# 🚀 多语言功能快速开始

## 5 分钟上手多语言 GitAI

### 步骤 1: 选择你的使用模式

#### 选项 A: 单一语言（最简单）

在 `.gitcommit.yaml` 中设置：

```yaml
language: "zh"  # 或 en, ja, ko, de, fr, es, pt, ru, it
```

#### 选项 B: 自动检测（智能推荐）

```yaml
auto_detect_language: true
```

GitAI 会自动分析你的项目并选择最合适的语言。

#### 选项 C: 多语言（团队协作）

```yaml
languages:
  - "en"  # 主语言（用于主题行）
  - "zh"  # 翻译
```

### 步骤 2: 使用 GitAI

正常使用 `gitai commit`，提交消息会自动以你配置的语言生成！

```bash
git add .
gitai commit
```

### 步骤 3: 查看结果

**单语言模式 (中文):**
```
feat: 添加用户认证接口

- 实现JWT认证
- 添加登录登出功能
```

**多语言模式 (英文+中文):**
```
feat: add user authentication

- Implement JWT authentication
- Add login/logout functionality

Translations:
- [中文] feat: 添加用户认证接口
```

## 支持的语言

| 代码 | 语言 | 示例 |
|-----|------|-----|
| `en` | English | feat: add new feature |
| `zh` | 中文 | feat: 添加新功能 |
| `ja` | 日本語 | feat: 新機能を追加 |
| `ko` | 한국어 | feat: 새로운 기능 추가 |
| `de` | Deutsch | feat: neue Funktion hinzufügen |
| `fr` | Français | feat: ajouter nouvelle fonctionnalité |
| `es` | Español | feat: añadir nueva función |
| `pt` | Português | feat: adicionar nova função |
| `ru` | Русский | feat: добавить новую функцию |
| `it` | Italiano | feat: aggiungere nuova funzione |

## 常见问题

### Q: 如何更改语言？

编辑 `.gitcommit.yaml`：

```yaml
language: "ja"  # 改为日语
```

### Q: 多语言模式会影响性能吗？

影响很小。2-3 种语言只增加 10-20% 的 token 使用。

### Q: 自动检测准确吗？

非常准确！基于：
- README 文件内容
- 最近的提交历史
- 文档目录结构

### Q: 可以混合使用吗？

可以！例如：

```yaml
auto_detect_language: true  # 自动检测主语言
languages:
  - "en"  # 但总是包含英语翻译
```

## 高级配置

### 国际团队配置

```yaml
model: "qwen2.5-coder:7b"
languages:
  - "en"  # 英语（国际通用）
  - "zh"  # 中文团队
  - "ja"  # 日本团队
detailed_commit: true
```

### 自动检测 + 备选

```yaml
auto_detect_language: true
language: "en"  # 检测失败时使用英语
detailed_commit: true
```

## 更多信息

- 📖 [完整多语言指南](../MULTILINGUAL_GUIDE.md)
- ⚙️ [配置示例](../.gitcommit.example.yaml)
- 🔧 [实现细节](../MULTILINGUAL_FEATURES.md)

---

开始使用多语言 GitAI，让你的团队沟通更流畅！🌍
