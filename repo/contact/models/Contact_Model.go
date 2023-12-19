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

type Contact struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true" sysCol:"true"`
	Name fields.ValText `json:"name" required:"true" alias:"Наименование"`
	Post_id fields.ValInt `json:"post_id" alias:"Должность"`
	Email fields.ValText `json:"email" alias:"Email"`
	Tel fields.ValText `json:"tel" alias:"Телефон"`
	Tel_ext fields.ValText `json:"tel_ext" alias:"Добавочный номер"`
	Descr fields.ValText `json:"descr" alias:"Описание для поиска" defOrder:"ASC"`
	Comment_text fields.ValText `json:"comment_text" alias:"Комментарий"`
	Email_confirmed fields.ValBool `json:"email_confirmed" alias:"Адрес электр.почты подтвержден"`
	Tel_confirmed fields.ValBool `json:"tel_confirmed" alias:"Номер телефона подтвержден"`
}

func (o *Contact) SetNull() {
	o.Id.SetNull()
	o.Name.SetNull()
	o.Post_id.SetNull()
	o.Email.SetNull()
	o.Tel.SetNull()
	o.Tel_ext.SetNull()
	o.Descr.SetNull()
	o.Comment_text.SetNull()
	o.Email_confirmed.SetNull()
	o.Tel_confirmed.SetNull()
}

func NewModelMD_Contact() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(Contact{})),
		ID: "Contact_Model",
		Relation: "contacts",
	}
}
//for insert
type Contact_argv struct {
	Argv *Contact `json:"argv"`	
}

//Keys for delete/get object
type Contact_keys struct {
	Id fields.ValInt `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type Contact_keys_argv struct {
	Argv *Contact_keys `json:"argv"`	
}

//old keys for update
type Contact_old_keys struct {
	Old_id fields.ValInt `json:"old_id"`
	Id fields.ValInt `json:"id"`
	Name fields.ValText `json:"name" alias:"Наименование"`
	Post_id fields.ValInt `json:"post_id" alias:"Должность"`
	Email fields.ValText `json:"email" alias:"Email"`
	Tel fields.ValText `json:"tel" alias:"Телефон"`
	Tel_ext fields.ValText `json:"tel_ext" alias:"Добавочный номер"`
	Descr fields.ValText `json:"descr" alias:"Описание для поиска"`
	Comment_text fields.ValText `json:"comment_text" alias:"Комментарий"`
	Email_confirmed fields.ValBool `json:"email_confirmed" alias:"Адрес электр.почты подтвержден"`
	Tel_confirmed fields.ValBool `json:"tel_confirmed" alias:"Номер телефона подтвержден"`
}

type Contact_old_keys_argv struct {
	Argv *Contact_old_keys `json:"argv"`	
}

