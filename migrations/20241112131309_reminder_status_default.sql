-- +goose Up
ALTER TABLE reminders RENAME COLUMN status TO status_old;
ALTER TABLE reminders ADD COLUMN status TEXT CHECK(status IN ('pending', 'completed', 'missed')) DEFAULT 'pending';
UPDATE reminders SET status = status_old;
ALTER TABLE reminders DROP COLUMN status_old;

-- +goose Down
ALTER TABLE reminders RENAME COLUMN status TO status_old;
ALTER TABLE reminders ADD COLUMN status TEXT CHECK(status IN ('pending', 'completed', 'missed'));
UPDATE reminders SET status = status_old;
ALTER TABLE reminders DROP COLUMN status_old;
