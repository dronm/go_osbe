package models

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

type SessionVal struct {
	Id fields.ValText `json:"id" primaryKey:"true"`
	Accessed_time fields.ValDateTimeTZ `json:"accessed_time"`
	Create_time fields.ValDateTimeTZ `json:"create_time"`
	Val fields.ValBytea `json:"val"`
}

func (o *SessionVal) SetNull() {
	o.Id.SetNull()
	o.Accessed_time.SetNull()
	o.Create_time.SetNull()
	o.Val.SetNull()
}

func NewModelMD_SessionVal() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(SessionVal{})),
		ID: "SessionVal_Model",
		Relation: "session_vals",
	}
}
//for insert
type SessionVal_argv struct {
	Argv *SessionVal `json:"argv"`	
}

//Keys for delete/get object
type SessionVal_keys struct {
	Id fields.ValText `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type SessionVal_keys_argv struct {
	Argv *SessionVal_keys `json:"argv"`	
}

//old keys for update
type SessionVal_old_keys struct {
	Old_id fields.ValText `json:"old_id"`
	Id fields.ValText `json:"id"`
	Accessed_time fields.ValDateTimeTZ `json:"accessed_time"`
	Create_time fields.ValDateTimeTZ `json:"create_time"`
	Val fields.ValBytea `json:"val"`
}

type SessionVal_old_keys_argv struct {
	Argv *SessionVal_old_keys `json:"argv"`	
}

