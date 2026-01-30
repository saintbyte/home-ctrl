-- Migration 002: Add default data

-- Check if we have any API keys
-- If not, create a default one
INSERT INTO api_keys (key, name)
SELECT 'default-api-key-12345', 'Default API Key'
WHERE NOT EXISTS (SELECT 1 FROM api_keys);