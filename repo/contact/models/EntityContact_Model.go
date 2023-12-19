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

type EntityContact struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true" sysCol:"true"`
	Entity_type fields.ValText `json:"entity_type" required:"true"`
	Entity_id fields.ValInt `json:"entity_id" required:"true"`
	Contact_id fields.ValInt `json:"contact_id" required:"true"`
	Mod_date_time fields.ValDateTimeTZ `json:"mod_date_time"`
}

func (o *EntityContact) SetNull() {
	o.Id.SetNull()
	o.Entity_type.SetNull()
	o.Entity_id.SetNull()
	o.Contact_id.SetNull()
	o.Mod_date_time.SetNull()
}

func NewModelMD_EntityContact() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(EntityContact{})),
		ID: "EntityContact_Model",
		Relation: "entity_contacts",
	}
}
//for insert
type EntityContact_argv struct {
	Argv *EntityContact `json:"argv"`	
}

//Keys for delete/get object
type EntityContact_keys struct {
	Id fields.ValInt `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type EntityContact_keys_argv struct {
	Argv *EntityContact_keys `json:"argv"`	
}

//old keys for update
type EntityContact_old_keys struct {
	Old_id fields.ValInt `json:"old_id"`
	Id fields.ValInt `json:"id"`
	Entity_type fields.ValText `json:"entity_type"`
	Entity_id fields.ValInt `json:"entity_id"`
	Contact_id fields.ValInt `json:"contact_id"`
	Mod_date_time fields.ValDateTimeTZ `json:"mod_date_time"`
}

type EntityContact_old_keys_argv struct {
	Argv *EntityContact_old_keys `json:"argv"`	
}

