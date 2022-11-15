CREATE TABLE test (
	id serial NOT NULL,
	f1 int,
	f2 text,
	f3 numeric(15,2),
	f4 bool,
	CONSTRAINT test_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.test OWNER to test_proj;

