BEGIN;

ALTER TABLE public.products DROP COLUMN price;
ALTER TABLE public.products DROP COLUMN quantity;

COMMIT;