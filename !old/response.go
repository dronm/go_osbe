package osbe

import (
)

const (
	RESPONSE_MODEL_ID ModelID = "ModelServResponse"

	RESP_OK = 0
	RESP_ER_AUTH = 100
	RESP_ER_PARSE = 2
	RESP_ER_VALID = 3
	RESP_ER_EXEC = 4
	RESP_ER_INTERNAL = 5
	
)

type Response struct {
	Models ModelCollection `json:"models"`
}

func (r *Response) AddModel(model *Model) {	
	r.Models[model.ID] = model
}

func (r *Response) GetModelCount() int{	
	return len(r.Models)
}

func (r *Response) SetError(code int, descr string) {	
	if m,ok := r.Models[RESPONSE_MODEL_ID]; ok && len(m.Rows) > 0 {
		fields := m.Rows[0].(*Response_Model_row)
		fields.Code = code
		fields.Descr = descr
	}
}

func (r *Response) GetQueryID() string {	
	if m,ok := r.Models[RESPONSE_MODEL_ID]; ok && len(m.Rows) > 0 {
		fields := m.Rows[0].(*Response_Model_row)
		return fields.QueryID
	}
	return ""
}

type Response_Model_row struct {
	Code int `json:"result"`
	Descr string `json:"descr"`
	QueryID string `json:"query_id"`
	AppVersion string `json:"app_version"`
}
/*func (m Response_Model) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}*/

func NewResponse(queryId, appVersion string) *Response{
	resp := &Response{Models: make(ModelCollection)}
	resp.Models[RESPONSE_MODEL_ID] = &Model{ID: RESPONSE_MODEL_ID, Rows:make([]ModelRow,1)}
	resp.Models[RESPONSE_MODEL_ID].Rows[0] = &Response_Model_row{
		Code: RESP_OK,
		Descr: "",
		QueryID: queryId,
		AppVersion: appVersion,
	}
	return resp
}

