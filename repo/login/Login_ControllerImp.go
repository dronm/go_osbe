package login

import (
	"reflect"	
	
	"osbe/repo/login/models"
	
	"osbe"
	"osbe/srv"
	"osbe/evnt"
	"osbe/socket"
	"osbe/response"	
)



//Method implemenation get_object
func (pm *Login_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginList"], &models.LoginList{}, sock.GetPresetFilter("LoginList"))	
}

//Method implemenation get_list
func (pm *Login_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginList"], &models.LoginList{}, sock.GetPresetFilter("LoginList"))	
}

func (pm *Login_Controller_destroy_session) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	args := rfltArgs.Interface().(*evnt.Event)
	session_id_i, ok := args.Params["session_id"]
	if !ok {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, "Login_Controller_destroy_session session_id parameter is missing")
	}
	session_id, ok := session_id_i.(string)
	if !ok {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, "Login_Controller_destroy_session session_id parameter is not a string")
	}
	app.GetSessManager().SessionDestroy(session_id)	
	return nil
}

