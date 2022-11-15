package response

import (
	"osbe/model"
)

const (
	RESPONSE_MODEL_ID model.ModelID = "Response"
	//"ModelServResponse"	

	RESP_OK = 0
	RESP_ER_AUTH = 100
	RESP_ER_PARSE = 2
	RESP_ER_VALID = 3
	RESP_ER_INTERNAL = 5
	
)

type Response struct {
	Models model.ModelCollection `json:"models"`
}

//*model.Model
func (r *Response) AddModel(model model.Modeler) {	
	r.Models[model.GetID()] = model
}

func (r *Response) AddModelFromStruct(modelID model.ModelID, data interface{}) {	
	rows := make([]model.ModelRow, 1)
	rows[0] = data
	r.AddModel(&model.Model{ID: modelID, Rows: rows})
}

func (r *Response) GetModelCount() int{	
	return len(r.Models)
}

func (r *Response) SetError(code int, descr string) {	
	if m,ok := r.Models[RESPONSE_MODEL_ID]; ok && m.GetRowCount() > 0 {
		fields := m.GetRow(0).(*Response_ModelRow)
		fields.Code = code
		fields.Descr = descr
	}
}

func (r *Response) GetQueryID() string {	
	if m,ok := r.Models[RESPONSE_MODEL_ID]; ok && m.GetRowCount() > 0 {
		fields := m.GetRow(0).(*Response_ModelRow)
		return fields.QueryID
	}
	return ""
}

func (r *Response) GetCode() int {	
	if m,ok := r.Models[RESPONSE_MODEL_ID]; ok && m.GetRowCount() > 0 {
		fields := m.GetRow(0).(*Response_ModelRow)
		return fields.Code
	}
	return 0
}

func (r *Response) GetDescr() string {	
	if m,ok := r.Models[RESPONSE_MODEL_ID]; ok && m.GetRowCount() > 0 {
		fields := m.GetRow(0).(*Response_ModelRow)
		return fields.Descr
	}
	return ""
}

type Response_ModelRow struct {
	Code int `json:"result"`
	Descr string `json:"descr"`
	QueryID string `json:"query_id"`
	AppVersion string `json:"app_version"`
}
/*func (m Response_Model) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}*/

func NewResponse(queryId, appVersion string) *Response{
	resp := &Response{Models: make(model.ModelCollection)}
	rows := make([]model.ModelRow, 1)
	rows[0] = &Response_ModelRow{
		Code: RESP_OK,
		Descr: "",
		QueryID: queryId,
		AppVersion: appVersion,
	}
	resp.Models[RESPONSE_MODEL_ID] = &model.Model{ID: RESPONSE_MODEL_ID, Rows: rows, SysModel: true}
	return resp
}


