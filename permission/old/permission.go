package permission

/**
 * table permissions has one column rules of type json
 * {
 *	"role":{
 *		"controller":{
 *			"insert": false,
 *			"delete": false
 *			"update": true
 *			"select": false
 *		}
 *	}
 * }
 */

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

//@ToDo make it subscribe to Permission.update and reload automaticaly

const DEFAULT_ROLE = "guest"

type permMethod map[string]bool
type permController map[string]permMethod
type permRules map[string]permController

type Manager struct {
	DbPool *pgxpool.Pool
	mx sync.RWMutex
	rules permRules
}

func (mng *Manager) Reload() error{
	mng.mx.Lock()
	defer mng.mx.Unlock()

	mng.rules = make(permRules)
	return mng.DbPool.QueryRow(context.Background(), `SELECT rules FROM permissions LIMIT 1`).Scan(&mng.rules)
}

func (mng *Manager) IsAllowed(role, controller, method string) bool{

	if role == "" {
		role = DEFAULT_ROLE
	}
	if mng_contr, ok := mng.rules[role]; ok {
		if mng_meth, ok := mng_contr[controller]; ok {
			if mng_allowed, ok := mng_meth[method]; ok {
				return mng_allowed
			}
		}
	}
	return false
}

func NewManager(dbPool *pgxpool.Pool) *Manager{
	m := &Manager{DbPool: dbPool}
	m.Reload()
	return m
}

