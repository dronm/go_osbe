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

type Constant_set_value_argv struct {
	Argv *Constant_set_value `json:"argv"`	
}

//Exported model metadata
var Constant_set_value_md fields.FieldCollection

func Constant_set_value_Model_init() {	
	Constant_set_value_md = fields.GenModelMD(reflect.ValueOf(Constant_set_value{}))
}

//
type Constant_set_value struct {
	Id fields.ValText `json:"id" required:"true"`
	Val fields.ValText `json:"val"`
}
