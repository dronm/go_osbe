package httpSrv

import(
	"sync"
	"context"
	"fmt"
	"os"
	"io/ioutil"
	//"strings"

	"github.com/dronm/ds/pgds"
	"osbe"
	"osbe/socket"
	"osbe/response"
	"osbe/model"
	"osbe/constants"
	
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HTTPApplication struct {
	osbe.Application
	
	UserTmplDir string
	UserTmplExtension string
	
	//ServerVariables *model.Model
	JavaScriptModel *model.Model	
	CSSModel *model.Model
	
	mx sync.RWMutex
	constantQuery string
	
	//cashed server templates for roles
	serverTemplates map[string]*model.Model
}

//@ToDo: store menu common and for user if it exists
func (a *HTTPApplication) AddMainMenuModel(sock socket.ClientSocketer, resp *response.Response, conn *pgx.Conn) error {
//fmt.Println("AddMainMenuModel")	
	sess := sock.GetSession()		
	role := sess.GetString(osbe.SESS_ROLE)

	menu := &model.Model{ID: model.ModelID("MainMenu_Model"), SysModel: true, Rows: make([]model.ModelRow, 1)}
	q := `SELECT model_content
		FROM (
			SELECT model_content FROM main_menus WHERE user_id = $1
			UNION ALL
			`
	if role != "" {
		q+= fmt.Sprintf("SELECT model_content FROM main_menus WHERE role_id='%s'::role_types AND user_id IS NULL", role)
	}else{
		q+= "SELECT model_content FROM main_menus WHERE role_id IS NULL AND user_id IS NULL"
	}
	q+= ") AS s LIMIT 1"	
	err := conn.QueryRow(context.Background(), q, sess.GetInt("USER_ID")).Scan(&menu.RawData)
//fmt.Println(q)	
	if err != nil && err != pgx.ErrNoRows {
		return err		
	}
	menu.RawData = []byte(`<model id="MainMenu_Model" sysModel="1">` + string(menu.RawData) + `</model>`)
	resp.AddModel(menu)
	
	return nil
}

func (a *HTTPApplication) InitServerTemplateCache() {
	a.serverTemplates = make(map[string]*model.Model, 0)
}

//adds view template to all views
func (a *HTTPApplication) AddServerTemplate(sock *HTTPSocket, resp *response.Response) error {
	sess := sock.GetSession()
	role_id := ""
	if sess != nil {
		role_id = sess.GetString(osbe.SESS_ROLE)
	}
	//add view template
	if sock.ViewTemplateID != "" && a.UserTmplDir != "" && a.UserTmplExtension != "" {
		if a.serverTemplates == nil {
			a.InitServerTemplateCache()
		}
		
		//server template cache
		srv_tmpl_id := sock.ViewTemplateID + role_id
		if m, ok := a.serverTemplates[srv_tmpl_id]; ok {
			resp.AddModel(m)
			return nil
		}
				
		template_file := ""		
		//for specific role
		if role_id != "" {
			fl := a.UserTmplDir+"/"+fmt.Sprintf("%s.%s.%s", sock.ViewTemplateID, role_id, a.UserTmplExtension)
			if _, err := os.Stat(fl); err == nil || !os.IsNotExist(err) {
				template_file = fl
			}
		}
		
		//for all roles
		if template_file == "" {
			fl := a.UserTmplDir+"/"+fmt.Sprintf("%s.%s", sock.ViewTemplateID, a.UserTmplExtension)
			if _, err := os.Stat(fl); err == nil || !os.IsNotExist(err) {
				template_file = fl
			}
		}
		
		if template_file != "" {
			cont_b, err := ioutil.ReadFile(template_file)
			if err != nil {
				return err					
			}
//fmt.Println("AddServerTemplate template_file=",template_file)			
			model_data := fmt.Sprintf(`<model id="%s-template" templateId="%s" sysModel="1">%s</model>`,
				sock.ViewTemplateID,
				sock.ViewTemplateID,
				//strings.Replace(string(cont_b), "{{id}}", sock.ViewTemplateID, -1))
				string(cont_b))
			
			a.serverTemplates[srv_tmpl_id] = &model.Model{ID: model.ModelID(sock.ViewTemplateID), SysModel: true, RawData: []byte(model_data)}
			resp.AddModel(a.serverTemplates[srv_tmpl_id])
		}			
	}
	return nil
}

//adds constants marked as autoload
//cashe values from constant objects
func (a *HTTPApplication) AddAutoloadConstants(resp *response.Response, conn *pgx.Conn) error {
	query_id := "autoload_constants"
	if a.constantQuery == "" {
		a.mx.Lock()
		for const_id, const_o := range a.GetMD().Constants {
			if const_o.GetAutoload() {
				if a.constantQuery != "" {
					a.constantQuery += " UNION ALL "
				}
				a.constantQuery += fmt.Sprintf(`SELECT
						'%s' AS id,
						const_%s_val()::text AS val,
						(SELECT c.val_type FROM const_%s c) AS val_type`, const_id, const_id, const_id)
			}
		}
		a.mx.Unlock()
	}
	if _, err := conn.Prepare(context.Background(), query_id, a.constantQuery); err != nil {
		return err
	}
	//"ConstantValue_Model"
	return osbe.AddQueryResult(resp, a.GetMD().Models["ConstantValue"], &constants.ConstantValue{}, query_id, "", nil, conn, true)	
}

func (a *HTTPApplication) OnBeforeRenderXML(sock socket.ClientSocketer, resp *response.Response) error {
	if http_sock, ok := sock.(*HTTPSocket); ok {
		return a.AddServerTemplate(http_sock, resp)
	}
	return nil 
}

func (a *HTTPApplication) BeforeRenderHTML(sock *HTTPSocket, resp *response.Response) error {
	d_store,_ := a.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn	
	pool_conn, conn_id, err := d_store.GetSecondary("")
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()

	sess := sock.GetSession()
		
	if err := a.AddServerTemplate(sock, resp); err != nil {
		return err
	}

	if sess.GetBool("LOGGED") {
		//+main menu, not for child!!!
		if err := a.AddMainMenuModel(sock, resp, conn); err != nil {
			return err
		}
		//+constants
		if err := a.AddAutoloadConstants(resp, conn); err != nil {
			return err
		}
	}						
	
	//+javascript
	if a.JavaScriptModel != nil {
		resp.AddModel(a.JavaScriptModel)
	}
	
	//+css
	if a.CSSModel != nil {
		resp.AddModel(a.CSSModel)
	}
	
	return nil
}

func (a *HTTPApplication) Reload() {
	//server templates
	a.serverTemplates = make(map[string]*model.Model, 0)
} 
