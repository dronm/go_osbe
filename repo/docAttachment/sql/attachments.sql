-- Table: public.attachments

-- DROP TABLE IF EXISTS public.attachments;

CREATE TABLE IF NOT EXISTS public.attachments
(
    id serial NOT NULL,
    date_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    ref jsonb,
    content_info jsonb,
    content_data bytea,
    content_preview bytea,
    CONSTRAINT attachments_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.attachments
    OWNER TO ;
-- Index: attachments_idx

-- DROP INDEX IF EXISTS public.attachments_idx;

CREATE INDEX IF NOT EXISTS attachments_idx
    ON public.attachments USING btree
    ((ref ->> 'dataType'::text) COLLATE pg_catalog."default" ASC NULLS LAST, (((ref -> 'keys'::text) ->> 'id'::text)::integer) ASC NULLS LAST)
    TABLESPACE pg_default;
