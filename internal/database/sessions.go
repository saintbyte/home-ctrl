package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Session represents a user session
type Session struct {
	ID        int       `json:"id"`
	SessionID string    `json:"session_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreateSession creates a new session
func (d *Database) CreateSession(sessionID, username string, expiresAt time.Time) (*Session, error) {
	result, err := d.db.Exec(
		"INSERT INTO sessions (session_id, username, expires_at) VALUES (?, ?, ?)",
		sessionID, username, expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &Session{
		ID:        int(id),
		SessionID: sessionID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}, nil
}

// GetSessionByID retrieves a session by its ID
func (d *Database) GetSessionByID(sessionID string) (*Session, error) {
	var session Session

	err := d.db.QueryRow(
		"SELECT id, session_id, username, created_at, expires_at FROM sessions WHERE session_id = ?",
		sessionID,
	).Scan(&session.ID, &session.SessionID, &session.Username, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

// DeleteSession deletes a session by its ID
func (d *Database) DeleteSession(sessionID string) error {
	_, err := d.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// ValidateSession checks if a session is valid
func (d *Database) ValidateSession(sessionID string) bool {
	session, err := d.GetSessionByID(sessionID)
	if err != nil || session == nil {
		return false
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}

// CleanupExpiredSessions removes expired sessions
func (d *Database) CleanupExpiredSessions() error {
	_, err := d.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	return nil
}