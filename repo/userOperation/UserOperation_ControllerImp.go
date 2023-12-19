package userOperation

import (
	"reflect"	
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
)


//Method implemenation delete
func (pm *UserOperation_Controller_delete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["UserOperation"], sock.GetPresetFilter("UserOperation"))	
}

//Method implemenation get_object
func (pm *UserOperation_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["UserOperationDialog"], &UserOperationDialog{}, sock.GetPresetFilter("UserOperationDialog"))	
}




