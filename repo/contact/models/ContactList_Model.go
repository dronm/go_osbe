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

type ContactList struct {
	Id fields.ValInt `json:"id" primaryKey:"true"`
	Name fields.ValText `json:"name" alias:"Наименование"`
	Posts_ref fields.ValJSON `json:"posts_ref"`
	Email fields.ValText `json:"email" alias:"Email"`
	Tel fields.ValText `json:"tel" alias:"Телефон"`
	Tel_ext fields.ValText `json:"tel_ext" alias:"Добавочный номер"`
	Descr fields.ValText `json:"descr" alias:"Описание для поиска" defOrder:"ASC"`
	Comment_text fields.ValText `json:"comment_text" alias:"Комментарий"`
}

func (o *ContactList) SetNull() {
	o.Id.SetNull()
	o.Name.SetNull()
	o.Posts_ref.SetNull()
	o.Email.SetNull()
	o.Tel.SetNull()
	o.Tel_ext.SetNull()
	o.Descr.SetNull()
	o.Comment_text.SetNull()
}

func NewModelMD_ContactList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(ContactList{})),
		ID: "ContactList_Model",
		Relation: "contacts_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
		
	}
}
