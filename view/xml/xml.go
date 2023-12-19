package viewXML

import (
	"errors"
	
	"osbe/xml"
	"osbe/view"
	"osbe/response"
	"osbe/socket"
)

const (
	VIEW_ID = "ViewXML"
)

type OnBeforeRenderProto = func(socket.ClientSocketer, *response.Response) error

var v = &ViewXML{}

type ViewXML struct {
	BeforeRender OnBeforeRenderProto
}
//Parameters:
//		BeforeRender(OnBeforeRenderProto)
func (v *ViewXML) Init(params map[string]interface {}) error {
	for id, val := range params {
		if err := v.SetParam(id, val); err != nil {
			return err
		}
	}
	return nil
}

func (v *ViewXML)  SetParam(paramID string, val interface{}) error {
	ok := false
	switch paramID {
	case "BeforeRender":
		if v.BeforeRender, ok = val.(OnBeforeRenderProto); !ok {
			return errors.New("parameter BeforeRender must be of OnBeforeRenderProto type")
		}
	}
	return nil
}

//All models from Response to xml
func (v *ViewXML) Render(sock socket.ClientSocketer, resp *response.Response) ([]byte, error){
	if v.BeforeRender != nil {
		if err := v.BeforeRender(sock, resp); err != nil {
			return nil, err
		}
	}
	
	return xml.Marshal(resp.Models, false)
}

func init() {
	view.Register(VIEW_ID, v)
}

