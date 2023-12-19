package notif

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type NotifTemplateList struct {
	Id fields.ValInt `json:"id" primaryKey:"true"`
	Notif_provider fields.ValText `json:"notif_provider"`
	Notif_type fields.ValText `json:"notif_type"`
	Template fields.ValText `json:"template"`
}

func (o *NotifTemplateList) SetNull() {
	o.Id.SetNull()
	o.Notif_provider.SetNull()
	o.Notif_type.SetNull()
	o.Template.SetNull()
}

func NewModelMD_NotifTemplateList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(NotifTemplateList{})),
		ID: "NotifTemplateList_Model",
		Relation: "notif_templates_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
	}
}
