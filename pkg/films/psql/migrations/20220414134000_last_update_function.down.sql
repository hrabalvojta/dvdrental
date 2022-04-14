DROP TRIGGER IF EXISTS last_updated ON public.actor CASCADE;

DROP TRIGGER IF EXISTS last_updated ON public.category CASCADE;

DROP TRIGGER IF EXISTS last_updated ON public.film CASCADE;

DROP TRIGGER IF EXISTS last_updated ON public.film_actor CASCADE;

DROP TRIGGER IF EXISTS last_updated ON public.language CASCADE;

DROP TRIGGER IF EXISTS last_updated ON public.film_category CASCADE;

DROP FUNCTION IF EXISTS public.last_updated();