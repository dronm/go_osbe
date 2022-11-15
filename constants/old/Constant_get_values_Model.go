package constants

/**
 * Andrey Mikhalevich 16/12/21
 * This file is part of the OSBE framework
 */

//Controller method model
import (
	"reflect"
	
		
	"osbe/fields"
)

type Constant_get_values_argv struct {
	Argv *Constant_get_values `json:"argv"`	
}

//Exported model metadata
var Constant_get_values_md fields.FieldCollection

func Constant_get_values_Model_init() {	
	Constant_get_values_md = fields.GenModelMD(reflect.ValueOf(Constant_get_values{}))
}

//
type Constant_get_values struct {
	Id_list fields.ValText `json:"id_list" required:"true"`
	Field_sep fields.ValText `json:"field_sep"`
}
