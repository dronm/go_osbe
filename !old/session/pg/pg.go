package pg

/** Requirements:
 *	jackc connection to PG Sql https://github.com/jackc/pgx
 *
 *	PG_CRYPTO extension must be installed CREATE EXTENSION pgrypto, PGP_SYM_DECRYPT, PGP_SYM_ENCRYPT functions are used,
 *		if encryption is not necessary - correct sql in SessionRead/SessionClose functions
 *
 *	Some SQL scripts:
 *	session_vals.sql contains table for holding session values
 *	session_vals_process.sql trigger function for updating login information (logins table must be present in database)
 *	session_vals_trigger.sql creating trigger script
*/

import (
	"container/list"	
	"sync"
	"time"
	"context"
	"errors"
	//"fmt"

	"encoding/gob"
	"bytes"
	"encoding/base64"
	
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
	
	"osbe/session"
	"osbe/fields"
)

var pder = &Provider{list: list.New()}

//*** SessionStore ***
type storeValue map[string]interface{}

type SessionStore struct {
	sid          string                     //session id
	mx sync.RWMutex
	timeAccessed time.Time                  //last modified
	timeCreated  time.Time                  //when created
	value        storeValue 		//key-value pair
	valueModified bool
}

//just set inmemory value
func (st *SessionStore) Set(key string, value interface{}) error {
	if st.value[key] != value {
		st.mx.Lock()
		st.value[key] = value
		st.valueModified = true
		st.mx.Unlock()
		
		pder.sessionUpdate(st.sid)
	}
	
	return nil
}

func (st *SessionStore) Flush() error {
//fmt.Println("Flush() ", st, "ID=",st.sid, "VAL=", st.value)
	val, err := getForDb(&st.value)
	if err!= nil {
		return err
	}
//fmt.Println("Flush() error", err)
	//flush val only if it's been modified
	if st.valueModified {
		//modified
//fmt.Println("Flush() actual update, ID=",st.sid)		
		_,err = pder.dbpool.Exec(context.Background(),
			`UPDATE session_vals
			SET
				val = PGP_SYM_ENCRYPT($1,$2),
				accessed_time = now()
			WHERE id = $3`,
			val,
			pder.encrkey,
			st.sid)			
		st.mx.Lock()
		st.valueModified = false
		st.mx.Unlock()
	}
	
	return nil
}

//get from memory
func (st *SessionStore) Get(key string) interface{} {
	pder.sessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *SessionStore) GetBool(key string) bool {
	v := st.Get(key)
	if v != nil {
		if v_bool, ok := v.(bool); ok {
			return v_bool
		}		
	}
	return false
}

func (st *SessionStore) GetString(key string) string {
	v := st.Get(key)
	if v != nil {
		if v_str, ok := v.(string); ok {
			return v_str
			
		}else if v_str, ok := v.([]byte); ok {
			return string(v_str)
		}
	}
	return ""
}

func (st *SessionStore) GetInt(key string) int64 {
	v := st.Get(key)
	if v != nil {
		if v_i, ok := v.(int64); ok {
			return v_i
			
		}else if v_i, ok := v.(int); ok {
			return int64(v_i)
		}
	}
	return 0
}

func (st *SessionStore) GetFloat(key string) float64 {
	v := st.Get(key)
	if v != nil {
		if v_f, ok := v.(float64); ok {
			return v_f
			
		}else if v_f, ok := v.(float32); ok {
			return float64(v_f)
		}
	}
	return 0
}

func (st *SessionStore) GetDate(key string) time.Time {
	v := st.Get(key)
	if v != nil {
		if v_t, ok := v.(time.Time); ok {
			return v_t			
		}
	}
	return time.Time{}
}

//delete from memmory
func (st *SessionStore) Delete(key string) error {

	delete(st.value, key)
	pder.sessionUpdate(st.sid)
	
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

func (st *SessionStore) TimeCreated() time.Time {
	return st.timeCreated
}

//*** Provider ***
type Provider struct {
	lock     sync.Mutex               
	sessions map[string]*list.Element 
	list     *list.List
	dbpool	 *pgxpool.Pool
	encrkey  string	
}

func (pder *Provider) sessionMemInit(sid string) (element *list.Element, newSess session.Session) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	
	v := make(map[string]interface{}, 0)
	
	newSess = &SessionStore{
		sid: sid,
		timeAccessed: time.Now(),
		timeCreated: time.Now(),
		value: v,
	}
	element = pder.list.PushBack(newSess)
	pder.sessions[sid] = element
	return 
}

func (pder *Provider) SessionInit(sid string) (session.Session, error) {
	if pder.dbpool == nil {
		return nil,errors.New("Provider not initialized")
	}
	
	_,new_sess := pder.sessionMemInit(sid)
	
	_, err := pder.dbpool.Exec(context.Background(),
		"INSERT INTO session_vals(id) VALUES($1)",
		sid,
	)
	
	return new_sess, err
}

func setFromDb(strucVal *storeValue, dbVal string) error{
	var b bytes.Buffer
	valb, err := base64.URLEncoding.DecodeString(dbVal)
	if err != nil {
		return err
	}		
	b.Write(valb)		
	dec := gob.NewDecoder(&b)
	dec.Decode(strucVal)
	
	return nil
}

func getForDb(strucVal *storeValue) (string, error){
	var b bytes.Buffer
	
	enc := gob.NewEncoder(&b)
	err := enc.Encode(strucVal)
	if err != nil {
		return "",err
	}
	res := base64.URLEncoding.EncodeToString(b.Bytes())
	return res,nil
}	
	
func (pder *Provider) SessionRead(sid string) (session.Session, error) {

	element,_ := pder.sessionMemInit(sid)
	var val fields.ValText
	
	err := pder.dbpool.QueryRow(context.Background(),
		`SELECT
			accessed_time,
			create_time,
			PGP_SYM_DECRYPT(val,$1)
		FROM session_vals
		WHERE id=$2`,
	pder.encrkey,
	sid).Scan(&element.Value.(*SessionStore).timeAccessed,
		&element.Value.(*SessionStore).timeCreated,
		&val)
	if err == pgx.ErrNoRows {
		//no such session
		return pder.SessionInit(sid)
		
	}else if err != nil {
		return nil, err
	}
	
	if err := setFromDb(&element.Value.(*SessionStore).value, val.GetValue()); err != nil {
		return nil, err
	}				
	return element.Value.(*SessionStore), nil
}

//write session data to db
func (pder *Provider) SessionClose(sid string) (err error) {
	
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).Flush()
		/*		
		var val string
		
		val,err = getForDb(&element.Value.(*SessionStore).value)
		if err!= nil {
			return err
		}
		
		//flush val only if it's been modified
		if !element.Value.(*SessionStore).valueModified {
			_,err = pder.dbpool.Exec(context.Background(),
				`UPDATE session_vals
				SET
					accessed_time = now()
				WHERE id = $1`,
				sid)	
		}else{
			//modified
			_,err = pder.dbpool.Exec(context.Background(),
				`UPDATE session_vals
				SET
					val = PGP_SYM_ENCRYPT($1,$2),
					accessed_time = now()
				WHERE id = $3`,
				val,
				pder.encrkey,
				sid)			
		}
		*/
	}

	return err
}

func (pder *Provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		return pder.removeSession(element, sid)
	}
	return nil
}

func (pder *Provider) SessionGC(maxLifeTime int64, maxIdleTime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		tm := time.Now().Unix()
		if ((element.Value.(*SessionStore).timeCreated.Unix() + maxLifeTime) < tm) || ((element.Value.(*SessionStore).timeAccessed.Unix() + maxIdleTime) < tm) {
			//pder.list.Remove(element)
			//delete(pder.sessions, element.Value.(*SessionStore).sid)
			pder.removeSession(element, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

//first parameter: ConnectionString string
//second parameter: encryptKey application unique,if to set no encryption used
func (pder *Provider) InitProvider(provParams []interface{}) (err error) {
	if len(provParams)<2 {
		return errors.New("Missing parameters: pgxpool.Pool, encryptKey")
	}	
	//pder.dbpool, err = pgxpool.Connect(context.Background(), provParams[0].(string))
	pder.dbpool = provParams[0].(*pgxpool.Pool)
	
	pder.encrkey = provParams[1].(string)
	
	return err
}


//helper function for SessionDestroy and SessionGC
func (pder *Provider) removeSession(el *list.Element, sid string) (err error) {
	_, err = pder.dbpool.Exec(context.Background(),
		`DELETE FROM session_vals WHERE id=$1`,sid)
	if err != nil {
		return err
	}

	delete(pder.sessions, sid)
	pder.list.Remove(el)
	return err
}

//protected
func (pder *Provider) sessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("pg", pder)
}
