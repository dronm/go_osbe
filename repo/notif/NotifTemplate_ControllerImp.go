package notif

import (
	"reflect"	
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
	
)

//Method implemenation insert
func (pm *NotifTemplate_Controller_insert) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["NotifTemplate"], &NotifTemplate_keys{}, sock.GetPresetFilter("NotifTemplate"))	
}

//Method implemenation delete
func (pm *NotifTemplate_Controller_delete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["NotifTemplate"], sock.GetPresetFilter("NotifTemplate"))	
}

//Method implemenation get_object
func (pm *NotifTemplate_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["NotifTemplate"], &NotifTemplate{}, sock.GetPresetFilter("NotifTemplate"))	
}

//Method implemenation get_list
func (pm *NotifTemplate_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["NotifTemplateList"], &NotifTemplateList{}, sock.GetPresetFilter("NotifTemplateList"))	
}

//Method implemenation update
func (pm *NotifTemplate_Controller_update) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["NotifTemplate"], sock.GetPresetFilter("NotifTemplate"))	
}


