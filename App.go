package osbe

import (
	"errors"
	"reflect"
	"time"
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"io"
	"encoding/json"
	"path/filepath"
	"io/ioutil"
	"strings"
	//"runtime/debug"
		
	"osbe/socket"
	"osbe/srv"
	"osbe/view"
	"osbe/response"
	"osbe/logger"
	"osbe/fields"
	"osbe/config"
	
	"github.com/dronm/session"
)

const (
	DEFAULT_ROLE = "guest"		// default user role if not logged
	SESS_ROLE = "USER_ROLE"		// session variable name
	SESS_LOCALE = "USER_LOCALE"	// session variable name
	
	VERSION_FILE_NAME = "version.txt"// for storing varion number
	
	FRAME_WORK_VERSION = "1.0.0.10"	// framwork version
)

var DebugQueries bool // if true, all queries will be logged

type OnPublishEventProto = func(string, string)
type OnReloadConfigProto = func()

// ServerList is a list of running servers.
// Defined on application startup.
type ServerList map[string]srv.Server

// Application is the main object holding
// application parameters. It is passed to all
// controllers. Derived applications must
// include this structure to their app objects.
type Application struct {
	Config AppConfiger			// app config
	Logger logger.Logger			// app logger
	SessManager *session.Manager		// app sessions
	MD *Metadata				// app description of all controllers/models
	Servers ServerList			// list of running servers
	PermisManager Permissioner		// handles permission rules (controllers to roles)
	OnPublishEvent OnPublishEventProto	
	DataStorage interface{}
	EncryptKey string			// application encryption key
	BaseDir string				// application directory
	ConfigFileName string			//
	OnReloadConfig OnReloadConfigProto
	AppVersion string
}

// GetConfig returns application config.
func (a *Application) GetConfig()  AppConfiger{
	return a.Config
}

// GetConfig returns application logger.
func (a *Application) GetLogger()  logger.Logger{
	return a.Logger
}

// GetConfig returns application metadata.
func (a *Application) GetMD()  *Metadata{
	return a.MD
}

// GetConfig returns a running service by its name.
func (a *Application) GetServer(ID string)  srv.Server{
	if s, ok := a.Servers[ID]; ok {
		return s
	}
	return nil
}

// GetPermisManager returns application permission manager.
func (a *Application) GetPermisManager()  Permissioner{
	return a.PermisManager
}

// GetConfig returns a list of running service.
func (a *Application) GetServers()  ServerList{
	return a.Servers
}

// GetOnPublishEvent returns on publish event function.
func (a *Application) GetOnPublishEvent()  OnPublishEventProto{
	return a.OnPublishEvent
}

// AddServer adds a new service to the list.
func (a *Application) AddServer(ID string, s srv.Server) {
	if a.Servers == nil {
		a.Servers = make(ServerList)
	}
	a.Servers[ID] = s
}

// GetSessManager returns application session manager.
func (a *Application) GetSessManager() *session.Manager{
	return a.SessManager
}

// HandleSession is run on a new client call. It finds an appropriate
// session by client token or creates a new one if it is dead.
func (a *Application) HandleSession(sock socket.ClientSocketer) error {
	if a.SessManager == nil {
		return nil
	}
	//session
	tok := sock.GetToken()
	if  tok!= "" && len(tok) > a.SessManager.GetSessionIDLen() {
		//wrong token
		sock.SetToken("")
	}
	
	//preset filter type for serialization
	socket.RegisterPresetFilter()
	
	sess, err := a.SessManager.SessionStart(tok)
	if err != nil {
		return err
	}
	sock.SetSession(sess)
	if sock.GetToken() == "" {
		//session not found: create a new one
		sock.SetToken(sess.SessionID())
		sock.SetTokenExpires( sess.TimeCreated().Add(time.Second * time.Duration(a.Config.GetSession().MaxLifeTime)) )
		//default role
		err := sess.Set(SESS_ROLE, DEFAULT_ROLE)
		if err != nil {
			return err
		}
	}
			
	return nil
}

// DestroySession deletes a session by its ID.
func (a *Application) DestroySession(sessID string) {
	a.GetLogger().Debugf("DestroySession sessID=%s", sessID)
	a.GetSessManager().SessionDestroy(sessID)
}

// HandlePermission is called to determine permission for Controller->method. If not permitted
// an error is returned.
// TODO: it is better to return bool value instead.
func (a *Application) HandlePermission(sock socket.ClientSocketer, controllerID string, methodID string) error {
	if a.PermisManager == nil {
		return nil
	}
	sess := sock.GetSession()	
	role_id := sess.GetString(SESS_ROLE)
	if !a.PermisManager.IsAllowed(role_id, controllerID, methodID) {
		a.GetLogger().Errorf("Method '%s.%s' not allowed for role '%s'", controllerID, methodID, role_id)
		return errors.New(ER_COM_METH_PROHIB)
	}
	return nil
}

// HandleServerError is called on any server error. The error is sent to client 
// with a specified view as an internal server error.
func (a *Application) HandleServerError(serv srv.Server, sock socket.ClientSocketer, queryID string, viewID string){
	resp := response.NewResponse(queryID, a.MD.Version.Value)
	resp.SetError(response.RESP_ER_INTERNAL, ER_INTERNAL)
	a.SendToClient(serv, sock, resp, viewID)
}

//TODO: whe is it called, from server?
// HandleProhibError is called on not allowed error.
func (a *Application) HandleProhibError(serv srv.Server, sock socket.ClientSocketer, queryID string, viewID string){
	resp := response.NewResponse(queryID, a.MD.Version.Value)
	resp.SetError(response.RESP_ER_INTERNAL, ER_COM_METH_PROHIB)
	a.SendToClient(serv, sock, resp, viewID)
}

// HandleRequestCont continues handling a client request. Checks permission for controller->method
// for a given role.
func (a *Application) HandleRequestCont(serv srv.Server, sock socket.ClientSocketer, pm PublicMethod, contr Controller, argv reflect.Value, resp *response.Response, viewID string) {
	if contr != nil && pm != nil {
		//permission
		if sock != nil {
			err := a.HandlePermission(sock, string(contr.GetID()), string(pm.GetID()))
			if serv != nil && err != nil {
				resp.SetError(response.RESP_ER_AUTH, ER_COM_METH_PROHIB)
				a.SendToClient(serv, sock, resp, viewID)
				//Block!
				return
				
			}else if err != nil {
				return
			}
		}
		
		//validate client arguments
		err := a.validateExtArgs(pm, contr, argv)
		if serv != nil && err != nil {
			resp.SetError(response.RESP_ER_VALID, err.Error())
			a.SendToClient(serv, sock, resp, viewID)
			return
		}
		//run controller method.
		err = pm.Run(a, serv, sock, resp, argv)
		if serv != nil && err != nil {
			var err_code int
			var err_txt string
			if pm_err, ok := err.(*PublicMethodError); ok {
				err_code = pm_err.Code
				err_txt = pm_err.Err.Error()
			}else{			
				err_code = response.RESP_ER_INTERNAL
				err_txt = err.Error()			
			}
			
			//log real error
			a.Logger.Errorf("Application.Run() %s.%s: %d:%s", contr.GetID(), pm.GetID(), err_code, err_txt)			
			
			if !a.Config.GetReportErrors() && err_code == response.RESP_ER_INTERNAL{
				//short to client				
				err_txt = ER_PM_INTERNAL
			}
			
			resp.SetError(err_code, err_txt)
			
			a.SendToClient(serv, sock, resp, viewID)
			return
			
		}else if err != nil {
			return
		}
	}	
	if serv != nil && resp != nil && (resp.GetQueryID() != "" || resp.GetModelCount() > 1) {
		//response is expected
		a.SendToClient(serv, sock, resp, viewID)
	}
}

// HandleEvent is an event handler. Called from event server.
// Event is in format "ControllerID.MethodID".
// No response is expected, so the function returns none.
// All errors are logged.
func (a *Application) HandleEvent(fn string, args []byte) {
	contr, pm, argv, err := a.MD.Controllers.ParseFunctionCommand(fn, args)
	if err != nil {
		a.Logger.Errorf("Application.HandleLocalEvent ParseFunctionCommand() failed: %v", err)
		return
	}
	if err := pm.Run(a, nil, nil, nil, argv); err != nil {
		a.Logger.Errorf("Application.HandleLocalEvent Run() failed: %v, %s.%s", err, contr.GetID(), pm.GetID())
	}
}

// HandleRequest handles net request. It has controller and method as separated parameters.
// If error occurs when parsing command it is sent to client.
// Otherwise, on success handling is continued in HandleRequestCont() function 
func (a *Application) HandleRequest(serv srv.Server, sock socket.ClientSocketer, controllerID string, methodID string, queryID string, argsPayload []byte, viewID string) {
	if a.Config.GetLogLevel() == "debug" {
		a.Logger.Debugf("HTTPServer HandleRequest(): controllerID=%s, methodID=%s, queryID=%s, argsPayload=%s, viewID=%s", controllerID, methodID, queryID, string(argsPayload), viewID)
	}
	
	var contr Controller
	var pm PublicMethod
	var argv reflect.Value
	var err error
	if controllerID !="" && methodID !="" {
		contr, pm, argv, err = a.MD.Controllers.ParseCommand(controllerID, methodID, argsPayload)	
	}
	resp := response.NewResponse(queryID, a.MD.Version.Value)
	if serv != nil && err != nil {
		resp.SetError(response.RESP_ER_PARSE, err.Error())
		a.Logger.Errorf("Application.HandleRequest ControllerCollection.ParseCommand(): %v", err)
		a.SendToClient(serv, sock, resp, viewID)
		return
		
	}else if err != nil {
		a.Logger.Errorf("Application.HandleRequest ControllerCollection.ParseCommand(): %v", err)
		return
	}
	a.HandleRequestCont(serv, sock, pm, contr, argv, resp, viewID)
}

// HandleJSONRequest handles requests in json format. Payload argument contains json data.
// It parses incoming arguments, then on success HandleRequestCont() is called.
func (a *Application) HandleJSONRequest(serv srv.Server, sock socket.ClientSocketer, payload []byte, viewID string) {	
	contr, pm, argv, query_id, view_id, err := a.MD.Controllers.ParseJSONCommand(payload)	
	if view_id == "" {
		view_id = viewID
	}
	resp := response.NewResponse(query_id, a.MD.Version.Value)
	if serv != nil && err != nil {
		resp.SetError(response.RESP_ER_PARSE, err.Error())
		a.Logger.Errorf("Application.HandleJSONRequest NewResponse() failed: %v", err)
		a.SendToClient(serv, sock, resp, view_id)
		return
		
	}else if err != nil {
		return
	}
	a.HandleRequestCont(serv, sock, pm, contr, argv, resp, view_id)
}

// SendToClient sends back to client response object rendered with a specific View.
// TODO: measure request time?
func (a *Application) SendToClient(serv srv.Server, sock socket.ClientSocketer, resp *response.Response, viewID string) error {
	msg, err := view.Render(viewID, sock, resp)	
	if err != nil {
		a.Logger.Errorf("Application.Render() failed: %v", err)
		//debug.PrintStack()
		msg = []byte(err.Error())
	}
	err = serv.SendToClient(sock, msg)
//a.GetLogger().Debugf("Query execution time: %v", time.Since(sock.GetLastActivity()))	
	return err 
}

// GetDataStorage returns application data storage object.
func (a *Application) GetDataStorage() interface{}{
	return a.DataStorage
}

func (a *Application) GetEncryptKey() string{
	return a.EncryptKey
}
func (a *Application) GetBaseDir() string{
	return a.BaseDir
}

// GetTempDir returns application temp directory.
// TODO: more platform independant way.
func (a *Application) GetTempDir() string {
	return "/tmp"
}

// XSLTransform applies XSLT rules from file xslFileName to data byte slice
// or data read from inFileName argument if given an uninitialized slice
// and inFileName is not an empty string.
// TODO:Make read from stdin.
// If outFileName is set then output will be written to this file.
// Othervise []byte will be returned.
// This function is used from util_xml && view/html.go, view/excel.go
// The function uses xalan as an XSL transformer.
// TODO: make a possibility to choose a transformer. Put transformer
// arguments some where else.
func (a *Application) XSLTransform(data []byte, inFileName string, xslFileName string, outFileName string) ([]byte, error) {
	//if a.OnXSLTransform
	//xalan transformation by default
	
	if (data == nil && inFileName == "") || xslFileName == "" {
		return nil, errors.New(ER_XSL_TRANSFORM)
	}	
	
	var out_b []byte
	var errb bytes.Buffer
	params := []string{"-q", "-xsl", xslFileName}
	
	if outFileName != "" {
		params = append(params, "-out", outFileName)
	}	
	if data == nil {
		params = append(params, "-in", inFileName)
	}
	cmd := exec.Command("xalan", params...)
	cmd.Stderr = &errb

	if data != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil { 
			return nil, err
		}

		go func() {
			defer stdin.Close()
			io.Copy(stdin, bytes.NewReader(data))
		}()

		if outFileName != "" {
			err := cmd.Run()
			if err != nil { 
				return nil, err
				//errors.New(string(out_b))
			}		
		}else{
			out_b, err = cmd.Output()		
			if err != nil {
				return nil, err
				//errors.New(string(out_b))
			}
		}

	}else{
		err := cmd.Run()
		if err != nil { 
			return nil, errors.New(string(out_b))
			//errors.New(errb.String())
		}
	}
	
	if outFileName != "" {
		_, err := os.Stat(outFileName)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
			
		}else if err != nil {
			return nil, errors.New(ER_XSL_TRANSFORM+" "+errb.String())
		}
		return nil, nil
	}else{
		return out_b, nil
	}
}

// XSLToPDFTransform transforms given data in form of a byte slice or from inFileName.
// Styles from xslFileName are applied.
// Transformed result is saved to file if outFileName is not empty or is returned
// as slice of bytes.
// The function uses Apache FOP as a transformer. 
// TODO: Make read from stdin.
// This function is used from view/pdf.go
func (a *Application) XSLToPDFTransform(fop string, confFile string, params []string, data []byte, inFileName string, xslFileName string, outFileName string) ([]byte, error) {
	//fop transformation by default	
	if (data == nil && inFileName == "") || xslFileName == "" {
		return nil, errors.New(ER_XSL_TRANSFORM)
	}	
	
	if params == nil || len(params) == 0 {
		params = []string{"-q"} //default param
	}
	if confFile != "" {
		params = append(params, "-c", confFile)
	}
	if fop == "" {
		fop = "fop"
	}
	
	var out_b []byte
	var errb bytes.Buffer
	
	if outFileName != "" {
		params = append(params, "-pdf", outFileName)
	}else{
		params = append(params, "-pdf", "-")
	}	
	if data == nil {
		params = append(params, "-xml", inFileName)
	}else{
		params = append(params, "-xml", "-")
	}
	params = append(params, "-xsl", xslFileName)
	
	a.Logger.Debugf("XSLToPDFTransform: %s %v\n", fop, params)	
	
	cmd := exec.Command(fop, params...)
	cmd.Stderr = &errb

	if data != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil { 
			return nil, err
		}

		go func() {
			defer stdin.Close()
			io.Copy(stdin, bytes.NewReader(data))
		}()

		if outFileName != "" {
			err := cmd.Run()
			if err != nil { 
				return nil, err
				//errors.New(string(out_b))
			}		
		}else{
			out_b, err = cmd.Output()		
			if err != nil {
				return nil, err
				//errors.New(string(out_b))
			}
		}

	}else{
		err := cmd.Run()
		if err != nil { 
			return nil, errors.New(string(out_b))
			//errors.New(errb.String())
		}
	}
	
	if outFileName != "" {
		_, err := os.Stat(outFileName)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
			
		}else if err != nil {
			return nil, errors.New(ER_XSL_TRANSFORM+" "+errb.String())
		}
		return nil, nil
	}else{
		return out_b, nil
	}
}

// GetFrameworkVersion returns current framwork version number.
func (a *Application) GetFrameworkVersion() string {
	return FRAME_WORK_VERSION
}

// validateExtArgs is a private function for validating external user arguments.
func (a *Application) validateExtArgs(pm PublicMethod, contr Controller, argv reflect.Value) error {

	md_model := pm.GetFields()
	if md_model == nil {
		return nil
	}
	
	//combines all errors in one string	
	var valid_err strings.Builder
	
	var arg_fld reflect.Value
	
	var argv_empty = argv.IsZero()
	
	for fid, fld := range md_model {
		var arg_fld_v reflect.Value		
		if !argv_empty {
			//Indirect always returns object!
			arg_fld = reflect.Indirect(argv).FieldByName(fid)
			if arg_fld.Kind() == reflect.Struct {
				arg_fld_v = arg_fld.FieldByName("TypedValue")
			}
		}
		
		if !argv_empty && arg_fld_v == (reflect.Value{}) {
			//custom structure
			if fld.GetRequired() && arg_fld.IsZero() {
				appendError(&valid_err, fmt.Sprintf(ER_PARSE_NOT_VALID_EMPTY, fld.GetDescr()) ) 
			}// or no validation here
		
		//GetRequired is implemented by all fields
		}else if fld.GetRequired() && (argv_empty || (arg_fld.IsValid() && arg_fld.Kind() == reflect.Struct && (!arg_fld.FieldByName("IsSet").Bool() || arg_fld.FieldByName("IsNull").Bool()) ) ) {
			//required field has no value
			appendError(&valid_err, fmt.Sprintf(ER_PARSE_NOT_VALID_EMPTY, fld.GetDescr()) ) 
			
		}else if !argv_empty && arg_fld.IsValid() && arg_fld.Kind() == reflect.Struct {
			//check if metadata field implements certain interfaces
			//if it does, call methods of these interfaces
			//fmt.Printf("fid=%s, arg_fld=%v\n",fid, arg_fld)	
			
			var err error			
			switch fld.GetDataType() {
			case fields.FIELD_TYPE_FLOAT:
				err = fields.ValidateFloat(fld.(fields.FielderFloat), arg_fld_v.Float())				
				
			case fields.FIELD_TYPE_INT:
				err = fields.ValidateInt(fld.(fields.FielderInt), arg_fld_v.Int())				
				
			case fields.FIELD_TYPE_TEXT:
				err = fields.ValidateText(fld.(fields.FielderText), arg_fld_v.String())				

			case fields.FIELD_TYPE_JSON:
				err = fields.ValidateJSON(fld.(fields.FielderJSON), []byte(arg_fld_v.String()))


			case fields.FIELD_TYPE_TIME:
				err = fields.ValidateTime(fld.(fields.Fielder), arg_fld_v.String())				

			case fields.FIELD_TYPE_DATE:
				err = fields.ValidateDate(fld.(fields.Fielder), arg_fld_v.String())				

			case fields.FIELD_TYPE_DATETIME:
				err = fields.ValidateDateTime(fld.(fields.Fielder), arg_fld_v.String())				

			case fields.FIELD_TYPE_DATETIMETZ:
				err = fields.ValidateDateTimeTZ(fld.(fields.Fielder), arg_fld_v.String())				

			case fields.FIELD_TYPE_ENUM:
				err = fields.ValidateEnum(fld.(fields.FielderEnum), arg_fld_v.String())				
				
			/*default:
				appendError(&valid_err, "osbe.ValidateExtArgs: unsupported field type" ) 
			*/
			}
			if err != nil {
				appendError(&valid_err, err.Error() ) 
			}
		//}else if !argv_empty {
			//field is present in ext argg but is not in metadata
		//	a.GetLogger().Warnf("External argument %s is not present in metadata of %s.%s", fid, contr.GetID(), pm.GetID())
			//fmt.Println("Field",fid, "arg_fld=",arg_fld)
		//}else{
			//fmt.Println("Otherwise")
		}
		
		//fmt.Println("Field",fid,"IsSet=",arg_fld.FieldByName("IsSet"),"IsNull=",arg_fld.FieldByName("IsNull"),"Value=",arg_fld.FieldByName("TypedValue"))
	}
	
	if valid_err.Len() > 0 {
		return errors.New(valid_err.String())
	}
	
	return nil
}

// PublishPublicMethodEvents publishes events from public method if there are any.
func (a *Application) PublishPublicMethodEvents(pm PublicMethod, params map[string]interface{}) {
	//params["lsn"] = (SELECT pg_current_wal_lsn())
	//SELECT pg_notify('%s','%s') ev_id, params_s
	on_ev := a.GetOnPublishEvent()
	if on_ev != nil {
		l := pm.GetEventList()
		if l != nil {
			params_s := "null"
			if params != nil && len(params) > 0 {				
				if par, err := json.Marshal(params); err == nil {
					params_s = string(par)
				}
			}
			for _, ev_id := range l {				
				if ev_id != "" {
					on_ev(ev_id, `"params":`+params_s)
				}
			}
		}
	}
}

// ReloadAppConfig reloads application configuration from file.
// After reload OnReloadConfig() is called if it points to a function.
func (a *Application) ReloadAppConfig() error {	
	if a.ConfigFileName == "" {
		return errors.New(ER_CONFIG_FILE_NOT_DEFINED)
	}
	if err := config.ReadConf(a.ConfigFileName, a.Config); err != nil {
		return err
	}
	if a.OnReloadConfig != nil {
		a.OnReloadConfig()
	}
	return nil
}

// LoadAppVersion loads application version from file defined by
// VERSION_FILE_NAME constant. File is searched in base directory.
func (a *Application) LoadAppVersion() error {	
	f_n := filepath.Join(a.BaseDir, VERSION_FILE_NAME)
	ver, err := ioutil.ReadFile(f_n)
	if err != nil {
		return err
	}
	if len(ver) == 0 {
		return errors.New(ER_VERSION_FILE_EMPTY)
	}
	if []rune(string(ver))[len(ver)-1] == 10 {
		ver = ver[0:len(ver)-1]
	}
	a.AppVersion = string(ver)
	
	if a.Logger != nil {
		a.Logger.Warnf("Version file loaded, current version: %s", a.AppVersion)
	}
	return nil
}

// GetAppVersion returns application version. Version is
// retrieved from file. It panics if version file is not found.
// Should be called on application startup. 
func (a *Application) GetAppVersion() string {	
	if a.AppVersion == "" {
		if err := a.LoadAppVersion(); err != nil {
			err_s := fmt.Sprintf("LoadAppVersion() failed: %v", err)
			if a.Logger != nil {
				a.Logger.Error(err_s)
				return ""
			}else{
				panic(err_s)
			}		
		}
	}
	return a.AppVersion
}

