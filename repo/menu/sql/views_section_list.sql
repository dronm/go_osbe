-- View: public.views_section_list

-- DROP VIEW public.views_section_list;

CREATE OR REPLACE VIEW public.views_section_list AS 
 SELECT DISTINCT views.section
   FROM views
  ORDER BY views.section;

ALTER TABLE public.views_section_list OWNER TO ;

