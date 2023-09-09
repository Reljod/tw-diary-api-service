-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
	id SERIAL,
	username VARCHAR(255) NOT NULL UNIQUE,
	password VARCHAR(255) NOT NULL,
	disabled smallint DEFAULT '0',
	PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts;
-- +goose StatementEnd
