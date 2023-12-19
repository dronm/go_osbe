package userOperation

import (
	"reflect"	
	"encoding/json"
	
	"osbe"
	"osbe/fields"
	
)

//Controller
type UserOperation_Controller struct {
	osbe.Base_Controller
}

func NewController_UserOperation() *UserOperation_Controller{
	c := &UserOperation_Controller{osbe.Base_Controller{ID: "UserOperation", PublicMethods: make(osbe.PublicMethodCollection)}}	
	keys_fields := fields.GenModelMD(reflect.ValueOf(UserOperation_keys{}))
	
	
	//************************** method delete *************************************
	c.PublicMethods["delete"] = &UserOperation_Controller_delete{
		osbe.Base_PublicMethod{
			ID: "delete",
			Fields: keys_fields,
			EventList: osbe.PublicMethodEventList{"UserOperation.delete"},
		},
	}
	
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &UserOperation_Controller_get_object{
		osbe.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},
	}
	
	
			
	
	return c
}

type UserOperation_Controller_keys_argv struct {
	Argv UserOperation_keys `json:"argv"`	
}


//************************* DELETE **********************************************
type UserOperation_Controller_delete struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *UserOperation_Controller_delete) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &UserOperation_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET OBJECT **********************************************
type UserOperation_Controller_get_object struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *UserOperation_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &UserOperation_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}



