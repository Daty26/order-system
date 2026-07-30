CREATE INDEX IF NOT EXISTS notifications_user_id_created_at_idx
ON notifications(user_id, created_at DESC)

CREATE INDEX IF NOT EXISTS notifications_user_id_status_created_at_idx
ON notifications(user_id, status, created_at DESC)
