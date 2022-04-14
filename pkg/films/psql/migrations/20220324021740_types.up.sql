CREATE TYPE public.mpaa_rating AS ENUM (
    'G',
    'PG',
    'PG-13',
    'R',
    'NC-17'
);
ALTER TYPE public.mpaa_rating OWNER TO postgres;

CREATE DOMAIN public.year AS integer
	CONSTRAINT year_check CHECK (((VALUE >= 1901) AND (VALUE <= 2155)));

ALTER DOMAIN public.year OWNER TO postgres;
