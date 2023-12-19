package socket

import(
	"sync"
	"time"
	"net"
		
	"osbe/sql"	
	"github.com/dronm/session"
)

//Interface for client sockets
type ClientSocketer interface {
	GetID() string
	GetDescr() string
	Close()
	GetConn() net.Conn	
	SetToken(string)
	GetToken() string
	SetTokenExpires(time.Time)
	GetTokenExpires() time.Time
	UpdateLastActivity()
	SetSession(session.Session)
	GetSession() session.Session
	GetPacketID() uint32
	GetIP() string
	GetPresetFilter(string) sql.FilterCondCollection
	SetPresetFilter(PresetFilter) error
	GetLastActivity() time.Time
}

//Structure for managing client sockets
type ClientSocketList struct {
	mx sync.RWMutex
	m []ClientSocketer 
}

func (l *ClientSocketList) Append(socket ClientSocketer){
	l.mx.Lock()	
	l.m = append(l.m, socket)
	l.mx.Unlock()
}

func (l *ClientSocketList) Remove(socket ClientSocketer){
	l.mx.Lock()
	for i, v := range l.m {
		if v == socket {
			l.m[i] = l.m[len(l.m) - 1]
			l.m[len(l.m) - 1] = nil
			l.m = l.m[:len(l.m) - 1]
			break
		}
	}
	l.mx.Unlock()
}

func (l ClientSocketList) Len() int{
	l.mx.Lock()
	defer l.mx.Unlock()
	return len(l.m)
}

// Iterates over the events in the concurrent slice
func (l *ClientSocketList) Iter() <-chan ClientSocketer {
	c := make(chan ClientSocketer)
	f := func() {
		l.mx.Lock()
		defer l.mx.Unlock()
		for _, v := range l.m {
			c <- v
		}
		close(c)
	}
	go f()
	return c
}

func NewClientSocketList() *ClientSocketList{
	return &ClientSocketList{m: make([]ClientSocketer, 0)}
}

