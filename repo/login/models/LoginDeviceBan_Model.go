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

type LoginDeviceBan struct {
	User_id fields.ValInt `json:"user_id" primaryKey:"true" sysCol:"true"`
	Hash fields.ValText `json:"hash" primaryKey:"true"`
	Create_dt fields.ValDateTimeTZ `json:"create_dt" alias:"Дата создания"`
}

func (o *LoginDeviceBan) SetNull() {
	o.User_id.SetNull()
	o.Hash.SetNull()
	o.Create_dt.SetNull()
}

func NewModelMD_LoginDeviceBan() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(LoginDeviceBan{})),
		ID: "LoginDeviceBan_Model",
		Relation: "login_device_bans",
	}
}
//for insert
type LoginDeviceBan_argv struct {
	Argv *LoginDeviceBan `json:"argv"`	
}

//Keys for delete/get object
type LoginDeviceBan_keys struct {
	User_id fields.ValInt `json:"user_id"`
	Hash fields.ValText `json:"hash"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type LoginDeviceBan_keys_argv struct {
	Argv *LoginDeviceBan_keys `json:"argv"`	
}

//old keys for update
type LoginDeviceBan_old_keys struct {
	Old_user_id fields.ValInt `json:"old_user_id"`
	User_id fields.ValInt `json:"user_id"`
	Old_hash fields.ValText `json:"old_hash"`
	Hash fields.ValText `json:"hash"`
	Create_dt fields.ValDateTimeTZ `json:"create_dt" alias:"Дата создания"`
}

type LoginDeviceBan_old_keys_argv struct {
	Argv *LoginDeviceBan_old_keys `json:"argv"`	
}

