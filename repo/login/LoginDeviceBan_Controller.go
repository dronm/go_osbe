package login

import (
	"reflect"	
	"encoding/json"
	
	"osbe/repo/login/models"
	
	"osbe"
	"osbe/fields"
	"osbe/model"
)

//Controller
type LoginDeviceBan_Controller struct {
	osbe.Base_Controller
}

func NewController_LoginDeviceBan() *LoginDeviceBan_Controller{
	c := &LoginDeviceBan_Controller{osbe.Base_Controller{ID: "LoginDeviceBan", PublicMethods: make(osbe.PublicMethodCollection)}}	
	keys_fields := fields.GenModelMD(reflect.ValueOf(models.LoginDeviceBan_keys{}))
	
	//************************** method insert **********************************
	c.PublicMethods["insert"] = &LoginDeviceBan_Controller_insert{
		osbe.Base_PublicMethod{
			ID: "insert",
			Fields: fields.GenModelMD(reflect.ValueOf(models.LoginDeviceBan{})),
			EventList: osbe.PublicMethodEventList{"LoginDeviceBan.insert"},
		},
	}
	
	//************************** method delete *************************************
	c.PublicMethods["delete"] = &LoginDeviceBan_Controller_delete{
		osbe.Base_PublicMethod{
			ID: "delete",
			Fields: keys_fields,
			EventList: osbe.PublicMethodEventList{"LoginDeviceBan.delete"},
		},
	}
	
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &LoginDeviceBan_Controller_get_object{
		osbe.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},
	}
	
	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &LoginDeviceBan_Controller_get_list{
		osbe.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}
	
			
	
	return c
}

type LoginDeviceBan_Controller_keys_argv struct {
	Argv models.LoginDeviceBan_keys `json:"argv"`	
}

//************************* INSERT **********************************************
//Public method: insert
type LoginDeviceBan_Controller_insert struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *LoginDeviceBan_Controller_insert) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &models.LoginDeviceBan_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* DELETE **********************************************
type LoginDeviceBan_Controller_delete struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *LoginDeviceBan_Controller_delete) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &models.LoginDeviceBan_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET OBJECT **********************************************
type LoginDeviceBan_Controller_get_object struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *LoginDeviceBan_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &models.LoginDeviceBan_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET LIST **********************************************
//Public method: get_list
type LoginDeviceBan_Controller_get_list struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *LoginDeviceBan_Controller_get_list) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &model.Controller_get_list_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

