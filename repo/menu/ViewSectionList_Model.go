package menu

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type ViewSectionList struct {
	Section fields.ValText `json:"section"`
}

func (o *ViewSectionList) SetNull() {
	o.Section.SetNull()
}

func NewModelMD_ViewSectionList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(ViewSectionList{})),
		ID: "ViewSectionList_Model",
		Relation: "views_section_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
		LimitConstant: "doc_per_page_count",
	}
}
