package validRusBanckAcc

import (
	"context"
	"testing"
	"os"
	
	"github.com/jackc/pgx/v4/pgxpool"	
)

const (
	TEST_VAR_PG_CONN = "PG_CONN"
	TEST_VAR_BIK = "BIK"
	TEST_VAR_ACC_NUM = "ACC_NUM"
)

func getTestVar(t *testing.T, n string) *string {
	v := os.Getenv(n)
	if v == "" {
		t.Fatalf("getTestVar() failed: %s environment variable is not set", n)
	}
	return &v
}

func getPGPool(t *testing.T) (*pgxpool.Pool, error) {
	conn_conf, err := pgxpool.ParseConfig(*getTestVar(t, TEST_VAR_PG_CONN))
	if err != nil {
		return nil, err
	}
	return pgxpool.ConnectConfig(context.Background(), conn_conf)	
}

func TestCheckBik(t *testing.T) {	
	pool, err := getPGPool(t)
	if err != nil {
		t.Fatalf("getPGPool() failed: %v", err)
	}	
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("pool.Acquire() failed: %v", err)	
	}
	
	bik := *getTestVar(t, TEST_VAR_BIK)
	if !CheckBik(bik, conn.Conn(), ctx) {
		t.Fatalf("Bik: %s expeced to exist, but not found", bik)	
	}
}

func TestCheckAcc(t *testing.T) {	
	bik := *getTestVar(t, TEST_VAR_BIK)
	acc_num := *getTestVar(t, TEST_VAR_ACC_NUM)
	if !CheckAcc(bik, acc_num) {
		t.Fatalf("Account num: %s, bik %s expeced to be correct, but it is not", acc_num, bik)	
	}
}
