package evnt

import(
	"sync"
	"net"
	"time"

	"osbe/socket"
)

type AddEventProto = func (string, chan string)
type RemoveEventProto = func (string, chan string)

type EventServer interface {
	CloseSocket(*EvntSocket)
}

//unique events for a socket
type uniqSocketEvents struct {
	mx sync.Mutex
	m map[string]bool
}

func (sev *uniqSocketEvents) AddEvent(eventID string, qChan chan string, addEvent AddEventProto) {
	sev.mx.Lock()
	if ok := sev.m[eventID]; !ok {
		sev.m[eventID] = true
		addEvent(eventID, qChan)
	}
	sev.mx.Unlock()	
}
func (sev *uniqSocketEvents) RemoveEvent(eventID string, qChan chan string, removeEvent RemoveEventProto) {
	sev.mx.Lock()
	if ok := sev.m[eventID]; ok {
		delete(sev.m, eventID)
		removeEvent(eventID, qChan)
	}	
	sev.mx.Unlock()	
}

func (sev *uniqSocketEvents) HasEvent(eventID string) bool{
	sev.mx.Lock()
	defer sev.mx.Unlock()	
	_,ok := sev.m[eventID]
	return ok
}

// Iterates over the events
func (sev *uniqSocketEvents) Iter() <-chan string {
	c := make(chan string)

	f := func() {
		sev.mx.Lock()
		defer sev.mx.Unlock()
		for k, v := range sev.m {
			if v {
				c <- k
			}
		}
		close(c)
	}
	go f()

	return c
}


//Socket for one client events
type EvntSocket struct {
	socket.ClientSocket
	Events *uniqSocketEvents //unique events for one client connection	
	Srv EventServer
}

func (s EvntSocket) GetDescr() string {
	return s.Conn.RemoteAddr().String()
}

func (s *EvntSocket) Close() {
	s.Srv.CloseSocket(s)
	s.Conn.Close()
}

func (s *EvntSocket) GetConn() net.Conn{
	return s.Conn
}

func (s EvntSocket) GetToken() string{
	return s.Token
}

//id string, ID: id,
func NewClientSocket(conn net.Conn, token string, tokenExp time.Time, srv EventServer) socket.ClientSocketer{
	return &EvntSocket{socket.ClientSocket{			
			Conn: conn,
			Token: token,
			TokenExpires: tokenExp,
			StartTime: time.Now(),
		},
		&uniqSocketEvents{m: make(map[string]bool)},
		srv,
	}	
}
