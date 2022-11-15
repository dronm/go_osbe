package about

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 *
 * Custom controller
 */

import (
	"reflect"
	
	"osbe"
	"osbe/fields"
	"osbe/srv"
	"osbe/socket"
	"osbe/response"	
)

//Controller
type About_Controller struct {
	osbe.Base_Controller
}

func (c *About_Controller) InitPublicMethods() {
	c.ID = "About"
	About_Model_init()

	c.PublicMethods = make(osbe.PublicMethodCollection)
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &About_Controller_get_object{
		osbe.Base_PublicMethod{
			ID: "get_object",
		},
	}
}

//************************* GET OBJECT **********************************************
type About_Controller_get_object struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *About_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {	
	var res reflect.Value
	return res, nil
}

//Method implemenation
func (pm *About_Controller_get_object) Run(app osbe.Applicationer, serv srv.Server, sock socket.ClientSocketer, resp *response.Response, rfltArgs reflect.Value) error {
	conf := app.GetConfig()	
	m_row := &About{Author: fields.ValText{TypedValue: conf.GetAuthor()},
		Tech_mail: fields.ValText{TypedValue: conf.GetTechMail()},
		App_name: fields.ValText{TypedValue: conf.GetAppID()},
		Fw_version: fields.ValText{TypedValue: app.GetFrameworkVersion()},
		App_version: fields.ValText{TypedValue: app.GetMD().Version.Value},
		Db_name: fields.ValText{},
	}
	resp.AddModelFromStruct(m_row.GetID(), m_row)
	return nil
}



