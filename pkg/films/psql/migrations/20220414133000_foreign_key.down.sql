ALTER TABLE ONLY public.film
    DROP CONSTRAINT IF EXISTS film_language_id_fkey;

ALTER TABLE ONLY public.film_actor
    DROP CONSTRAINT IF EXISTS film_actor_actor_id_fkey;

ALTER TABLE ONLY public.film_actor
    DROP CONSTRAINT IF EXISTS film_actor_film_id_fkey;

ALTER TABLE ONLY public.film_category
    DROP CONSTRAINT IF EXISTS film_category_category_id_fkey;

ALTER TABLE ONLY public.film_category
    DROP CONSTRAINT IF EXISTS film_category_film_id_fkey;
