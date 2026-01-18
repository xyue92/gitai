# Multilingual Support Guide

GitAI supports generating commit messages in multiple languages, making it ideal for international teams and projects with diverse contributors.

## Supported Languages

GitAI currently supports the following languages:

| Code | Language | Native Name |
|------|----------|-------------|
| `en` | English | English |
| `zh` | Chinese (Simplified) | 中文 |
| `ja` | Japanese | 日本語 |
| `ko` | Korean | 한국어 |
| `de` | German | Deutsch |
| `fr` | French | Français |
| `es` | Spanish | Español |
| `pt` | Portuguese | Português |
| `ru` | Russian | Русский |
| `it` | Italian | Italiano |

## Usage Modes

### 1. Single Language Mode (Default)

Generate commit messages in a single language.

**Configuration:**

```yaml
# .gitcommit.yaml
language: "en"  # or zh, ja, ko, de, fr, es, pt, ru, it
```

**Example output:**

```
feat(auth): add user authentication endpoint

- Implement JWT-based authentication
- Add login and logout endpoints
- Include token validation middleware
```

### 2. Multilingual Mode

Generate commit messages with translations in multiple languages. Perfect for international teams where developers speak different languages.

**Configuration:**

```yaml
# .gitcommit.yaml
languages:
  - "en"  # Primary language (used for subject line)
  - "zh"  # Chinese translation
  - "ja"  # Japanese translation
```

**Example output:**

```
feat(auth): add user authentication endpoint

- Implement JWT-based authentication
- Add login and logout endpoints
- Include token validation middleware

Translations:
- [中文] feat(auth): 添加用户认证接口
- [日本語] feat(auth): ユーザー認証エンドポイントを追加
```

### 3. Auto-Detection Mode

Let GitAI automatically detect the appropriate language based on your project's README, commit history, and documentation.

**Configuration:**

```yaml
# .gitcommit.yaml
auto_detect_language: true
```

**How it works:**

1. Analyzes README files for language patterns
2. Examines recent commit messages
3. Checks documentation directories (e.g., `docs/zh`, `docs/ja`)
4. Selects the language with the highest confidence score

## Language-Specific Examples

### English (en)

```yaml
language: "en"
```

Output:
```
feat: add user authentication endpoint

- Implement JWT-based authentication
- Add login and logout endpoints
```

### Chinese (zh)

```yaml
language: "zh"
```

Output:
```
feat: 添加用户认证接口

- 实现基于JWT的认证功能
- 添加登录和登出接口
```

### Japanese (ja)

```yaml
language: "ja"
```

Output:
```
feat: ユーザー認証エンドポイントを追加

- JWT ベースの認証を実装
- ログインとログアウトのエンドポイントを追加
```

### Korean (ko)

```yaml
language: "ko"
```

Output:
```
feat: 사용자 인증 엔드포인트 추가

- JWT 기반 인증 구현
- 로그인 및 로그아웃 엔드포인트 추가
```

### German (de)

```yaml
language: "de"
```

Output:
```
feat: Benutzerauthentifizierungs-Endpoint hinzufügen

- JWT-basierte Authentifizierung implementiert
- Login- und Logout-Endpunkte hinzugefügt
```

### French (fr)

```yaml
language: "fr"
```

Output:
```
feat: ajouter l'endpoint d'authentification utilisateur

- Implémentation de l'authentification JWT
- Ajout des endpoints de connexion et déconnexion
```

### Spanish (es)

```yaml
language: "es"
```

Output:
```
feat: añadir endpoint de autenticación de usuario

- Implementación de autenticación JWT
- Añadidos endpoints de inicio y cierre de sesión
```

### Portuguese (pt)

```yaml
language: "pt"
```

Output:
```
feat: adicionar endpoint de autenticação de usuário

- Implementação de autenticação JWT
- Adicionados endpoints de login e logout
```

### Russian (ru)

```yaml
language: "ru"
```

Output:
```
feat: добавить эндпоинт аутентификации пользователя

- Реализована JWT-аутентификация
- Добавлены эндпоинты входа и выхода
```

### Italian (it)

```yaml
language: "it"
```

Output:
```
feat: aggiungere endpoint di autenticazione utente

- Implementazione dell'autenticazione JWT
- Aggiunti endpoint di login e logout
```

## Configuration Examples

### Example 1: International Team (English + Chinese)

```yaml
# .gitcommit.yaml
languages:
  - "en"
  - "zh"
detailed_commit: true
```

### Example 2: European Project (English + German + French)

```yaml
# .gitcommit.yaml
languages:
  - "en"
  - "de"
  - "fr"
detailed_commit: true
```

### Example 3: Auto-detect with Fallback

```yaml
# .gitcommit.yaml
auto_detect_language: true
language: "en"  # Fallback if detection fails
```

### Example 4: Japanese Open Source Project

```yaml
# .gitcommit.yaml
languages:
  - "ja"  # Primary language
  - "en"  # For international contributors
detailed_commit: true
```

## Best Practices

### 1. Choose Primary Language Carefully

The first language in the `languages` array is used for the subject line. Choose the language that most team members understand.

### 2. Limit Translation Languages

For readability, limit translations to 2-3 languages maximum:

```yaml
# Good: 2-3 languages
languages:
  - "en"
  - "zh"
  - "ja"

# Avoid: Too many languages
languages:
  - "en"
  - "zh"
  - "ja"
  - "ko"
  - "de"
  - "fr"  # Commit message becomes too long
```

### 3. Use Auto-detection for Mixed Projects

If your project has contributions in multiple languages, use auto-detection:

```yaml
auto_detect_language: true
```

### 4. Consistent Team Configuration

Use a project-level `.gitcommit.yaml` file so all team members use the same language settings:

```bash
# Project root
.gitcommit.yaml  # Team configuration
```

### 5. Combine with Ticket Integration

Multilingual commits work great with ticket numbers:

```yaml
languages:
  - "en"
  - "zh"
require_ticket: true
ticket_prefix: "JIRA"
```

Output:
```
feat(auth): [JIRA-123] add user authentication

- Implement JWT-based authentication

Translations:
- [中文] feat(auth): [JIRA-123] 添加用户认证
```

## Command-Line Overrides

You can override language settings via command-line flags (if implemented):

```bash
# Use specific language for this commit
gitai commit --language ja

# Force multilingual mode
gitai commit --languages en,zh,ja

# Disable auto-detection
gitai commit --no-auto-detect
```

## Technical Details

### Language Detection Algorithm

GitAI uses a scoring system to detect project language:

1. **README Analysis (Weight: 3x)**
   - Detects character sets (CJK, Cyrillic, Latin)
   - Identifies common words and phrases
   - Analyzes language-specific punctuation

2. **Commit History (Weight: 2x)**
   - Analyzes last 5 commit messages
   - Detects language patterns

3. **Documentation Structure (Weight: 1x)**
   - Checks for language-specific directories
   - Detects locale patterns (e.g., `docs/zh-CN`)

### Language Normalization

GitAI normalizes various language code formats:

- `zh-CN`, `zh-Hans`, `Chinese` → `zh`
- `ja`, `jp`, `Japanese` → `ja`
- `en-US`, `en-GB`, `English` → `en`

### AI Model Compatibility

All supported Ollama models work with multilingual prompts:

- **Recommended for Multilingual:**
  - `qwen2.5-coder:7b` - Excellent Chinese/English support
  - `mistral:7b` - Good European language support
  - `codellama:7b` - Decent multilingual capabilities

- **Best Single-Language:**
  - English: All models
  - Chinese: `qwen2.5-coder:7b`
  - Japanese: `qwen2.5-coder:7b`, `mistral:7b`
  - European: `mistral:7b`

## Troubleshooting

### Issue: AI generates wrong language

**Solution:** Explicitly set the language in config:

```yaml
language: "zh"  # Force Chinese
# Remove or comment out auto_detect_language
```

### Issue: Translations are low quality

**Solution:** Use a better model or reduce number of languages:

```yaml
model: "qwen2.5-coder:7b"  # Better multilingual support
languages:
  - "en"
  - "zh"  # Limit to 2 languages
```

### Issue: Auto-detection picks wrong language

**Solution:** Add more context to your project:

- Add a detailed README in your preferred language
- Use consistent language in commit messages
- Create language-specific documentation directories

### Issue: Commit messages too long with translations

**Solution:** Use concise mode or reduce translations:

```yaml
detailed_commit: false  # Concise mode
languages:
  - "en"  # Primary only
  - "zh"  # Just one translation
```

## Contributing New Languages

To add support for a new language:

1. Update `internal/i18n/language.go`:
   - Add language to `SupportedLanguages` map
   - Add detection patterns to `detectLanguageInText()`

2. Update `internal/i18n/templates.go`:
   - Add language-specific examples
   - Implement example generator function

3. Add tests in `internal/i18n/language_test.go`

4. Update documentation

See the codebase for implementation details.

## Examples in the Wild

### Open Source Project (English + Chinese)

```yaml
# .gitcommit.yaml
model: "qwen2.5-coder:7b"
languages:
  - "en"
  - "zh"
detailed_commit: true
```

### Corporate Project (English + Japanese + Korean)

```yaml
# .gitcommit.yaml
model: "qwen2.5-coder:7b"
languages:
  - "en"
  - "ja"
  - "ko"
require_ticket: true
ticket_prefix: "PROJ"
detailed_commit: true
```

### European Startup (English + German + French)

```yaml
# .gitcommit.yaml
model: "mistral:7b"
languages:
  - "en"
  - "de"
  - "fr"
detailed_commit: true
subject_length: "normal"
```

---

For more information, see:
- [Main README](README.md)
- [Configuration Guide](.gitcommit.example.yaml)
- [Ticket Integration Guide](TICKET_INTEGRATION_GUIDE.md)
