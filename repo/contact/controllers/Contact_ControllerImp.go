package contact

import (
	"reflect"	
	
	models "osbe/repo/contact/models"
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
	
	//"github.com/jackc/pgx/v4"
)

//********************
//Method implemenation insert
func (pm *Contact_Controller_insert) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["Contact"], &models.Contact_keys{}, sock.GetPresetFilter("Contact"))	
}

//Method implemenation delete
func (pm *Contact_Controller_delete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["Contact"], sock.GetPresetFilter("Contact"))	
}

//Method implemenation get_object
func (pm *Contact_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["ContactDialog"], &models.ContactDialog{}, sock.GetPresetFilter("ContactDialog"))	
}

//Method implemenation get_list
func (pm *Contact_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["ContactList"], &models.ContactList{}, sock.GetPresetFilter("ContactList"))	
}

//Method implemenation update
func (pm *Contact_Controller_update) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["Contact"], sock.GetPresetFilter("Contact"))	
}

func (pm *Contact_Controller_complete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.CompleteOnArgs(app, resp, rfltArgs, app.GetMD().Models["ContactList"], &models.ContactList{}, sock.GetPresetFilter("ContactList"))	
}


