package pg_pool

import (
	"sync"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgconn"
)

type OnDbNotificationProto = func(*pgconn.PgConn, *pgconn.Notification)

//Holds db instances
type Db struct {
	connStr string
	onNotification OnDbNotificationProto
	Pool *pgxpool.Pool
	mx sync.RWMutex
	active bool
	ref int
}

func (d *Db) getRefCount() int {
	d.mx.Lock()
	defer d.mx.Unlock()
	return d.ref
}

func (d *Db) addRef() (error) {
	d.mx.Lock()
	defer d.mx.Unlock()
	
	if !d.active {
		conn_conf, err := pgxpool.ParseConfig(d.connStr)
		if err != nil {
			return err
		}

		conn_conf.ConnConfig.OnNotification = d.onNotification
		
		d.Pool, err = pgxpool.ConnectConfig(context.Background(), conn_conf)
		if err != nil {
			return err
		}
		
		d.active = true
	}	
	d.ref++
	
	return nil
}

func (d *Db) Release() {
	d.mx.Lock()
	d.ref--
	if d.ref < 0 {
		d.ref = 0
	}
	d.mx.Unlock()
}

type ServerID string

//server pool holds primary and secondary instances
type ServerPool struct {
	Primary *Db
	Secondaries map[ServerID]*Db
}

func (p *ServerPool) GetPrimary() (*Db, error) {
	err := p.Primary.addRef()
	if err != nil {
		return nil , err
	}else{
		return p.Primary, nil
	}
}

//Looks for a secondary with less ref count
//if nothing found returns primary
func (p *ServerPool) GetSecondary() (*Db, error) {
	var exclude_id []ServerID
	var tmp_db *Db
	var tmp_id ServerID
	
srv_loop:	
	if p.Secondaries != nil {
		min_ref := 9999999		
		for sec_id, sec := range p.Secondaries {
			do_exclude := false
			for _,j := range exclude_id {
				if j == sec_id {
					do_exclude = true
					break
				}
			}
			if !do_exclude {
				cnt := sec.getRefCount()
				if min_ref > cnt {
					min_ref = cnt
					tmp_id = sec_id
					tmp_db = sec
				}
			}
		}
	}
	if tmp_db != nil {
		err := tmp_db.addRef()
		if err != nil {
			exclude_id = append(exclude_id, tmp_id)
			tmp_db = nil
			tmp_id = ""
			goto srv_loop
		}else{
			return tmp_db ,nil
		}
	}else{
		//no secondary available
		return p.GetPrimary()
	}
}

func NewServerPool(primaryConnStr string, onDbNotification OnDbNotificationProto, secondaries map[string]string) *ServerPool{
	p := &ServerPool{}
	p.Primary = &Db{connStr: primaryConnStr, onNotification: onDbNotification}	
	
	if secondaries != nil {
		p.Secondaries = make(map[ServerID]*Db, 0)
		for id,conn := range secondaries {
			p.Secondaries[ServerID(id)] = &Db{connStr: conn}
		}
	}	
	return p
}
