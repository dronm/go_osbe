package wsSrv

import (
	"net"
	"net/http"
	"crypto/tls"
	"net/url"
	"time"
	"strings"
	"io"
	
	"osbe/srv"
	"osbe/tokenBlock"
	"osbe/socket"
	"osbe/stat"
	"osbe/view/json"
	
	"github.com/gobwas/ws"	
	"github.com/gobwas/httphead"
	"github.com/gobwas/ws/wsutil"
)

/**
 * Returns json to the client
 * View parameter is not considered
 */

type OnCloseSocketProto = func(socket.ClientSocketer)

type WSServer struct {
	srv.BaseServer
	Statistics stat.SrvStater
	OnCloseSocket OnCloseSocketProto
}

func (s *WSServer) Run() {
	var err error
	var ln net.Listener

	if s.OnHandleJSONRequest == nil {
		s.Logger.Fatal("WSServer.OnHandleJSONRequest not defined")
	}
	if s.OnConstructSocket == nil {
		s.Logger.Fatal("WSServer.OnConstructSocket not defined")
	}
	if s.OnHandleSession == nil {
		s.Logger.Fatal("WSServer.OnHandleSession not defined")
	}
	if s.OnDestroySession == nil {
		s.Logger.Fatal("WSServer.OnDestroySession not defined")
	}
	
	if s.OnHandleServerError == nil {
		s.Logger.Fatal("WSServer.OnHandleServerError not defined")
	}

	//TLS if nedded
	tls_start := (s.TlsAddress != "" && s.TlsCert != "" && s.TlsKey != "")
	ws_start := (s.Address!= "")
	
	if tls_start {
		var ln_sec net.Listener
		var cer tls.Certificate
	
		cer, err = tls.LoadX509KeyPair(s.TlsCert, s.TlsKey)
		if err != nil {
			s.Logger.Fatalf("tls.LoadX509KeyPair: %v",err)
		}
		
		config := &tls.Config{
			Certificates: []tls.Certificate{cer},
		}
		ln_sec, err = tls.Listen("tcp", s.TlsAddress, config)
	
		if err != nil {
			s.Logger.Fatalf("tls.Listen: %v",err)
		}
		
		s.Logger.Infof("Starting secured websocket server: %s", s.TlsAddress)
		
		if !ws_start {
			//main loop
			s.listenLoop(ln_sec);	
		}else{
			//2 servers
			go s.listenLoop(ln_sec);	
		}
	}
	
	
	if ws_start {
		ln, err = net.Listen("tcp", s.Address)
		if err != nil {
			s.Logger.Fatalf("net.Listen: %v",err)
		}
		
		s.Logger.Infof("Starting websocket server: %s", s.Address)
		
		s.listenLoop(ln);	
	}
}

func (s *WSServer) listenLoop(ln net.Listener) {
	defer ln.Close()
	
	// Prepare handshake header writer from http.Header mapping.
	header := ws.HandshakeHeaderHTTP(http.Header{
		"X-AppSrv-Version": []string{"1.0"},
		//"X-Go-Version": []string{runtime.Version()},
	})	
	
	s.BlockedTokens = tokenBlock.NewTokenBlockList()
	s.ClientSockets = socket.NewClientSocketList()
	s.Statistics = stat.NewSrvStat()
	//Token cam be send
	//URL= CLIENT_APP_ID/TOKEN
	for {		
		conn, err := ln.Accept()
		if err != nil {
			s.Logger.Errorf("ln.Accept: %v",err)
			continue
		}
		
		//, client_origin
		var client_uri, conn_token, client_app_id string
		var conn_token_exp time.Time
		
		//struct description: https://github.com/gobwas/ws/blob/master/server.go
		u := ws.Upgrader{
			OnHeader: func(key, value []byte) error {
				
				//if string(key) == "Origin" {
				//	client_origin = string(value)
				
				//}else
				if string(key) != "Cookie" {
					return nil
				}
				ok := httphead.ScanCookie(value, func(key, value []byte) bool {
					if conn_token=="" && string(key)=="token"{
						conn_token,_ = url.QueryUnescape(string(value))
						
					}else if string(key)=="tokenExpires" {
						exp_str,_ := url.QueryUnescape(string(value))
						conn_token_exp, err = time.Parse(srv.TIME_LAYOUT, exp_str)
						if err != nil {
							s.Logger.Errorf("tokenExpires time.Parse: %s %v", exp_str, err)
							return false
						}
						
					}else if client_app_id=="" && string(key)=="appId" {
						client_app_id,_ = url.QueryUnescape(string(value))
					}
					
					return true
				})
				if ok {
					return nil
				}
				s.Logger.Error("u.Upgrade bad cookie")
				return ws.RejectConnectionError(
					ws.RejectionReason("bad cookie"),
					ws.RejectionStatus(400),
				)
			},
			OnBeforeUpgrade: func() (ws.HandshakeHeader, error) {				
				if client_app_id == "" && len(client_uri)>=2 {
					client_app_id = client_uri[1:]
					p := strings.Index(client_app_id,"/")
					if p >= 0 {
						//token
						conn_token = client_app_id[p+1:]
						client_app_id = client_app_id[:p]
					}
				}
				
				if client_app_id != "" && client_app_id != s.AppID {
					s.Logger.Warnf("u.Upgrade client_app_id %s <> s.AppID %s", client_app_id, s.AppID)
					return nil, ws.RejectConnectionError(
						ws.RejectionReason("bad appID"),
						ws.RejectionStatus(401),
					)
				}
				
				/*if app.AllowedOrigin != "" && client_origin != "" && app.AllowedOrigin != client_origin {
					return nil, ws.RejectConnectionError(
						ws.RejectionReason(ER_ACCESS_DENIED),
						ws.RejectionStatus(403),
					)
				}*/
				
				if conn_token != "" && s.BlockedTokens.Contains(conn_token) {
					s.Logger.Warnf("u.Upgrade BlockedTokens.Contains: %s", conn_token)
					return nil, ws.RejectConnectionError(
						ws.RejectionReason(srv.ER_ACCESS_DENIED),
						ws.RejectionStatus(401),
					)
				}
						
				return header, nil
			},
			OnRequest: func(uri []byte) error {
				
				client_uri,_ = url.QueryUnescape(string(uri))
				
				return nil
			},
		}
		
		_, err = u.Upgrade(conn)		
		if err != nil {
			s.Logger.Errorf("u.Upgrade: %v", err)
			conn.Close()
			continue
		}

		/*id := ""
		id, err = srv.GenID()
		if err != nil {
			s.Logger.Errorf("srv.GenID: %v", err)
			conn.Close()
			continue
		}
		*/
		socket := s.OnConstructSocket(conn, conn_token, conn_token_exp)
		s.ClientSockets.Append(socket)
		
		go s.HandleConnection(socket)
				
	}	
}

func (s *WSServer) HandleConnection(socket socket.ClientSocketer) {
	s.Logger.Debugf("Got connection from: %s", socket.GetDescr())
	
	s.Statistics.IncHandshakes()
	
	defer s.CloseSocket(socket)

	//check token for validity and reject if not valid with error code
	//clientToken can be empty in case of publish event
	//@toDo: handle this case
	/*
	token := socket.GetToken()
	if token != "" && s.OnCheckToken != nil{
		err := s.OnCheckToken(token)
		if err != nil {
			s.BlockedTokens.Append(token)
			s.Logger.Warnf("%s, %v", socket.GetDescr(), err)				
			
			if s.OnAccessDenied != nil {
				s.OnAccessDenied(socket)
			}
			return
		}
	}
	*/
	
	//session
	err := s.OnHandleSession(socket)
	if err != nil {
		s.Logger.Errorf("WSServer.HandleConnection OnHandleSession(): %v", err)
		s.OnHandleServerError(s, socket, "", viewJSON.VIEW_ID)//error always in JSON
		return
	}
	
	for {	
		conn := socket.GetConn()
		if conn == nil {
			return;
		}
		
		header, err := ws.ReadHeader(conn)
		
		socket.UpdateLastActivity()
		
		switch err {
		case nil:			
			payload := make([]byte, header.Length)
			_, err = io.ReadFull(conn, payload)
			if err != nil {
				s.Logger.Errorf("%s, io.ReadFull: %v", socket.GetDescr(), err)
				return
			}
			if header.Masked {
				ws.Cipher(payload, header.Mask, 0)
			}
			if len(payload)<12 {
				//min query={"func":"a"}
				s.Logger.Errorf("Wrong payload: %s", payload)
				return
			}
			
			s.Logger.Debugf("%s, Data payload: %s", socket.GetDescr(), string(payload))
			
			s.Statistics.IncDownloadedBytes(uint64(len(payload)))
			
			go s.OnHandleJSONRequest(s, socket, payload, viewJSON.VIEW_ID)//default view if no v param in query
			
			if header.OpCode == ws.OpClose {					
				return
			}
			
		
		case io.EOF:
			s.Logger.Warnf("%s, Closed on timeout", socket.GetDescr())
			return
			
		default:
			s.Logger.Errorf("%s, conn.Read: %v", socket.GetDescr(), err)
			return
		}		
	}
}


func (s *WSServer) SendToClient(sock socket.ClientSocketer, msg []byte) error {
//s.Logger.Debugf("WSServer SendToClient sock.ID=%s, msg=%s", sock.GetID(), string(msg))	
	err := wsutil.WriteServerText(sock.GetConn(), msg)
	if err != nil {
		s.Logger.Errorf("%s, wsutil.WriteServerText: %v", sock.GetDescr(), err)
		return err
	}	
	s.Statistics.IncUploadedBytes(uint64(len(msg)))
	return nil
}

func (s *WSServer) CloseSocket(sock socket.ClientSocketer){
	id := sock.GetID()
	token := sock.GetToken()
	if s.OnCloseSocket != nil {
		s.OnCloseSocket(sock)
	}
	//s.OnDestroySession(sock.GetToken())
	s.ClientSockets.Remove(sock)
	s.Statistics.OnClientDisconnceted()
	s.Logger.Debugf("CloseSocket, id=%s, token = %s, socket count:%d", id, token, s.ClientSockets.Len())	
}

func (s *WSServer) GetClientSockets() *socket.ClientSocketList{
	return s.ClientSockets
}

func (s *WSServer) GetStatistics() stat.SrvStater {
	return s.Statistics
}
