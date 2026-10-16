/* Permissions */
-- name: GetAllPermissions :many
SELECT roles.name as roleName, actions.name as actionsName FROM permissions
JOIN roles ON permissions.role_id = roles.id
JOIN actions ON permissions.action_id = actions.id;

/* Actions */
-- name: GetAllActions :many
SELECT * FROM actions;

-- name: GetActionById :one
SELECT * FROM actions
WHERE id = ? LIMIT 1;

/* Roles */
-- name: GetAllRoles :many
SELECT * FROM roles;

-- name: GetRoleById :one
SELECT * FROM roles
WHERE id = ? LIMIT 1;

/* Users */
-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: AddUser :execresult
INSERT INTO users(name, email, password_hash, phone_number, role_id) VALUES (?,?,?,?,?);

-- name: UpdateUser :exec
UPDATE users 
SET name = ?, email = ?, password_hash = ?, phone_number = ?
WHERE id = ?;

-- name: ChangeRoles :exec
UPDATE users 
SET role_id = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: CheckAndGetInfoFromEmail :one
SELECT * FROM users 
WHERE email = ? LIMIT 1;

-- name: HasAccess :one
SELECT COUNT(*) FROM users
JOIN roles ON users.role_id = roles.id
JOIN permissions ON roles.id = permissions.role_id
JOIN actions ON permissions.action_id = actions.id
WHERE users.id = ? AND actions.name = ?;

/* Notes */
-- name: GetAllNotes :many
SELECT * FROM notes
JOIN users 
ON notes.user_id = users.id;

-- name: GetAllNotesByUserId :many
SELECT * FROM notes
WHERE user_id = ?;

-- name: GetNoteById :one
SELECT * FROM notes
WHERE id = ? LIMIT 1;

-- name: AddNote :execresult
INSERT INTO notes(text, priority, user_id) VALUES (?,?,?);

-- name: UpdateNote :exec
UPDATE notes 
SET text = ?, priority = ?
WHERE id = ?;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = ?;