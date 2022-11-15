package session

import (
	"crypto/rand"
	"fmt"	
	"sync"
	"time"
)

type Session interface {
	Set(key string, value interface{}) error//set session value
	Get(key string) interface{}		//get session value
	GetBool(key string) bool		//get bool session value
	GetString(key string) string		//get string session value
	GetInt(key string) int64		//get int64 session value
	Delete(key string) error     		//delete session value
	SessionID() string           		//back current sessionID
	Flush() error
	TimeCreated() time.Time
}

type Provider interface {
	InitProvider(provParams []interface{}) error
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionClose(sid string) error
	SessionGC(maxLifeTime int64,maxIdleTime int64)
}

var provides = make(map[string]Provider)

// Register makes a session provider available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provide Provider) {
	if provide == nil {
		panic("session: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provide
}

type Manager struct {
	lock        sync.Mutex 
	provider    Provider
	maxLifeTime int64
	maxIdleTime int64
}

func NewManager(provideName string, maxLifeTime int64, maxIdleTime int64, provParams ...interface{}) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	manager := &Manager{provider: provider, maxLifeTime: maxLifeTime, maxIdleTime: maxIdleTime}
	if len(provParams) > 0 {
		er := manager.provider.InitProvider(provParams)
		if er != nil {
			return nil,er
		}
	}
	return manager, nil
}

//open Session
func (manager *Manager) SessionStart(sid string) (session Session, er error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	
	if sid == "" {
		sid := manager.genSessionID()
		session, er = manager.provider.SessionInit(sid)
	} else {
		session, er = manager.provider.SessionRead(sid)
	}
	return
}

//close Session
func (manager *Manager) SessionClose(sid string) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if sid != "" {
		return manager.provider.SessionClose(sid)
	}
	return nil
}

func (manager *Manager) InitProvider(provParams []interface{}) error {
	return manager.provider.InitProvider(provParams)
}

//Destroy sessionid
func (manager *Manager) SessionDestroy(sid string) {
	if sid == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(sid)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime,manager.maxIdleTime)
	time.AfterFunc(time.Duration(manager.maxIdleTime)*time.Second, func() { manager.GC() })
}

func (manager *Manager) genSessionID() string {
	/*b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
	*/
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
	    return ""
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid	
}

