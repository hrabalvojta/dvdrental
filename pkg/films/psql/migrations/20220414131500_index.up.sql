CREATE INDEX idx_actor_last_name ON public.actor USING btree (last_name);

CREATE INDEX film_fulltext_idx ON public.film USING gist (fulltext);

CREATE INDEX idx_fk_language_id ON public.film USING btree (language_id);

CREATE INDEX idx_title ON public.film USING btree (title);

CREATE INDEX idx_fk_film_id ON public.film_actor USING btree (film_id);