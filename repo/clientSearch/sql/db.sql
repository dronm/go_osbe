-- Table: client_search.dadata_cache

-- DROP TABLE IF EXISTS client_search.dadata_cache;

CREATE TABLE IF NOT EXISTS client_search.dadata_cache
(
    query text COLLATE pg_catalog."default" NOT NULL,
    response json,
    CONSTRAINT dadata_cache_pkey PRIMARY KEY (query)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS client_search.dadata_cache
    OWNER to ;
