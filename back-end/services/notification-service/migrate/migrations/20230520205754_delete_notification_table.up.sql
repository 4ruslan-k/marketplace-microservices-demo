ALTER TABLE notifications_by_users DROP CONSTRAINT notifications_by_users_fk;
DROP TABLE IF  EXISTS notifications;
ALTER TABLE notifications_by_users DROP COLUMN notification_id;
DELETE from notifications_by_users;
ALTER TABLE notifications_by_users ADD notification_type_id varchar NOT NULL;
ALTER TABLE notifications_by_users RENAME TO notifications;

