package database

import (
	"database/sql"
	"fmt"
	"time"
)

// APIKey represents an API key in the database
type APIKey struct {
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// CreateAPIKey creates a new API key
func (d *Database) CreateAPIKey(key, name string, expiresAt *time.Time) (*APIKey, error) {
	result, err := d.db.Exec(
		"INSERT INTO api_keys (key, name, expires_at) VALUES (?, ?, ?)",
		key, name, expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &APIKey{
		ID:        int(id),
		Key:       key,
		Name:      name,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}, nil
}

// GetAPIKeyByKey retrieves an API key by its key value
func (d *Database) GetAPIKeyByKey(key string) (*APIKey, error) {
	var apiKey APIKey
	var expiresAt sql.NullTime

	err := d.db.QueryRow(
		"SELECT id, key, name, created_at, expires_at FROM api_keys WHERE key = ?",
		key,
	).Scan(&apiKey.ID, &apiKey.Key, &apiKey.Name, &apiKey.CreatedAt, &expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if expiresAt.Valid {
		apiKey.ExpiresAt = &expiresAt.Time
	}

	return &apiKey, nil
}

// ListAPIKeys lists all API keys
func (d *Database) ListAPIKeys() ([]APIKey, error) {
	rows, err := d.db.Query("SELECT id, key, name, created_at, expires_at FROM api_keys")
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer rows.Close()

	var keys []APIKey
	for rows.Next() {
		var key APIKey
		var expiresAt sql.NullTime

		if err := rows.Scan(&key.ID, &key.Key, &key.Name, &key.CreatedAt, &expiresAt); err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}

		if expiresAt.Valid {
			key.ExpiresAt = &expiresAt.Time
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// DeleteAPIKey deletes an API key by its key value
func (d *Database) DeleteAPIKey(key string) error {
	_, err := d.db.Exec("DELETE FROM api_keys WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}

// ValidateAPIKey checks if an API key is valid
func (d *Database) ValidateAPIKey(key string) bool {
	apiKey, err := d.GetAPIKeyByKey(key)
	if err != nil || apiKey == nil {
		return false
	}

	// Check if key is expired
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}