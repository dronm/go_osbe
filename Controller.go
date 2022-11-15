package osbe

import(
	"encoding/json"	
	"errors"	
	"strings"
	"reflect"
	"fmt"
)

//Command function
type ExtCommand struct {
	Func string `json:"func"` //function Controller.method
	Query_id string `json:"query_id"`
	View_id string `json:"view_id"`
}


type Controller interface {
	//InitPublicMethods()
	GetPublicMethod(PublicMethodID) (PublicMethod, error)
	GetID() string
}

type Base_Controller struct {
	ID string
	PublicMethods PublicMethodCollection	
}
func (c *Base_Controller) GetID() string {
	return c.ID
}
func (c *Base_Controller) GetPublicMethod(publicMethodID PublicMethodID) (PublicMethod, error) {
	if pm, ok := c.PublicMethods[publicMethodID]; ok {
		return pm, nil
	}
	return nil, errors.New(fmt.Sprintf(ER_CONTOLLER_METH_NOT_DEFINED, string(publicMethodID), string(c.GetID())))
}


type ControllerCollection map[string] Controller

//parses external command from JSON string of arguments
//argument "func" at least MUST exist!
//Returns:
//	public method interface,
//	argv - all parsed arguments not validated,
//	query id
//	view id
//	error if any
func (c *ControllerCollection) ParseJSONCommand(payload []byte) (Controller, PublicMethod, reflect.Value, string, string, error) {
	//1) unmarshal Controller-method to base function structure		
	var cmd_payload = ExtCommand{}	
	if err := json.Unmarshal(payload, &cmd_payload); err != nil {		
		return nil, nil, reflect.ValueOf(nil), "", "", err
	}		
	//func - Controller.method
	contr, pm, argv, err := c.ParseFunctionCommand(cmd_payload.Func, payload)
	return contr, pm, argv, cmd_payload.Query_id, cmd_payload.View_id, err
	
}

//parses command from function string of type ControllerID.MethodID
func (c *ControllerCollection) ParseFunctionCommand(fn string, argsPayload []byte) (contr Controller, pm PublicMethod, argv reflect.Value, err error) {
	//
	p := strings.Index(fn, ".")
	if p == -1 {
		return nil, nil, reflect.ValueOf(nil), errors.New(ER_PARSE_NO_METH)
	}
		
	return c.ParseCommand(fn[:p], fn[p+1:], argsPayload)
}

//parses command from separated controller,method IDs
func (c *ControllerCollection) ParseCommand(controllerID string, methodID string, argsPayload []byte) (contr Controller, pm PublicMethod, argv reflect.Value, err error) {
	//check controller
	ok := false
	contr, ok = (*c)[controllerID]
	if !ok {
		err = errors.New(fmt.Sprintf(ER_PARSE_CTRL_NOT_DEFINED, controllerID)) 
		return
	}
	
	//check method
	pm, err = contr.GetPublicMethod(PublicMethodID(methodID))
	if err != nil {
		return
	}
	
	//unmarshal params to structure
	argv, err = pm.Unmarshal(argsPayload)
	return
}


