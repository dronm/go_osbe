package menu

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type MainMenuContent struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true"`
	Descr fields.ValText `json:"descr"`
	Viewid fields.ValInt `json:"viewid"`
	Viewdescr fields.ValText `json:"viewdescr"`
	Default fields.ValBool `json:"default"`
	Glyphclass fields.ValText `json:"glyphclass"`
}

func (o *MainMenuContent) SetNull() {
	o.Id.SetNull()
	o.Descr.SetNull()
	o.Viewid.SetNull()
	o.Viewdescr.SetNull()
	o.Default.SetNull()
	o.Glyphclass.SetNull()
}

func NewModelMD_MainMenuContent() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(MainMenuContent{})),
		ID: "MainMenuContent_Model",
		Relation: "",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
	}
}
