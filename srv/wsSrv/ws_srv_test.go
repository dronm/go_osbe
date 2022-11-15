package wsSrv

import(
	"testing"
	"context"
	"time"
	
	"osbe/srv"
	
	"github.com/labstack/gommon/log"
	"github.com/gobwas/ws"	
	"github.com/gobwas/ws/wsutil"
)

func TestSrvStart(t *testing.T) {

	addr := "localhost:57000"
	app_id := "Test"
	
	logger := log.New("-")
	logger.SetHeader("${time_rfc3339_nano} ${short_file}:${line} ${level} -${message}")
	logger.SetLevel(log.DEBUG)

	ws_srv := &WSServer{srv.BaseServer{Address: addr,
		Logger: logger,
		AppID: app_id,
		OnHandleRequest:func(*socket.ClientSocket, payload []byte){
			logger.Debugf("OnHandleRequest payload=%s", string(payload))
			
			//SendToClient(sock *socket.ClientSocket, msg string)
		},
	},
	nil}
	
	go ws_srv.Run()
	time.Sleep(time.Duration(1)*time.Second)
	
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "ws://"+addr+"/"+app_id)	
	if err != nil {
		t.Fatalf("unexpected error: %v;\n", err)
		
	}
	payload := []byte("hello, server!")			
	
	err = wsutil.WriteClientMessage(conn, ws.OpText, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v;\n", err)
	}
	
	time.Sleep(time.Duration(2)*time.Second)	
		
	logger.Debugf("Client count: %d", ws_srv.Statistics.GetClientCount())
	logger.Debugf("MaxClient count: %d", ws_srv.Statistics.GetMaxClientCount())
	logger.Debugf("downloadedBytes count: %d", ws_srv.Statistics.GetDownloadedBytes())
	
	logger.Debugf("Total run seconds: %d", ws_srv.Statistics.GetRunSeconds())
	time.Sleep(time.Duration(2)*time.Second)
	logger.Debugf("Total run seconds: %d", ws_srv.Statistics.GetRunSeconds())	
}	

