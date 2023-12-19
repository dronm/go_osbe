package constants

/**
 * Andrey Mikhalevich 16/12/22
 *
 * Controller implimentation file
 *
 */

import (
	"fmt"
	"reflect"
	"context"
	"strings"
	
	"github.com/dronm/ds/pgds"
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
	
	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//Method implemenation
func (pm *Constant_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["ConstantList"], &ConstantList{}, sock.GetPresetFilter("ConstantList"))	
}

//Method implemenation
func (pm *Constant_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["ConstantList"], &ConstantList{}, sock.GetPresetFilter("ConstantList"))	
}


//Method implemenation
func (pm *Constant_Controller_set_value) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {

	args := rfltArgs.Interface().(*Constant_set_value)
	id := args.Id.GetValue()

	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err := d_store.GetPrimary()
	if err != nil {
		return err
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	
	if !app.GetMD().Constants.Exists(id) {
		return osbe.NewPublicMethodError(RESP_ER_NOT_FOUND, fmt.Sprintf(ER_CONST_NOT_DEFINED, id))
	}
//@ToDo sql injections!!!
//fmt.Println("OrigConstVal=", args.Val.GetValue())		
	const_val, err := app.GetMD().Constants[id].Sanatize(args.Val.GetValue())	
	if err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("Sanatize(): %v",err))
	}
//fmt.Println("SanatizeConstVal=", const_val)			
//fmt.Println(fmt.Sprintf(`SELECT const_%s_set_val(%s)`, id, const_val))
	if _, err := conn.Exec(context.Background(), fmt.Sprintf(`SELECT const_%s_set_val(%s)`, id, const_val)); err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Exec() 1: %v",err))
	}
	
	//+event Constant.update(id:"", val:"")
	if _, err := conn.Exec(context.Background(),
		fmt.Sprintf(`SELECT pg_notify('Constant.update',
					json_build_object('params',
						json_build_object(
							'id', '%s',
							'val',%s
						)
					)::text
			)`,
			id, const_val)); err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("pgx.Conn.Exec() 2: %v",err))
	}
	
	return nil
}

//Method implemenation
func (pm *Constant_Controller_get_values) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {

	args := rfltArgs.Interface().(*Constant_get_values)
	
	fld_sep := osbe.ArgsFieldSep(rfltArgs)
	ids_str := strings.Split(args.Id_list.GetValue(), fld_sep)
	query := ""
	for _, id := range ids_str {
		if !app.GetMD().Constants.Exists(id) {
			return osbe.NewPublicMethodError(RESP_ER_NOT_FOUND, fmt.Sprintf(ER_CONST_NOT_DEFINED, id))
		}
	
		if query != "" {
			query += " UNION ALL "
		}
		query += fmt.Sprintf(`SELECT
			'%s' AS id,
			const_%s_val()::text AS val,
			(SELECT c.val_type FROM const_%s c) AS val_type`,
			id, id, id);		
	}
	if query != "" {
		d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
		var conn_id pgds.ServerID
		var pool_conn *pgxpool.Conn
		pool_conn, conn_id, err := d_store.GetSecondary("")
		if err != nil {
			return err
		}
		defer d_store.Release(pool_conn, conn_id)
		conn := pool_conn.Conn()
	
		if err := osbe.AddQueryResult(resp, app.GetMD().Models["ConstantValue"], &ConstantValue{}, query, "", nil, conn, false); err != nil {
			return err
		}
	}
	
	return nil
}

