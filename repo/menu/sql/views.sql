-- Table: public.views

-- DROP TABLE IF EXISTS public.views;

CREATE TABLE IF NOT EXISTS public.views
(
    id serial NOT NULL,
    c text COLLATE pg_catalog."default",
    f text COLLATE pg_catalog."default",
    t text COLLATE pg_catalog."default",
    section text COLLATE pg_catalog."default" NOT NULL,
    descr text COLLATE pg_catalog."default" NOT NULL,
    limited boolean,
    CONSTRAINT views_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.views
    OWNER TO ;
-- Index: views_section_descr_idx

-- DROP INDEX IF EXISTS public.views_section_descr_idx;

CREATE UNIQUE INDEX IF NOT EXISTS views_section_descr_idx
    ON public.views USING btree
    (section COLLATE pg_catalog."default" ASC NULLS LAST, descr COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
