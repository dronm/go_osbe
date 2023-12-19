
CREATE TYPE template_value AS (
    field text,
    value text
);


-- Function: templates_text(template_value[],text)

-- DROP FUNCTION templates_text(template_value[],text);

CREATE OR REPLACE FUNCTION templates_text(template_value[],text)
	RETURNS text AS
$BODY$
DECLARE
   v_value template_value;
   v_text text;
BEGIN
	v_text = $2;
	FOREACH v_value IN ARRAY $1
	LOOP
		v_text = replace(v_text,
				'['||v_value.field||']',
				COALESCE(v_value.value,'')
		);
	END LOOP;
	
	RETURN v_text;
END
$BODY$
  LANGUAGE plpgsql COST 100;
  
ALTER FUNCTION templates_text(template_value[],text) OWNER TO ;
