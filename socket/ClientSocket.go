package socket

import (
	"net"
	"time"
	"sync"
	"strings"	
	"errors"
	"math/rand"
		
	"osbe/sql"
	"github.com/dronm/session"
)

const SOCKET_ID_LEN = 32

type ClientSocket struct {
	ID string	//emitter unique ID
	Conn net.Conn
	mx sync.RWMutex
	PacketID uint32
	Token string
	TokenExpires time.Time		
	LastActivity time.Time
	StartTime time.Time
	Session session.Session
}

func (s *ClientSocket) GetDescr() string {
	return s.Conn.RemoteAddr().String()
}

func (s *ClientSocket) Close() {
	s.mx.Lock()
	s.Conn.Close()
	s.Conn = nil
	s.mx.Unlock()
}

func (s *ClientSocket) GetConn() net.Conn{
	return s.Conn
}

func (s *ClientSocket) UpdateLastActivity(){
	s.mx.Lock()
	s.LastActivity = time.Now()
	s.mx.Unlock()
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
	
		f := PresetFilter{}
		if err := sess.Get(SESS_PRESET_FILTER, &f); err == nil {
			return f.Get(modelID)	
		}		
	}
	return nil
}

func (s *ClientSocket) GetIP() string{
	if s.Conn == nil {
		return ""
	}
	return GetRemoteAddrIP(s.Conn.RemoteAddr().String())
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

func (s *ClientSocket) GetID() string {
	return s.ID
}

//*************
//id string, ID: id, 
func NewClientSocket(conn net.Conn, token string, tokenExp time.Time) ClientSocketer{
	return &ClientSocket{
		ID: GenSocketID(),
		Conn: conn,
		Token: token,
		TokenExpires: tokenExp,
		StartTime: time.Now(),
	}
}

func GenSocketID() string{
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, SOCKET_ID_LEN)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

