package i18n

import "fmt"

// LanguageTemplate contains language-specific instructions for AI
type LanguageTemplate struct {
	LanguageInstruction string // How to instruct the AI about the language
	ExampleSubject      string // Example subject line in this language
	ExampleBody         []string // Example body bullet points in this language
}

// GetLanguageTemplate returns language-specific template for AI prompts
func GetLanguageTemplate(langCode, commitType, scope string) LanguageTemplate {
	lang, ok := GetLanguage(langCode)
	if !ok {
		lang = SupportedLanguages["en"] // Fallback to English
	}

	template := LanguageTemplate{}

	switch langCode {
	case "zh":
		template.LanguageInstruction = "Generate the commit message in Chinese (Simplified). Use professional technical Chinese terminology."
		template.ExampleSubject = getChineseExample(commitType, scope)
		template.ExampleBody = []string{
			"实现了用户认证的基本功能",
			"添加了登录和登出接口",
			"包含 JWT 令牌验证中间件",
		}

	case "ja":
		template.LanguageInstruction = "Generate the commit message in Japanese. Use professional technical Japanese terminology."
		template.ExampleSubject = getJapaneseExample(commitType, scope)
		template.ExampleBody = []string{
			"ユーザー認証機能を実装",
			"ログインとログアウトのエンドポイントを追加",
			"JWTトークン検証ミドルウェアを含む",
		}

	case "ko":
		template.LanguageInstruction = "Generate the commit message in Korean. Use professional technical Korean terminology."
		template.ExampleSubject = getKoreanExample(commitType, scope)
		template.ExampleBody = []string{
			"사용자 인증 기능 구현",
			"로그인 및 로그아웃 엔드포인트 추가",
			"JWT 토큰 검증 미들웨어 포함",
		}

	case "de":
		template.LanguageInstruction = "Generate the commit message in German. Use professional technical German terminology."
		template.ExampleSubject = getGermanExample(commitType, scope)
		template.ExampleBody = []string{
			"Benutzerauthentifizierung implementiert",
			"Login- und Logout-Endpunkte hinzugefügt",
			"JWT-Token-Validierungsmiddleware enthalten",
		}

	case "fr":
		template.LanguageInstruction = "Generate the commit message in French. Use professional technical French terminology."
		template.ExampleSubject = getFrenchExample(commitType, scope)
		template.ExampleBody = []string{
			"Implémentation de l'authentification utilisateur",
			"Ajout des endpoints de connexion et déconnexion",
			"Inclusion du middleware de validation JWT",
		}

	case "es":
		template.LanguageInstruction = "Generate the commit message in Spanish. Use professional technical Spanish terminology."
		template.ExampleSubject = getSpanishExample(commitType, scope)
		template.ExampleBody = []string{
			"Implementación de autenticación de usuario",
			"Añadidos endpoints de inicio y cierre de sesión",
			"Incluye middleware de validación JWT",
		}

	case "pt":
		template.LanguageInstruction = "Generate the commit message in Portuguese. Use professional technical Portuguese terminology."
		template.ExampleSubject = getPortugueseExample(commitType, scope)
		template.ExampleBody = []string{
			"Implementação de autenticação de usuário",
			"Adicionados endpoints de login e logout",
			"Inclui middleware de validação JWT",
		}

	case "ru":
		template.LanguageInstruction = "Generate the commit message in Russian. Use professional technical Russian terminology."
		template.ExampleSubject = getRussianExample(commitType, scope)
		template.ExampleBody = []string{
			"Реализована аутентификация пользователя",
			"Добавлены эндпоинты входа и выхода",
			"Включено промежуточное ПО для проверки JWT",
		}

	case "it":
		template.LanguageInstruction = "Generate the commit message in Italian. Use professional technical Italian terminology."
		template.ExampleSubject = getItalianExample(commitType, scope)
		template.ExampleBody = []string{
			"Implementazione dell'autenticazione utente",
			"Aggiunti endpoint di login e logout",
			"Include middleware di validazione JWT",
		}

	default: // English
		template.LanguageInstruction = fmt.Sprintf("Generate the commit message in %s. Use professional technical terminology.", lang.Name)
		template.ExampleSubject = getEnglishExample(commitType, scope)
		template.ExampleBody = []string{
			"Implement JWT-based authentication",
			"Add login and logout endpoints",
			"Include token validation middleware",
		}
	}

	return template
}

// Language-specific example generators
func getEnglishExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): add user authentication endpoint", commitType, scope)
	}
	return fmt.Sprintf("%s: add user authentication endpoint", commitType)
}

func getChineseExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): 添加用户认证接口", commitType, scope)
	}
	return fmt.Sprintf("%s: 添加用户认证接口", commitType)
}

func getJapaneseExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): ユーザー認証エンドポイントを追加", commitType, scope)
	}
	return fmt.Sprintf("%s: ユーザー認証エンドポイントを追加", commitType)
}

func getKoreanExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): 사용자 인증 엔드포인트 추가", commitType, scope)
	}
	return fmt.Sprintf("%s: 사용자 인증 엔드포인트 추가", commitType)
}

func getGermanExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): Benutzerauthentifizierungs-Endpoint hinzufügen", commitType, scope)
	}
	return fmt.Sprintf("%s: Benutzerauthentifizierungs-Endpoint hinzufügen", commitType)
}

func getFrenchExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): ajouter l'endpoint d'authentification utilisateur", commitType, scope)
	}
	return fmt.Sprintf("%s: ajouter l'endpoint d'authentification utilisateur", commitType)
}

func getSpanishExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): añadir endpoint de autenticación de usuario", commitType, scope)
	}
	return fmt.Sprintf("%s: añadir endpoint de autenticación de usuario", commitType)
}

func getPortugueseExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): adicionar endpoint de autenticação de usuário", commitType, scope)
	}
	return fmt.Sprintf("%s: adicionar endpoint de autenticação de usuário", commitType)
}

func getRussianExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): добавить эндпоинт аутентификации пользователя", commitType, scope)
	}
	return fmt.Sprintf("%s: добавить эндпоинт аутентификации пользователя", commitType)
}

func getItalianExample(commitType, scope string) string {
	if scope != "" {
		return fmt.Sprintf("%s(%s): aggiungere endpoint di autenticazione utente", commitType, scope)
	}
	return fmt.Sprintf("%s: aggiungere endpoint di autenticazione utente", commitType)
}

// GetMultilingualInstructions generates instructions for multilingual commit messages
func GetMultilingualInstructions(languages []string) string {
	if len(languages) == 0 {
		return ""
	}

	if len(languages) == 1 {
		lang, ok := GetLanguage(languages[0])
		if !ok {
			return ""
		}
		return fmt.Sprintf("Generate the commit message in %s.", lang.Name)
	}

	// Multiple languages
	langNames := make([]string, len(languages))
	for i, code := range languages {
		if lang, ok := GetLanguage(code); ok {
			langNames[i] = lang.Name
		} else {
			langNames[i] = code
		}
	}

	instruction := "Generate a MULTILINGUAL commit message with translations in the following languages:\n"
	for i, name := range langNames {
		instruction += fmt.Sprintf("%d. %s\n", i+1, name)
	}

	instruction += "\nFormat:\n"
	instruction += "Subject line in primary language (" + langNames[0] + ")\n\n"
	instruction += "Body:\n"
	instruction += "- Bullet points explaining changes in " + langNames[0] + "\n\n"

	if len(languages) > 1 {
		instruction += "Translations:\n"
		for i := 1; i < len(langNames); i++ {
			instruction += fmt.Sprintf("- [%s] Translation of the subject line\n", langNames[i])
		}
	}

	return instruction
}
