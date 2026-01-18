package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetLanguage(t *testing.T) {
	tests := []struct {
		code    string
		wantOk  bool
		wantCode string
	}{
		{"en", true, "en"},
		{"zh", true, "zh"},
		{"ja", true, "ja"},
		{"ko", true, "ko"},
		{"de", true, "de"},
		{"fr", true, "fr"},
		{"invalid", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			lang, ok := GetLanguage(tt.code)
			if ok != tt.wantOk {
				t.Errorf("GetLanguage(%q) ok = %v, want %v", tt.code, ok, tt.wantOk)
			}
			if ok && lang.Code != tt.wantCode {
				t.Errorf("GetLanguage(%q).Code = %v, want %v", tt.code, lang.Code, tt.wantCode)
			}
		})
	}
}

func TestIsSupported(t *testing.T) {
	tests := []struct {
		code string
		want bool
	}{
		{"en", true},
		{"zh", true},
		{"ja", true},
		{"ko", true},
		{"de", true},
		{"fr", true},
		{"es", true},
		{"pt", true},
		{"ru", true},
		{"it", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if got := IsSupported(tt.code); got != tt.want {
				t.Errorf("IsSupported(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestNormalizeLanguageCode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"en", "en"},
		{"EN", "en"},
		{"zh-cn", "zh"},
		{"zh-hans", "zh"},
		{"chinese", "zh"},
		{"ja", "ja"},
		{"japanese", "ja"},
		{"jp", "ja"},
		{"ko", "ko"},
		{"korean", "ko"},
		{"kr", "ko"},
		{"de", "de"},
		{"german", "de"},
		{"fr", "fr"},
		{"french", "fr"},
		{"es", "es"},
		{"spanish", "es"},
		{"pt", "pt"},
		{"pt-br", "pt"},
		{"ru", "ru"},
		{"russian", "ru"},
		{"it", "it"},
		{"italian", "it"},
		{"en-us", "en"},
		{"en-gb", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := NormalizeLanguageCode(tt.input); got != tt.want {
				t.Errorf("NormalizeLanguageCode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDetectLanguageInText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		wantLang string // Expected language with highest score
	}{
		{
			name:     "English text",
			text:     "The quick brown fox jumps over the lazy dog. This is a test.",
			wantLang: "en",
		},
		{
			name:     "Chinese text",
			text:     "这是一个测试文本。我们正在测试中文检测功能。",
			wantLang: "zh",
		},
		{
			name:     "Japanese text",
			text:     "これはテストです。日本語の検出をテストしています。",
			wantLang: "ja",
		},
		{
			name:     "Korean text",
			text:     "이것은 테스트입니다. 한국어 감지를 테스트하고 있습니다.",
			wantLang: "ko",
		},
		{
			name:     "German text",
			text:     "Das ist ein Test. Wir testen die deutsche Spracherkennung mit ä ö ü ß.",
			wantLang: "de",
		},
		{
			name:     "French text",
			text:     "C'est un test. Nous testons la détection de la langue française avec é è ê.",
			wantLang: "fr",
		},
		{
			name:     "Spanish text",
			text:     "Esta es una prueba. Estamos probando la detección del idioma español con ñ.",
			wantLang: "es",
		},
		{
			name:     "Portuguese text",
			text:     "Este é um teste. Estamos testando a detecção da língua portuguesa com ã õ.",
			wantLang: "pt",
		},
		{
			name:     "Russian text",
			text:     "Это тест. Мы тестируем обнаружение русского языка.",
			wantLang: "ru",
		},
		{
			name:     "Italian text",
			text:     "Questo è un test. Stiamo testando il rilevamento della lingua italiana.",
			wantLang: "it",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scores := detectLanguageInText(tt.text)

			// Find language with highest score
			maxScore := 0
			detectedLang := ""
			for lang, score := range scores {
				if score > maxScore {
					maxScore = score
					detectedLang = lang
				}
			}

			if detectedLang != tt.wantLang {
				t.Errorf("detectLanguageInText() detected %q, want %q (scores: %v)", detectedLang, tt.wantLang, scores)
			}
		})
	}
}

func TestDetectProjectLanguage(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		readmeContent  string
		commits        []string
		wantLang       string
	}{
		{
			name:          "English project",
			readmeContent: "# Test Project\n\nThis is a test project for testing language detection.",
			commits: []string{
				"feat: add new feature",
				"fix: resolve bug in authentication",
				"docs: update README",
			},
			wantLang: "en",
		},
		{
			name:          "Chinese project",
			readmeContent: "# 测试项目\n\n这是一个用于测试语言检测的测试项目。",
			commits: []string{
				"feat: 添加新功能",
				"fix: 修复认证中的错误",
				"docs: 更新文档",
			},
			wantLang: "zh",
		},
		{
			name:          "Japanese project",
			readmeContent: "# テストプロジェクト\n\nこれは言語検出をテストするためのテストプロジェクトです。",
			commits: []string{
				"feat: 新機能を追加",
				"fix: 認証のバグを修正",
				"docs: ドキュメントを更新",
			},
			wantLang: "ja",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create README file
			readmePath := filepath.Join(tmpDir, "README.md")
			err := os.WriteFile(readmePath, []byte(tt.readmeContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write README: %v", err)
			}

			// Detect language
			detected := DetectProjectLanguage(readmePath, tt.commits)

			if detected != tt.wantLang {
				t.Errorf("DetectProjectLanguage() = %q, want %q", detected, tt.wantLang)
			}

			// Cleanup
			os.Remove(readmePath)
		})
	}
}

func TestAnalyzeCommitLanguages(t *testing.T) {
	tests := []struct {
		name     string
		commits  []string
		wantLang string
	}{
		{
			name: "English commits",
			commits: []string{
				"feat: add user authentication",
				"fix: resolve login bug",
				"docs: update API documentation",
			},
			wantLang: "en",
		},
		{
			name: "Chinese commits",
			commits: []string{
				"feat: 添加用户认证",
				"fix: 修复登录错误",
				"docs: 更新API文档",
			},
			wantLang: "zh",
		},
		{
			name: "Mixed commits (English dominant)",
			commits: []string{
				"feat: add user authentication",
				"feat: 添加功能",
				"fix: resolve bug in authentication system",
				"docs: update documentation for API endpoints",
				"test: add unit tests for authentication",
			},
			wantLang: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scores := analyzeCommitLanguages(tt.commits)

			// Find language with highest score
			maxScore := 0
			detectedLang := ""
			for lang, score := range scores {
				if score > maxScore {
					maxScore = score
					detectedLang = lang
				}
			}

			if detectedLang != tt.wantLang {
				t.Errorf("analyzeCommitLanguages() detected %q, want %q (scores: %v)", detectedLang, tt.wantLang, scores)
			}
		})
	}
}
