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

func (pm *TimeZoneLocale_Controller_insert) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["TimeZoneLocale"], &models.TimeZoneLocale_keys{}, sock.GetPresetFilter("TimeZoneLocale"))
}

//Method implemenation
func (pm *TimeZoneLocale_Controller_delete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["TimeZoneLocale"], sock.GetPresetFilter("TimeZoneLocale"))
}

//Method implemenation
func (pm *TimeZoneLocale_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["TimeZoneLocale"], &models.TimeZoneLocale{}, sock.GetPresetFilter("TimeZoneLocale"))
}

//Method implemenation
func (pm *TimeZoneLocale_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["TimeZoneLocale"], &models.TimeZoneLocale{}, sock.GetPresetFilter("TimeZoneLocale"))
}

//Method implemenation
func (pm *TimeZoneLocale_Controller_update) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["TimeZoneLocale"], sock.GetPresetFilter("TimeZoneLocale"))
}

