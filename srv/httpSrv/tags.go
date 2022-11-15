package httpSrv

import(
	"osbe/model"
)

//HTML tags: TagScript (javascript), TagLink (css)

const (
	SCRIPT_MODEL_ID model.ModelID = "Script"
	LINK_MODEL_ID model.ModelID = "Link"
)

//script tag
type TagScript struct {
	Src string `json:"src"`
	Type string `json:"type" xml:"omitempty`
	Defer bool `json:"defer" xml:"omitempty`
	Language string `json:"language" xml:"omitempty`
}

//link tag
type TagLink struct {
	Charset string `json:"charset" xml:"omitempty`
	Href string `json:"href"`
	Media string `json:"media" xml:"omitempty`
	Rel string `json:"rel" xml:"omitempty`
	Sizes string `json:"sizes" xml:"omitempty` //widthxheight | widthXheight | any
	Type string `json:"type" xml:"omitempty`
}

func NewScriptModel(rowCount int) *model.Model{
	m := &model.Model{ID: SCRIPT_MODEL_ID, SysModel: true, Rows: make([]model.ModelRow, rowCount)}
	return m
}

func NewLinkModel(rowCount int) *model.Model{
	m := &model.Model{ID: LINK_MODEL_ID, SysModel: true, Rows: make([]model.ModelRow, rowCount)}
	return m
}

