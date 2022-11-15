package osbe

import(
	"fmt"
	"testing"
	"context"
	
	"osbe/config"	
	"osbe/srv"
	"osbe/srv/wsSrv"
	"osbe/evnt"
	"osbe/socket"
	//"session"
	//_ "session/pg"
	
	"github.com/labstack/gommon/log"
	"github.com/jackc/pgx/v4/pgxpool"
)

func construct_md() *Metadata{

	md := osbe.Metadata{Controllers: make(ControllerCollection)
	}
	md.Controllers["Test"] = &Test_Controller{}
	md.Controllers["Test"].InitPublicMethods()
	
	return md
}

func initApp() *Application {
	app := &Application{}
	app.Config = &config.AppConfig{}
	err := app.Config.ReadConf("test_app.json")
	if err != nil {
		panic(fmt.Sprintf("ReadConf: %v",err))
	}

	//app.MD := construct_md()

	err = app.InitDb()
	//app.DbConn, err = pgxpool.Connect(context.Background(), app.Config.GetStorageConnection())
	if err != nil {
		panic(fmt.Sprintf("app.InitDb: %v", err))
	}

	app.Logger = log.New("-")
	app.Logger.SetHeader("${time_rfc3339_nano} ${short_file}:${line} ${level} -${message}")
	app.Logger.SetLevel(app.GetLogLevel())
	
	return app
}

//tests web socket server
func TestWSServer(t *testing.T) {
	
	/*
	ENC_KEY := "2fgV65sh465Fh4054Nhryn4dty54nH5d4j6G41c356j4T5h1g76dtYj0"
	app.SessManager, err = session.NewManager("pg", "token", 3600, 3600, app.DbConn, ENC_KEY)
	if err != nil {
		panic(fmt.Sprintf("session.NewManager: %v", err))
	}
	*/
	app := initApp()	
	ws_server := &wsSrv.WSServer{}
	ws_server.Address = app.Config.GetWSServer()
	ws_server.TlsCert = app.Config.GetTLSCert()
	ws_server.TlsKey = app.Config.GetTLSKey()
	ws_server.TlsAddress = app.Config.GetTLSWSServer()
	ws_server.AppID = app.Config.GetAppID()
	ws_server.Logger = app.Logger
	ws_server.OnCheckToken = func(token string) error{
		return app.CheckToken(token)
	}
	ws_server.OnHandleRequest = (func(s *wsSrv.WSServer) srv.OnHandleRequestProto{
		return func(sock socket.ClientSocketer, payload []byte){
			app.HandleRequest(s, sock, payload)
		}
	})(ws_server)
	ws_server.OnConstructSocket = evnt.NewClientSocket
	
	app.AddServer(ws_server)
	ws_server.Run()
}

//tests parser
func TestParseExtCmd(t *testing.T) {
	cmd := `{"func":"Test_Controller.get_list"}`	
	
	
}
