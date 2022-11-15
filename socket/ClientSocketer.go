package socket

import (
	"net"
	"time"
	
	"session"
	"osbe/sql"
)
	
//Interface for client sockets
type ClientSocketer interface {
	GetDescr() string
	Close()
	GetConn() net.Conn	
	SetToken(string)
	GetToken() string
	SetTokenExpires(time.Time)
	GetTokenExpires() time.Time
	//GetID() string
	GetDemandLogout() chan bool
	UpdateLastActivity()
	SetSession(session.Session)
	GetSession() session.Session
	GetPacketID() uint32
	GetIP() string
	GetPresetFilter(string) sql.FilterCondCollection
	SetPresetFilter(PresetFilter) error
	GetLastActivity() time.Time
}

