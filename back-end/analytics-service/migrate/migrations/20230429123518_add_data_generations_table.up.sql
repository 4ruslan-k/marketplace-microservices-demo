CREATE TABLE IF NOT EXISTS data_generation_by_users (
	user_id uuid NULL,
	"data" jsonb NOT NULL,
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	session_id varchar NULL,
	CONSTRAINT data_generation_by_users_pk PRIMARY KEY (id)
);
