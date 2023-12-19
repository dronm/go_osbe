package service

import (
	"reflect"
	
	"osbe"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"
)

//Controller
type Service_Controller struct {
	osbe.Base_Controller
}

func NewController_Service() *Service_Controller{
	c := &Service_Controller{osbe.Base_Controller{ID: "Service", PublicMethods: make(osbe.PublicMethodCollection)}}
	
	//************************** method reload_config **********************************
	c.PublicMethods["reload_config"] = &Service_Controller_reload_config{
		osbe.Base_PublicMethod{
			ID: "reload_config",
			Fields: nil,
		},
	}

	//************************** method reload_version **********************************
	c.PublicMethods["reload_version"] = &Service_Controller_reload_version{
		osbe.Base_PublicMethod{
			ID: "reload_version",
			Fields: nil,
		},
	}
	
	return c
}

//**************************************************************************************
//Public method: reload_config
type Service_Controller_reload_config struct {	
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Service_Controller_reload_config) Unmarshal(payload []byte) (res reflect.Value, err error) {
	return res, nil
}

//custom method
func (pm *Service_Controller_reload_config) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	if err := app.ReloadAppConfig(); err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	return nil
}

//**************************************************************************************
//Public method: reload_version
type Service_Controller_reload_version struct {	
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Service_Controller_reload_version) Unmarshal(payload []byte) (res reflect.Value, err error) {
	return res, nil
}

//custom method
func (pm *Service_Controller_reload_version) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	if err := app.LoadAppVersion(); err != nil {
		return osbe.NewPublicMethodError(response.RESP_ER_INTERNAL, err.Error())
	}
	
	return nil
}

