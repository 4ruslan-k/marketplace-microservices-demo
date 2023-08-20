CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL,
    name varchar NOT NULL,
    email varchar NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NULL,
    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_un UNIQUE (email)
);   

CREATE TABLE IF NOT EXISTS notifications (
	id uuid NOT NULL DEFAULT uuid_generate_v1(),
	type_id varchar NOT NULL,
	template varchar NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT notifications_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS notifications_by_users (
	id uuid NOT NULL DEFAULT uuid_generate_v1(),
	user_id uuid NULL,
	viewed_at timestamptz NULL,
	"data" jsonb NULL,
	title varchar NOT NULL,
	message varchar NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NULL,
	notification_id uuid NOT NULL,
	CONSTRAINT notifications_by_users_pk PRIMARY KEY (id),
	CONSTRAINT notifications_by_users_fk FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE ON UPDATE CASCADE
);




