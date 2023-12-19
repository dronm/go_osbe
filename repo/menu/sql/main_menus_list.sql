-- View: public.main_menus_list

-- DROP VIEW public.main_menus_list;

CREATE OR REPLACE VIEW public.main_menus_list AS 
 SELECT m.id,
    m.role_id,
    m.user_id,
    u.name AS user_descr
   FROM main_menus m
     LEFT JOIN users u ON u.id = m.user_id
  ORDER BY m.role_id, u.name;

ALTER TABLE public.main_menus_list OWNER TO ;
