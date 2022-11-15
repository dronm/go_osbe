package pg

import (
	"context"
	"sync"
	"errors"

	"osbe/permission"

	"github.com/jackc/pgx/v4/pgxpool"
)

var manager = &Manager{}

type Manager struct {
	DbPool *pgxpool.Pool
	mx sync.RWMutex
	rules permission.PermRules
}

func (mng *Manager) Reload() error{
	mng.mx.Lock()
	defer mng.mx.Unlock()

	mng.rules = make(permission.PermRules)
	return mng.DbPool.QueryRow(context.Background(), `SELECT rules FROM permissions LIMIT 1`).Scan(&mng.rules)
}

//controller=no _Controller postfix!!!
func (mng *Manager) IsAllowed(role, controller, method string) bool{

	if role == "" {
		role = permission.DEFAULT_ROLE
	}
//fmt.Println("role=", role, "controller=", controller,"method=",method)		
	if mng_contr, ok := mng.rules[role]; ok {
		if mng_meth, ok := mng_contr[controller]; ok {
			if mng_allowed, ok := mng_meth[method]; ok {
				return mng_allowed
			}
		}
	}
	return false
}

//first parameter: *pgxpool.Pool
func (mng *Manager) InitManager(provParams []interface{}) (err error) {
	if len(provParams)<1 {
		return errors.New("Missing parameters: *pgxpool.Pool")
	}	
	var ok bool
	mng.DbPool, ok = provParams[0].(*pgxpool.Pool)
	if !ok {
		return errors.New("Parameter must be of type *pgxpool.Pool")
	}
	mng.Reload()
	return nil
}

func init() {
	permission.Register("pg", manager)
}
