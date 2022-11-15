package view

import (
	"fmt"
	"errors"
	"os"
	
	"osbe/response"
	"osbe/socket"
)

const ER_NOT_REGISTERED = "View '%s' not registered"

type View interface {
	Init(viewParams map[string]interface{}) error
	SetParam(paramID string, val interface{}) error
	Render(sock socket.ClientSocketer, resp *response.Response) ([]byte, error)
}

var views = make(map[string]View)

// Register makes a view available by the view name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, view View) {
	if view == nil {
		panic("view: Register provide is nil")
	}
	if _, dup := views[name]; dup {
		panic("view: Register called twice for view " + name)
	}
	views[name] = view
}

func Registered(viewID string) bool {
	if _, ok := views[viewID]; ok {
		return true
	}
	return false
}

func SetParam(viewID, paramID string, val interface{}) error {
	if v, ok := views[viewID]; ok {
		return v.SetParam(paramID, val)
	}else{
		return errors.New(fmt.Sprintf(ER_NOT_REGISTERED, viewID))
	}
}

func Render(viewID string, sock socket.ClientSocketer, resp *response.Response) ([]byte, error) {
	if v, ok := views[viewID]; ok {
		return v.Render(sock, resp)
	}else{
		return nil, errors.New(fmt.Sprintf(ER_NOT_REGISTERED, viewID))
	}
}

func Init(viewID string, viewParams map[string]interface{}) (error) {
	if v, ok := views[viewID]; ok {
		return v.Init(viewParams)
	}else{
		return errors.New(fmt.Sprintf(ER_NOT_REGISTERED, viewID))
	}
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil || !os.IsNotExist(err) {
		return true
	}
	return false
}
