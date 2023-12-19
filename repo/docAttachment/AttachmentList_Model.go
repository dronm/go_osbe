package docAttachment

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 *
 * THIS FILE IS GENERATED FROM TEMPLATE build/templates/models/Model.go.tmpl
 * ALL DIRECT MODIFICATIONS WILL BE LOST WITH THE NEXT BUILD PROCESS!!!
 */

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type AttachmentList struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true"`
	Date_time fields.ValDateTimeTZ `json:"date_time" defOrder:"DESC"`
	Ref fields.ValJSON `json:"ref"`
	Content_info fields.ValJSON `json:"content_info"`
	Content_preview fields.ValText `json:"content_preview"`
}

func (o *AttachmentList) SetNull() {
	o.Id.SetNull()
	o.Date_time.SetNull()
	o.Ref.SetNull()
	o.Content_info.SetNull()
	o.Content_preview.SetNull()
}

func NewModelMD_AttachmentList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(AttachmentList{})),
		ID: "AttachmentList_Model",
		Relation: "attachments_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
		
	}
}
