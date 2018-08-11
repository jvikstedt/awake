package job

var Schema = `
CREATE TABLE IF NOT EXISTS jobs (
	id integer PRIMARY KEY,
	name text,
	active integer DEFAULT 0,
	cron text,
	step_config_ids text,
	created_at timestamp,
	updated_at timestamp DEFAULT current_timestamp,
	deleted_at timestamp
);
`
