package bank

import (
	"reflect"	
	
	models "osbe/repo/bank/models"
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
	
	//"github.com/jackc/pgx/v4"
)



//Method implemenation get_object
func (pm *Bank_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["BankList"], &models.BankList{}, sock.GetPresetFilter("BankList"))	
}

//Method implemenation get_list
func (pm *Bank_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["BankList"], &models.BankList{}, sock.GetPresetFilter("BankList"))	
}


//Method implemenation complete
func (pm *Bank_Controller_complete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.CompleteOnArgs(app, resp, rfltArgs, app.GetMD().Models["BankList"], &models.BankList{}, sock.GetPresetFilter("BankList"))	
}

