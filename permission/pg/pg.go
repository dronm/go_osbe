package pg

import (
	"context"
	"sync"
	"errors"

	"osbe/permission"

	"github.com/jackc/pgx/v4"
)

var manager = &Manager{}

type Manager struct {
	DbConnStr string
	mx sync.RWMutex
	rules permission.PermRules
}

func (mng *Manager) Reload() error{
	conn, err := pgx.Connect(context.Background(), mng.DbConnStr)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	
	mng.mx.Lock()
	defer mng.mx.Unlock()

	mng.rules = make(permission.PermRules)
	if err := conn.QueryRow(context.Background(), `SELECT rules FROM permissions LIMIT 1`).Scan(&mng.rules); err != nil {
		return err
	}
	
	return nil
}

//controller=no _Controller postfix!!!
func (mng *Manager) IsAllowed(role, controller, method string) bool{
	return mng.rules.IsAllowed(role, controller, method)
}

// InitManager initializes permission manager
// First parameter: ConnectionString string in pg format: postgresql://{USER_NAME}@{HOST}:{PORT}/{DATABASE}
func (mng *Manager) InitManager(mngParams []interface{}) (err error) {
	if len(mngParams) < 1 {
		return errors.New("InitManager missing parameter: pg connection string")
	}	
	ok := false
	mng.DbConnStr, ok = mngParams[0].(string)
	if !ok {
		return errors.New("InitManager db connection parameter must be a string")
	}
	mng.Reload()
	return nil
}

func init() {
	permission.Register("pg", manager)
}
