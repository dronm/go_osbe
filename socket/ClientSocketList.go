package socket

import(
	"sync"
)
/*
type ClientSocketItem struct {
	ID string
	Socket ClientSocketer
}
*/
//Structure for managing client sockets
type ClientSocketList struct {
	mx sync.RWMutex
	m map[string]ClientSocketer //client connections: key=token		
}

func (l *ClientSocketList) Append(socket ClientSocketer){
	l.mx.Lock()
	l.m[socket.GetToken()] = socket
	l.mx.Unlock()
}
func (l *ClientSocketList) Remove(id string){
	l.mx.Lock()
	if _,ok := l.m[id]; ok {
		delete(l.m, id) 
	}
	l.mx.Unlock()
	return
}
/*
func (l *ClientSocketList) Get(id string) ClientSocketer {
	l.mx.Lock()
	defer l.mx.Unlock()
	
	if sock,ok := l.m[id]; ok {
		return sock
	}
	return nil
}
*/
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
			c <-v
			//ClientSocketItem{v}//k, 
		}
		close(c)
	}
	go f()

	return c
}

func NewClientSocketList() *ClientSocketList{
	return &ClientSocketList{m: make(map[string]ClientSocketer)}
}

