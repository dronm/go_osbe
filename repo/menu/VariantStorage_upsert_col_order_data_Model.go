package menu

//Controller method model
import (
	"reflect"
	
		
	"osbe/fields"
)

type VariantStorage_upsert_col_order_data_argv struct {
	Argv *VariantStorage_upsert_col_order_data `json:"argv"`	
}

//Exported model metadata
var VariantStorage_upsert_col_order_data_md fields.FieldCollection

func VariantStorage_upsert_col_order_data_Model_init() {	
	VariantStorage_upsert_col_order_data_md = fields.GenModelMD(reflect.ValueOf(VariantStorage_upsert_col_order_data{}))
}

//
type VariantStorage_upsert_col_order_data struct {
	Storage_name fields.ValText `json:"storage_name" required:"true"`
	Variant_name fields.ValText `json:"variant_name" required:"true"`
	Col_order fields.ValText `json:"col_order" required:"true"`
	Default_variant fields.ValBool `json:"default_variant" required:"true"`
}
