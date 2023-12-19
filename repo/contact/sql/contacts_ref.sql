-- FUNCTION: public.contacts_ref(contacts)

-- DROP FUNCTION IF EXISTS public.contacts_ref(contacts);

CREATE OR REPLACE FUNCTION public.contacts_ref(
	contacts)
    RETURNS json
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
	SELECT json_build_object(
		'keys',json_build_object(
			'id',$1.id    
			),	
		'descr',$1.descr
		,
		'dataType','contacts'
	);
$BODY$;

ALTER FUNCTION public.contacts_ref(contacts)
    OWNER TO ;

