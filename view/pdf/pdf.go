package viewPDF

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"osbe/xml"
	"osbe/view"
	"osbe/response"
	"osbe/socket"
	"osbe/srv/httpSrv"
)

const (
	VIEW_ID = "ViewPDF"
	TEMPLATE_EXT = "pdf.xsl"
	DEF_FILE_NAME = "model"
	QUERY_PARM_FNAME = "fname"
	OUT_FILE_EXT = "pdf"
	QUERY_PARAM_TRANSFORM_TMPL = "templ" //transformation template
)

var v = &ViewPDF{}

type OnTemplateToPDFTransformProto = func(string, string, []string, []byte, string, string, string) ([]byte, error)

type ViewPDF struct {
	SrvTemplateDir string
	UserViewDir string
	TemplateTransform OnTemplateToPDFTransformProto
	
	DebugDir string
	Debug bool //if true, xml data will be saved to DebugDir
	PDFDebug bool //if true, xml data will be saved to DebugDir
	Fop string
	ConfFile string
	ExecParams []string
}

//Parameters:
//		SrvTemplateDir (string),
//		UserViewDir (string),
func (v *ViewPDF) Init(params map[string]interface{}) error {
	for id, val := range params {
		if err := v.SetParam(id, val); err != nil {
			return err
		}
	}
	return nil
}

func (v *ViewPDF) SetParam(paramID string, val interface{}) error {
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
	case "TemplateTransform":
		if v.TemplateTransform, ok = val.(OnTemplateToPDFTransformProto); !ok {
			return errors.New("parameter TemplateTransform must be of OnTemplateToPDFTransformProto type")
		}
	case "DebugDir":
		if v.DebugDir, ok = val.(string); !ok {
			return errors.New("parameter DebugDir must be of string type")
		}
	case "Debug":
		if v.Debug, ok = val.(bool); !ok {
			return errors.New("parameter Debug must be of bool type")
		}
	case "Fop":
		if v.Fop, ok = val.(string); !ok {
			return errors.New("parameter Fop must be of string type")
		}
	case "ConfFile":
		if v.ConfFile, ok = val.(string); !ok {
			return errors.New("parameter ConfFile must be of string type")
		}
	case "ExecParams":
		if v.ExecParams, ok = val.([]string); !ok {
			return errors.New("parameter ExecParams must be of []string type")
		}
		
	}
	return nil
}

func (v *ViewPDF) Render(sock socket.ClientSocketer, resp *response.Response) ([]byte, error){
	var sock_http *httpSrv.HTTPSocket
	var is_sock_http bool
	if sock != nil{
		sock_http, is_sock_http = sock.(*httpSrv.HTTPSocket)
	}

	xml_data, err := xml.Marshal(resp.Models, true)
	if err != nil {
		return nil, err
	}
	//if no socket defined
	if sock == nil || !is_sock_http {
		return xml_data, nil
	}
	
	template_file := ""
	
	template_id := ""	
	if par_tmpl, ok := sock_http.QueryParams[QUERY_PARAM_TRANSFORM_TMPL]; ok && len(par_tmpl)>0 {
		template_id = par_tmpl[0]
	}
	
	if template_id != "" {
		if view.FileExists(v.UserViewDir + "/" +  fmt.Sprintf("%s.%s", template_id, TEMPLATE_EXT)) {
			template_file = v.UserViewDir + "/" +  fmt.Sprintf("%s.%s", template_id, TEMPLATE_EXT)
			
		}else if view.FileExists(v.SrvTemplateDir + "/" +  fmt.Sprintf("%s.%s", template_id, TEMPLATE_EXT)) {
			template_file = v.SrvTemplateDir + "/" +  fmt.Sprintf("%s.%s", template_id, TEMPLATE_EXT)
		}
	}else {
		if view.FileExists(v.UserViewDir + "/" +  fmt.Sprintf("%s.%s", sock_http.TransformClassID, TEMPLATE_EXT)) {
			template_file = v.UserViewDir + "/" +  fmt.Sprintf("%s.%s", sock_http.TransformClassID, TEMPLATE_EXT)
			
		}else if view.FileExists(v.SrvTemplateDir + "/" +  fmt.Sprintf("%s.%s", sock_http.TransformClassID, TEMPLATE_EXT)) {		
			template_file = v.SrvTemplateDir + "/" +  fmt.Sprintf("%s.%s", sock_http.TransformClassID, TEMPLATE_EXT)
		}
	}

	if template_file == "" {
		return nil, errors.New("default server template not found in server template directory")
	}
	if v.Debug && v.DebugDir != "" {
		//fmt.Println("ViewPDF->Render template_file=", template_file)	
		ioutil.WriteFile(v.DebugDir + "/xml_data.xml", xml_data, 0644)
	}

	view_data, err := v.TemplateTransform(v.Fop, v.ConfFile, v.ExecParams, xml_data, "", template_file, "")
	if err != nil {
		return nil, err
	}	
	if v.Debug && v.DebugDir != "" {
		ioutil.WriteFile(v.DebugDir + "/view_data.pdf", view_data, 0644)		
	}

	f_name := ""
	if fn_par, ok := sock_http.QueryParams[QUERY_PARM_FNAME]; ok && len(fn_par) > 0 {
		f_name = fn_par[0]
	}else{
		f_name = DEF_FILE_NAME
	}
	f_name+= "." + OUT_FILE_EXT

	resp = nil
	httpSrv.ServeContent(sock_http, &view_data, f_name, httpSrv.MIME_TYPE_pdf, time.Now(), httpSrv.CONTENT_DISPOSITION_ATTACHMENT)

	return nil, nil
}

func init() {
	view.Register(VIEW_ID, v)
}

