-- FUNCTION: public.posts_ref(posts)

-- DROP FUNCTION IF EXISTS public.posts_ref(posts);

CREATE OR REPLACE FUNCTION public.posts_ref(
	posts)
    RETURNS json
    LANGUAGE 'sql'
    COST 100
    VOLATILE PARALLEL UNSAFE
AS $BODY$
	SELECT json_build_object(
		'keys',json_build_object(
			'id',$1.id    
			),	
		'descr',$1.name,
		'dataType','posts'
	);
$BODY$;

ALTER FUNCTION public.posts_ref(posts)
    OWNER TO ;

