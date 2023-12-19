package notif

import(
	"reflect"
	"encoding/json"
	
	"osbe"
	"osbe/fields"
	"osbe/model"
)

//Controller
type NotifTemplate_Controller struct {
	osbe.Base_Controller
}

func NewController_NotifTemplate() *NotifTemplate_Controller{
	c := &NotifTemplate_Controller{osbe.Base_Controller{ID: "NotifTemplate", PublicMethods: make(osbe.PublicMethodCollection)}}	
	keys_fields := fields.GenModelMD(reflect.ValueOf(NotifTemplate_keys{}))
	
	//************************** method insert **********************************
	c.PublicMethods["insert"] = &NotifTemplate_Controller_insert{
		osbe.Base_PublicMethod{
			ID: "insert",
			Fields: fields.GenModelMD(reflect.ValueOf(NotifTemplate{})),
			EventList: osbe.PublicMethodEventList{"NotifTemplate.insert"},
		},
	}
	
	//************************** method delete *************************************
	c.PublicMethods["delete"] = &NotifTemplate_Controller_delete{
		osbe.Base_PublicMethod{
			ID: "delete",
			Fields: keys_fields,
			EventList: osbe.PublicMethodEventList{"NotifTemplate.delete"},
		},
	}
	
	//************************** method update *************************************
	c.PublicMethods["update"] = &NotifTemplate_Controller_update{
		osbe.Base_PublicMethod{
			ID: "update",
			Fields: fields.GenModelMD(reflect.ValueOf(NotifTemplate_old_keys{})),
			EventList: osbe.PublicMethodEventList{"NotifTemplate.update"},
		},
	}
	
	//************************** method get_object *************************************
	c.PublicMethods["get_object"] = &NotifTemplate_Controller_get_object{
		osbe.Base_PublicMethod{
			ID: "get_object",
			Fields: keys_fields,
		},
	}
	
	//************************** method get_list *************************************
	c.PublicMethods["get_list"] = &NotifTemplate_Controller_get_list{
		osbe.Base_PublicMethod{
			ID: "get_list",
			Fields: model.Cond_Model_fields,
		},
	}
	
			
	
	return c
}

type NotifTemplate_Controller_keys_argv struct {
	Argv NotifTemplate_keys `json:"argv"`	
}

//************************* INSERT **********************************************
//Public method: insert
type NotifTemplate_Controller_insert struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_insert) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &NotifTemplate_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* DELETE **********************************************
type NotifTemplate_Controller_delete struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_delete) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &NotifTemplate_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET OBJECT **********************************************
type NotifTemplate_Controller_get_object struct {
	osbe.Base_PublicMethod
}

//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_get_object) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &NotifTemplate_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

//************************* GET LIST **********************************************
//Public method: get_list
type NotifTemplate_Controller_get_list struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_get_list) Unmarshal(payload []byte) (reflect.Value, error) {
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
type NotifTemplate_Controller_update struct {
	osbe.Base_PublicMethod
}
//Public method Unmarshal to structure
func (pm *NotifTemplate_Controller_update) Unmarshal(payload []byte) (reflect.Value, error) {
	var res reflect.Value
	argv := &NotifTemplate_old_keys_argv{}
		
	if err := json.Unmarshal(payload, argv); err != nil {
		return res, err
	}	
	res = reflect.ValueOf(&argv.Argv).Elem()	
	return res, nil
}

