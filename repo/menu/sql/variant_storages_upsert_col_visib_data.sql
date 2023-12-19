-- Function: variant_storages_upsert_col_visib_data(in_user_id int, in_storage_name text, in_variant_name text, in_col_visib_data text, in_default_variant boolean)

--DROP FUNCTION variant_storages_upsert_col_visib_data(in_user_id int, in_storage_name text, in_variant_name text, in_col_visib_data text, in_default_variant boolean);

CREATE OR REPLACE FUNCTION variant_storages_upsert_col_visib_data(in_user_id int, in_storage_name text, in_variant_name text, in_col_visib_data text, in_default_variant boolean)
  RETURNS void AS
$BODY$  
BEGIN
	IF in_default_variant THEN
		UPDATE variant_storages
		SET
			default_variant = FALSE
		WHERE
			user_id = in_user_id
			AND storage_name = in_storage_name
		;	
	END IF;
	
	UPDATE variant_storages
	SET
		--set_time = now(),
		col_visib_data = in_col_visib_data,
		default_variant = in_default_variant
	WHERE
		user_id = in_user_id
		AND storage_name = in_storage_name
		AND variant_name = in_variant_name
	;
	
	IF FOUND THEN
		RETURN;
	END IF;
	
	BEGIN
		INSERT INTO variant_storages (user_id, storage_name, variant_name, col_visib_data, default_variant)
		VALUES(in_user_id, in_storage_name, in_variant_name, in_col_visib_data, in_default_variant);
		
	EXCEPTION WHEN OTHERS THEN
		UPDATE variant_storages
		SET
			--set_time = now(),
			col_visib_data = in_col_visib_data,
			default_variant = in_default_variant
		WHERE
			user_id = in_user_id
			AND storage_name = in_storage_name
			AND variant_name = in_variant_name
		;
	END;
	
	RETURN;

END;	
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION variant_storages_upsert_col_visib_data(in_user_id int, in_storage_name text, in_variant_name text, in_col_visib_data text, in_default_variant boolean) OWNER TO ;
