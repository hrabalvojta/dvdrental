DROP TABLE public.film_actor CASCADE;

ALTER TABLE ONLY public.film_actor
    DROP CONSTRAINT IF EXISTS film_actor_pkey;