-- Add column program_id and class_id
ALTER TABLE password_requests ADD COLUMN request_type SMALLINT NOT NULL AFTER token;
