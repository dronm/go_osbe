package menu

import (
	"reflect"	
	"nails/enums"	
	"osbe/fields"
	"osbe/model"
)

type MainMenuConstructor struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true" alias:"Код"`
	Role_id enums.ValEnum_role_types `json:"role_id" required:"true" defOrder:"ASC"`
	User_id fields.ValInt `json:"user_id"`
	Content fields.ValText `json:"content" required:"true" alias:"Содержание"`
	Model_content fields.ValText `json:"model_content" alias:"Содержание для модели,заполняется при записи из контроллера!"`
}

func (o *MainMenuConstructor) SetNull() {
	o.Id.SetNull()
	o.Role_id.SetNull()
	o.User_id.SetNull()
	o.Content.SetNull()
	o.Model_content.SetNull()
}

func NewModelMD_MainMenuConstructor() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(MainMenuConstructor{})),
		ID: "MainMenuConstructor_Model",
		Relation: "main_menus",
	}
}
//for insert
type MainMenuConstructor_argv struct {
	Argv *MainMenuConstructor `json:"argv"`	
}

//Keys for delete/get object
type MainMenuConstructor_keys struct {
	Id fields.ValInt `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type MainMenuConstructor_keys_argv struct {
	Argv *MainMenuConstructor_keys `json:"argv"`	
}

//old keys for update
type MainMenuConstructor_old_keys struct {
	Old_id fields.ValInt `json:"old_id" alias:"Код"`
	Id fields.ValInt `json:"id" alias:"Код"`
	Role_id enums.ValEnum_role_types `json:"role_id"`
	User_id fields.ValInt `json:"user_id"`
	Content fields.ValText `json:"content" alias:"Содержание"`
	Model_content fields.ValText `json:"model_content" alias:"Содержание для модели,заполняется при записи из контроллера!"`
}

type MainMenuConstructor_old_keys_argv struct {
	Argv *MainMenuConstructor_old_keys `json:"argv"`	
}

