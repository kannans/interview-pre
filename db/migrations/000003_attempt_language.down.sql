DROP INDEX IF EXISTS idx_attempts_language;

ALTER TABLE attempts DROP COLUMN IF EXISTS language;
