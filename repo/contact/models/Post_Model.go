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

type Post struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true" sysCol:"true"`
	Name fields.ValText `json:"name" required:"true" alias:"Наименование" defOrder:"ASC"`
}

func (o *Post) SetNull() {
	o.Id.SetNull()
	o.Name.SetNull()
}

func NewModelMD_Post() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(Post{})),
		ID: "Post_Model",
		Relation: "posts",
	}
}
//for insert
type Post_argv struct {
	Argv *Post `json:"argv"`	
}

//Keys for delete/get object
type Post_keys struct {
	Id fields.ValInt `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type Post_keys_argv struct {
	Argv *Post_keys `json:"argv"`	
}

//old keys for update
type Post_old_keys struct {
	Old_id fields.ValInt `json:"old_id"`
	Id fields.ValInt `json:"id"`
	Name fields.ValText `json:"name" alias:"Наименование"`
}

type Post_old_keys_argv struct {
	Argv *Post_old_keys `json:"argv"`	
}

