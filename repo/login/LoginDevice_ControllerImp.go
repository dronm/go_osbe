package login

import (
	"reflect"	
	"context"
	"fmt"
	
	"osbe/repo/login/models"
	
	"github.com/dronm/ds/pgds"
	
	"osbe"
	"osbe/fields"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
	
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
)




//Method implemenation get_list
func (pm *LoginDevice_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginDeviceList"], &models.LoginDeviceList{}, sock.GetPresetFilter("LoginDeviceList"))	
}


//Method implemenation
func (pm *LoginDevice_Controller_switch_banned) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {

	args := rfltArgs.Interface().(*models.LoginDevice_switch_banned)
	
	d_store,_ := app.GetDataStorage().(*pgds.PgProvider)
	var conn_id pgds.ServerID
	var pool_conn *pgxpool.Conn
	pool_conn, conn_id, err_conn := d_store.GetPrimary()
	if err_conn != nil {
		return err_conn
	}
	defer d_store.Release(pool_conn, conn_id)
	conn := pool_conn.Conn()
	if args.Banned.GetValue() {
		_, err := conn.Exec(context.Background(), "BEGIN")
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("LoginDevice_Controller_switch_banned pgx.Conn.Exec() BEGIN: %v",err))
		}	
		//DELETE FROM session_vals WHERE id IN
		var session_id fields.ValText
		if err := conn.QueryRow(context.Background(),
			`SELECT
				session_id
			FROM logins
			WHERE user_id = $1
			AND md5(login_devices_uniq(user_agent)) = $2
			AND date_time_out IS NULL
			ORDER BY date_time_in DESC
			LIMIT 1`,
			args.User_id.GetValue(), args.Hash.GetValue()).Scan(&session_id); err != nil && err != pgx.ErrNoRows {
			
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("LoginDevice_Controller_switch_banned pgx.Conn.QueryRow() SELECT session_id: %v",err))
			
		}else if err == nil {
			app.GetSessManager().SessionDestroy(session_id.GetValue())
		}
		_, err = conn.Exec(context.Background(),
			`INSERT INTO login_device_bans (user_id, hash) VALUES ($1, $2)`,
			args.User_id.GetValue(), args.Hash.GetValue())
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("LoginDevice_Controller_switch_banned pgx.Conn.Exec() INSERT login_device_bans: %v",err))
		}	
	
		_, err = conn.Exec(context.Background(), "COMMIT")
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("LoginDevice_Controller_switch_banned pgx.Conn.Exec() COMMIT: %v",err))
		}	
		//Send event!!!
		//deleted from logins after delete!!!
	
	}else{
		_, err := conn.Exec(context.Background(),
			"DELETE FROM login_device_bans WHERE user_id = $1 AND hash = $2",
			args.User_id.GetValue(), args.Hash.GetValue())
		if err != nil {
			return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, fmt.Sprintf("LoginDevice_Controller_switch_banned pgx.Conn.Exec() DELETE login_device_bans: %v",err))
		}	
	}
	
	return nil
}
