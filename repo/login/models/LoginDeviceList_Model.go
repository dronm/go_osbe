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

type LoginDeviceList struct {
	User_id fields.ValInt `json:"user_id" primaryKey:"true" sysCol:"true"`
	User_agent fields.ValText `json:"user_agent" primaryKey:"true"`
	User_descr fields.ValText `json:"user_descr" alias:"Имя пользователя"`
	Date_time_in fields.ValDateTimeTZ `json:"date_time_in" alias:"Дата входа"`
	Banned fields.ValBool `json:"banned" alias:"Вход запрещен"`
	Ban_hash fields.ValText `json:"ban_hash"`
}

func (o *LoginDeviceList) SetNull() {
	o.User_id.SetNull()
	o.User_agent.SetNull()
	o.User_descr.SetNull()
	o.Date_time_in.SetNull()
	o.Banned.SetNull()
	o.Ban_hash.SetNull()
}

func NewModelMD_LoginDeviceList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(LoginDeviceList{})),
		ID: "LoginDeviceList_Model",
		Relation: "login_devices_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
		
	}
}
