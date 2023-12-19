-- View: public.attachments_list

-- DROP VIEW public.attachments_list;

CREATE OR REPLACE VIEW public.attachments_list
 AS
 SELECT t.id,
    t.date_time,
    t.ref,
    t.content_info,
    encode(t.content_preview, 'base64'::text) AS content_preview
   FROM attachments t
  ORDER BY t.date_time DESC;

ALTER TABLE public.attachments_list
    OWNER TO ;


