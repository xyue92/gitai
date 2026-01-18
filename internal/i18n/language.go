package i18n

import (
	"os"
	"path/filepath"
	"strings"
)

// Language represents a supported language
type Language struct {
	Code        string // ISO 639-1 code (en, zh, ja, ko, de, fr, etc.)
	Name        string // English name
	NativeName  string // Native name
	AIPromptKey string // Key to use in AI prompts
}

// SupportedLanguages defines all languages supported by GitAI
var SupportedLanguages = map[string]Language{
	"en": {Code: "en", Name: "English", NativeName: "English", AIPromptKey: "English"},
	"zh": {Code: "zh", Name: "Chinese", NativeName: "中文", AIPromptKey: "Chinese (Simplified)"},
	"ja": {Code: "ja", Name: "Japanese", NativeName: "日本語", AIPromptKey: "Japanese"},
	"ko": {Code: "ko", Name: "Korean", NativeName: "한국어", AIPromptKey: "Korean"},
	"de": {Code: "de", Name: "German", NativeName: "Deutsch", AIPromptKey: "German"},
	"fr": {Code: "fr", Name: "French", NativeName: "Français", AIPromptKey: "French"},
	"es": {Code: "es", Name: "Spanish", NativeName: "Español", AIPromptKey: "Spanish"},
	"pt": {Code: "pt", Name: "Portuguese", NativeName: "Português", AIPromptKey: "Portuguese"},
	"ru": {Code: "ru", Name: "Russian", NativeName: "Русский", AIPromptKey: "Russian"},
	"it": {Code: "it", Name: "Italian", NativeName: "Italiano", AIPromptKey: "Italian"},
}

// GetLanguage returns the language metadata for a language code
func GetLanguage(code string) (Language, bool) {
	lang, ok := SupportedLanguages[code]
	return lang, ok
}

// IsSupported checks if a language code is supported
func IsSupported(code string) bool {
	_, ok := SupportedLanguages[code]
	return ok
}

// DetectProjectLanguage attempts to detect the primary language used in the project
// It analyzes README files, documentation, and recent commit messages
func DetectProjectLanguage(readmePath string, recentCommits []string) string {
	scores := make(map[string]int)

	// Analyze README files
	readmeScore := analyzeReadmeLanguage(readmePath)
	for lang, score := range readmeScore {
		scores[lang] += score * 3 // README has higher weight
	}

	// Analyze recent commit messages
	commitScore := analyzeCommitLanguages(recentCommits)
	for lang, score := range commitScore {
		scores[lang] += score * 2 // Commits have medium weight
	}

	// Analyze documentation files
	docScore := analyzeDocumentationLanguage()
	for lang, score := range docScore {
		scores[lang] += score
	}

	// Find the language with highest score
	maxScore := 0
	detectedLang := "en" // Default to English
	for lang, score := range scores {
		if score > maxScore {
			maxScore = score
			detectedLang = lang
		}
	}

	return detectedLang
}

// analyzeReadmeLanguage analyzes README content to detect language
func analyzeReadmeLanguage(readmePath string) map[string]int {
	scores := make(map[string]int)

	// Try common README file names
	readmeFiles := []string{"README.md", "README.MD", "readme.md", "Readme.md"}
	if readmePath != "" {
		readmeFiles = append([]string{readmePath}, readmeFiles...)
	}

	for _, filename := range readmeFiles {
		content, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		text := string(content)
		langScores := detectLanguageInText(text)
		for lang, score := range langScores {
			scores[lang] += score
		}
		break // Only analyze the first README found
	}

	return scores
}

// analyzeCommitLanguages analyzes recent commits to detect language patterns
func analyzeCommitLanguages(commits []string) map[string]int {
	scores := make(map[string]int)

	for _, commit := range commits {
		langScores := detectLanguageInText(commit)
		for lang, score := range langScores {
			scores[lang] += score
		}
	}

	return scores
}

// analyzeDocumentationLanguage analyzes documentation files
func analyzeDocumentationLanguage() map[string]int {
	scores := make(map[string]int)

	// Check for language-specific documentation directories
	docDirs := []string{"docs", "doc", "documentation", ".github"}

	for _, dir := range docDirs {
		if _, err := os.Stat(dir); err == nil {
			// Check for language-specific subdirectories
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					dirName := strings.ToLower(entry.Name())
					// Common language directory patterns
					switch {
					case dirName == "zh" || dirName == "zh-cn" || dirName == "chinese":
						scores["zh"]++
					case dirName == "ja" || dirName == "japanese":
						scores["ja"]++
					case dirName == "ko" || dirName == "korean":
						scores["ko"]++
					case dirName == "de" || dirName == "german":
						scores["de"]++
					case dirName == "fr" || dirName == "french":
						scores["fr"]++
					case dirName == "es" || dirName == "spanish":
						scores["es"]++
					case dirName == "pt" || dirName == "portuguese":
						scores["pt"]++
					case dirName == "ru" || dirName == "russian":
						scores["ru"]++
					case dirName == "it" || dirName == "italian":
						scores["it"]++
					case dirName == "en" || dirName == "english":
						scores["en"]++
					}
				} else if strings.HasSuffix(entry.Name(), ".md") {
					// Analyze markdown files in docs directory
					filePath := filepath.Join(dir, entry.Name())
					content, err := os.ReadFile(filePath)
					if err != nil {
						continue
					}
					langScores := detectLanguageInText(string(content))
					for lang, score := range langScores {
						scores[lang] += score / 5 // Lower weight for individual files
					}
				}
			}
		}
	}

	return scores
}

// detectLanguageInText uses heuristics to detect language in text
func detectLanguageInText(text string) map[string]int {
	scores := make(map[string]int)

	// Convert to lowercase for analysis
	lowerText := strings.ToLower(text)

	// Chinese detection (CJK characters)
	chineseChars := 0
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF { // Common CJK ideographs
			chineseChars++
		}
	}
	// Even a few Chinese characters indicate Chinese text
	// Use balanced scoring: each char counts as 1 point
	if chineseChars > 0 {
		scores["zh"] = chineseChars
	}

	// Japanese detection (Hiragana, Katakana)
	japaneseChars := 0
	for _, r := range text {
		if (r >= 0x3040 && r <= 0x309F) || // Hiragana
			(r >= 0x30A0 && r <= 0x30FF) { // Katakana
			japaneseChars++
		}
	}
	if japaneseChars > 0 {
		scores["ja"] = japaneseChars
	}

	// Korean detection (Hangul)
	koreanChars := 0
	for _, r := range text {
		if r >= 0xAC00 && r <= 0xD7AF { // Hangul syllables
			koreanChars++
		}
	}
	if koreanChars > 0 {
		scores["ko"] = koreanChars
	}

	// German detection (common German words and umlauts)
	germanIndicators := []string{"der ", "die ", "das ", "und ", "ein ", "eine ", "für ", "mit ", "über ", "ä", "ö", "ü", "ß"}
	germanScore := 0
	for _, indicator := range germanIndicators {
		germanScore += strings.Count(lowerText, indicator)
	}
	if germanScore > 0 {
		scores["de"] = germanScore
	}

	// French detection (common French words and accents)
	frenchIndicators := []string{"le ", "la ", "les ", "de ", "des ", "un ", "une ", "et ", "est ", "pour ", "avec ", "à ", "é", "è", "ê", "ç"}
	frenchScore := 0
	for _, indicator := range frenchIndicators {
		frenchScore += strings.Count(lowerText, indicator)
	}
	if frenchScore > 0 {
		scores["fr"] = frenchScore
	}

	// Spanish detection
	spanishIndicators := []string{"el ", "la ", "los ", "las ", "de ", "del ", "un ", "una ", "y ", "es ", "para ", "con ", "en ", "á", "é", "í", "ó", "ú", "ñ"}
	spanishScore := 0
	for _, indicator := range spanishIndicators {
		spanishScore += strings.Count(lowerText, indicator)
	}
	if spanishScore > 0 {
		scores["es"] = spanishScore
	}

	// Portuguese detection (more specific indicators to avoid confusion with Italian)
	portugueseIndicators := []string{" o ", " a ", " os ", " as ", " do ", " da ", " um ", " uma ", " para ", " com ", " em ", "ã", "õ"}
	portugueseScore := 0
	for _, indicator := range portugueseIndicators {
		portugueseScore += strings.Count(lowerText, indicator)
	}
	// Portuguese-specific diacritics have higher weight
	portugueseScore += strings.Count(lowerText, "ã") * 2
	portugueseScore += strings.Count(lowerText, "õ") * 2
	if portugueseScore > 0 {
		scores["pt"] = portugueseScore
	}

	// Russian detection (Cyrillic)
	russianChars := 0
	for _, r := range text {
		if r >= 0x0400 && r <= 0x04FF { // Cyrillic
			russianChars++
		}
	}
	if russianChars > 0 {
		scores["ru"] = russianChars
	}

	// Italian detection (more specific to avoid confusion with Portuguese/Spanish)
	italianIndicators := []string{" il ", " la ", " di ", " un ", " una ", " per ", " con ", " in ", " sono ", " della ", " dello "}
	italianScore := 0
	for _, indicator := range italianIndicators {
		italianScore += strings.Count(lowerText, indicator)
	}
	// Italian-specific diacritics
	italianScore += strings.Count(lowerText, "è") * 2
	italianScore += strings.Count(lowerText, "ò") * 2
	italianScore += strings.Count(lowerText, "ì") * 2
	if italianScore > 0 {
		scores["it"] = italianScore
	}

	// English detection (common English words)
	// Always check for English, but give it lower priority unless it's clearly dominant
	englishIndicators := []string{"the ", "is ", "are ", "was ", "were ", "and ", "or ", "but ", "for ", "with ", "this ", "that ", "have ", "has ", "been ", "can ", "will ", "add ", "update ", "fix ", "remove "}
	englishScore := 0
	for _, indicator := range englishIndicators {
		englishScore += strings.Count(lowerText, indicator)
	}
	// Set English score if detected (even small amounts for short texts like commits)
	if englishScore > 0 {
		scores["en"] = englishScore
	}

	return scores
}

// NormalizeLanguageCode normalizes various language code formats to standard codes
func NormalizeLanguageCode(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))

	// Handle common variations
	switch code {
	case "zh-cn", "zh-hans", "chinese", "simplified chinese":
		return "zh"
	case "ja", "japanese", "jp":
		return "ja"
	case "ko", "korean", "kr":
		return "ko"
	case "de", "german":
		return "de"
	case "fr", "french":
		return "fr"
	case "es", "spanish":
		return "es"
	case "pt", "portuguese", "pt-br":
		return "pt"
	case "ru", "russian":
		return "ru"
	case "it", "italian":
		return "it"
	case "en", "english", "en-us", "en-gb":
		return "en"
	default:
		// If not recognized, return as-is
		return code
	}
}
