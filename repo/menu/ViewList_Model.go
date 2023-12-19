package menu

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type ViewList struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true"`
	C fields.ValText `json:"c"`
	F fields.ValText `json:"f"`
	T fields.ValText `json:"t"`
	Href fields.ValText `json:"href"`
	User_descr fields.ValText `json:"user_descr" defOrder:"ASC"`
	Section fields.ValText `json:"section"`
}

func (o *ViewList) SetNull() {
	o.Id.SetNull()
	o.C.SetNull()
	o.F.SetNull()
	o.T.SetNull()
	o.Href.SetNull()
	o.User_descr.SetNull()
	o.Section.SetNull()
}

func NewModelMD_ViewList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(ViewList{})),
		ID: "ViewList_Model",
		Relation: "views_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
		LimitConstant: "doc_per_page_count",
	}
}
