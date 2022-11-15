package socket

import (
	"net"
	"time"
	"sync"
	"strings"	
	"errors"
//	"fmt"	
	
	"session"
	"osbe/sql"
)

type ClientSocket struct {
	//ID string
	DemandLogout chan bool
	Conn net.Conn
	mx sync.RWMutex
	PacketID uint32
	Token string
	TokenExpires time.Time		
	LastActivity time.Time
	StartTime time.Time
	Session session.Session
}

func (s ClientSocket) GetDescr() string {
	return s.Conn.RemoteAddr().String()
}

func (s *ClientSocket) Close() {
	s.Conn.Close()
}

func (s *ClientSocket) GetConn() net.Conn{
	return s.Conn
}

/*func (s ClientSocket) GetID() string{
	return s.ID
}*/

func (s ClientSocket) GetDemandLogout() chan bool{
	return s.DemandLogout
}

func (s *ClientSocket) UpdateLastActivity(){
	s.LastActivity = time.Now()
}

func (s *ClientSocket) SetToken(token string){
	s.mx.Lock()
	s.Token = token
	s.mx.Unlock()
}
func (s *ClientSocket) GetToken() string {
	return s.Token
}
func (s *ClientSocket) SetTokenExpires(t time.Time) {
	s.mx.Lock()
	s.TokenExpires = t
	s.mx.Unlock()
}
func (s *ClientSocket) GetTokenExpires() time.Time {
	return s.TokenExpires
}

func (s *ClientSocket) GetSession() session.Session{
	return s.Session
}

func (s *ClientSocket) SetSession(sess session.Session){
	s.Session = sess
}

func (s *ClientSocket) SetPresetFilter(f PresetFilter) error {
	sess := s.GetSession()
	if sess != nil {
		//for session serialization
		//registerPresetFilter()
	
		sess.Set(SESS_PRESET_FILTER, f)
		return sess.Flush()
	}
	return errors.New("Session not defined")
}

func (s *ClientSocket) GetPresetFilter(modelID string) sql.FilterCondCollection {
	sess := s.GetSession()
	if sess != nil {
		//for session serialization
		//registerPresetFilter()
	
		f := sess.Get(SESS_PRESET_FILTER)
//fmt.Printf("ClientSocket.GetPresetFilter=%v\n", f)		
		if v, ok := f.(PresetFilter); ok {
			return v.Get(modelID)
		}
	}
	return nil
}

func (s *ClientSocket) GetIP() string{
	if s.Conn == nil {
		return ""
	}
	return GetRemoteAddrIP(s.Conn.RemoteAddr().String())
	/*addr := s.Conn.RemoteAddr().String()
	if p := strings.Index(addr, ":"); p >= 0 {
		return addr[:p]
	}else{
		return addr
	}*/
}

func GetRemoteAddrIP(remoteAddr string) string{
	if p := strings.Index(remoteAddr, ":"); p >= 0 {
		return remoteAddr[:p]
	}else{
		return remoteAddr
	}
}

func (s *ClientSocket) GetPacketID() uint32{
	s.mx.Lock()
	id := s.PacketID
	s.PacketID++
	s.mx.Unlock()	
	return id
}
func (s *ClientSocket) GetLastActivity() time.Time {
	return s.LastActivity 
}

//*************
//id string, ID: id, 
func NewClientSocket(conn net.Conn, token string, tokenExp time.Time) ClientSocketer{
	return &ClientSocket{Conn: conn, Token: token, TokenExpires: tokenExp, StartTime: time.Now()}
}

