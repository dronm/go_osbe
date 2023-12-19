package menu

import (
	"reflect"	
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
)

//
func (pm *VariantStorage_Controller_insert) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.InsertOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["VariantStorage"], VariantStorage_keys{}, nil)
}

//Method implemenation
func (pm *VariantStorage_Controller_delete) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.DeleteOnArgKeys(app, pm, resp, sock, rfltArgs, app.GetMD().Models["VariantStorage"], nil)
}

//Method implemenation
func (pm *VariantStorage_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetObjectOnArgs(app, resp, rfltArgs, app.GetMD().Models["VariantStorage"], &VariantStorage{}, sock.GetPresetFilter("VariantStorage"))
}

//Method implemenation
func (pm *VariantStorage_Controller_get_list) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.GetListOnArgs(app, resp, rfltArgs, app.GetMD().Models["VariantStorage"], &VariantStorage{}, sock.GetPresetFilter("VariantStorage"))
}

//Method implemenation
func (pm *VariantStorage_Controller_update) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return osbe.UpdateOnArgs(app, pm, resp, sock, rfltArgs, app.GetMD().Models["VariantStorage"], nil)
}

//Method implemenation
func (pm *VariantStorage_Controller_upsert_filter_data) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return nil
}

//Method implemenation
func (pm *VariantStorage_Controller_upsert_col_visib_data) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return nil
}

//Method implemenation
func (pm *VariantStorage_Controller_upsert_col_order_data) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return nil
}

//Method implemenation
func (pm *VariantStorage_Controller_get_filter_data) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return nil
}

//Method implemenation
func (pm *VariantStorage_Controller_get_col_visib_data) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return nil
}


//Method implemenation
func (pm *VariantStorage_Controller_get_col_order_data) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	return nil
}
