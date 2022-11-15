package tcpSrv

import (
	"net"
	"crypto/tls"
	"time"
	"encoding/binary"
	"io"
	
	"osbe/srv"
	"osbe/tokenBlock"
	"osbe/socket"
	"osbe/stat"	
)

const (
	//Pref(2bytes) + length(4bytes) +ID(4bytes) + data(variable) + postf(2bytes)	
	PREF_PACKET_START byte = 0xFF
	PREF_PACKET_LAST byte = 0x0A
	PREF_PACKET_CONT byte = 0x0B
	POSTF_0 byte = 0x0A
	POSTF_1 byte = 0x0D 
	PREF_LEN uint32 = 10
	POSTF_LEN uint32 = 2
	MAX_DATA_LEN uint16 = 65535
	
	DEFAULT_VIEW = "json"
)

type TCPServer struct {
	srv.BaseServer
	Statistics stat.SrvStater
}

func (s *TCPServer) Run() {
	var err error
	var ln net.Listener

	if s.OnHandleJSONRequest == nil {
		s.Logger.Fatal("TCPServer.OnHandleJSONRequest not defined")
	}
	if s.OnConstructSocket == nil {
		s.Logger.Fatal("TCPServer.OnConstructSocket not defined")
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
		
		s.Logger.Infof("Starting secured tcp server: %s", s.TlsAddress)
		
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
		
		s.Logger.Infof("Starting tcp server: %s", s.Address)
		
		s.listenLoop(ln);	
	}
}

func (s *TCPServer) listenLoop(ln net.Listener) {
	defer ln.Close()
	
	s.BlockedTokens = tokenBlock.NewTokenBlockList()
	s.ClientSockets = socket.NewClientSocketList()
	s.Statistics = stat.NewSrvStat()
	
	for {		
		conn, err := ln.Accept()
		if err != nil {
			s.Logger.Errorf("ln.Accept: %v",err)
			continue
		}
		
		id := ""
		id, err = srv.GenID()
		if err != nil {
			s.Logger.Errorf("srv.GenID: %v", err)
			conn.Close()
			continue
		}
		
		socket := s.OnConstructSocket(id, conn, "", time.Time{})
		s.ClientSockets.Append(socket)
		
		go s.HandleConnection(socket)
				
	}	
}

func (s *TCPServer) HandleConnection(socket socket.ClientSocketer) {
	s.Logger.Debugf("Got connection from: %s", socket.GetDescr())
	
	s.Statistics.IncHandshakes()
	
	defer s.CloseSocket(socket)

	//session
	err := s.OnHandleSession(socket)
	if err != nil {
		s.Logger.Errorf("TCPServer.HandleConnection OnHandleSession: %v", err)
		s.OnHandleServerError(s, socket, "", DEFAULT_VIEW)
		return
	}

	conn := socket.GetConn()
	head_b := make([]byte, PREF_LEN)
	for {	
		select {
		case _= <-socket.GetDemandLogout():
			break
		default:
		}
			
		_, err := conn.Read(head_b)			
		switch err {		
		case nil:
		
			s.Logger.Debugf("Got incoming data: %v",head_b)
			
			//prefix check
			if head_b[0] != PREF_PACKET_START || (head_b[1] != PREF_PACKET_LAST && head_b[1] != PREF_PACKET_CONT) {
				//wrong structute
				s.Logger.Errorf("%s, TCPServer.HandleConnection() wrong packet structure", socket.GetDescr())
				return			
			}

			//Packet structure:
			//PREFIX(2 bytes) + data length(2 bytes) + JSON data (=data length), POSTF(2 bytes)
		
			packet_len := binary.LittleEndian.Uint32(head_b[2:6])
			packet_id := binary.LittleEndian.Uint32(head_b[6:10])			
			s.Logger.Debugf("Data length=%d, packetID=%d", packet_len, packet_id)
			payload := make([]byte, packet_len + POSTF_LEN) //Data + postfix
			payload_len := packet_len + POSTF_LEN
			full_payload := ""
			
			for payload_len>0 {
				b_cnt, err := conn.Read(payload)			
				//fmt.Println("ReadBytes=",b_cnt)
				payload_len = payload_len - uint32(b_cnt)
				
				switch err {
				case nil:
					//got message
					if payload_len == 0 {
						//postfix check
						if payload[b_cnt-2] != POSTF_0 || payload[b_cnt-1] != POSTF_1 {
							s.Logger.Errorf("%s, TCPServer.HandleConnection() wrong packet structure", socket.GetDescr())
							return
						}
						full_payload += string(payload[:uint32(b_cnt)-POSTF_LEN])
						
						s.Statistics.IncDownloadedBytes(uint64(PREF_LEN + packet_len + POSTF_LEN))						
						s.Logger.Debugf("Starting message parsing %s", full_payload)
						go s.OnHandleJSONRequest(s, socket, []byte(full_payload), DEFAULT_VIEW)
						
					}else{
						full_payload += string(payload[:b_cnt])
					}
				case io.EOF:
					s.Logger.Warnf("%s, TCPServer.HandleConnection() gracefully closed", socket.GetDescr())
					return
					
				default:
					s.Logger.Errorf("%s, TCPServer.HandleConnection() conn.Read: %v", socket.GetDescr(), err)
					return
				}		
			}					
		case io.EOF:
			s.Logger.Warnf("%s, TCPServer.HandleConnection() gracefully closed", socket.GetDescr())
			return
			
		default:
			s.Logger.Errorf("%s, TCPServer.HandleConnection() conn.Read: %v", socket.GetDescr(), err)
			return
		}		
		
	}
}


func (s *TCPServer) SendToClient(sock socket.ClientSocketer, msg []byte) error {
	conn := sock.GetConn()
	
	packet_id := sock.GetPacketID()
	packet_len := uint32(len(msg))	
	
	bf := make([]byte, PREF_LEN + packet_len + POSTF_LEN)
	bf[0] = PREF_PACKET_START
	bf[1] = PREF_PACKET_LAST
	binary.LittleEndian.PutUint32(bf[2:6], packet_len)
	binary.LittleEndian.PutUint32(bf[6:10], packet_id)		
	copy(bf[PREF_LEN : PREF_LEN+packet_len], msg)
	bf[PREF_LEN+packet_len] = POSTF_0
	bf[PREF_LEN+packet_len+1] = POSTF_1
s.Logger.Debugf("Sending message ID=%d, len=%d, msg=%s", packet_id, packet_len, string(msg))
//s.Logger.Debug(string(msg))		
	//bf := append([]byte{PREF_PACKET_START}, packet_type, packet_len_b, msg[data_offset:data_offset+packet_len], []byte{POSTF_0}, []byte{POSTF_0}...)
	_, err := conn.Write(bf)
	if err != nil {
		s.Logger.Errorf("%s, TCPServer.SendToClient() Write: %v", sock.GetDescr(), err)
		return err		
	}
	
	s.Statistics.IncUploadedBytes(uint64(packet_len + PREF_LEN + POSTF_LEN))
	return nil	
}

func (s *TCPServer) CloseSocket(socket socket.ClientSocketer){
	s.ClientSockets.Remove(socket.GetID())
	socket.Close()
	s.Statistics.OnClientDisconnceted()
}

func (s *TCPServer) GetClientSockets() *socket.ClientSocketList{
	return s.ClientSockets
}
