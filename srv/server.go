package srv

import (
	//"fmt"
	//"crypto/rand"
	"net"
	"time"
	
	"osbe/tokenBlock"
	"osbe/socket"
	"osbe/logger"
)

const TIME_LAYOUT = "2006-01-02T15:04:05.000-07"

//Here controller and method as separate strings, patlad is a string {argv:[]}
type OnHandleRequestProto = func(Server, socket.ClientSocketer, string, string, string, []byte, string)

//here payload is a json string {"func":"Controller.method","queryId":"",argv:[]}
type OnHandleJSONRequestProto = func(Server, socket.ClientSocketer, []byte, string)

type OnHandleSessionProto = func(socket.ClientSocketer) error
type OnDestroySessionProto = func(string)

type OnConstructSocketProto = func(net.Conn, string, time.Time) socket.ClientSocketer

type OnHandleServerErrorProto = func(Server, socket.ClientSocketer, string, string)
//type OnHandleProhibErrorProto = func(Server, socket.ClientSocketer, string, string)

//type OnHandlePermissionProto = func(socket.ClientSocketer, string, string) error

type Server interface {
	Run() //params ...interface{}
	SendToClient(socket.ClientSocketer, []byte) error
	GetClientSockets() *socket.ClientSocketList
}

type BaseServer struct {
	TlsCert string
	TlsKey string
	TlsAddress string
	Address string //Host:Port
	Logger logger.Logger
	AppID string
	BlockedTokens *tokenBlock.TokenBlockList
	ClientSockets *socket.ClientSocketList
	
	OnHandleRequest OnHandleRequestProto
	OnHandleJSONRequest OnHandleJSONRequestProto
	OnHandleSession OnHandleSessionProto
	OnDestroySession OnDestroySessionProto
	OnConstructSocket OnConstructSocketProto
	//OnHandlePermission OnHandlePermissionProto
	
	OnHandleServerError OnHandleServerErrorProto
	//OnHandleProhibError OnHandleProhibErrorProto
}

// generates unique ID SESSION
/*
func GenID() (string,error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
	    return "",err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid,nil
}
*/
