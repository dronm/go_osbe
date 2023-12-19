package menu

import (
	"reflect"	
	"encoding/json"
	
	"osbe"
	"osbe/fields"
	"osbe/model"
)

//Controller
type MainMenuConstructor_Controller struct {
	osbe.Base_Controller
}

func NewController_MainMenuConstructor() *MainMenuConstructor_Controller{
	c := &MainMenuConstructor_Controller{osbe.Base_Controller{ID: "MainMenuConstructor", PublicMethods: make(osbe.PublicMethodCollection)}}	
	keys_fields := fields.GenModelMD(reflect.ValueOf(MainMenuConstructor_keys{}))
	
	//************************** method insert **********************************
	c.PublicMethods["insert"] = &MainMenuConstructor_Controller_insert{
		osbe.Base_PublicMethod{
			ID: "insert",
			Fields: fields.GenModelMD(reflect.ValueOf(MainMenuConstructor{})),
			EventList: osbe.PublicMethodEventList{"MainMenuConstructor.insert"},
		},
	}
	
	//************************** method delete *************************************
	c.PublicMethods["delete"] = &MainMenuConstructor_Controller_delete{
		osbe.Base_PublicMethod{
			ID: "delete",
			Fields: keys_fields,
			EventList: osbe.PublicMethodEventList{"MainMenuConstructor.delete"},
		},
	}
	
	//************************** method update *************************************
	c.PublicMethods["update"] = &MainMenuConstructor_Controller_update{
		osbe.Base_PublicMethod{
			ID: "update",
			Fields: fields.GenModelMD(reflect.ValueOf(MainMenuConstructor_old_keys{})),
			EventList: osbe.PublicMethodEventList{"MainMenuConstructor.update"},
		},
	}
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &MainMenuConstructor_Controller_get_object{
		osbe.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},
	}
	
	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &MainMenuConstructor_Controller_get_list{
		osbe.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}
	
	return c
}

type MainMenuConstructor_Controller_keys_argv struct {
	Argv MainMenuConstructor_keys `json:"argv"`	
}

//************************* INSERT **********************************************
//Public method: insert
type MainMenuConstructor_Controller_insert struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *MainMenuConstructor_Controller_insert) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &MainMenuConstructor_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* DELETE **********************************************
type MainMenuConstructor_Controller_delete struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *MainMenuConstructor_Controller_delete) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &MainMenuConstructor_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET OBJECT **********************************************
type MainMenuConstructor_Controller_get_object struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *MainMenuConstructor_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &MainMenuConstructor_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET LIST **********************************************
//Public method: get_list
type MainMenuConstructor_Controller_get_list struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *MainMenuConstructor_Controller_get_list) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &model.Controller_get_list_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* UPDATE **********************************************
//Public method: update
type MainMenuConstructor_Controller_update struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *MainMenuConstructor_Controller_update) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &MainMenuConstructor_old_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

