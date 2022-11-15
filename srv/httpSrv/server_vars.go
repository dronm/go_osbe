package httpSrv

import(
	"osbe/model"
)

const SERVER_VARS_MODEL_ID model.ModelID = "ServerVars"

type ServerVars struct {
	//BasePath string `json:"basePath"`
	ScriptId string `json:"scriptId"`
	Title string `json:"title"`
	Author string `json:"author"`
	Keywords string `json:"keywords"`
	Description string `json:"description"`
	CurDate int64 `json:"curDate"`
	Version string `json:"version"`
	Locale_id string `json:"locale_id"`
	Debug int `json:"debug"`
}


func NewServerVarsModel(row model.ModelRow) *model.Model{
	m := &model.Model{ID: SERVER_VARS_MODEL_ID, SysModel: true, Rows: make([]model.ModelRow, 1)}
	m.Rows[0] = row
	return m
}
