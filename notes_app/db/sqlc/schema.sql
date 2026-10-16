use notesdb;

CREATE TABLE roles(
	id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
	name VARCHAR(50) NOT NULL
);

CREATE TABLE actions(
	id int AUTO_INCREMENT PRIMARY KEY NOT NULL,
	name VARCHAR(50) NOT NULL
);

CREATE TABLE permissions(
	role_id int NOT NULL,
	action_id int NOT NULL, 
	CONSTRAINT fk_roles_in_permissions FOREIGN KEY(role_id) REFERENCES roles(id),
	CONSTRAINT fk_actions_in_permissions FOREIGN KEY(action_id) REFERENCES actions(id),
	PRIMARY KEY(role_id, action_id)
);

CREATE TABLE users(
	id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
	name VARCHAR(50) NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	phone_number VARCHAR(11),
	created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW() NOT NULL,
	role_id INT NOT NULL,
	CONSTRAINT fk_roles_in_users FOREIGN KEY(role_id) REFERENCES roles(id)
);

CREATE TABLE notes(
	id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
	text TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW() NOT NULL,
	priority TINYINT CHECK (priority BETWEEN 1 AND 3),
	user_id INT NOT NULL,
	CONSTRAINT fk_users_in_notes FOREIGN KEY(user_id) REFERENCES users(id)
);