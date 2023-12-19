package menu

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type MainMenuConstructorList struct {
	Id fields.ValInt `json:"id" primaryKey:"true" sysCol:"true"`
	Role_id fields.ValText `json:"role_id"`
	User_id fields.ValText `json:"user_id"`
	User_descr fields.ValText `json:"user_descr"`
}

func (o *MainMenuConstructorList) SetNull() {
	o.Id.SetNull()
	o.Role_id.SetNull()
	o.User_id.SetNull()
	o.User_descr.SetNull()
}

func NewModelMD_MainMenuConstructorList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(MainMenuConstructorList{})),
		ID: "MainMenuConstructorList_Model",
		Relation: "main_menus_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},		
	}
}
