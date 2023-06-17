
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.cart_products (
	customer_id uuid NOT NULL,
	product_id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NULL,
	quantity int NOT NULL,
	CONSTRAINT cart_products_pk PRIMARY KEY (customer_id,product_id)
);



