package osbe

import (
	"reflect"
	"fmt"
	"errors"	
	
	"osbe/fields"
	"osbe/socket"
	"osbe/srv"
	"osbe/response"
)

//public method
type PublicMethodID string
type PublicMethodCollection map[PublicMethodID] PublicMethod
type PublicMethodEventList []string

type PublicMethod interface {
	GetID() PublicMethodID
	GetFields() fields.FieldCollection
	Unmarshal(payload []byte) (reflect.Value, error)
	Run(Applicationer, srv.Server, socket.ClientSocketer, *response.Response, reflect.Value) error
	GetEventList() PublicMethodEventList
	//AddEvent(string)
}

type Base_PublicMethod struct {
	ID string
	Fields fields.FieldCollection
	EventList PublicMethodEventList
}

func (pm *Base_PublicMethod) GetID() PublicMethodID {
	return PublicMethodID(pm.ID)
}

func (pm *Base_PublicMethod) GetEventList() PublicMethodEventList {
	return pm.EventList
}

func (pm *Base_PublicMethod) AddEvent(evId string) {
	pm.EventList[len(pm.EventList)-1] = evId
}

func (pm *Base_PublicMethod) GetFields() fields.FieldCollection {
	return pm.Fields
}

type PublicMethodWithEvent interface {
	GetEventList() PublicMethodEventList
	AddEvent(string)
}

type PublicMethodError struct {
    Code int
    Err error
}
func (e PublicMethodError) Error() string {
	//e.Code
	return fmt.Sprintf("%v", e.Err)
}

func NewPublicMethodError(code int, err string) *PublicMethodError{
	return &PublicMethodError{Code: code, Err: errors.New(err)}
}


