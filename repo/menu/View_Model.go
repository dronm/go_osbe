package menu

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type View struct {
	Id fields.ValInt `json:"id" primaryKey:"true" autoInc:"true" alias:"Код"`
	C fields.ValText `json:"c"`
	F fields.ValText `json:"f"`
	T fields.ValText `json:"t"`
	Section fields.ValText `json:"section" required:"true" defOrder:"ASC"`
	Descr fields.ValText `json:"descr" required:"true" defOrder:"ASC"`
	Limited fields.ValBool `json:"limited"`
}

func (o *View) SetNull() {
	o.Id.SetNull()
	o.C.SetNull()
	o.F.SetNull()
	o.T.SetNull()
	o.Section.SetNull()
	o.Descr.SetNull()
	o.Limited.SetNull()
}

func NewModelMD_View() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(View{})),
		ID: "View_Model",
		Relation: "views",
	}
}
//for insert
type View_argv struct {
	Argv *View `json:"argv"`	
}

//Keys for delete/get object
type View_keys struct {
	Id fields.ValInt `json:"id"`
	Mode string `json:"mode" openMode:"true"` //open mode insert|copy|edit
}
type View_keys_argv struct {
	Argv *View_keys `json:"argv"`	
}

//old keys for update
type View_old_keys struct {
	Old_id fields.ValInt `json:"old_id" alias:"Код"`
	Id fields.ValInt `json:"id" alias:"Код"`
	C fields.ValText `json:"c"`
	F fields.ValText `json:"f"`
	T fields.ValText `json:"t"`
	Section fields.ValText `json:"section"`
	Descr fields.ValText `json:"descr"`
	Limited fields.ValBool `json:"limited"`
}

type View_old_keys_argv struct {
	Argv *View_old_keys `json:"argv"`	
}

