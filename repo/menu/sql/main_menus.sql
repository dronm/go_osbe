-- Table: public.main_menus

-- DROP TABLE IF EXISTS public.main_menus;

CREATE TABLE IF NOT EXISTS public.main_menus
(
    id serial NOT NULL,
    role_id role_types NOT NULL,
    user_id integer,
    content text COLLATE pg_catalog."default" NOT NULL,
    model_content text COLLATE pg_catalog."default",
    CONSTRAINT main_menus_pkey PRIMARY KEY (id),
    CONSTRAINT main_menus_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.main_menus
    OWNER TO ;
-- Index: main_menus_role_user_idx

-- DROP INDEX IF EXISTS public.main_menus_role_user_idx;

CREATE UNIQUE INDEX IF NOT EXISTS main_menus_role_user_idx
    ON public.main_menus USING btree
    (role_id ASC NULLS LAST, user_id ASC NULLS LAST)
    TABLESPACE pg_default;
