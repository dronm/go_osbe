package contact

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

type EntityContactList struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true" sysCol:"true"`
	Entity_type fields.ValText `json:"entity_type"`
	Entity_id fields.ValInt `json:"entity_id"`
	Entities_ref fields.ValJSON `json:"entities_ref"`
	Contact_id fields.ValInt `json:"contact_id"`
	Contact_attrs fields.ValJSON `json:"contact_attrs"`
	Contacts_ref fields.ValJSON `json:"contacts_ref"`
	Tm_exists fields.ValBool `json:"tm_exists"`
	Tm_activated fields.ValBool `json:"tm_activated"`
}

func (o *EntityContactList) SetNull() {
	o.Id.SetNull()
	o.Entity_type.SetNull()
	o.Entity_id.SetNull()
	o.Entities_ref.SetNull()
	o.Contact_id.SetNull()
	o.Contact_attrs.SetNull()
	o.Contacts_ref.SetNull()
	o.Tm_exists.SetNull()
	o.Tm_activated.SetNull()
}

func NewModelMD_EntityContactList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(EntityContactList{})),
		ID: "EntityContactList_Model",
		Relation: "entity_contacts_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
		
	}
}
