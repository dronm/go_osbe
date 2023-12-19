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
//Protocol structure
//1) Header. Fixed size 10 bytes.
//	Pref(2bytes) + length(4bytes) +ID(4bytes) 
//
//2) Data. Variable size:
//	data(variable)
//
//3) Footer. Fixe size 2 bytes: POSTF_0+POSTF_1
//	postf(2bytes)	

const (
	PREF_PACKET_START byte = 0xFF
	PREF_PACKET_LAST byte = 0x0A
	PREF_PACKET_CONT byte = 0x0B
	POSTF_0 byte = 0x0A
	POSTF_1 byte = 0x0D 
	PREF_LEN uint32 = 10
	POSTF_LEN uint32 = 2
	MAX_DATA_LEN uint16 = 65535
	
	DEFAULT_VIEW = "ViewJSON"
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
	if s.OnHandleSession != nil && s.OnHandleServerError == nil {
		s.Logger.Fatal("TCPServer.OnHandleSession defined, but OnHandleServerError not defined")
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

func (s *TCPServer) HandleConnection(sock socket.ClientSocketer) {
	s.Logger.Debugf("Got connection from: %s", sock.GetDescr())
	
	s.Statistics.IncHandshakes()
	
	defer s.CloseSocket(sock)

	//session
	if s.OnHandleSession != nil {
		err := s.OnHandleSession(sock)
		if err != nil {
			s.Logger.Errorf("TCPServer.HandleConnection OnHandleSession: %v", err)
			s.OnHandleServerError(s, sock, "", DEFAULT_VIEW)
			return
		}
	}
	
	conn := sock.GetConn()
	head_b := make([]byte, PREF_LEN)
	s.Logger.Debugf("Starting listen loop for: %s", sock.GetDescr())
	for {	
		select {
		case _= <-sock.GetDemandLogout():
			s.Logger.Error("sock.GetDemandLogout()")
			return
		default:
		}
			
		_, err := conn.Read(head_b)			
		switch err {		
		case nil:
		
			s.Logger.Debugf("Got incoming data for:%s, %v",sock.GetDescr(), head_b)
			
			//prefix check
			if head_b[0] != PREF_PACKET_START || (head_b[1] != PREF_PACKET_LAST && head_b[1] != PREF_PACKET_CONT) {
				//wrong structute
				s.Logger.Errorf("%s, TCPServer.HandleConnection() wrong packet structure", sock.GetDescr())
				return			
			}

			//Packet structure:
			//PREFIX(2 bytes) + data length(2 bytes) + JSON data (=data length), POSTF(2 bytes)
		
			packet_len := binary.LittleEndian.Uint32(head_b[2:6])
			packet_id := binary.LittleEndian.Uint32(head_b[6:10])			
			s.Logger.Debugf("Data length=%d, packetID=%d", packet_len, packet_id)
			
			tot_packet_len := packet_len + POSTF_LEN
			payload := make([]byte, tot_packet_len) //Data + postfix			
			to_read := tot_packet_len
			read_cnt := 0
			//b, err := io.ReadAll(conn)
			var payload_full []byte //на случай, если все не прочиается за раз
			for to_read > 0 {
				b_cnt, err := conn.Read(payload)							
				//s.Logger.Debugf("Read=%d", b_cnt)
				to_read -= uint32(b_cnt)
				
				switch err {
				case nil:
					//got message
					if to_read == 0 {
						//got full message
						//may be one byte long!
						if (b_cnt==1 && (payload[0] != POSTF_1 || payload_full[len(payload_full)-1] != POSTF_0) ) ||
						(payload[b_cnt-2] != POSTF_0 || payload[b_cnt-1] != POSTF_1) {
							s.Logger.Errorf("read wrong packet postfix: %d %d", payload[b_cnt-2], payload[b_cnt-1])
							break
						}
						var response *[]byte
						if read_cnt == 0 {
							//one buffer
							response = &payload
						}else{
							//concat data
							payload_full = append(payload_full, payload[:b_cnt]...)
							response = &payload_full
						}
						s.Statistics.IncDownloadedBytes(uint64(tot_packet_len))
						*response = (*response)[:packet_len]
						s.Logger.Debugf("Starting message parsing %s", string(*response))
						go s.OnHandleJSONRequest(s, sock, *response, DEFAULT_VIEW)
					}else{
						//message part
						if payload_full == nil {
							payload_full = make([]byte, 0)
						}
						//concat data
						payload_full = append(payload_full, payload[:b_cnt]...)
						payload = make([]byte, to_read)
						//s.Logger.Debugf("ToRead=%d", to_read)
					}
				case io.EOF:
					s.Logger.Warnf("%s, TCPServer.HandleConnection() gracefully closed", sock.GetDescr())
					return
					
				default:
					s.Logger.Errorf("%s, TCPServer.HandleConnection() conn.Read: %v", sock.GetDescr(), err)
					return
				}
				read_cnt++		
			}					
		case io.EOF:
			s.Logger.Warnf("%s, TCPServer.HandleConnection() gracefully closed", sock.GetDescr())
			return
			
		default:
			s.Logger.Errorf("%s, TCPServer.HandleConnection() conn.Read: %v", sock.GetDescr(), err)
			return
		}				
	}
}

//@ToDo: send msg []byte directly to conn.Write without copying
func (s *TCPServer) SendToClient(sock socket.ClientSocketer, msg []byte) error {
	conn := sock.GetConn()
	
	packet_id := sock.GetPacketID()
	packet_len := uint32(len(msg))	
	
	bf_header := make([]byte, PREF_LEN + packet_len + POSTF_LEN)
	bf_header[0] = PREF_PACKET_START
	bf_header[1] = PREF_PACKET_LAST
	binary.LittleEndian.PutUint32(bf_header[2:6], packet_len)
	binary.LittleEndian.PutUint32(bf_header[6:10], packet_id)		
	//header
	_, err := conn.Write(bf_header)
	if err != nil {
		s.Logger.Errorf("%s, TCPServer.SendToClient() header Write: %v", sock.GetDescr(), err)
		return err		
	}
	//body
	_, err := conn.Write(msg)
	if err != nil {
		s.Logger.Errorf("%s, TCPServer.SendToClient() body Write: %v", sock.GetDescr(), err)
		return err		
	}
	//footer
	_, err := conn.Write([]byte{POSTF_0, POSTF_1})
	if err != nil {
		s.Logger.Errorf("%s, TCPServer.SendToClient() footer Write: %v", sock.GetDescr(), err)
		return err		
	}
	
	s.Statistics.IncUploadedBytes(uint64(packet_len + PREF_LEN + POSTF_LEN))
	return nil	
}

/*
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
	if packet_len<20000{
		s.Logger.Debugf("Sending message ID=%d, len=%d, msg=%s", packet_id, packet_len, string(msg))
	}else{
		s.Logger.Debugf("Sending BIG message ID=%d, len=%d", packet_id, packet_len)
	}
	_, err := conn.Write(bf)
	if err != nil {
		s.Logger.Errorf("%s, TCPServer.SendToClient() Write: %v", sock.GetDescr(), err)
		return err		
	}
	
	s.Statistics.IncUploadedBytes(uint64(packet_len + PREF_LEN + POSTF_LEN))
	return nil	
}
*/
func (s *TCPServer) CloseSocket(sock socket.ClientSocketer){
	s.ClientSockets.Remove(sock.GetID())
	sock.Close()
	s.Statistics.OnClientDisconnceted()
	s.Logger.Debugf("Socket closed: %s", sock.GetDescr())
}

func (s *TCPServer) GetClientSockets() *socket.ClientSocketList{
	return s.ClientSockets
}

func (s *TCPServer) GetStatistics() stat.SrvStater {
	return s.Statistics
}
