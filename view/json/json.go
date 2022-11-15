package viewJSON

import (
	"encoding/json"

	"osbe/view"
	"osbe/response"
	"osbe/socket"
)

const VIEW_ID = "ViewJSON"

var v = &ViewJSON{}

type ViewJSON struct {
}

func (v *ViewJSON) Init(map[string]interface{}) (err error) {
	return nil
}

func (v *ViewJSON)  SetParam(string, interface{}) error {
	return nil
}

func (v *ViewJSON) Render(_ socket.ClientSocketer, resp *response.Response) ([]byte, error){
	return json.Marshal(resp)
}

func init() {
	view.Register(VIEW_ID, v)
}
