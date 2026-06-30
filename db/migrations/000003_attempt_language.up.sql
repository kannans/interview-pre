ALTER TABLE attempts
    ADD COLUMN IF NOT EXISTS language VARCHAR(20) NOT NULL DEFAULT 'javascript';

CREATE INDEX IF NOT EXISTS idx_attempts_language ON attempts(language);
