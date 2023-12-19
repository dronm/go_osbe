-- Function: contacts_upsert(in_name text, in_tel text, in_email text, in_tel_ext text)

-- DROP FUNCTION contacts_upsert(in_name text, in_tel text, in_email text, in_tel_ext text);

CREATE OR REPLACE FUNCTION contacts_upsert(in_name text, in_tel text, in_email text, in_tel_ext text)
  RETURNS json AS
$$  
DECLARE
	v_contacts_ref json;
BEGIN  
	BEGIN
		INSERT INTO contacts (name, tel, email, tel_ext)
		VALUES (
			in_name,
			in_tel,
			CASE WHEN coalesce(in_email,'') = '' THEN NULL ELSE in_email END,
			CASE WHEN coalesce(in_tel_ext,'') = '' THEN NULL ELSE in_tel_ext END
		)
		RETURNING contacts_ref(contacts) AS contacts_ref INTO v_contacts_ref;
		
	EXCEPTION WHEN SQLSTATE '23505' THEN
		SELECT
			contacts_ref(contacts) AS contacts_ref
		INTO v_contacts_ref
		FROM contacts
		WHERE tel=in_tel;
	END;
	
	RETURN v_contacts_ref;
END;
$$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION contacts_upsert(in_name text, in_tel text, in_email text, in_tel_ext text) OWNER TO ;
