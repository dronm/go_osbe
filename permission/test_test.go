package permission

import(
	"testing"
	"context"
	"fmt"
	
	"github.com/jackc/pgx/v4/pgxpool"
)

const DB_CONN = "postgresql://postgres@localhost:5432/test_proj"


func TestPermissions(t *testing.T) {
	conn_conf, err := pgxpool.ParseConfig(DB_CONN)
	if err != nil {
		panic(fmt.Sprintf("pgxpool.ParseConfig: %v", err))
	}
	
	conn, err := pgxpool.ConnectConfig(context.Background(), conn_conf)

	perm_mng := &Manager{DbConn: conn}
	perm_mng.Reload()
	
	if !perm_mng.IsAllowed("admin", "Test", "insert") {
		t.Fail()
		t.Log("admin Test.insert Expected: true, got false")
	}
	
	if !perm_mng.IsAllowed("admin", "Test", "select") {
		t.Fail()
		t.Log("admin Test.select Expected: true, got false")
	}

	if perm_mng.IsAllowed("admin", "Test", "get_list") {
		t.Fail()
		t.Log("admin Test.get_list Expected: false, got true")
	}

	if perm_mng.IsAllowed("manager", "Test", "insert") {
		t.Fail()
		t.Log("manager Test.get_list Expected: false, got true")
	}
	
}	

