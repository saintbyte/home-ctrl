-- Test data for database testing

-- Delete all existing data
DELETE FROM key_value;
DELETE FROM api_keys;
DELETE FROM sessions;

-- Insert test API keys
INSERT INTO api_keys (key, name) VALUES
  ('test-api-key-12345', 'Test API Key'),
  ('admin-api-key-67890', 'Admin API Key');

-- Insert test sessions
INSERT INTO sessions (session_id, username, expires_at) VALUES
  ('valid-session-id', 'testuser', datetime('now', '+8 hours')),
  ('expired-session-id', 'testuser', datetime('now', '-1 hour'));

-- Insert test key-value pairs
INSERT INTO key_value (key, value, status, hidden) VALUES
  ('test-key-1', 'test-value-1', 'unread', 0),
  ('test-key-2', 'test-value-2', 'read', 1),
  ('test-key-3', 'test-value-3', 'archived', 0);
