package src

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Translator handles translation operations
type Translator struct {
	config *Config
	client *Client
	cache  *Cache
	logger *slog.Logger
}

// NewTranslator creates a new translator instance
func NewTranslator(cfg *Config) (*Translator, error) {
	level := slog.LevelInfo
	if cfg.Advanced.Debug {
		level = slog.LevelDebug
	}
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

	return &Translator{
		config: cfg,
		client: NewClient(cfg),
		cache:  NewCache(cfg),
		logger: log,
	}, nil
}

// Translate translates text to the specified language(s)
func (t *Translator) Translate(text string, targetLang string) (string, error) {
	useColor := shouldUseColor()

	// If target language specified, translate to that language only
	if targetLang != "" {
		translation, err := t.translateToLanguage(text, targetLang)
		if err != nil {
			return "", err
		}
		return formatSingleOutput(text, translation, targetLang, useColor), nil
	}

	// Otherwise, translate to priority languages
	var results []string
	for _, lang := range t.config.Languages.Priority {
		translation, err := t.translateToLanguage(text, lang)
		if err != nil {
			t.logger.Error("failed to translate", "lang", lang, "err", err)
			continue
		}

		langName := getLanguageName(lang)
		results = append(results, formatLine(langName, translation, useColor))
	}

	if len(results) == 0 {
		return "", fmt.Errorf("failed to translate to any language")
	}

	output := formatMultiOutput(text, results, useColor)
	return output, nil
}

// translateToLanguage translates text to a specific language
func (t *Translator) translateToLanguage(text, lang string) (string, error) {
	// Check cache first
	if cached, found := t.cache.Get(text, lang); found {
		t.logger.Debug("cache hit", "lang", lang)
		return cached, nil
	}

	// Translate via API
	t.logger.Debug("calling API", "lang", lang)
	translation, err := t.client.Translate(text, lang)
	if err != nil {
		return "", fmt.Errorf("translation failed: %w", err)
	}

	// Cache the result
	if err := t.cache.Set(text, lang, translation); err != nil {
		t.logger.Error("failed to cache translation", "lang", lang, "err", err)
		// Don't fail on cache errors
	}

	return translation, nil
}

// Close closes the translator and its resources
func (t *Translator) Close() error {
	return nil
}

// getLanguageName returns the full name of a language code
func getLanguageName(code string) string {
	names := map[string]string{
		"zh":  "Chinese",
		"en":  "English",
		"ja":  "Japanese",
		"ko":  "Korean",
		"es":  "Spanish",
		"fr":  "French",
		"de":  "German",
		"ru":  "Russian",
		"pt":  "Portuguese",
		"it":  "Italian",
		"ar":  "Arabic",
		"hi":  "Hindi",
		"nl":  "Dutch",
		"pl":  "Polish",
		"tr":  "Turkish",
		"vi":  "Vietnamese",
		"th":  "Thai",
		"id":  "Indonesian",
		"ms":  "Malay",
		"fil": "Filipino",
		"he":  "Hebrew",
		"sv":  "Swedish",
		"no":  "Norwegian",
		"da":  "Danish",
		"fi":  "Finnish",
		"cs":  "Czech",
		"el":  "Greek",
		"ro":  "Romanian",
		"hu":  "Hungarian",
		"uk":  "Ukrainian",
	}

	if name, ok := names[code]; ok {
		return name
	}
	return strings.ToUpper(code)
}

// -------- Formatting helpers --------

const (
	colorReset   = "\033[0m"
	colorGreen   = "\033[32m"
	colorMagenta = "\033[35m"
	colorGray    = "\033[90m"
	colorBold    = "\033[1m"
)

func shouldUseColor() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("FANYI_NO_COLOR") != "" {
		return false
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func c(s, code string, enable bool) string {
	if !enable {
		return s
	}
	return code + s + colorReset
}

func formatSingleOutput(original, translation, lang string, color bool) string {
	langName := getLanguageName(lang)
	var b strings.Builder
	b.WriteString(c("Original", colorGray+colorBold, color))
	b.WriteString(": ")
	b.WriteString(original)
	b.WriteString("\n")
	b.WriteString(c(langName+":", colorGreen+colorBold, color))
	b.WriteString(" ")
	b.WriteString(translation)
	return b.String()
}

func formatMultiOutput(original string, lines []string, color bool) string {
	var b strings.Builder
	b.WriteString(c("Original", colorGray+colorBold, color))
	b.WriteString(": ")
	b.WriteString(original)
	b.WriteString("\n")
	b.WriteString(c(strings.Repeat("-", 6), colorGray, color))
	b.WriteString("\n")
	b.WriteString(strings.Join(lines, "\n"))
	return b.String()
}

func formatLine(langName, translation string, color bool) string {
	return fmt.Sprintf("%s %s %s", c("•", colorMagenta, color), c(langName+":", colorGreen+colorBold, color), translation)
}
