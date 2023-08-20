ALTER TABLE notifications RENAME TO notifications_by_users;
ALTER TABLE notifications_by_users DROP COLUMN notification_type_id;
ALTER TABLE notifications_by_users ADD notification_id uuid NOT NULL;


CREATE TABLE IF NOT EXISTS notifications (
	id uuid NOT NULL DEFAULT uuid_generate_v1(),
	type_id varchar NOT NULL,
	template varchar NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT notifications_pk PRIMARY KEY (id)
);

ALTER TABLE notifications_by_users ADD CONSTRAINT notifications_by_users_fk FOREIGN KEY (id)
 REFERENCES notifications(id) ON DELETE CASCADE ON UPDATE CASCADE;
