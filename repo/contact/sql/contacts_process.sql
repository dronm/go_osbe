-- Function: contacts_process()

-- DROP FUNCTION contacts_process();

CREATE OR REPLACE FUNCTION contacts_process()
  RETURNS trigger AS
$BODY$
BEGIN
	IF (TG_WHEN='BEFORE' AND TG_OP='UPDATE' OR TG_OP='INSERT') THEN
		/*IF TG_OP='INSERT'
		OR coalesce(NEW.name,'') <> coalesce(OLD.name,'')
		OR coalesce(NEW.post_id,0) <> coalesce(OLD.post_id,0)
		OR coalesce(NEW.tel,'') <> coalesce(OLD.tel,'')
		OR coalesce(NEW.email,'') <> coalesce(OLD.email,'')
		OR coalesce(NEW.tel_ext,'') <> coalesce(OLD.tel_ext,'')
		THEN*/
			NEW.descr = coalesce(NEW.name,'')||
				CASE
					WHEN NEW.post_id IS NOT NULL THEN
						(SELECT '('||posts.name||') ' FROM posts WHERE posts.id = NEW.post_id)
					ELSE ''
				END||
				CASE
					WHEN coalesce(NEW.email,'')<>'' THEN ', '||NEW.email
					ELSE ''
				END||
				CASE
					--||format_cel_standart(
					WHEN coalesce(NEW.tel,'')<>'' THEN ', +7'||NEW.tel||
						CASE
							WHEN coalesce(NEW.tel_ext,'')<>'' THEN ' ('||NEW.tel_ext||')'
							ELSE ''
						END
					ELSE ''
				END				
			;
		--END IF;
		
		RETURN NEW;
	END IF;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION contacts_process() OWNER TO ;

