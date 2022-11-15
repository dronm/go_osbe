package session_test

import(
	"fmt"
	"testing"
	"context"
	
	"session"
	_ "session/pg"
	
	"github.com/jackc/pgx/v4/pgxpool"
)

const CONN_STR = "postgresql://postgres@localhost:5432/test_proj"
const ENC_KEY = "er4gf1tc43t84gbcfw5e4c6we5c40x5r4wfc2wrt4h54677b1hvg4e854xdwgbyujk467bv46er5014gr4w4gctk78k34r"

func TestSessionStart(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), CONN_STR)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	defer dbpool.Close()

	SessManager, er := session.NewManager("pg", 3600, 3600, dbpool, ENC_KEY)
	if er != nil {
		panic(er)
	}
	
	currentSession,er := SessManager.SessionStart("")
	//defer SessManager.SessionClose(currentSession.SessionID())
	
	if er != nil {
		panic(er)
	}	
	
	fmt.Println("SessionID=", currentSession.SessionID() )
	currentSession.Set("strVal", "Some string")
	currentSession.Set("intVal", 125)
	currentSession.Set("floatVal", 35.85)
	
	SessManager.SessionClose(currentSession.SessionID())
	
}

func TestSessionRead(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), CONN_STR)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	defer dbpool.Close()

	SessManager, er := session.NewManager("pg", 3600, 3600, dbpool, ENC_KEY)
	if er != nil {
		panic(er)
	}
	
	currentSession,er := SessManager.SessionStart("bf7db9e6-e855-f370-2b86-8c254b901102")	
	if er != nil {
		panic(er)
	}	
	
	fmt.Println("strVal=",currentSession.Get("strVal"))
	fmt.Println("intVal=",currentSession.Get("intVal"))
	fmt.Println("floatVal=",currentSession.Get("floatVal"))
	
	currentSession.Set("strVal", "Modified string2")
	fmt.Println("strVal=",currentSession.Get("strVal"))
}
