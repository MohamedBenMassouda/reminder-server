-- +goose Up
CREATE TABLE reminders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    category_id INTEGER,
    due_date DATETIME,
    priority TEXT CHECK(priority IN ('high', 'medium', 'low')) DEFAULT 'medium',
    status TEXT CHECK(status IN ('pending', 'completed', 'missed')) DEFAULT 'pending',
    is_recurring BOOLEAN DEFAULT FALSE,
    recurring_pattern TEXT,
    user_id INTEGER,
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_reminders_user_id ON reminders(user_id);
CREATE INDEX idx_reminders_category_id ON reminders(category_id);
CREATE INDEX idx_reminders_due_date ON reminders(due_date);

-- +goose Down
DROP INDEX IF EXISTS idx_reminders_due_date;
DROP INDEX IF EXISTS idx_reminders_category_id;
DROP INDEX IF EXISTS idx_reminders_user_id;
DROP TABLE reminders;
