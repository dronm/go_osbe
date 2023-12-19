package login

import (
	"reflect"	
	
	"osbe/repo/login/models"
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
	
	//"github.com/jackc/pgx/v4"
)

func (pm *LoginDeviceBan_Controller_insert) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["LoginDeviceBan"], &models.LoginDeviceBan_keys{}, sock.GetPresetFilter("LoginDeviceBan"))	
}

//Method implemenation
func (pm *LoginDeviceBan_Controller_delete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["LoginDeviceBan"], sock.GetPresetFilter("LoginDeviceBan"))	
}

//Method implemenation
func (pm *LoginDeviceBan_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginDeviceBan"], &models.LoginDeviceBan{}, sock.GetPresetFilter("LoginDeviceBan"))	
}

//Method implemenation
func (pm *LoginDeviceBan_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["LoginDeviceBan"], &models.LoginDeviceBan{}, sock.GetPresetFilter("LoginDeviceBan"))	
}

