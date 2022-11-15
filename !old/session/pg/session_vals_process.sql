-- FUNCTION: public.session_vals_process()

-- DROP FUNCTION public.session_vals_process();

CREATE FUNCTION public.session_vals_process()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
	IF (TG_WHEN='AFTER' AND TG_OP='DELETE') THEN
	ELSE 
		UPDATE logins SET date_time_out = now() WHERE session_id=OLD.id;
		
		RETURN OLD;
	END IF;
END;
$BODY$;

--ALTER FUNCTION public.session_vals_process() OWNER TO ;

