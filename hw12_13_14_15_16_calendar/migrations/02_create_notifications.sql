-- +goose Up
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id TEXT NOT NULL,
    userId TEXT NOT NULL,
    title TEXT NOT NULL,
    event_date TIMESTAMPTZ NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL DEFAULT 'sent'
);

CREATE INDEX IF NOT EXISTS idx_notifications_event_id ON notifications (event_id);
CREATE INDEX IF NOT EXISTS idx_notifications_userId ON notifications (userId);
CREATE INDEX IF NOT EXISTS idx_notifications_sent_at ON notifications (sent_at);

-- +goose Down
DROP TABLE IF EXISTS notifications;