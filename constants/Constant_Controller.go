// Constants package manages application constants.
// The package includes a controller and models. The constant controller
// supports methods:
//	get_list - for getting list of all constants
//	get_object - detailed information about one constant
//	set_value for updating constant value. 
//	get_values - for fetching values of several constants
//
// OSBE project by Andrey Mikhalevich
package constants

import (
	"encoding/json"
	"reflect"
	
	"osbe"
	"osbe/fields"
	"osbe/model"
)

const (
	RESP_ER_NOT_FOUND = 1000	
)

//Controller
type Constant_Controller struct {
	osbe.Base_Controller
}

func NewController_Constant() *Constant_Controller{
	c := &Constant_Controller{osbe.Base_Controller{ID: "Constant", PublicMethods: make(osbe.PublicMethodCollection)}}

	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &Constant_Controller_get_list{
		osbe.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}

	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &Constant_Controller_get_object{
		osbe.Base_PublicMethod{
			ID: "get_object",
			Fields: fields.GenModelMD(reflect.ValueOf(Constant_keys{})),
		},
	}
	
	//************************** method set_value *************************************
	c.PublicMethods["set_value"] = &Constant_Controller_set_value{
		osbe.Base_PublicMethod{
			ID: "set_value",	
			Fields: fields.GenModelMD(reflect.ValueOf(Constant_set_value{})),
		},
	}
	
	//************************** method get_values *************************************
	c.PublicMethods["get_values"] = &Constant_Controller_get_values{
		osbe.Base_PublicMethod{
			ID: "get_values",
			Fields: fields.GenModelMD(reflect.ValueOf(Constant_get_values{})),
		},
	}
	
	return c
}

type Constant_keys_argv struct {
	Argv Constant_keys `json:"argv"`	
}

//************************* GET LIST **********************************************
//Public method: get_list
type Constant_Controller_get_list struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Constant_Controller_get_list) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &model.Controller_get_list_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//Public method: get_object
type Constant_Controller_get_object struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Constant_Controller_get_object) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Constant_keys_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//*******************************************************************************************************
//Public method: set_value
type Constant_Controller_set_value struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Constant_Controller_set_value) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Constant_set_value_argv{}
	
	//json values will raise errors!
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

//*******************************************************************************************************
//Public method: get_values
type Constant_Controller_get_values struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Constant_Controller_get_values) Unmarshal(payload []byte) (res reflect.Value, err error) {

	//argument structrure
	argv := &Constant_get_values_argv{}
	
	err = json.Unmarshal(payload, argv)
	if err != nil {
		return 
	}
	
	res = reflect.ValueOf(&argv.Argv).Elem()
	
	return
}

