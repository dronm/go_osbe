package constants

/**
 * Andrey Mikhalevich 15/12/21
 * This file is part of the OSBE framework
 */

import (
	"reflect"
	
	"osbe/fields"
)	

//Exported model metadata
var (
	Constant_keys_md fields.FieldCollection
)

func Constant_Model_init() {	
	Constant_keys_md = fields.GenModelMD(reflect.ValueOf(Constant_keys{}))
}

//Key
type Constant_keys struct {
	Id fields.ValText `json:"id"`
}

