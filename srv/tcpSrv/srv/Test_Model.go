package main

import (
	"reflect"
	
	"osbe/fields"
)

//Exported model metadata
var (
	Test_obj_md fields.FieldCollection
	Test_keys_md fields.FieldCollection
	Test_old_keys_md fields.FieldCollection
	Test_cond_md fields.FieldCollection
)

//Object model for insert
type Test_obj struct {
	Id fields.ValInt `json:"id" alias:"Ид" descr:"Object ID"`
	F1 fields.ValInt `json:"f1" alias:"Ф1" descr:"Поле int" maxValue:100 maxValue:100`
	F2 fields.ValText `json:"f2" alias:"Ф2" descr:"Поле text, length:20" length:20`
	F3 fields.ValFloat `json:"f3" alias:"Ф3" descr:"Поле float, 15,2" length:15 precision:2`
	F4 fields.ValBool `json:"f4" alias:"Ф4" descr:"Поле bool2"`
}
func Get_Test_obj_md() fields.FieldCollection {
	if Test_obj_md == nil {
		Test_obj_md = fields.GenModelMD(reflect.ValueOf(Test_obj{}))
	}
	return Test_obj_md
}

//object key model
type Test_keys struct {
	Id fields.ValInt `json:"id"`
}

func Get_Test_keys_md() fields.FieldCollection {
	if Test_keys_md == nil {
		Test_keys_md = fields.GenModelMD(reflect.ValueOf(Test_keys{}))
	}
	return Test_keys_md
}

type Test_old_keys struct {
	Old_id fields.ValInt `json:"old_id" required:true notZero:true primaryKey:true`
	Id fields.ValInt `json:"id" alias:"Ид" descr:"Object ID"`
	F1 fields.ValInt `json:"f1" alias:"Ф1" descr:"Поле int" maxValue:100 maxValue:100`
	F2 fields.ValText `json:"f2" alias:"Ф2" descr:"Поле text, length:20" length:20`
	F3 fields.ValFloat `json:"f3" alias:"Ф3" descr:"Поле float, 15,2" length:15 precision:2`
	F4 fields.ValBool `json:"f4" alias:"Ф4" descr:"Поле bool2"`
}
func Get_Test_old_keys_md() fields.FieldCollection {
	if Test_old_keys_md == nil {
		Test_old_keys_md = fields.GenModelMD(reflect.ValueOf(Test_old_keys{}))
	}
	return Test_old_keys_md
}

//Condition model
type Test_cond struct {
	Count fields.ValInt `json:"count" notZero:true`
	From fields.ValInt `json:"from"`
	Cond_fields fields.ValText `json:"cond_fields" length:1000`
	Cond_sgns fields.ValText `json:"cond_sgns" length:1000`
	Cond_vals fields.ValText `json:"cond_vals" length:1000`
	Cond_ic fields.ValText `json:"cond_ic" length:1000`
	Ord_fields fields.ValText `json:"ord_fields" length:1000`
	Ord_directs fields.ValText `json:"ord_directs" length:1000`
	Field_sep fields.ValText `json:"field_sep" length:2`
}
func Get_Test_cond_md() fields.FieldCollection {
	if Test_cond_md == nil {
		Test_cond_md = fields.GenModelMD(reflect.ValueOf(Test_cond{}))
	}
	return Test_cond_md
}

