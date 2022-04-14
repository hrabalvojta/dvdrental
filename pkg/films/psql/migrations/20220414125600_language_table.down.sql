DROP TABLE public.language CASCADE;

DROP SEQUENCE public.language_language_id_seq CASCADE;

ALTER TABLE ONLY public.language
    DROP CONSTRAINT IF EXISTS language_pkey;