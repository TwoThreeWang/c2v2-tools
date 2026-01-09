package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Manager struct {
	translations map[string]map[string]string
	mu           sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		translations: make(map[string]map[string]string),
	}
}

// LoadTranslations reads all json files from the locales directory
func (m *Manager) LoadTranslations(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			lang := entry.Name()[0 : len(entry.Name())-5] // remove .json
			content, err := os.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				return err
			}

			var data map[string]string
			if err := json.Unmarshal(content, &data); err != nil {
				return fmt.Errorf("failed to parse %s: %v", entry.Name(), err)
			}
			m.translations[lang] = data
		}
	}
	return nil
}

// Translate returns the translation for a key in the given language
// If not found, it falls back to English, then to the key itself
func (m *Manager) Translate(lang, key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Try requested language
	if trans, ok := m.translations[lang]; ok {
		if val, ok := trans[key]; ok {
			return val
		}
	}

	// Fallback to "en"
	if lang != "en" {
		if trans, ok := m.translations["en"]; ok {
			if val, ok := trans[key]; ok {
				return val
			}
		}
	}

	// Fallback to key
	return key
}
