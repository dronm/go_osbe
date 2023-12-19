package contact

import (
	"reflect"	
	
	models "osbe/repo/contact/models"
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
	
	"github.com/dronm/session"	
)

type ClientSockSessioner interface {
	GetSession() session.Session
}

//Method implemenation insert
func (pm *EntityContact_Controller_insert) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["EntityContact"], &models.EntityContact_keys{}, sock.GetPresetFilter("EntityContact"))	
}

//Method implemenation delete
func (pm *EntityContact_Controller_delete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["EntityContact"], sock.GetPresetFilter("EntityContact"))	
}

//Method implemenation get_object
func (pm *EntityContact_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["EntityContactList"], &models.EntityContactList{}, sock.GetPresetFilter("EntityContactList"))	
}

//Method implemenation get_list
func (pm *EntityContact_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["EntityContactList"], &models.EntityContactList{}, sock.GetPresetFilter("EntityContactList"))	
}

//Method implemenation update
func (pm *EntityContact_Controller_update) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["EntityContact"], sock.GetPresetFilter("EntityContact"))	
}

//Method implemenation complete

