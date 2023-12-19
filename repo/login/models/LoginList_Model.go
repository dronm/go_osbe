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

type LoginList struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true" sysCol:"true"`
	Date_time_in fields.ValDateTimeTZ `json:"date_time_in" alias:"Дата входа"`
	Date_time_out fields.ValDateTimeTZ `json:"date_time_out" alias:"Дата выхода"`
	Ip fields.ValText `json:"ip" alias:"IP адрес"`
	User_id fields.ValInt `json:"user_id" sysCol:"true"`
	Users_ref fields.ValJSON `json:"users_ref" alias:"Пользователь"`
	Pub_key fields.ValText `json:"pub_key" sysCol:"true"`
	Set_date_time fields.ValDateTimeTZ `json:"set_date_time" alias:"Время последнего обращения"`
	User_agent_descr fields.ValText `json:"user_agent_descr"`
	User_agent fields.ValJSON `json:"user_agent"`
	Headers fields.ValJSON `json:"headers"`
}

func (o *LoginList) SetNull() {
	o.Id.SetNull()
	o.Date_time_in.SetNull()
	o.Date_time_out.SetNull()
	o.Ip.SetNull()
	o.User_id.SetNull()
	o.Users_ref.SetNull()
	o.Pub_key.SetNull()
	o.Set_date_time.SetNull()
	o.User_agent_descr.SetNull()
	o.User_agent.SetNull()
	o.Headers.SetNull()
}

func NewModelMD_LoginList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(LoginList{})),
		ID: "LoginList_Model",
		Relation: "logins_list",
	}
}
