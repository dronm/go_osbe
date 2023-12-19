package menu

import(
	"encoding/xml"
	"strings"
	"fmt"
	"context"
	
	"osbe"
	
	"github.com/jackc/pgx/v4"
)

type MenuItem struct {
	XMLName xml.Name `xml:"menuitem"`
	ViewID    string   `xml:"viewid,attr"`
	MenuItems   []MenuItem   `xml:"menuitem"`
}
type Menu struct {
    XMLName xml.Name `xml:"menu"`
    MenuItems   []MenuItem   `xml:"menuitem"`
}

func add_item(sql *strings.Builder, items []MenuItem, view_ids *[]string) {
	for _,it := range items {
		if it.ViewID == "" {
			if it.MenuItems!= nil && len(it.MenuItems) > 0 {
				add_item(sql, it.MenuItems, view_ids)
			}
			continue
		}
		if sql.Len() > 0 {
			sql.WriteString(" UNION ALL ")
		}
		sql.WriteString(fmt.Sprintf(
			`SELECT
				CASE WHEN v.c IS NOT NULL THEN 'c="' || v.c|| '"' ELSE '' END
				||CASE WHEN v.f IS NOT NULL THEN CASE WHEN v.c IS NULL THEN '' ELSE ' ' END|| 'f="' || v.f || '"' ELSE '' END
				||CASE WHEN v.t IS NOT NULL THEN CASE WHEN v.c IS NULL AND v.f IS NULL THEN '' ELSE ' ' END|| 't="' || v.t || '"' ELSE '' END
				||CASE WHEN v.limited IS NOT NULL AND v.limited THEN CASE WHEN v.c IS NULL AND v.f IS NULL AND v.t IS NULL THEN '' ELSE ' ' END|| 'limit="TRUE"' ELSE '' END
			FROM views v WHERE v.id=%s`,
			it.ViewID))		
		*view_ids = append(*view_ids, it.ViewID)
	}
}

func gen_user_menu(app osbe.Applicationer, conn *pgx.Conn, content string) (string,error) {
	content = strings.ReplaceAll(content, `xmlns="http://www.w3.org/1999/xhtml"`, "")
	content = strings.ReplaceAll(content, `xmlns="http://www.katren.org/crm/doc/mainmenu"`, "")
	
	var menu Menu
	if err := xml.Unmarshal([]byte(content), &menu); err != nil {
		return "",err
	}
	if menu.MenuItems == nil || len(menu.MenuItems) == 0 {
		return "",nil
	}
	var sql strings.Builder;
	view_ids := []string{}
	add_item(&sql, menu.MenuItems, &view_ids)
	
	sql_q := sql.String()
	if app.GetConfig().GetDebugQueries() {
		app.GetLogger().Debugf("Query debug gen_user_menu: %s", sql_q)
	}	
	rows, err := conn.Query(context.Background(), sql_q)
	if err != nil {
		return "",err
	}
	
	view_ind := 0	
	for rows.Next() {
		cmd := ""
		if err := rows.Scan(&cmd); err != nil {		
			return "",err
		}
		content = strings.ReplaceAll(content, fmt.Sprintf(`viewid="%s"`, view_ids[view_ind]), cmd);
		content = strings.ReplaceAll(content, fmt.Sprintf(`viewid ="%s"`, view_ids[view_ind]), cmd);
		content = strings.ReplaceAll(content, fmt.Sprintf(`viewid= "%s"`, view_ids[view_ind]), cmd);
		content = strings.ReplaceAll(content, fmt.Sprintf(`viewid = "%s"`, view_ids[view_ind]), cmd);
		view_ind++
	}
	if err := rows.Err(); err != nil {
		return "",err
	}

	return content, nil
}

