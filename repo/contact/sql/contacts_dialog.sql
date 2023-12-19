-- VIEW: contacts_dialog

DROP VIEW contacts_dialog;

CREATE OR REPLACE VIEW contacts_dialog AS
	SELECT
		ct.id,
		ct.name,
		posts_ref(p) AS posts_ref,
		ct.email,
		ct.tel,
		ct.descr,
		ct.tel_ext,
		ct.comment_text,
		ct.email_confirmed,
		ct.tel_confirmed

		/*
		(e_us.id IS NOT NULL) AS tm_exists,
		(e_us.tm_id IS NOT NULL) AS tm_activated,
		e_us.tm_photo,
		e_us.tm_first_name,
		e_us.id AS ext_id
		*/
		/*
		FALSE AS tm_exists,
		FALSE AS tm_activated,
		NULL AS tm_photo,
		NULL AS ext_id,
		NULL AS tm_first_name
		*/
		
		
	FROM contacts AS ct
	LEFT JOIN posts AS p ON p.id = ct.post_id
	--LEFT JOIN notifications.ext_users_photo_list AS e_us ON e_us.ext_contact_id = ct.id
	ORDER BY ct.name
	;
	
ALTER VIEW contacts_dialog OWNER TO ;
