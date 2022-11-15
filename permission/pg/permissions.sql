CREATE TABLE permissions (
    rules json NOT NULL
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.permissions OWNER to ;

INSERT INTO permissions VALUES ('{"admin":{
	"Test":{
		"insert":true,"delete":true,"update":true,"get_object":true
	}								
}}');

update permissions set rules = '{
	"admin":{
		"Test":{
			"insert":true,"delete":true,"update":true,"get_object":true
		}								
	},
	"manager":{
		"Test":{
			"insert":false,"delete":false,"update":true,"get_object":true
		}								
	},
	"guest":{
		"Test":{
			"insert":true,"delete":true,"update":true,"get_object":true,"get_list":true
		}								
	}
	
}';
