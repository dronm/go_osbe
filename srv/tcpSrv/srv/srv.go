package main

import (
	"fmt"
	
	"osbe"	
	"osbe/config"	
	"osbe/srv"
	"osbe/srv/tcpSrv"
	"osbe/evnt"
	"osbe/socket"
	"osbe/permission"
	"osbe/session"
	"osbe/db"
	_ "osbe/session/pg"	
	_ "osbe/view/json"
	
	"github.com/labstack/gommon/log"

)

func initApp() *osbe.Application {
	app := &osbe.Application{}
	app.Config = &config.AppConfig{}
	err := app.Config.ReadConf("srv.json")
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

func main() {
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
