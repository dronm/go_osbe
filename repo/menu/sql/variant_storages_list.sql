-- VIEW: variant_storages_list

--DROP VIEW variant_storages_list;

CREATE OR REPLACE VIEW variant_storages_list AS
	SELECT
		user_id,
		storage_name,
		variant_name
	FROM variant_storages
	;
	
ALTER VIEW variant_storages_list OWNER TO ;
