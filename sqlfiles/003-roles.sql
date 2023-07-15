--
--
-- Role Fact Table
--
--	
-- Table contained list of roles available in the system
CREATE TABLE roles (
	id           VARCHAR(36)  NOT NULL,
 	code         VARCHAR(100) NOT NULL,
	name         VARCHAR(100) NOT NULL,
	status       VARCHAR(10)  NOT NULL DEFAULT 'inactive',
	created_by   VARCHAR(36)  NOT NULL,
	created_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by   VARCHAR(36)  NOT NULL,
	updated_time TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);
CREATE UNIQUE INDEX idx_roles_unique_code 
	ON roles (code);

