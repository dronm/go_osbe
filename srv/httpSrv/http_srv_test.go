package httpSrv

import(
	"testing"
	//"net/http"
	
	"osbe/srv"
	"osbe/socket"
	
	"github.com/labstack/gommon/log"
)

func TestSrvStart(t *testing.T) {

	addr := "localhost:7070"
	app_id := "Test"
	
	//client_cmd := localhost:7070/?c=Test_Controller&f=get_data&v=ViewJSON
	
	logger := log.New("-")
	logger.SetHeader("${time_rfc3339_nano} ${short_file}:${line} ${level} -${message}")
	logger.SetLevel(log.DEBUG)

	http_srv := &HTTPServer{srv.BaseServer{Address: addr,
		Logger: logger,
		AppID: app_id,
	}}
		
	http_srv.OnHandleRequest = func(serv srv.Server, sock socket.ClientSocketer, controllerID string, methodID string, queryID string, payload []byte){
		s.SendToClient(sock, []byte(fmt.Sprintf(`<h1>Test page</h1>
			<div>controllerID=<strong>%s</strong></div>
			<div>methodID=<strong>%s</strong></div>
			<div>queryID=<strong>%s</strong></div>
			<div>Other arguments=<strong>%s</strong></div>
			`, controllerID, methodID, queryID, string(payload))))
	}
	
	http_srv.Run()
}	

