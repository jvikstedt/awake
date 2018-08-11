package stepconfig

var Schema = `
CREATE TABLE IF NOT EXISTS step_configs (
	id integer PRIMARY KEY,
	tag text,
	variables text,
	created_at timestamp
);
`
