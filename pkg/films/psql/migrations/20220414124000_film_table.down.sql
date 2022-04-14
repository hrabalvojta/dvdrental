--
-- Name: film; Type: TABLE; Schema: public; Owner: postgres
--

DROP TABLE public.film CASCADE;

DROP SEQUENCE public.film_film_id_seq CASCADE;

ALTER TABLE ONLY public.film
    DROP CONSTRAINT IF EXISTS film_pkey;

DROP TRIGGER IF EXISTS film_fulltext_trigger ON public.film CASCADE;

