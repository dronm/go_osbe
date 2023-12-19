-- View: public.views_list

-- DROP VIEW public.views_list;

CREATE OR REPLACE VIEW public.views_list AS 
 SELECT v.id,
    v.c,
    v.f,
    v.t,
    v.section,
    v.descr,
    v.limited,
    (v.section || ' '::text) || v.descr AS user_descr,
    ((
        CASE
            WHEN v.c IS NOT NULL THEN ('c="'::text || v.c) || '"'::text
            ELSE ''::text
        END ||
        CASE
            WHEN v.f IS NOT NULL THEN ((
            CASE
                WHEN v.c IS NULL THEN ''::text
                ELSE ' '::text
            END || 'f="'::text) || v.f) || '"'::text
            ELSE ''::text
        END) ||
        CASE
            WHEN v.t IS NOT NULL THEN ((
            CASE
                WHEN v.c IS NULL AND v.f IS NULL THEN ''::text
                ELSE ' '::text
            END || 't="'::text) || v.t) || '"'::text
            ELSE ''::text
        END) ||
        CASE
            WHEN v.limited IS NOT NULL AND v.limited THEN
            CASE
                WHEN v.c IS NULL AND v.f IS NULL AND v.t IS NULL THEN ''::text
                ELSE ' '::text
            END || 'limit="TRUE"'::text
            ELSE ''::text
        END AS href
   FROM views v
  ORDER BY v.section, v.descr;

ALTER TABLE public.views_list OWNER TO ;

