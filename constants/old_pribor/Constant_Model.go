package constants

import (
	"reflect"
	
	"osbe/fields"
	"osbe/model"
)

//Exported model metadata
var (
	ConstantList_md fields.FieldCollection
	Constant_keys_md fields.FieldCollection
	Constant_get_values_md fields.FieldCollection
	Constant_set_value_md fields.FieldCollection
)

func Constant_Model_init() {
	ConstantList_md = fields.GenModelMD(reflect.ValueOf(ConstantList{}))
	Constant_keys_md = fields.GenModelMD(reflect.ValueOf(Constant_keys{}))
	Constant_get_values_md = fields.GenModelMD(reflect.ValueOf(Constant_get_values{}))
	Constant_set_value_md = fields.GenModelMD(reflect.ValueOf(Constant_set_value{}))
}

type Constant_keys struct {
	Id fields.ValText `json:"id" primaryKey:"true" required:"true"`
}

type Constant_keys_argv struct {
	Argv *Constant_keys `json:"argv"`	
}

//Object model for insert
type ConstantList struct {
	Id fields.ValText `json:"id" primaryKey:"true"`
	Name fields.ValText `json:"name"`
	Descr fields.ValText `json:"descr"`
	Val fields.ValText `json:"val"`
}
func (m *ConstantList) GetDataTable() string{
	return "constants_list_view"
}
func (m *ConstantList) GetID() model.ModelID{
	return model.ModelID("ConstantList")
}

type Constant_get_values struct {
	Id_list fields.ValText `json:"id_list" required:"true"`
	Field_sep fields.ValText `json:"field_sep" length:"2"`
}
type Constant_get_values_argv struct {
	Argv *Constant_get_values `json:"argv"`	
}

//
type ConstantValue struct {
	Id fields.ValText `json:"id" primaryKey:"true"`
	Val fields.ValText `json:"val"`
	Val_type fields.ValText `json:"val_type"`
}
func (m *ConstantValue) GetID() model.ModelID{
	return model.ModelID("ConstantValueList")
}

type Constant_set_value struct {
	Id fields.ValText `json:"id" required:"true"`
	Val fields.ValText `json:"val"`
}
type Constant_set_value_argv struct {
	Argv *Constant_set_value `json:"argv"`	
}

