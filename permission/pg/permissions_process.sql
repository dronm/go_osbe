-- Function: permissions_process()

-- DROP FUNCTION permissions_process();

CREATE OR REPLACE FUNCTION permissions_process()
  RETURNS trigger AS
$BODY$
BEGIN
	PERFORM pg_notify('Permission.change', NULL);
	
	IF TG_WHEN='AFTER' AND (TG_OP='UPDATE' OR TG_OP='INSERT') THEN
		
		RETURN NEW;
		
	ELSIF TG_WHEN='AFTER' AND TG_OP='DELETE' THEN
		RETURN OLD;
		
	END IF;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION permissions_process() OWNER TO ;

