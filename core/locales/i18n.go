package i18n

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

// Default default locale for i18n
var Default = "en-US"

// I18n struct that hold all translations
type I18n struct {
	Resource   *admin.Resource
	scope      string
	value      string
	Backends   []Backend
	CacheStore cache.CacheStoreInterface
}

// ResourceName change display name in qor admin
func (I18n) ResourceName() string {
	return "Translation"
}

// Backend defined methods that needs for translation backend
type Backend interface {
	LoadTranslations() []*Translation
	SaveTranslation(*Translation) error
	DeleteTranslation(*Translation) error
}

// Translation is a struct for translations, including Translation Key, Locale, Value
type Translation struct {
	Key     string
	Locale  string
	Value   string
	Backend Backend `json:"-"`
}

// New initialize I18n with backends
func New(backends ...Backend) *I18n {
	i18n := &I18n{Backends: backends, CacheStore: memory.New()}
	for i := len(backends) - 1; i >= 0; i-- {
		var backend = backends[i]
		for _, translation := range backend.LoadTranslations() {
			i18n.AddTranslation(translation)
		}
	}
	return i18n
}

func (i18n *I18n) LoadTranslations() map[string]map[string]*Translation {
	var translations = map[string]map[string]*Translation{}

	for _, backend := range i18n.Backends {
		for _, translation := range backend.LoadTranslations() {
			if translations[translation.Locale] == nil {
				translations[translation.Locale] = map[string]*Translation{}
			}
			translations[translation.Locale][translation.Key] = translation
		}
	}
	return translations
}

// AddTranslation add translation
func (i18n *I18n) AddTranslation(translation *Translation) error {
	return i18n.CacheStore.Set(cacheKey(translation.Locale, translation.Key), translation)
}

// SaveTranslation save translation
func (i18n *I18n) SaveTranslation(translation *Translation) error {
	for _, backend := range i18n.Backends {
		if backend.SaveTranslation(translation) == nil {
			i18n.AddTranslation(translation)
			return nil
		}
	}

	return errors.New("failed to save translation")
}

// DeleteTranslation delete translation
func (i18n *I18n) DeleteTranslation(translation *Translation) (err error) {
	for _, backend := range i18n.Backends {
		backend.DeleteTranslation(translation)
	}

	return i18n.CacheStore.Delete(cacheKey(translation.Locale, translation.Key))
}

// Scope i18n scope
func (i18n *I18n) Scope(scope string) admin.I18n {
	return &I18n{CacheStore: i18n.CacheStore, scope: scope, value: i18n.value, Backends: i18n.Backends, Resource: i18n.Resource}
}

// Default default value of translation if key is missing
func (i18n *I18n) Default(value string) admin.I18n {
	return &I18n{CacheStore: i18n.CacheStore, scope: i18n.scope, value: value, Backends: i18n.Backends, Resource: i18n.Resource}
}

// T translate with locale, key and arguments
func (i18n *I18n) T(locale, key string, args ...interface{}) template.HTML {
	var (
		value          = i18n.value
		translationKey = key
	)

	if locale == "" {
		locale = Default
	}

	if i18n.scope != "" {
		translationKey = strings.Join([]string{i18n.scope, key}, ".")
	}

	var translation Translation
	if err := i18n.CacheStore.Unmarshal(cacheKey(locale, key), &translation); err != nil || translation.Value == "" {
		// Get default translation if not translated
		if err := i18n.CacheStore.Unmarshal(cacheKey(Default, key), &translation); err != nil || translation.Value == "" {
			// If not initialized
			translation = Translation{Key: translationKey, Value: value, Locale: locale, Backend: i18n.Backends[0]}

			// Save translation
			i18n.SaveTranslation(&translation)
		}
	}

	if translation.Value != "" {
		value = translation.Value
	}

	if str, err := cldr.Parse(locale, value, args...); err == nil {
		value = str
	}

	return template.HTML(value)
}

func cacheKey(strs ...string) string {
	return strings.Join(strs, "/")
}
