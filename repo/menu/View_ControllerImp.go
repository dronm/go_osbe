package menu

import (
	"reflect"	
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
)

//Method implemenation
func (pm *View_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["ViewList"], &ViewList{}, nil)
}

//Method implemenation
func (pm *View_Controller_get_section_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return nil
}

//Method implemenation
func (pm *View_Controller_complete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.CompleteOnArgs(app, resp, rfltArgs, app.GetMD().Models["ViewList"], &ViewList{}, nil)
}


