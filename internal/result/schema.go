package result

var Schema = `
CREATE TABLE IF NOT EXISTS results (
	id integer PRIMARY KEY,
	job_id integer,
	steps text,
	created_at timestamp,
	updated_at timestamp DEFAULT current_timestamp,
	deleted_at timestamp,
	FOREIGN KEY(job_id) REFERENCES jobs(id)
);
`
