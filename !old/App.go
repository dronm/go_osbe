package osbe

import (
	"errors"
	"reflect"
	"time"
	"fmt"
	"context"
	"os"
	"os/exec"
	"bytes"
	"io"
	
	"osbe/session"	
	"osbe/socket"
	"osbe/srv"
	"osbe/view"
	"osbe/response"
	"osbe/db"

	"github.com/labstack/gommon/log"
	
	//"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DEFAULT_ROLE = "guest"
	SESS_ROLE = "USER_ROLE"
	SESS_LOCALE = "USER_LOCALE"
	
	FRAME_WORK_VERSION = "1.0.0.1"
)

//Not used mooved to db/pool
//type OnDbNotificationProto = func(*pgconn.PgConn, *pgconn.Notification)

type OnPublishEventProto = func(string, string)

type Applicationer interface {
	GetConfig() AppConfiger
	GetLogger() *log.Logger
	GetMD() *Metadata
	GetServer(string) srv.Server
	GetServers() ServerList
	GetServerPool() *db.ServerPool
	SendToClient(srv.Server, socket.ClientSocketer, *response.Response, string) error
	HandleRequest(srv.Server, socket.ClientSocketer, string, string, string, []byte, string)
	HandleJSONRequest(srv.Server, socket.ClientSocketer, []byte, string)
	HandleSession(socket.ClientSocketer) error
	HandlePermission(socket.ClientSocketer, string, string) error
	HandleServerError(srv.Server, socket.ClientSocketer, string, string)
	HandleProhibError(srv.Server, socket.ClientSocketer, string, string)
	GetPrimaryPoolConn() (*pgxpool.Conn, *PublicMethodError)
	GetSecondaryPoolConn() (*pgxpool.Conn, *PublicMethodError)
	GetSessManager() *session.Manager
	XSLTransform([]byte, string, string, string) ([]byte, error)	
	GetFrameworkVersion() string
	GetPermisManager() Permissioner
	GetOnPublishEvent() OnPublishEventProto
}

type ServerList map[string]srv.Server

type Application struct {
	Config AppConfiger
	Logger *log.Logger
	ServerPool *db.ServerPool
	SessManager *session.Manager
	MD *Metadata
	Servers ServerList
	PermisManager Permissioner
	OnPublishEvent OnPublishEventProto
}

func (a *Application) GetConfig()  AppConfiger{
	return a.Config
}

func (a *Application) GetLogger()  *log.Logger{
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

func (a *Application) GetServerPool() *db.ServerPool{
	return a.ServerPool
}

func (a *Application) HandleSession(sock socket.ClientSocketer) error {
	if a.SessManager == nil {
		return nil
	}
	//session
	sess, err := a.SessManager.SessionStart(sock.GetToken())
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

		err := ValidateExtArgs(a, pm, contr, argv)
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
			
			if !a.Config.GetReportErrors() && err_code == response.RESP_ER_INTERNAL{
				//log real error, short to client
				a.Logger.Errorf("Application.Run() failed %s.%s: %d:%s", contr.GetID(), pm.GetID(), err_code, err_txt)			
				err_txt = ER_PM_INTERNAL
			}
			
			resp.SetError(err_code, err_txt)
			
			a.SendToClient(serv, sock, resp, viewID)
			return
			
		}else if err != nil {
			return
		}
	}	
	if serv != nil && (resp.GetQueryID() != "" || resp.GetModelCount() > 1) {
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
		a.Logger.Errorf("Application.HandleLocalEvent ParseFunctionCommand() failed: %v", err)
		return
	}
	if err := pm.Run(a, nil, nil, nil, argv); err != nil {
		a.Logger.Errorf("Application.HandleLocalEvent Run() failed: %v, %s.%s", err, contr.GetID(), pm.GetID())
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
		a.Logger.Errorf("Application.HandleRequest NewResponse() failed: %v", err)
		a.SendToClient(serv, sock, resp, viewID)
		return
		
	}else if err != nil {
		a.Logger.Errorf("Application.HandleRequest ParseCommand() failed: %v", err)
		return
	}
//a.Logger.Debug("Application.HandleRequest is called")	
	a.HandleRequestCont(serv, sock, pm, contr, argv, resp, viewID)
}

//Parsing incoming arguments, Controller method calling
//payload contains json request
func (a *Application) HandleJSONRequest(serv srv.Server, sock socket.ClientSocketer, payload []byte, viewID string) {	
	contr, pm, argv, query_id, err := a.MD.Controllers.ParseJSONCommand(payload)	
	resp := response.NewResponse(query_id, a.MD.Version.Value)
	if serv != nil && err != nil {
		resp.SetError(response.RESP_ER_PARSE, err.Error())
		a.Logger.Errorf("Application.HandleJSONRequest NewResponse() failed: %v", err)
		a.SendToClient(serv, sock, resp, viewID)
		return
		
	}else if err != nil {
		return
	}
a.Logger.Debug("Application.HandleJSONRequest is called")	
	a.HandleRequestCont(serv, sock, pm, contr, argv, resp, viewID)
}

func (a *Application) SendToClient(serv srv.Server, sock socket.ClientSocketer, resp *response.Response, viewID string) error {
	msg, err := view.Render(viewID, sock, resp)	
	if err != nil {
		a.Logger.Errorf("Application.Render() failed: %v", err)
		msg = []byte(err.Error())
	}
	err = serv.SendToClient(sock, msg)
	return err 
}

func (a *Application) GetLogLevel() log.Lvl {
	var lvl log.Lvl

	switch a.Config.GetLogLevel() {
	case "debug":
		lvl = log.DEBUG
		break
	case "info":
		lvl = log.INFO
		break
	case "warn":
		lvl = log.WARN
		break
	case "error":
		lvl = log.ERROR
		break
	default:
		lvl = log.INFO
	}
	return lvl
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

func (a *Application) GetTempDir() string {
	return "/tmp"
}

//Make read from stdin
//now used from util_xml && html.go
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
				return nil, errors.New(string(out_b))
				//errors.New(errb.String())
			}		
		}else{
			out_b, err = cmd.Output()		
			if err != nil {
				return nil, errors.New(string(out_b))
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

