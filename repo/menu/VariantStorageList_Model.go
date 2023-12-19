package menu

import (
	"reflect"	
		
	"osbe/fields"
	"osbe/model"
)

type VariantStorageList struct {
	Id fields.ValInt `json:"id" primaryKey:"true"`
	User_id fields.ValInt `json:"user_id"`
	Storage_name fields.ValText `json:"storage_name"`
	Default_variant fields.ValBool `json:"default_variant"`
	Variant_name fields.ValText `json:"variant_name"`
}

func (o *VariantStorageList) SetNull() {
	o.Id.SetNull()
	o.User_id.SetNull()
	o.Storage_name.SetNull()
	o.Default_variant.SetNull()
	o.Variant_name.SetNull()
}

func NewModelMD_VariantStorageList() *model.ModelMD{
	return &model.ModelMD{Fields: fields.GenModelMD(reflect.ValueOf(VariantStorageList{})),
		ID: "VariantStorageList_Model",
		Relation: "variant_storages_list",
		AggFunctions: []*model.AggFunction{
			&model.AggFunction{Alias: "totalCount", Expr: "count(*)"},
		},
		LimitConstant: "doc_per_page_count",
	}
}
