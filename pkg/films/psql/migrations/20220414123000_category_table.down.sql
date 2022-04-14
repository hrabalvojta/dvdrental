DROP TABLE public.category CASCADE;

DROP SEQUENCE public.category_category_id_seq CASCADE;

ALTER TABLE ONLY public.category
    DROP CONSTRAINT IF EXISTS category_pkey;
