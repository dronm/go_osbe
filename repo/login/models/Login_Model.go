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

type Login struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true"`
	Date_time_in fields.ValDateTimeTZ `json:"date_time_in"`
	Date_time_out fields.ValDateTimeTZ `json:"date_time_out"`
	Ip fields.ValText `json:"ip"`
	Session_id fields.ValText `json:"session_id"`
	User_id fields.ValInt `json:"user_id"`
	Pub_key fields.ValText `json:"pub_key"`
	Set_date_time fields.ValDateTimeTZ `json:"set_date_time"`
	Headers fields.ValJSON `json:"headers"`
	User_agent fields.ValJSON `json:"user_agent"`
}

func (o *Login) SetNull() {
	o.Id.SetNull()
	o.Date_time_in.SetNull()
	o.Date_time_out.SetNull()
	o.Ip.SetNull()
	o.Session_id.SetNull()
	o.User_id.SetNull()
	o.Pub_key.SetNull()
	o.Set_date_time.SetNull()
	o.Headers.SetNull()
	o.User_agent.SetNull()
}

func NewModelMD_Login() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(Login{})),
		ID: "Login_Model",
		Relation: "logins",
	}
}
//for insert
type Login_argv struct {
	Argv *Login `json:"argv"`	
}

//Keys for delete/get object
type Login_keys struct {
	Id fields.ValInt `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type Login_keys_argv struct {
	Argv *Login_keys `json:"argv"`	
}

//old keys for update
type Login_old_keys struct {
	Old_id fields.ValInt `json:"old_id"`
	Id fields.ValInt `json:"id"`
	Date_time_in fields.ValDateTimeTZ `json:"date_time_in"`
	Date_time_out fields.ValDateTimeTZ `json:"date_time_out"`
	Ip fields.ValText `json:"ip"`
	Session_id fields.ValText `json:"session_id"`
	User_id fields.ValInt `json:"user_id"`
	Pub_key fields.ValText `json:"pub_key"`
	Set_date_time fields.ValDateTimeTZ `json:"set_date_time"`
	Headers fields.ValJSON `json:"headers"`
	User_agent fields.ValJSON `json:"user_agent"`
}

type Login_old_keys_argv struct {
	Argv *Login_old_keys `json:"argv"`	
}

