package httpSrv

import (
	"net/http"
	"net/url"
	"strings"
	"fmt"
	"time"
	"encoding/json"
	"bytes"
	"os"
	"errors"
	
	"osbe/srv"
	"osbe/socket"
	"osbe/view"
	//"osbe/tokenBlock"	
	"osbe/stat"	
	"osbe/response"
)

//HTTP server, OnHandleRequest, must be defined

const (
	PARAM_TOKEN = "token"
	
	PARAM_CONTROLLER = "c"
	PARAM_METH = "f"
	PARAM_VIEW = "v"
	PARAM_VIEW_TMPL = "t" //view template to send with response, added in http_app
	PARAM_QUERY_ID = "query_id"
	
	PARAM_TRANSFORM_TMPL = "templ" //transformation template
	
	CONTROLLER_QUERY_POSF = "_Controller"
	
	DEF_USER_TRANSFORM_CLASS_ID =  "ViewBase"
	DEF_GUEST_TRANSFORM_CLASS_ID =  "Login"
	
	DEF_MULTYPART_MAX_MEM = 256
)

type requestHandlerProto = func(w http.ResponseWriter, r *http.Request)
type OnBeforeHandleRequestProto func(socket.ClientSocketer)

type argvType map[string]string//[]byte
//

type methodParams struct {
	Argv argvType `json:"argv"`
}

type URLShortcut struct {
	ControllerID string
	MethodID string
	ViewID string
	Params map[string]string
}

type HTTPServer struct {
	srv.BaseServer
	Statistics stat.SrvStater
	
	OnBeforeHandleRequest OnBeforeHandleRequestProto
	HTTPDir string	
	AllowedExtensions []string
	Headers map[string]string
	UserTransformClassID string
	GuestTransformClassID string
	URLShortcuts map[string]URLShortcut
	viewContentTypes map[string]string
	
	MultypartMaxMemory int64 //bytes
}

func (s *HTTPServer) AddURLShortcut(ID, cID, mID, vID string, params map[string]string) {
	if s.URLShortcuts == nil {
		s.URLShortcuts = make(map[string]URLShortcut)
	}
	s.URLShortcuts[ID] = URLShortcut{ControllerID: cID, MethodID: mID, ViewID: vID, Params: params}
}


func (s *HTTPServer) Run() {

	if s.OnHandleRequest == nil {
		s.Logger.Fatal("HTTPServer.OnHandleRequest not defined")
	}
	/*if s.OnHandlePermission != nil && s.OnHandleProhibError == nil {
		s.Logger.Fatal("HTTPServer.OnHandlePermission defined, but OnHandleProhibError not defined")
	}*/
	if s.OnHandleSession != nil && s.OnHandleServerError == nil {
		s.Logger.Fatal("HTTPServer.OnHandleSession defined, but OnHandleServerError not defined")
	}

	//TLS if nedded
	tls_start := (s.TlsAddress != "" && s.TlsCert != "" && s.TlsKey != "")
	ws_start := (s.Address!= "")
	
	http.HandleFunc("/", s.HandleRequest)
	
	s.Statistics = stat.NewSrvStat()
	
	if tls_start {
		s.Logger.Infof("Starting secured web server: %s", s.TlsAddress)		
		if !ws_start {
			//main loop
			http.ListenAndServeTLS(s.TlsAddress, s.TlsCert, s.TlsKey, nil)
		}else{
			//2 servers
			go http.ListenAndServeTLS(s.TlsAddress, s.TlsCert, s.TlsKey, nil)
		}
	}
	
	
	if ws_start {		
		s.Logger.Infof("Starting web server: %s", s.Address)		
		http.ListenAndServe(s.Address, nil)	
	}
}

//Retrieves parameter, first looks into cookies, then QueryParams
func extractParam(r *http.Request, queryParams url.Values, param string) (val string, exp time.Time) {
	//analyse cookies	
	if v_cookie, err := r.Cookie(param); v_cookie != nil && err == nil {
		val = v_cookie.Value
		exp = v_cookie.Expires
	}	
	
	//if param is not present in cookies, trying to get it from query params
	if val == "" {
		if val_par, ok := queryParams[param]; ok && len(val_par)>0 {
			val = val_par[0]
		}	
	}

	return
}

//parses query params based on query method, queryParams always non-nil map
func (s *HTTPServer) parseQueryParams(r *http.Request, queryParams *url.Values) {
	if r.Method == http.MethodGet {
		*queryParams = r.URL.Query()
	}else{
		ct := r.Header.Get("Content-Type")
		if strings.Contains(ct, "multipart/form-data") {
			var mem  int64
			if s.MultypartMaxMemory == 0 {
				mem = DEF_MULTYPART_MAX_MEM
			}else{
				mem = s.MultypartMaxMemory
			}
			r.ParseMultipartForm(mem)
			*queryParams = r.MultipartForm.Value
		}else{
			r.ParseForm()
			*queryParams = r.Form
		}
	}
}

func (s *HTTPServer) checkExtension(ext string) bool {
	for _,s_ext := range s.AllowedExtensions {
		if ext == s_ext {
			return true
		}
	}
	return false
}

func (s *HTTPServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	
	var query_params url.Values
	if r.URL.Path != "/" {
		path_parts := strings.Split(r.URL.Path, "/")

		sh_cut_found := false
		if s.URLShortcuts != nil {
			path := r.URL.Path
			if path[len(path)-1:] == "/" {
				path = path[:len(path)-1]
			}
			if sh_cut, ok := s.URLShortcuts[path]; ok {
				//Shortcuts - predefined paths
				s.parseQueryParams(r, &query_params)
							
				query_params.Add(PARAM_CONTROLLER, sh_cut.ControllerID)
				query_params.Add(PARAM_METH, sh_cut.MethodID)
				query_params.Add(PARAM_VIEW, sh_cut.ViewID)
				sh_cut_found = true
			}
		}
		
		if !sh_cut_found {
			file_parts := strings.Split(path_parts[len(path_parts)-1], ".")		
			n := len(file_parts)
			if n > 0 && s.checkExtension(file_parts[n-1]) && view.FileExists(s.HTTPDir + r.URL.Path) {
				//file serving
				http.ServeFile(w, r, s.HTTPDir + r.URL.Path)
				return
			}
			
			if len(path_parts) >= 2 {		
				//schema: controller/method/view
				s.parseQueryParams(r, &query_params)			
				query_params.Add(PARAM_CONTROLLER, path_parts[1])
				if len(path_parts) >= 3 {
					query_params.Add(PARAM_METH, path_parts[2])
				}	
				if len(path_parts) >= 4 {
					query_params.Add(PARAM_VIEW, path_parts[3])
				}
			}else{
				
				//not found
				s.Logger.Errorf("HTTPServer.OnHandleRequest %s file with extension %s not found", s.HTTPDir + r.URL.Path,file_parts[n-1])
				w.WriteHeader(http.StatusNotFound)
				//+ if ViewHTML return NotFound page???
				return
			}
		}
	}else{
		s.parseQueryParams(r, &query_params)
	}	
	sock := NewHTTPSocket(w, r)
	sock.Token, sock.TokenExpires = extractParam(r, query_params, PARAM_TOKEN)
	token_from_query := (sock.Token != "")
	
	//turn query/body parameters to json payload
	var query_id, view_id string	
	meth_params := methodParams{Argv: make(argvType)}
	
	for par_key, par_val:= range query_params {
		if par_key == PARAM_CONTROLLER && len(par_val)>0 {
			sock.ControllerID = par_val[0]
			//extract postfix if any
			posf_pos := len(sock.ControllerID)-len(CONTROLLER_QUERY_POSF)
			if posf_pos > 0 && sock.ControllerID[posf_pos:] == CONTROLLER_QUERY_POSF {
				sock.ControllerID = sock.ControllerID[:posf_pos]
			}
		
		}else if par_key == PARAM_METH && len(par_val)>0 {
			sock.MethodID = par_val[0]
			
		}else if par_key == PARAM_VIEW && len(par_val)>0 {
			sock.TransformClassID = par_val[0]
			view_id = par_val[0] 

		}else if par_key == PARAM_VIEW_TMPL && len(par_val)>0 {
			sock.ViewTemplateID = par_val[0]
			
		}else if par_key == PARAM_QUERY_ID && len(par_val)>0 {
			query_id = par_val[0]

		}else if par_key == PARAM_TRANSFORM_TMPL && len(par_val)>0 {
			sock.TransformTemplateID = par_val[0]
						
		}else if len(par_val)>0 {
			//[]byte(
			//meth_params.Argv[par_key] = par_val[0]
			meth_params.Argv[par_key] = par_val[0]
fmt.Println("par_val=", par_val[0], "byte=", string(meth_params.Argv[par_key]))					
		}
	}
	
	//session
	if s.OnHandleSession != nil {
		err := s.OnHandleSession(sock)
		if err != nil {
			s.Logger.Errorf("HTTPServer HandleRequest OnHandleSession: %v", err)
			s.OnHandleServerError(s, sock, query_id, view_id)
			return
		}

		if sock.Token == "" {
			//new session started
			sess := sock.GetSession()
			sock.Token = sess.SessionID()
			
			s.Statistics.IncHandshakes()
		}
		
		if !token_from_query {			
			//sock.TokenExpires = sess.
			// Make sure the session cookie is not accessable via javascript.
			http.SetCookie(w, &http.Cookie{Name: PARAM_TOKEN,
					Value: sock.Token,
					HttpOnly: true,
					//Expires= sock.TokenExpires,
					//Path:
					//Domain
					//Expires
					//MaxAge
				})
		}		
	}

	if sock.TransformClassID == "" {
		//defaults
		sess := sock.GetSession()	
		if sess.GetBool("LOGGED") {
			if s.UserTransformClassID != "" {
				sock.TransformClassID = s.UserTransformClassID
			}else{
				sock.TransformClassID = DEF_USER_TRANSFORM_CLASS_ID
			}
		}else{
			if s.GuestTransformClassID != "" {
				sock.TransformClassID = s.GuestTransformClassID
			}else{
				sock.TransformClassID = DEF_GUEST_TRANSFORM_CLASS_ID
			}			
		}
	}
	view_id = sock.TransformClassID
	if !view.Registered(view_id) {
		view_id = "ViewHTML"
	}	

	if query_id =="" {
		//http always expects result
		query_id = "1"
	}

fmt.Println("meth_params=", meth_params)			
	argv_s, err := json.Marshal(meth_params)
fmt.Println("argv_s=", string(argv_s))	
	if err != nil {
		s.Logger.Errorf("HTTPServer json.Marshal: %v", err)
		s.OnHandleServerError(s, sock, query_id, view_id)
		return
	}
	
	//header
	cont_tp := s.GetViewContentType(view_id)
	if cont_tp != "" {
		w.Header().Set("Content-Type", cont_tp)
	}else{
		s.Logger.Warnf("Content type for view %s not defined", view_id)
	}
	
	if s.Headers != nil {
		for key, val := range s.Headers {
			w.Header().Set(key, val)
		}
	}
	
	if s.OnBeforeHandleRequest != nil {
		s.OnBeforeHandleRequest(sock)
	}
	s.Logger.Debugf("HTTPServer calling OnHandleRequest ControllerID=%s, MethodID=%s, query_id=%s, argv_s=%s, view_id=%s", sock.ControllerID, sock.MethodID, query_id, argv_s, view_id)	
	
	s.OnHandleRequest(s, sock, sock.ControllerID, sock.MethodID, query_id, argv_s, view_id)
}

func (s *HTTPServer) SendToClient(sock socket.ClientSocketer, msg []byte) error {
	if http_sock, ok := sock.(*HTTPSocket); ok {
		fmt.Fprint(http_sock.Response, string(msg))
	}
	return nil
}

//empty stub, ClientSockets is not used
func (s *HTTPServer) GetClientSockets() *socket.ClientSocketList{
	return nil
}

func (s *HTTPServer) AddViewContentType(viewID, mimeType, charset string) {
	if s.viewContentTypes == nil {
		s.viewContentTypes = make(map[string]string)
	}
	s.viewContentTypes[viewID] = mimeType
	if charset != "" {
		s.viewContentTypes[viewID] += "; "+charset
	}
}

func (s *HTTPServer) GetViewContentType(viewID string) string {
	if tp, ok := s.viewContentTypes[viewID]; ok {
		return tp
	}
	return ""
}

func (s *HTTPServer) AddFile(viewID, mimeType, charset string) {
	if s.viewContentTypes == nil {
		s.viewContentTypes = make(map[string]string)
	}
	s.viewContentTypes[viewID] = mimeType
	if charset != "" {
		s.viewContentTypes[viewID] += "; "+charset
	}
}

func (s *HTTPServer) GetStatistics() stat.SrvStater {
	return s.Statistics
}

//mimetype default is GetMimeTypeOnFileExt()
//contDisposition (attachment|inline) default is attachment
func DownloadFile(resp *response.Response, sock socket.ClientSocketer, f *os.File, fName, mimetype, contDisposition string) error {
	sock_http, ok := sock.(*HTTPSocket)
	if !ok {
		return errors.New("sock must be *HTTPSocket")
	}

	if mimetype == "" {
		mimetype = GetMimeTypeOnFileExt(fName)
	}
	if contDisposition == "" {
		contDisposition = "attachment"
	}
	
	file_info, _ := f.Stat()
	f_size := file_info.Size()
	f_mod := file_info.ModTime()
	
	buffer := make([]byte, f_size)
	f.Read(buffer)
	f_bytes := bytes.NewReader(buffer) // converted to io.ReadSeeker type
	
	sock_http.Response.Header().Set("Content-Type", mimetype)
	sock_http.Response.Header().Set("Content-Disposition", contDisposition+";filename="+fName)
	http.ServeContent(sock_http.Response, sock_http.Request, fName, f_mod, f_bytes)
	resp = nil
	
	return nil
}

/*
func file_exists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil || !os.IsNotExist(err) {
		return true
	}
	return false
}
*/

