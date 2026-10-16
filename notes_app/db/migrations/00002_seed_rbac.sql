-- +goose Up
INSERT IGNORE INTO roles(id, name) VALUES
(1, 'Administrator'),
(2, 'Manager'),
(3, 'User');

INSERT IGNORE INTO actions(id, name) VALUES
(1, 'Admin:Add'),
(2, 'Admin:Read'),
(3, 'Admin:Update'),
(4, 'Admin:Delete'),
(5, 'Manager:Add'),
(6, 'Manager:Read'),
(7, 'Manager:Update'),
(8, 'User:Add'),
(9, 'User:Read');

INSERT IGNORE INTO permissions(role_id, action_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4),
(2, 5), (2, 6), (2, 7),
(3, 8), (3, 9);

-- +goose Down
DELETE FROM permissions WHERE (role_id, action_id) IN ((1,1),(1,2),(1,3),(1,4),(2,5),(2,6),(2,7),(3,8),(3,9));
DELETE FROM actions WHERE id IN (1,2,3,4,5,6,7,8,9);
DELETE FROM roles WHERE id IN (1,2,3);
