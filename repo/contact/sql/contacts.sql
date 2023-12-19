-- Table: public.contacts

-- DROP TABLE IF EXISTS public.contacts;

CREATE TABLE IF NOT EXISTS public.contacts
(
    id serial NOT NULL,
    name character varying(250) COLLATE pg_catalog."default" NOT NULL,
    post_id integer,
    email character varying(100) COLLATE pg_catalog."default",
    tel character varying(11) COLLATE pg_catalog."default",
    tel_ext character varying(20) COLLATE pg_catalog."default",
    descr text COLLATE pg_catalog."default",
    comment_text text COLLATE pg_catalog."default",
    email_confirmed boolean DEFAULT false,
    tel_confirmed boolean DEFAULT false,
    CONSTRAINT contacts_pkey PRIMARY KEY (id),
    CONSTRAINT contacts_post_id_fkey FOREIGN KEY (post_id)
        REFERENCES public.posts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.contacts
    OWNER TO ;
-- Index: contacts_descr_idx

-- DROP INDEX IF EXISTS public.contacts_descr_idx;

CREATE UNIQUE INDEX IF NOT EXISTS contacts_descr_idx
    ON public.contacts USING btree
    (lower(descr) COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: contacts_tel_idx

-- DROP INDEX IF EXISTS public.contacts_tel_idx;

CREATE UNIQUE INDEX IF NOT EXISTS contacts_tel_idx
    ON public.contacts USING btree
    (tel COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- Trigger: contacts_trigger_before

-- DROP TRIGGER IF EXISTS contacts_trigger_before ON public.contacts;

CREATE OR REPLACE TRIGGER contacts_trigger_before
    BEFORE INSERT OR UPDATE 
    ON public.contacts
    FOR EACH ROW
    EXECUTE FUNCTION public.contacts_process();
