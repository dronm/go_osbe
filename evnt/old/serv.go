package evnt

import (
	"sync"
	"fmt"
	"context"
	"time"
	
	"waitStrat"
	
	"github.com/labstack/gommon/log"
	"github.com/jackc/pgx/v4/pgxpool"	
	"github.com/jackc/pgconn"
)

const (
	CMD_PUBLISH_ARGS = `{"argv":{"id":"%s"%s}}`
	CMD_PUBLISH_FN = "Event.publish"
	QUERY_PAUSE_MS = 50
	
	EVNT_PARAM_EMITTER_ID = "emitterId"		
	EVNT_MODEL_ID = "Event"	
)

type UniqEvents struct {
	mx sync.RWMutex
	m map[string]int //unique event counter
}

func (e *UniqEvents) AddEvent(dbEventID string, qChan chan string){
	e.mx.Lock()
	cnt, ok := e.m[dbEventID]
	if !ok {
		qChan <- `LISTEN "`+dbEventID+`"`
		e.m[dbEventID] = 1
	}else{
		e.m[dbEventID] = cnt + 1
	}
	e.mx.Unlock()
}

func (e *UniqEvents) RemoveEvent(dbEventID string, qChan chan string){
	e.mx.Lock()
	if cnt, ok := e.m[dbEventID]; ok {		
		if cnt == 1 {
			qChan <- `UNLISTEN "`+dbEventID+`"`
			delete(e.m, dbEventID)
		}else{
			e.m[dbEventID] = cnt - 1
		}		
	}
	e.mx.Unlock()
}

func (e *UniqEvents) GetEvent(eventID string) (int, bool) {
	e.mx.Lock()
	defer e.mx.Unlock()

	value, ok := e.m[eventID]
	
	return value, ok
}

func (e *UniqEvents) GetEventCount() (int) {
	e.mx.Lock()
	defer e.mx.Unlock()

	return len(e.m)
}

type OnEventProto = func(string, []byte) //"ControllerID.MethodId", arguments

//event server params
type EvntSrv struct {
	DbPool *pgxpool.Pool	//	
	DbQuery chan string	//for notification queries
	Events *UniqEvents	//count of unique events for db
	Logger *log.Logger
	LocalEvents map[string]bool
	
	QueryPauseMS int
	ReconnectParams waitStrat.WaitStrategy
	
	OnEvent OnEventProto
}

func (s *EvntSrv) OnNotification(_ *pgconn.PgConn, n *pgconn.Notification) {
	s.Logger.Debugf("OnNotification Channel:%s, Payload:%s", n.Channel, n.Payload)		
	
	//calling Event.publish
	params := n.Payload
	if s.LocalEvents != nil {
		if _,ok := s.LocalEvents[n.Channel]; ok {
			//local event, direct call
			if len(params) > 0 {
				params = fmt.Sprintf(`{"argv":%s}`, params)
			}			
			go s.OnEvent(n.Channel, []byte(params))
			return
		}
	}
		
	if len(params) > 0 {
		//strip curly braces
		params = n.Payload[1:len(n.Payload)-1]
	}
	s.PublishEvent(n.Channel, params)
}

//calls Event.Publish(args)
//params - comma separated string of json params
func (s *EvntSrv) PublishEvent(evId, params string) {	
	if params != "" {
		params = ","+params
	}
	payload := fmt.Sprintf(CMD_PUBLISH_ARGS, evId, params)
	s.Logger.Debugf("PublishEvent payload=%s", payload)
	
	go s.OnEvent(CMD_PUBLISH_FN, []byte(payload))
}

func (s *EvntSrv) Run() {
	if s.OnEvent == nil {
		s.Logger.Fatal("EvntSrv.OnEvent not defined")
		return
	}
	
	s.DbQuery = make(chan string)
	s.Events = &UniqEvents{m: make(map[string]int, 0)}
	
	/*
	//default event for closing connections
	s.Events.m[EVNT_LOGIN_OUT] = 1			
	*/	
	//local events
	if s.LocalEvents != nil {
		for evnt_id,_ := range s.LocalEvents {
			s.Events.m[evnt_id] = 1
		}
	}	
new_conn:	
	conn, err := s.DbPool.Acquire(context.Background())
	if err != nil {		
		time.Sleep(time.Duration(s.ReconnectParams.NextWait()) * time.Millisecond)		
		s.Logger.Errorf("EvntSrv faled acquiring connection: %v", err)
		
		goto new_conn
	}
	s.ReconnectParams.Init()
	s.Logger.Debug("EvntSrv connected to db")
	
	//all events, concorrency here if reconnecting on failur
	for evnt,_ := range s.Events.m {
		conn.Exec(context.Background(),  `LISTEN "`+evnt+`"`)
		s.Logger.Debugf("EvntSrv executing 'LISTEN %s'", evnt)		
	}
	
	var q string
	for {
		select {
		case q = <-s.DbQuery:
			s.Logger.Debugf("EvntSrv executing query: %s", q)
			
		default:
			q = ";"
		}
		
		if _, err := conn.Exec(context.Background(), q); err != nil {
			s.Logger.Errorf("EventSrv failed to execute db query '%s': %v", q, err)
			
			conn.Release()
			goto new_conn
		}			
		time.Sleep(time.Duration(s.QueryPauseMS) * time.Millisecond)
	}	
}

// adds listener
func (s *EvntSrv) AddDbListener(dbEventID string, socket *EvntSocket){	
	s.Logger.Debugf("EventSrv AddDbListener Event: %s", dbEventID)
	//s.Events.AddEvent(dbEventID, s.DbQuery, )
	socket.Events.AddEvent(dbEventID, s.DbQuery, s.Events.AddEvent)
}

// removes listener
func (s *EvntSrv) RemoveDbListener(dbEventID string, socket *EvntSocket){	
	//s.Events.RemoveEvent(dbEventID, s.DbQuery)	
	s.Logger.Debugf("EventSrv RemoveDbListener Event: %s", dbEventID)
	socket.Events.RemoveEvent(dbEventID, s.DbQuery, s.Events.RemoveEvent)
}

func (s *EvntSrv) CloseSocket(socket *EvntSocket){
	for ev_id := range socket.Events.Iter() {			
		s.Events.RemoveEvent(ev_id, s.DbQuery)
	}	
}

func NewEvntSrv(logger *log.Logger, onEvent OnEventProto, localEvents []string) *EvntSrv{
	ev := &EvntSrv{
		Logger: logger,
		QueryPauseMS: QUERY_PAUSE_MS,
		OnEvent: onEvent,
		ReconnectParams: waitStrat.WaitStrategy{
			Strategies: []waitStrat.WaitStrategyValues{
				waitStrat.WaitStrategyValues{10,1000},
				waitStrat.WaitStrategyValues{12,10000},
				waitStrat.WaitStrategyValues{0,30000},
			}},
	}
	if localEvents != nil && len(localEvents) > 0 {
		ev.LocalEvents = make(map[string]bool)
		for _, ev_id := range localEvents {
			ev.LocalEvents[ev_id] = true
		}
	}
	return ev
}
