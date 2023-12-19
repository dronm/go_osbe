package viewHTML

import (
	"errors"
	"fmt"
	"io/ioutil"

	"osbe/xml"
	"osbe/view"
	//"osbe/view/xml"
	"osbe/response"
	"osbe/socket"
	"osbe/srv/httpSrv"
)

const (
	VIEW_ID = "ViewHTML"
	QUERY_PARAM_TRANSFORM_TMPL = "templ" //transformation template
)

type OnTemplateTransformProto = func([]byte, string, string, string) ([]byte, error)
type OnBeforeRenderProto = func(*httpSrv.HTTPSocket, *response.Response) error

var v = &ViewHTML{}

type ViewHTML struct {
	SrvTemplateDir string
	UserViewDir string
	TemplateExtension string
	TemplateTransform OnTemplateTransformProto
	BeforeRender OnBeforeRenderProto
	
	DebugDir string
	Debug bool //if true, prints debug information
	XMLDebug bool //if true, xml data will be saved to DebugDir
	HTMLDebug bool //if true, html data will be saved to DebugDir
}

//Parameters:
//		SrvTemplateDir (string),
//		UserViewDir (string),
//		TemplateExtension (string),
//		TemplateTransform(OnTemplateTransformProto),
//		BeforeRender(OnBeforeRenderProto)
func (v *ViewHTML) Init(params map[string]interface{}) error {
	for id, val := range params {
		if err := v.SetParam(id, val); err != nil {
			return err
		}
	}
	return nil
}

func (v *ViewHTML) SetParam(paramID string, val interface{}) error {
	ok := false
	switch paramID {
	case "SrvTemplateDir":
		if v.SrvTemplateDir, ok = val.(string); !ok {
			return errors.New("parameter SrvTemplateDir must be a string")
		}
	case "UserViewDir":
		if v.UserViewDir, ok = val.(string); !ok {
			return errors.New("parameter UserViewDir must be a string")
		}	
	case "TemplateExtension":
		if v.TemplateExtension, ok = val.(string); !ok {
			return errors.New("parameter TemplateExtension must be a string")
		}	
	case "TemplateTransform":
		if v.TemplateTransform, ok = val.(OnTemplateTransformProto); !ok {
			return errors.New("parameter TemplateTransform must be of OnTemplateTransformProto type")
		}
	case "BeforeRender":
		if v.BeforeRender, ok = val.(OnBeforeRenderProto); !ok {
			return errors.New("parameter BeforeRender must be of OnBeforeRenderProto type")
		}
		
	case "DebugDir":
		if v.DebugDir, ok = val.(string); !ok {
			return errors.New("parameter DebugDir must be of string type")
		}
	case "XMLDebug":
		if v.XMLDebug, ok = val.(bool); !ok {
			return errors.New("parameter XMLDebug must be of bool type")
		}
	case "HTMLDebug":
		if v.HTMLDebug, ok = val.(bool); !ok {
			return errors.New("parameter HTMLDebug must be of bool type")
		}

	case "Debug":
		if v.Debug, ok = val.(bool); !ok {
			return errors.New("parameter Debug must be of bool type")
		}
		
	}
	return nil
}

//template is resolved in the following order:
//1) templ for error, user and role
//2) templ for error and user
//3) templ for user and role
//4) templ for user
//5) error + role + v
//6) error + v
//7) role + v
//8) v
func (v *ViewHTML) Render(sock socket.ClientSocketer, resp *response.Response) ([]byte, error){
	var sock_http *httpSrv.HTTPSocket
	var is_sock_http bool
	if sock != nil{
		sock_http, is_sock_http = sock.(*httpSrv.HTTPSocket)
	}
	if is_sock_http && v.BeforeRender != nil {
		//add extra models
		if err := v.BeforeRender(sock_http, resp); err != nil {
			return nil, err
		}
	}

	//render xml
	//xml_data, err := view.Render(viewXML.VIEW_ID, sock, resp)
	xml_data, err := xml.Marshal(resp.Models, false)
	if err != nil {
		return nil, err
	}
	
	//if no socket defined
	if sock == nil || !is_sock_http {
		return xml_data, nil
	}
	
	//+header
	//if sock_http.Response != nil {
		//sock_http.Response.Header().Set("Expires", ) //"Mon, 26 Jul 1997 05:00:00 GMT"
		//sock_http.Response.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		//sock_http.Response.Header().Set("Date", "time.Now().UTC().Format(http.TimeFormat))
		//sock_http.Response.Header().Set("Server", "")
		//sock_http.Response.Header().Set("Pragma", "no-cache")
		//Cache-Control: max-age=3600;
		//Cache-Control: no-cache, no-store
	//}
	
	//adding template
	err_code := resp.GetCode()
	sess := sock.GetSession()
	role_id := sess.GetString("USER_ROLE")
	
	template_file := ""
	
	template_id := ""	
	if par_tmpl, ok := sock_http.QueryParams[QUERY_PARAM_TRANSFORM_TMPL]; ok && len(par_tmpl)>0 {
		template_id = par_tmpl[0]
	}
	
	//error + role + templ		
	if err_code != 0 && template_id != "" && role_id != "" {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.Ex.%d.%s.%s", template_id, err_code, role_id, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: error + role + templ")
			}
		}
	}
	
	//error + templ
	if template_file == "" && err_code != 0 && template_id != "" {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.Ex.%d.%s", template_id, err_code, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: error + templ")
			}
		}
	}
	
	//templ + role
	if template_file == "" && err_code == 0 && template_id != "" && role_id != "" {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.%s.%s", template_id, role_id, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: templ + role")
			}
		}
	}

	//templ
	if template_file == "" && err_code == 0 && template_id != "" {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.%s", template_id, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: templ")
			}			
		}
	}

	//error + v + role
	if template_file == "" && err_code != 0 && role_id != "" {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.Ex.%d.%s.%s", sock_http.TransformClassID, err_code, role_id, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: error + v + role")
			}						
		}
	}

	//error + v
	if template_file == "" && err_code != 0 {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.Ex.%d.%s", sock_http.TransformClassID, err_code, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: error + v")
			}									
		}
	}

	//v + role
	if template_file == "" && err_code == 0 && role_id != "" {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.%s.%s", sock_http.TransformClassID, role_id, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: v + role")
			}									
		}
	}

	//v
	if template_file == "" {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.%s", sock_http.TransformClassID, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: v")
			}									
		}
	}

	//on error
	if template_file == "" && err_code != 0 && template_id != "" {
		fl := v.UserViewDir + "/" +  fmt.Sprintf("%s.%s", template_id, v.TemplateExtension)
		if view.FileExists(fl) {
			template_file = fl
			if v.Debug {
				fmt.Println("ViewHTML debug: template_id")
			}									
		}
	}

	//default server template for html is not used any more
	//if template_file == ""  && view.FileExists(v.SrvTemplateDir + "/" + DEF_SRV_TEMPLATE + "." + v.TemplateExtension) {
	//	template_file = v.SrvTemplateDir + "/" + DEF_SRV_TEMPLATE + "." + v.TemplateExtension
	//}
	
	if template_file == "" {
		return nil, errors.New("default server template not found in server template directory")
	}
	
	if v.Debug {
		fmt.Println("ViewHTML debug: template_file=", template_file)
	}									

	//transformation	
	html_data, err := v.TemplateTransform(xml_data, "", template_file, "")
	if err != nil {
		return nil, err
	}	
	if v.XMLDebug && v.DebugDir != "" {
		ioutil.WriteFile(v.DebugDir + "/xml_data.xml", xml_data, 0644)
	}
	if v.HTMLDebug && v.DebugDir != "" {
		ioutil.WriteFile(v.DebugDir + "/html_data.html", html_data, 0644)		
	}

	return html_data, nil
}

func init() {
	view.Register("ViewHTML", v)
}

