package test_app

import(
	"fmt"
	"testing"
	"encoding/json"
	//"context"
	
	"osbe"	
	"osbe/config"	
	"osbe/srv"
	"osbe/srv/wsSrv"
	"osbe/srv/httpSrv"
	"osbe/srv/tcpSrv"
	"osbe/evnt"
	"osbe/socket"
	"osbe/permission"
	"osbe/response"
	"osbe/session"
	"osbe/db"
	_ "osbe/session/pg"
	
	_ "osbe/view/json"
	_ "osbe/view/xml"
	
	"github.com/labstack/gommon/log"
	//"github.com/jackc/pgx/v4/pgxpool"
)

func initApp() *osbe.Application {
	app := &osbe.Application{}
	app.Config = &config.AppConfig{}
	err := app.Config.ReadConf("test_app.json")
	if err != nil {
		panic(fmt.Sprintf("ReadConf: %v",err))
	}

	//Db
	db_conf := app.Config.GetDb()
	app.ServerPool = db.NewServerPool(db_conf.Primary, nil, db_conf.Secondaries)

	//metadata
	app.MD= &osbe.Metadata{Controllers: make(osbe.ControllerCollection)}
	app.MD.Controllers["Test"] = &Test_Controller{}
	app.MD.Controllers["Test"].InitPublicMethods()
		
	app.Logger = log.New("-")
	app.Logger.SetHeader("${time_rfc3339_nano} ${short_file}:${line} ${level} -${message}")
	app.Logger.SetLevel(app.GetLogLevel())
		
	srv_pool := app.GetServerPool()
	pr_pool, err := srv_pool.GetPrimary()
	if err != nil {
		panic(fmt.Sprintf("srv_pool.GetPrimary: %v", err))
	}
	
	app.PermisManager = permission.NewManager(pr_pool.Pool)
	sess_conf := app.Config.GetSession()
	app.SessManager, err = session.NewManager("pg", sess_conf.MaxLifeTime, sess_conf.MaxIdleTime, pr_pool.Pool, sess_conf.EncKey)
	if err != nil {
		panic(fmt.Sprintf("session.NewManager: %v", err))
	}
//	defer app.SessManager.SessionClose(currentSession.SessionID())
	
	//Event server
	app.EvntServer = &evnt.EvntSrv{
		DbPool: pr_pool.Pool,
		Logger: app.Logger,
		WaitBeforeReconnectMS: 1000,
		OnHandleRequest: (func(a osbe.Applicationer) srv.OnHandleRequestProto{
			return func(serv srv.Server, sock socket.ClientSocketer, controllerID string, methodID string, queryID string, argsPayload []byte, viewID string){
				a.HandleRequest(serv, sock, controllerID, methodID, queryID, argsPayload, viewID)
			}
		})(app),
	}
	
	return app
}

func runCMD(app *osbe.Application, cmd string, t *testing.T) string{
	app.Logger.Debugf("Running cmd: %s", cmd)
	contr, pm, argv, query_id, err := app.MD.Controllers.ParseJSONCommand([]byte(cmd))
	if err != nil {
		t.Fatalf("Fatal: %v", err)
	}	
	err = osbe.ValidateExtArgs(app, pm, contr, argv)
	if err != nil {
		t.Errorf("Validation error: %v", err)
	}
	
	resp := response.NewResponse(query_id, "1.001")
	err = pm.Run(app, nil, nil, resp, argv)
	if err != nil {
		t.Fatalf("Fatal: %v", err)
	}	
	
	if resp.GetModelCount() > 1 {
		msg, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Fatal: %v", err)
		}
		return string(msg)
	}	
	
	return ""	
}

//tests web socket server with one Text_Controller
func TestWSServer(t *testing.T) {
	
	app := initApp()
		
	ws_server := &wsSrv.WSServer{}
	ws_server.Address = app.Config.GetWSServer()
	ws_server.TlsCert = app.Config.GetTLSCert()
	ws_server.TlsKey = app.Config.GetTLSKey()
	ws_server.TlsAddress = app.Config.GetTLSWSServer()
	ws_server.AppID = app.Config.GetAppID()
	ws_server.Logger = app.Logger
	ws_server.OnHandleJSONRequest = app.HandleJSONRequest
	ws_server.OnConstructSocket = socket.NewClientSocket
	
	app.AddServer(ws_server)
	ws_server.Run()
}

//tests web server with one Text_Controller
//http://localhost:59000/?c=Test_Controller&f=get_object&id=15
func TestWebServer(t *testing.T) {
	
	app := initApp()
	
	http_server := &httpSrv.HTTPServer{}
	http_server.Address = app.Config.GetWSServer()
	http_server.AppID = app.Config.GetAppID()
	http_server.Logger = app.Logger
	http_server.UnloggedDefaulController = "Logger"
	
	http_server.OnHandleRequest = (func(a osbe.Applicationer) srv.OnHandleRequestProto{
		return func(serv srv.Server, sock socket.ClientSocketer, controllerID string, methodID string, queryID string, args []byte, viewID string){
			a.HandleRequest(serv, sock, controllerID, methodID, queryID, args, viewID)
		}
	})(app)
		
	http_server.OnHandleSession = (func(a osbe.Applicationer) srv.OnHandleSessionProto{
		return func(sock socket.ClientSocketer) error{
			return a.HandleSession(sock)
		}
	})(app)
	
	http_server.OnHandlePermission = (func(a osbe.Applicationer) srv.OnHandlePermissionProto{
		return func(sock socket.ClientSocketer, controllerID string, methodID string) error{
			return a.HandlePermission(sock, controllerID, methodID)
		}
	})(app)
	
	http_server.OnHandleServerError = (func(a osbe.Applicationer) srv.OnHandleServerErrorProto{
		return func(serv srv.Server, sock socket.ClientSocketer, queryID string, viewID string){
			a.HandleServerError(serv, sock, queryID, viewID)
		}
	})(app)
	
	http_server.OnHandleProhibError = (func(a osbe.Applicationer) srv.OnHandleProhibErrorProto{
		return func(serv srv.Server, sock socket.ClientSocketer, queryID string, viewID string){
			a.HandleProhibError(serv, sock, queryID, viewID)
		}
	})(app)
	
	app.AddServer(http_server)
	http_server.Run()
}

//tests parser
func TestParseExtCmd(t *testing.T) {
	
	app := initApp()
	
	cmd := `{"func":"Test.insert"}`	
	_, pm, argv, queryId, err := app.MD.Controllers.ParseJSONCommand([]byte(cmd))
	if err != nil {
		panic(fmt.Sprintf("Fatal: %v", err))
	}
	app.Logger.Debug("Simple parsing, no validation")
	app.Logger.Debugf("Pm=%v, argv=%v, queryId=%s, err=%v", pm, argv, queryId, err)	
}

//tests parser Requirted argument is missing
func TestValidateExtCmd1(t *testing.T) {	
	
	app := initApp()
	
	app.Logger.Debug("1) Requirted argument is missing")
	cmd_fail1 := `{"func":"Test.get_object"}`	
	contr, pm, argv, queryId, err := app.MD.Controllers.ParseJSONCommand([]byte(cmd_fail1))
	if err != nil {
		t.Fatalf("Fatal: %v", err)
	}	
	err = osbe.ValidateExtArgs(app, pm, contr, argv)
	if err != nil {
		t.Errorf("Validation error: %v", err)
	}
	
	app.Logger.Debugf("Pm=%v, argv=%v, queryId=%s, err=%v", pm, argv, queryId, err)	
}

//tests parser With one required argument
func TestValidateExtCmd2(t *testing.T) {	
	
	app := initApp()
	
	//2) 
	app.Logger.Debug("2) With one required argument")
	cmd_fail2 := `{"func":"Test.test_func", "argv":{"a":"Какой-то текст"}}`	
	contr, pm, argv, queryId, err := app.MD.Controllers.ParseJSONCommand([]byte(cmd_fail2))
	if err != nil {
		t.Fatalf("Fatal: %v", err)
	}	
	err = osbe.ValidateExtArgs(app, pm, contr, argv)
	if err != nil {
		t.Errorf("Validation error: %v", err)
	}
		
	app.Logger.Debugf("Pm=%v, argv=%v, queryId=%s, err=%v", pm, argv, queryId, err)	
}

//runs external command
func TestRunExtCmd(t *testing.T) {		
	app := initApp()
	
	cmd := `{"func":"Test.insert", "argv":{"f1":1, "f3":"Text"}}`	
	runCMD(app, cmd, t)
}

//runs external command with undefined controller
func TestRunNoCtrlExtCmd(t *testing.T) {		
	app := initApp()
	
	cmd := `{"func":"Test.get_list"}`	
	runCMD(app, cmd, t)
}
//runs external command with undefined method
func TestRunNoMethodExtCmd(t *testing.T) {		
	app := initApp()
	
	cmd := `{"func":"Test.get_object", "argv":{"id":15}}`	
	runCMD(app, cmd, t)
}

//runs external command expecting result
func TestRunWithResultExtCmd(t *testing.T) {		
	app := initApp()
	
	cmd := `{"func":"Test.get_list", "argv":{"id":18}}`	
	msg := runCMD(app, cmd, t)
	
	app.Logger.Debugf("Server response: %v", string(msg))
}

func TestTCPServer(t *testing.T) {
	
	app := initApp()
	
	server := &tcpSrv.TCPServer{}
	server.Address = app.Config.GetWSServer()
	server.AppID = app.Config.GetAppID()
	server.Logger = app.Logger
	server.OnConstructSocket = socket.NewClientSocket
	
	server.OnHandleJSONRequest = (func(a osbe.Applicationer) srv.OnHandleJSONRequestProto{
		return func(serv srv.Server, sock socket.ClientSocketer, args []byte, viewID string){
			a.HandleJSONRequest(serv, sock, args, viewID)
		}
	})(app)
		
	server.OnHandleSession = (func(a osbe.Applicationer) srv.OnHandleSessionProto{
		return func(sock socket.ClientSocketer) error{
			return a.HandleSession(sock)
		}
	})(app)
	
	server.OnHandlePermission = (func(a osbe.Applicationer) srv.OnHandlePermissionProto{
		return func(sock socket.ClientSocketer, controllerID string, methodID string) error{
			return a.HandlePermission(sock, controllerID, methodID)
		}
	})(app)
	
	server.OnHandleServerError = (func(a osbe.Applicationer) srv.OnHandleServerErrorProto{
		return func(serv srv.Server, sock socket.ClientSocketer, queryID string, viewID string){
			a.HandleServerError(serv, sock, queryID, viewID)
		}
	})(app)
	
	server.OnHandleProhibError = (func(a osbe.Applicationer) srv.OnHandleProhibErrorProto{
		return func(serv srv.Server, sock socket.ClientSocketer, queryID string, viewID string){
			a.HandleProhibError(serv, sock, queryID, viewID)
		}
	})(app)
	
	app.AddServer(server)
	server.Run()
}

