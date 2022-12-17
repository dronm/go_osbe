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
	//"runtime/debug"
	
	"session"		
	"osbe/socket"
	"osbe/srv"
	"osbe/view"
	"osbe/response"
	"osbe/logger"
	"osbe/fields"
)

const (
	DEFAULT_ROLE = "guest"
	SESS_ROLE = "USER_ROLE"
	SESS_LOCALE = "USER_LOCALE"
	
	FRAME_WORK_VERSION = "1.0.0.1"
)

type OnPublishEventProto = func(string, string)


type Applicationer interface {
	GetConfig() AppConfiger
	GetLogger() logger.Logger
	GetMD() *Metadata
	GetServer(string) srv.Server
	GetServers() ServerList
	SendToClient(srv.Server, socket.ClientSocketer, *response.Response, string) error
	HandleRequest(srv.Server, socket.ClientSocketer, string, string, string, []byte, string)
	HandleJSONRequest(srv.Server, socket.ClientSocketer, []byte, string)
	HandleSession(socket.ClientSocketer) error
	DestroySession(sessID string)
	HandlePermission(socket.ClientSocketer, string, string) error
	HandleServerError(srv.Server, socket.ClientSocketer, string, string)
	HandleProhibError(srv.Server, socket.ClientSocketer, string, string)
	GetSessManager() *session.Manager
	XSLTransform([]byte, string, string, string) ([]byte, error)	
	GetFrameworkVersion() string
	GetPermisManager() Permissioner
	GetOnPublishEvent() OnPublishEventProto
	GetDataStorage() interface{}
	PublishPublicMethodEvents(PublicMethod, map[string]interface{})
	GetEncryptKey() string
	GetBaseDir() string
}

type ServerList map[string]srv.Server

type Application struct {
	Config AppConfiger
	Logger logger.Logger
	//ServerPool *db.ServerPool
	SessManager *session.Manager
	MD *Metadata
	Servers ServerList
	PermisManager Permissioner
	OnPublishEvent OnPublishEventProto
	DataStorage interface{}
	EncryptKey string
	BaseDir string
}

func (a *Application) GetConfig()  AppConfiger{
	return a.Config
}

func (a *Application) GetLogger()  logger.Logger{
	return a.Logger
}

func (a *Application) GetMD()  *Metadata{
	return a.MD
}

func (a *Application) GetServer(ID string)  srv.Server{
	if s, ok := a.Servers[ID]; ok {
		return s
	}
	return nil
}

func (a *Application) GetPermisManager()  Permissioner{
	return a.PermisManager
}

func (a *Application) GetServers()  ServerList{
	return a.Servers
}
func (a *Application) GetOnPublishEvent()  OnPublishEventProto{
	return a.OnPublishEvent
}

func (a *Application) AddServer(ID string, s srv.Server) {
	if a.Servers == nil {
		a.Servers = make(ServerList)
	}
	a.Servers[ID] = s
}

func (a *Application) GetSessManager() *session.Manager{
	return a.SessManager
}

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
		//new session
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

func (a *Application) DestroySession(sessID string) {
a.GetLogger().Debugf("DestroySession sessID=%s", sessID)
	a.GetSessManager().SessionDestroy(sessID)
}

func (a *Application) HandlePermission(sock socket.ClientSocketer, controllerID string, methodID string) error{
	if a.PermisManager == nil {
		return nil
	}
	sess := sock.GetSession()	
	if !a.PermisManager.IsAllowed(sess.GetString(SESS_ROLE), controllerID, methodID) {
		return errors.New(ER_COM_METH_PROHIB)
	}
	return nil
}

func (a *Application) HandleServerError(serv srv.Server, sock socket.ClientSocketer, queryID string, viewID string){
	resp := response.NewResponse(queryID, a.MD.Version.Value)
	resp.SetError(response.RESP_ER_INTERNAL, ER_INTERNAL)
	a.SendToClient(serv, sock, resp, viewID)
}

func (a *Application) HandleProhibError(serv srv.Server, sock socket.ClientSocketer, queryID string, viewID string){
	resp := response.NewResponse(queryID, a.MD.Version.Value)
	resp.SetError(response.RESP_ER_INTERNAL, ER_COM_METH_PROHIB)
	a.SendToClient(serv, sock, resp, viewID)
}

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

		err := a.validateExtArgs(pm, contr, argv)
		if serv != nil && err != nil {
			resp.SetError(response.RESP_ER_VALID, err.Error())
			a.SendToClient(serv, sock, resp, viewID)
			return
		}
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

//event is of type ControllerID.MethodID
//no response is expected. All errors are logged
func (a *Application) HandleEvent(fn string, args []byte) {
	contr, pm, argv, err := a.MD.Controllers.ParseFunctionCommand(fn, args)
	if err != nil {
		//log error
		a.Logger.Errorf("Application.HandleLocalEvent ParseFunctionCommand(): %v", err)
		return
	}
	if err := pm.Run(a, nil, nil, nil, argv); err != nil {
		a.Logger.Errorf("Application.HandleLocalEvent Run(): %v, %s.%s", err, contr.GetID(), pm.GetID())
	}
}

//handles html request (controller, method - strings, no need to parse from struct)
func (a *Application) HandleRequest(serv srv.Server, sock socket.ClientSocketer, controllerID string, methodID string, queryID string, argsPayload []byte, viewID string) {
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
//a.Logger.Debug("Application.HandleRequest is called")	
	a.HandleRequestCont(serv, sock, pm, contr, argv, resp, viewID)
}

//Parsing incoming arguments, Controller method calling
//payload contains json request
func (a *Application) HandleJSONRequest(serv srv.Server, sock socket.ClientSocketer, payload []byte, viewID string) {	
	contr, pm, argv, query_id, view_id, err := a.MD.Controllers.ParseJSONCommand(payload)	
	if view_id == "" {
		view_id = viewID
	}
	resp := response.NewResponse(query_id, a.MD.Version.Value)
	if serv != nil && err != nil {
		resp.SetError(response.RESP_ER_PARSE, err.Error())
		a.Logger.Errorf("Application.HandleJSONRequest NewResponse(): %v", err)
		a.SendToClient(serv, sock, resp, view_id)
		return
		
	}else if err != nil {
		return
	}
	a.HandleRequestCont(serv, sock, pm, contr, argv, resp, view_id)
}

func (a *Application) SendToClient(serv srv.Server, sock socket.ClientSocketer, resp *response.Response, viewID string) error {
	msg, err := view.Render(viewID, sock, resp)	
	if err != nil {
		a.Logger.Errorf("Application.Render(): %v", err)
		//debug.PrintStack()
		msg = []byte(err.Error())
	}
	err = serv.SendToClient(sock, msg)
//a.GetLogger().Debugf("Query execution time: %v", time.Since(sock.GetLastActivity()))	
	return err 
}

func (a *Application) GetDataStorage() interface{}{
	return a.DataStorage
}

func (a *Application) GetEncryptKey() string{
	return a.EncryptKey
}
func (a *Application) GetBaseDir() string{
	return a.BaseDir
}

/*func (a *Application) GetServerPool() *db.ServerPool{
	return a.ServerPool
}

func (a *Application) GetPrimaryPoolConn() (*pgxpool.Conn, *PublicMethodError){
	pool,err := a.GetServerPool().GetPrimary()
	if err != nil {
		return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("App.GetPrimaryPoolConn() db.ServerPool.GetPrimary(): %v",err))
	}
	
	pool_conn, err := pool.Pool.Acquire(context.Background())
	if err != nil {
		return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("App.GetPrimaryPoolConn() pgxpool.Pool.Acquire(): %v",err))
	}
	return pool_conn, nil
}

func (a *Application) GetSecondaryPoolConn() (*pgxpool.Conn, *PublicMethodError){
	pool,err := a.GetServerPool().GetSecondary()
	if err != nil {
		return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("db.ServerPool.GetSecondary(): %v",err))
	}
	
	pool_conn, err := pool.Pool.Acquire(context.Background())
	if err != nil {
		return nil, NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgxpool.Pool.Acquire(): %v",err))
	}
	return pool_conn, nil
}
*/
func (a *Application) GetTempDir() string {
	return "/tmp"
}

//Make read from stdin
//now used from util_xml && view/html.go, view/excel.go
//incoming data may come from []byte or can be taken from inFileName
//if outFileName is set then output will be written to this file
//othervise []byte will be returned
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

func (a *Application) GetFrameworkVersion() string {
	return FRAME_WORK_VERSION
}

//External argument validation
func (a *Application) validateExtArgs(pm PublicMethod, contr Controller, argv reflect.Value) error {

	md_model := pm.GetFields()
	if md_model == nil {
		return nil
	}
	
	//combines all errors in one string	
	valid_err := ""
	
	var arg_fld reflect.Value
	
	var argv_empty = argv.IsZero()
	
	for fid, fld := range md_model {
		var arg_fld_v reflect.Value
		
		//fmt.Println("fid=", fid, "GetRequired=", fld.GetRequired(), "argv_empty=", argv_empty, "IsValid=", arg_fld.IsValid())		
		//,"IsSet=", arg_fld.FieldByName("IsSet").Bool(),"IsNull=", arg_fld.FieldByName("IsNull").Bool())
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
			//fmt.Println("!argv_empty && arg_fld.IsValid()")
			
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
	
	if valid_err != "" {
		return errors.New(valid_err)
	}
	
	return nil
}

//notifies all servers through database event
//
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
