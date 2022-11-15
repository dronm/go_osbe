package constants

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 */

import (
	"reflect"
	
		
	"osbe/fields"
	"osbe/model"
)

//
type ConstantList struct {
	Id fields.ValText `json:"id" primaryKey:"true"`
	Name fields.ValText `json:"name"`
	Descr fields.ValText `json:"descr"`
	Val fields.ValText `json:"val"`
	Val_type fields.ValText `json:"val_type"`
	Ctrl_class fields.ValText `json:"ctrl_class"`
	Ctrl_options fields.ValJSON `json:"ctrl_options"`
	View_class fields.ValText `json:"view_class"`
	View_options fields.ValJSON `json:"view_options"`
}

func NewModelMD_ConstantList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(ConstantList{})),
		ID: "ConstantList_Model",
		Relation: "constants_list_view",
	}
}

