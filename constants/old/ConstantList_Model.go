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

const ConstantList_DATA_TABLE = "constants_list_view"

//Exported model metadata
var (
	ConstantList_md fields.FieldCollection	
)

func ConstantList_Model_init() {	
	ConstantList_md = fields.GenModelMD(reflect.ValueOf(ConstantList{}))
}

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
func (m *ConstantList) GetDataTable() string{
	return ConstantList_DATA_TABLE
}
func (m *ConstantList) GetID() model.ModelID{
	return model.ModelID("ConstantList_Model")
}
