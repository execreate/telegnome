package locale

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// These tests check the locale files that ship with the bot rather than fixtures. A key
// missing from one language would silently fall through to an empty string at runtime,
// so every language must define exactly the keys English does.

// shippedLanguages lists the languages that have a texts file in this directory.
func shippedLanguages(t *testing.T) []string {
	t.Helper()

	paths, err := filepath.Glob("*.yaml")
	if err != nil {
		t.Fatalf("failed to list locale files: %v", err)
	}

	var languages []string
	for _, path := range paths {
		name := strings.TrimSuffix(path, ".yaml")
		if !strings.HasSuffix(name, "_commands") {
			languages = append(languages, name)
		}
	}
	if !slices.Contains(languages, FallbackLanguage) {
		t.Fatalf("no %s.yaml among the shipped locales %v", FallbackLanguage, languages)
	}

	return languages
}

func readLocaleFile(t *testing.T, path string) *viper.Viper {
	t.Helper()

	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}

	return v
}

func TestShippedLocalesDefineTheSameKeys(t *testing.T) {
	for _, suffix := range []string{".yaml", "_commands.yaml"} {
		want := readLocaleFile(t, FallbackLanguage+suffix)
		wantKeys := want.AllKeys()
		slices.Sort(wantKeys)

		for _, lang := range shippedLanguages(t) {
			path := lang + suffix
			t.Run(path, func(t *testing.T) {
				got := readLocaleFile(t, path)
				gotKeys := got.AllKeys()
				slices.Sort(gotKeys)

				for _, key := range wantKeys {
					if !slices.Contains(gotKeys, key) {
						t.Errorf("%s is missing %q", path, key)
					} else if got.GetString(key) == "" {
						t.Errorf("%s has an empty %q", path, key)
					}
				}
				for _, key := range gotKeys {
					if !slices.Contains(wantKeys, key) {
						t.Errorf("%s has %q, which %s%s does not", path, key, FallbackLanguage, suffix)
					}
				}
			})
		}
	}
}

func TestShippedLocalesHaveATermsPage(t *testing.T) {
	for _, lang := range shippedLanguages(t) {
		page := filepath.Join("..", "static", "terms_and_conditions."+lang+".html")
		if _, err := os.Stat(page); err != nil {
			t.Errorf("no terms and conditions page for %s: %v", lang, err)
		}
	}
}
