package models

import "time"

// KeyValue represents a key-value pair with status and visibility flags
type KeyValue struct {
	ID        int       `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	Status    string    `json:"status" db:"status"`       // "unread", "read", "archived"
	IsHidden  bool      `json:"is_hidden" db:"is_hidden"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// KeyValueStatus constants
const (
	StatusUnread   = "unread"
	StatusRead     = "read"
	StatusArchived = "archived"
)

// NewKeyValue creates a new KeyValue instance
func NewKeyValue(key, value string) *KeyValue {
	return &KeyValue{
		Key:       key,
		Value:     value,
		Status:    StatusUnread,
		IsHidden:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SetStatus updates the status of the key-value pair
func (kv *KeyValue) SetStatus(status string) {
	kv.Status = status
	kv.UpdatedAt = time.Now()
}

// SetHidden updates the hidden flag
func (kv *KeyValue) SetHidden(hidden bool) {
	kv.IsHidden = hidden
	kv.UpdatedAt = time.Now()
}

// UpdateValue updates the value
func (kv *KeyValue) UpdateValue(value string) {
	kv.Value = value
	kv.UpdatedAt = time.Now()
}