package test_app

import (
	//"osbe/model"	
	"osbe"
)

const (
	FIELD_id_ALIAS = "Object identifier"
)

func Get_Test_Model_obj_md() osbe.FieldCollection{

	md := make(osbe.FieldCollection, 5)
	
	f_Id := &osbe.FieldInt{}
	f_Id.Id = "Id"
	f_Id.Alias = FIELD_id_ALIAS
	f_Id.AutoInc = true
	md["Id"] = f_Id
	
	f_F1 := &osbe.FieldInt{}
	f_F1.Id = "F1"
	f_F1.Alias = "integer field (max=100)"
	f_F1.MaxValue = osbe.NewParamInt64(100)
	md["F1"] = f_F1

	f_F2 := &osbe.FieldText{}
	f_F2.Id = "F2"
	f_F2.Alias = "text field (20 chars)"
	f_F2.Length = osbe.NewParamInt(20)
	md["F2"] = f_F2

	f_F3 := &osbe.FieldFloat{}
	f_F3.Id = "F3"
	f_F3.Alias = "float field (15,2)"
	f_F3.Length = osbe.NewParamInt(15)
	f_F3.Precision = osbe.NewParamInt(2)
	md["F3"] = f_F3

	f_F4 := &osbe.FieldBool{}
	f_F4.Id = "F4"
	f_F4.Alias = "boolean field (required)"
	f_F4.Required = true
	md["F4"] = f_F4
		
	return md
}

func Get_Test_Model_get_list_md() osbe.FieldCollection{

	md := make(osbe.FieldCollection, 5)
	
	f_From := &osbe.FieldInt{}
	f_From.Id = "From"
	md["From"] = f_From
	
	f_Count := &osbe.FieldInt{}
	f_Count.Id = "Count"
	md["Count"] = f_Count

	f_Cond_fields := &osbe.FieldText{}
	f_Cond_fields.Id = "Cond_fields"
	md["Cond_fields"] = f_Cond_fields

	f_Cond_sgns := &osbe.FieldText{}
	f_Cond_sgns.Id = "Cond_sgns"
	md["cond_sgns"] = f_Cond_fields

	f_Cond_vals := &osbe.FieldText{}
	f_Cond_vals.Id = "Cond_vals"
	md["Cond_vals"] = f_Cond_vals

	f_Cond_ic := &osbe.FieldText{}
	f_Cond_ic.Id = "Cond_ic"
	md["Cond_ic"] = f_Cond_ic

	f_Ord_fields := &osbe.FieldText{}
	f_Ord_fields.Id = "Ord_fields"
	md["Ord_fields"] = f_Ord_fields

	f_Ord_directs := &osbe.FieldText{}
	f_Ord_directs.Id = "Ord_directs"
	md["Ord_directs"] = f_Ord_directs

	f_Field_sep := &osbe.FieldText{}
	f_Field_sep.Id = "Field_sep"
	f_Field_sep.Length = osbe.NewParamInt(1)
	md["Field_sep"] = f_Field_sep
		
	return md
}

func Get_Test_Model_object_md() osbe.FieldCollection{

	md := make(osbe.FieldCollection, 1)
	
	f_Id := &osbe.FieldInt{}
	f_Id.Id = "Id"
	f_Id.Alias = FIELD_id_ALIAS
	f_Id.Required = true
	md["Id"] = f_Id
		
	return md
}

func Get_Test_Model_update_md() osbe.FieldCollection{

	md := make(osbe.FieldCollection, 1)
	
	f_Id := &osbe.FieldInt{}
	f_Id.Id = "Old_id"
	f_Id.Alias = FIELD_id_ALIAS
	f_Id.Required = true
	md["Old_id"] = f_Id
		
	return md
}

//full object for insert
type Test_Model_obj struct {
	Id osbe.ValInt `json:"id"`
	F1 osbe.ValInt `json:"f1"`
	F2 osbe.ValText `json:"f2"`
	F3 osbe.ValFloat `json:"f3"`
	F4 osbe.ValBool `json:"f4"`
}

type Test_Model_keys struct {
	Id osbe.ValInt `json:"id"`
}

type Test_Model_cond struct {
	Count osbe.ValInt `json:"count"`
	From osbe.ValInt `json:"from"`
	Cond_fields osbe.ValText `json:"cond_fields"`
	Cond_sgns osbe.ValText `json:"cond_sgns"`
	Cond_vals osbe.ValText `json:"cond_vals"`
	Cond_ic osbe.ValText `json:"cond_ic"`
	Ord_fields osbe.ValText `json:"ord_fields"`
	Ord_directs osbe.ValText `json:"ord_directs"`
	Field_sep osbe.ValText `json:"field_sep"`
}
