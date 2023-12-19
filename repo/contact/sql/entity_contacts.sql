-- Table: public.entity_contacts

-- DROP TABLE IF EXISTS public.entity_contacts;

CREATE TABLE IF NOT EXISTS public.entity_contacts
(
    id serial NOT NULL,
    entity_type data_types NOT NULL,
    entity_id integer NOT NULL,
    contact_id integer NOT NULL,
    mod_date_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT entity_contacts_pkey PRIMARY KEY (id),
    CONSTRAINT entity_contacts_contact_id_fkey FOREIGN KEY (contact_id)
        REFERENCES public.contacts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.entity_contacts
    OWNER TO ;
-- Index: entity_contacts_contact_idx

-- DROP INDEX IF EXISTS public.entity_contacts_contact_idx;

CREATE INDEX IF NOT EXISTS entity_contacts_contact_idx
    ON public.entity_contacts USING btree
    (entity_type ASC NULLS LAST, contact_id ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: entity_contacts_id_idx

-- DROP INDEX IF EXISTS public.entity_contacts_id_idx;

CREATE UNIQUE INDEX IF NOT EXISTS entity_contacts_id_idx
    ON public.entity_contacts USING btree
    (entity_type ASC NULLS LAST, entity_id ASC NULLS LAST, contact_id ASC NULLS LAST)
    TABLESPACE pg_default;
