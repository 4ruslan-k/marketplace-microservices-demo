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

CREATE TABLE IF NOT EXISTS messages (
	id uuid OT NULL,
	"text" varchar NOT NULL,
	user_id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NULL,
	CONSTRAINT messages_pk PRIMARY KEY (id),
	CONSTRAINT messages_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE
);
