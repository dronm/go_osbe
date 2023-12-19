package menu

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type MainMenuConstructorDialog struct {
	Id fields.ValInt `json:"id" primaryKey:"true" sysCol:"true"`
	Role_id fields.ValText `json:"role_id"`
	User_id fields.ValText `json:"user_id"`
	User_descr fields.ValText `json:"user_descr"`
	Content fields.ValText `json:"content"`
}

func (o *MainMenuConstructorDialog) SetNull() {
	o.Id.SetNull()
	o.Role_id.SetNull()
	o.User_id.SetNull()
	o.User_descr.SetNull()
	o.Content.SetNull()
}

func NewModelMD_MainMenuConstructorDialog() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(MainMenuConstructorDialog{})),
		ID: "MainMenuConstructorDialog_Model",
		Relation: "main_menus_dialog",
	}
}
