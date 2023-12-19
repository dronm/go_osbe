-- Table: public.posts

-- DROP TABLE IF EXISTS public.posts;

CREATE TABLE IF NOT EXISTS public.posts
(
    id serial NOT NULL,
    name character varying(250) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT posts_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.posts
    OWNER TO ;
-- Index: posts_name_idx

-- DROP INDEX IF EXISTS public.posts_name_idx;

CREATE UNIQUE INDEX IF NOT EXISTS posts_name_idx
    ON public.posts USING btree
    (lower(name::text) COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
