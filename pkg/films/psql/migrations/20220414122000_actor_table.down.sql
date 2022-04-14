DROP TABLE public.actor CASCADE;

DROP SEQUENCE public.actor_actor_id_seq CASCADE;

ALTER TABLE ONLY public.actor
    DROP CONSTRAINT IF EXISTS actor_pkey;
