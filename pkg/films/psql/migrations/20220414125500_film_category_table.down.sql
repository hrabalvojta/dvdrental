DROP TABLE public.film_category CASCADE;

ALTER TABLE ONLY public.film_category
    DROP CONSTRAINT IF EXISTS film_category_pkey;