package bank

import (
	"reflect"	
	"encoding/json"
	
	models "osbe/repo/bank/models"
	
	"osbe"
	"osbe/fields"
	"osbe/model"
)

//Controller
type Bank_Controller struct {
	osbe.Base_Controller
}

func NewController_Bank() *Bank_Controller{
	c := &Bank_Controller{osbe.Base_Controller{ID: "Bank", PublicMethods: make(osbe.PublicMethodCollection)}}	
	keys_fields := fields.GenModelMD(reflect.ValueOf(models.Bank_keys{}))
	
	
	
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &Bank_Controller_get_object{
		osbe.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},
	}
	
	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &Bank_Controller_get_list{
		osbe.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}
	
	//************************** method complete *************************************
	c.PublicMethods["complete"] = &Bank_Controller_complete{
		osbe.Base_PublicMethod{
			ID: "complete",
			Fields: fields.GenModelMD(reflect.ValueOf(models.Bank_complete{})),
		},
	}
			
	
	return c
}

type Bank_Controller_keys_argv struct {
	Argv models.Bank_keys `json:"argv"`	
}



//************************* GET OBJECT **********************************************
type Bank_Controller_get_object struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *Bank_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &models.Bank_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET LIST **********************************************
//Public method: get_list
type Bank_Controller_get_list struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Bank_Controller_get_list) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &model.Controller_get_list_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}


//************************** COMPLETE ********************************************************
//Public method: complete
type Bank_Controller_complete struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *Bank_Controller_complete) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &models.Bank_complete_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}
