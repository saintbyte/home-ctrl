package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/saintbyte/home-ctrl/internal/database/models"
)

// CreateKeyValueTable creates the key_value table if it doesn't exist
func (d *Database) CreateKeyValueTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS key_values (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT UNIQUE NOT NULL,
		value TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'unread',
		is_hidden BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create key_values table: %w", err)
	}

	// Create index for key column
	_, err = d.db.Exec("CREATE INDEX IF NOT EXISTS idx_key_values_key ON key_values(key)")
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

// CreateKeyValue creates a new key-value pair
func (d *Database) CreateKeyValue(key, value string) (*models.KeyValue, error) {
	kv := models.NewKeyValue(key, value)

	result, err := d.db.Exec(
		"INSERT INTO key_values (key, value, status, is_hidden, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		kv.Key, kv.Value, kv.Status, kv.IsHidden, kv.CreatedAt, kv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create key-value: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	kv.ID = int(id)
	return kv, nil
}

// GetKeyValue retrieves a key-value pair by key
func (d *Database) GetKeyValue(key string) (*models.KeyValue, error) {
	var kv models.KeyValue

	err := d.db.QueryRow(
		"SELECT id, key, value, status, is_hidden, created_at, updated_at FROM key_values WHERE key = ?",
		key,
	).Scan(&kv.ID, &kv.Key, &kv.Value, &kv.Status, &kv.IsHidden, &kv.CreatedAt, &kv.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get key-value: %w", err)
	}

	return &kv, nil
}

// UpdateKeyValue updates an existing key-value pair
func (d *Database) UpdateKeyValue(key, value string) (*models.KeyValue, error) {
	kv, err := d.GetKeyValue(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key-value: %w", err)
	}
	if kv == nil {
		return nil, fmt.Errorf("key not found")
	}

	kv.UpdateValue(value)

	_, err = d.db.Exec(
		"UPDATE key_values SET value = ?, updated_at = ? WHERE key = ?",
		kv.Value, kv.UpdatedAt, kv.Key,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update key-value: %w", err)
	}

	return kv, nil
}

// UpdateKeyValueStatus updates the status of a key-value pair
func (d *Database) UpdateKeyValueStatus(key, status string) (*models.KeyValue, error) {
	kv, err := d.GetKeyValue(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key-value: %w", err)
	}
	if kv == nil {
		return nil, fmt.Errorf("key not found")
	}

	kv.SetStatus(status)

	_, err = d.db.Exec(
		"UPDATE key_values SET status = ?, updated_at = ? WHERE key = ?",
		kv.Status, kv.UpdatedAt, kv.Key,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update key-value status: %w", err)
	}

	return kv, nil
}

// UpdateKeyValueHidden updates the hidden flag of a key-value pair
func (d *Database) UpdateKeyValueHidden(key string, hidden bool) (*models.KeyValue, error) {
	kv, err := d.GetKeyValue(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key-value: %w", err)
	}
	if kv == nil {
		return nil, fmt.Errorf("key not found")
	}

	kv.SetHidden(hidden)

	_, err = d.db.Exec(
		"UPDATE key_values SET is_hidden = ?, updated_at = ? WHERE key = ?",
		kv.IsHidden, kv.UpdatedAt, kv.Key,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update key-value hidden flag: %w", err)
	}

	return kv, nil
}

// ListKeyValues lists all key-value pairs with optional filters
func (d *Database) ListKeyValues(includeHidden bool) ([]models.KeyValue, error) {
	var query string
	var rows *sql.Rows
	var err error

	if includeHidden {
		query = "SELECT id, key, value, status, is_hidden, created_at, updated_at FROM key_values ORDER BY created_at DESC"
		rows, err = d.db.Query(query)
	} else {
		query = "SELECT id, key, value, status, is_hidden, created_at, updated_at FROM key_values WHERE is_hidden = FALSE ORDER BY created_at DESC"
		rows, err = d.db.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list key-values: %w", err)
	}
	defer rows.Close()

	var keyValues []models.KeyValue
	for rows.Next() {
		var kv models.KeyValue
		if err := rows.Scan(&kv.ID, &kv.Key, &kv.Value, &kv.Status, &kv.IsHidden, &kv.CreatedAt, &kv.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan key-value: %w", err)
		}
		keyValues = append(keyValues, kv)
	}

	return keyValues, nil
}

// DeleteKeyValue deletes a key-value pair
func (d *Database) DeleteKeyValue(key string) error {
	_, err := d.db.Exec("DELETE FROM key_values WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("failed to delete key-value: %w", err)
	}
	return nil
}

// CheckKeyValueStatus checks if a key exists and returns its status
func (d *Database) CheckKeyValueStatus(key string) (string, bool, error) {
	var status string

	err := d.db.QueryRow(
		"SELECT status FROM key_values WHERE key = ?",
		key,
	).Scan(&status)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("failed to check key-value status: %w", err)
	}

	return status, true, nil
}

// CheckKeyValueExists checks if a key exists
func (d *Database) CheckKeyValueExists(key string) (bool, error) {
	var exists bool
	err := d.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM key_values WHERE key = ?)",
		key,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check key-value existence: %w", err)
	}

	return exists, nil
}

// CleanupKeyValues cleans up old or archived key-value pairs
func (d *Database) CleanupKeyValues(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	_, err := d.db.Exec(
		"DELETE FROM key_values WHERE status = ? AND updated_at < ?",
		models.StatusArchived, cutoff,
	)
	if err != nil {
		return fmt.Errorf("failed to cleanup key-values: %w", err)
	}
	return nil
}
