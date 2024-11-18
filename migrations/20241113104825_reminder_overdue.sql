-- +goose Up
ALTER TABLE reminders ADD COLUMN overdue BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE reminders DROP COLUMN overdue;
