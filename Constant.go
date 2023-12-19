package osbe

// TODO: correct validation of certain types.

import(
	"fmt"
	"time"
	"strings"
	"context"
	
	"github.com/dronm/ds/pgds"	
	"osbe/fields"	
	
	"github.com/jackc/pgx/v4/pgxpool"
)

// Constant describes application constant object.
type Constant interface {
	GetAutoload() bool		// if constant is sent to client at first request.
	Sanatize(string) (string,error) // manages value validation.
}

// ConstantCollection is a collection of all application constants.
type ConstantCollection map[string] Constant

// Exists returns true if a given constant exists.
func (c ConstantCollection) Exists(ID string) bool {
	_, ok := c[ID]
	return ok
/*	for const_id, _ := range c {
		if const_id == ID {
			return true
		}
	
	}
	return false*/
}

// GetValue fetches value of constant constID from store to constVal.
func RetrieveValue(dStore *pgds.PgProvider, constID string, constVal interface{}) error {
	//from data base
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := dStore.GetSecondary("")
	if err != nil {
		return err
	}
	defer dStore.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
		
	if err := conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT const_%s_val()`, constID)).Scan(constVal); err != nil {
		return err
	}
	return nil
}

// Structures for specific constants.

// ConstantInt is a constant of integer type.
type ConstantInt struct {
	ID string		// constant ID
	Autoload bool		// if autoload
	Value fields.ValInt	// constant value
}
func (c *ConstantInt) GetAutoload() bool {
	return c.Autoload
}

// GetValue return constant value.
// TODO: caching. Clear value on updates through events.
func (c *ConstantInt) GetValue(app Applicationer) (int64, error) {
	if c.Value.GetIsSet() {
		return c.Value.GetValue(), nil
	}
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	if err := RetrieveValue(d_store, c.ID, &c.Value); err != nil {
		return 0, err
	}
/*
	//from data base
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary("")
	if err != nil {
		return 0, err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
		
	if err := conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT const_%s_val()`,c.ID)).Scan(&c.Value); err != nil {
		return 0, err
	}
*/
	return c.Value.GetValue(), nil
}

// Sanatize sanatizes value for db.
func (c *ConstantInt) Sanatize(val string) (string, error) {	
	i, err := fields.StrToInt(val)
	if err != nil {
		return "",err
	}
	return fmt.Sprintf("%d::int", i), nil
}

// ConstantText is string value constant
type ConstantText struct {
	ID string
	Autoload bool
	Value fields.ValText
}
func (c *ConstantText) GetAutoload() bool {
	return c.Autoload
}

func (c *ConstantText) GetValue(app Applicationer) (string, error) {
	if c.Value.GetIsSet() {
		return c.Value.GetValue(), nil
	}
	//from data base
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary("")
	if err != nil {
		return "", err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
		
	if err := conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT const_%s_val()`,c.ID)).Scan(&c.Value); err != nil {
		return "", err
	}
	return c.Value.GetValue(), nil
}
// TODO: string validating.
func (c *ConstantText) Sanatize(val string) (string,error) {	
	return "'"+strings.ReplaceAll(val, "'", `\'`)+"'::text", nil
}

//******************************
type ConstantTime struct {
	ID string
	Autoload bool
	Value fields.ValTime
}
func (c *ConstantTime) GetAutoload() bool {
	return c.Autoload
}

func (c *ConstantTime) GetValue(app Applicationer) (time.Time, error) {
	if c.Value.GetIsSet() {
		return c.Value.GetValue(), nil
	}
	//from data base
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary("")
	if err != nil {
		return time.Time{}, err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
		
	if err := conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT const_%s_val()`,c.ID)).Scan(&c.Value); err != nil {
		return time.Time{}, err
	}
	return c.Value.GetValue(), nil
}
func (c *ConstantTime) Sanatize(val string) (string,error) {	
	return "'"+strings.ReplaceAll(val, "'", `\'`)+"'::interval", nil
}

//******************************
type ConstantFloat struct {
	ID string
	Autoload bool
	Value fields.ValFloat
}
func (c *ConstantFloat) GetAutoload() bool {
	return c.Autoload
}

func (c *ConstantFloat) GetValue(app Applicationer) (float64, error) {
	if c.Value.GetIsSet() {
		return c.Value.GetValue(), nil
	}
	//from data base
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary("")
	if err != nil {
		return 0, err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
		
	if err := conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT const_%s_val()`,c.ID)).Scan(&c.Value); err != nil {
		return 0, err
	}
	return c.Value.GetValue(), nil
}
func (c *ConstantFloat) Sanatize(val string) (string,error) {	
	f, err := fields.StrToFloat(val)
	if err != nil {
		return "",err
	}
	return fmt.Sprintf("%f::numeric", f), nil
}

//******************************
type ConstantJSON struct {
	ID string
	Autoload bool
	Value fields.ValJSON
}
func (c *ConstantJSON) GetAutoload() bool {
	return c.Autoload
}

func (c *ConstantJSON) GetValue(app Applicationer) ([]byte, error) {
	if c.Value.GetIsSet() {
		return c.Value.GetValue(), nil
	}
	//from data base
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary("")
	if err != nil {
		return nil, err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
		
	if err := conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT const_%s_val()`,c.ID)).Scan(&c.Value); err != nil {
		return nil, err
	}
	return c.Value.GetValue(), nil
}
func (c *ConstantJSON) Sanatize(val string) (string, error) {	
	return "'"+strings.ReplaceAll(val, "'", `\'`)+"'::json", nil
}

//******************************
type ConstantBytea struct {
	ID string
	Autoload bool
	Value fields.ValBytea
}
func (c *ConstantBytea) GetAutoload() bool {
	return c.Autoload
}

func (c *ConstantBytea) GetValue(app Applicationer) ([]byte, error) {
	if c.Value.GetIsSet() {
		return c.Value.GetValue(), nil
	}
	//from data base
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetSecondary("")
	if err != nil {
		return []byte{}, err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
		
	if err := conn.QueryRow(context.Background(), fmt.Sprintf(`SELECT const_%s_val()`,c.ID)).Scan(&c.Value); err != nil {
		return []byte{}, err
	}
	return c.Value.GetValue(), nil
}
func (c *ConstantBytea) Sanatize(val string) (string,error) {	
	return val, nil
}

