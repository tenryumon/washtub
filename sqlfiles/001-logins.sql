--
--
-- User Part
--
--	
-- Table contained password for email login and token for google/facebook login
CREATE TABLE IF NOT EXISTS logins (
	id           VARCHAR(36)  NOT NULL,
	user_id      VARCHAR(36)  NOT NULL,
	login_type   VARCHAR(10)  NOT NULL DEFAULT 'password',
	token        VARCHAR(255) NOT NULL,
	status       VARCHAR(10)  NOT NULL DEFAULT 'inactive',
	created_by   VARCHAR(36)  NOT NULL,
	created_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by   VARCHAR(36)  NOT NULL,
	updated_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);
CREATE UNIQUE INDEX idx_logins_unique_user_id_login_type 
	ON logins (user_id, login_type);

--
--
-- Table contained unified UserID
CREATE TABLE IF NOT EXISTS users (
	id           VARCHAR(36)  NOT NULL,
	name         VARCHAR(100) NOT NULL,
	photo        VARCHAR(255) NOT NULL,
	email        VARCHAR(100) NOT NULL,
	phone        VARCHAR(50)  NOT NULL,
	status       VARCHAR(10)  NOT NULL DEFAULT 'inactive',
	user_type    VARCHAR(10)  NOT NULL DEFAULT 'user',
	last_action  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by   VARCHAR(36)  NOT NULL,
	created_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by   VARCHAR(36)  NOT NULL,
	updated_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);
CREATE UNIQUE INDEX idx_users_unique_email ON users (email, user_type);
CREATE UNIQUE INDEX idx_users_unique_phone ON users (phone, user_type);

--
--
-- Table contained email verification code for password creation
CREATE TABLE IF NOT EXISTS password_requests (
	id           VARCHAR(36)  NOT NULL,
	user_id      VARCHAR(36)  NOT NULL,
	token        VARCHAR(255) NOT NULL,
	status       VARCHAR(10)  NOT NULL DEFAULT 'inactive',
	expired_time TIMESTAMP    NOT NULL,
	created_by   VARCHAR(36)  NOT NULL,
	created_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by   VARCHAR(36)  NOT NULL,
	updated_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);
