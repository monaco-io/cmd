package src

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Cache represents a translation cache
type Cache struct {
	config *Config
	dir    string
}

// CacheEntry represents a cached translation
type CacheEntry struct {
	Text        string    `json:"text"`
	Language    string    `json:"language"`
	Translation string    `json:"translation"`
	Timestamp   time.Time `json:"timestamp"`
}

// NewCache creates a new cache instance
func NewCache(cfg *Config) *Cache {
	return &Cache{
		config: cfg,
		dir:    cfg.GetCacheDir(),
	}
}

// Get retrieves a translation from cache if available and not expired
func (c *Cache) Get(text, language string) (string, bool) {
	if !c.config.Cache.Enabled {
		return "", false
	}

	key := c.generateKey(text, language)
	path := c.getCachePath(key)

	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false
	}

	// Check if cache is expired
	if c.config.Cache.TTL > 0 {
		expiryTime := entry.Timestamp.Add(time.Duration(c.config.Cache.TTL) * time.Hour)
		if time.Now().After(expiryTime) {
			// Cache expired, remove it
			os.Remove(path)
			return "", false
		}
	}

	if c.config.Advanced.Debug {
		fmt.Printf("[DEBUG] Cache hit for key: %s\n", key)
	}

	return entry.Translation, true
}

// Set stores a translation in cache
func (c *Cache) Set(text, language, translation string) error {
	if !c.config.Cache.Enabled {
		return nil
	}

	key := c.generateKey(text, language)
	path := c.getCachePath(key)

	// Ensure cache directory exists
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	entry := CacheEntry{
		Text:        text,
		Language:    language,
		Translation: translation,
		Timestamp:   time.Now(),
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	if c.config.Advanced.Debug {
		fmt.Printf("[DEBUG] Cached translation for key: %s\n", key)
	}

	return nil
}

// Clear removes all cached translations
func (c *Cache) Clear() error {
	if err := os.RemoveAll(c.dir); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}

// Status returns cache statistics
func (c *Cache) Status() (int64, int, error) {
	var totalSize int64
	fileCount := 0

	if _, err := os.Stat(c.dir); os.IsNotExist(err) {
		return 0, 0, nil
	}

	err := filepath.Walk(c.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	if err != nil {
		return 0, 0, fmt.Errorf("failed to calculate cache status: %w", err)
	}

	return totalSize, fileCount, nil
}

// generateKey generates a cache key from text and language
func (c *Cache) generateKey(text, language string) string {
	combined := fmt.Sprintf("%s:%s", language, text)
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// getCachePath returns the full path to a cache file
func (c *Cache) getCachePath(key string) string {
	// Use subdirectories to avoid too many files in one directory
	subdir := key[:2]
	return filepath.Join(c.dir, subdir, key+".json")
}
