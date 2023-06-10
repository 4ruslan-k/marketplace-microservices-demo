BEGIN;
ALTER TABLE public.products ADD price numeric NOT NULL;
ALTER TABLE public.products ADD quantity int NOT NULL;
COMMIT;