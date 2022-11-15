package model

import (
	"reflect"
	
	"osbe/fields"
)

//Exported model metadata
var Cond_Model_fields fields.FieldCollection

//Condition model
type Cond_Model struct {
	Count fields.ValInt `json:"count" notZero:true`
	From fields.ValInt `json:"from"`
	Cond_fields fields.ValText `json:"cond_fields" length:1000`
	Cond_sgns fields.ValText `json:"cond_sgns" length:1000`
	Cond_vals fields.ValText `json:"cond_vals" length:1000`
	Cond_ic fields.ValText `json:"cond_ic" length:1000`
	Cond_joins fields.ValText `json:"cond_joins" length:1000`
	Ord_fields fields.ValText `json:"ord_fields" length:1000`
	Ord_directs fields.ValText `json:"ord_directs" length:1000`
	Field_sep fields.ValText `json:"field_sep" length:2`
	Lsn fields.ValText `json:"lsn" length:50`
}

func Cond_Model_init() {
	Cond_Model_fields = fields.GenModelMD(reflect.ValueOf(Cond_Model{}))
}

type Controller_get_list_argv struct {
	Argv *Cond_Model `json:"argv"`	
}

