-- +goose Up
-- Default administrator (email admin@gmail.com, password "admin").
-- Development seed only: the password is intentionally weak. Change it (or
-- delete this user) before any shared/staging environment.
INSERT IGNORE INTO users (id, name, email, password_hash, role_id) VALUES
(1, 'Admin', 'admin@gmail.com', '$2a$12$xIsy4tcOdfR8vsu5PYxYFuUyP1XsRagqEiV6MfwBB6htbtgCpThmC', 1);

-- +goose Down
DELETE FROM users WHERE id = 1;
