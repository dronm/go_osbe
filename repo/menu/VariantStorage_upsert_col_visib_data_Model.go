package menu

//Controller method model
import (
	"reflect"
	
		
	"osbe/fields"
)

type VariantStorage_upsert_col_visib_data_argv struct {
	Argv *VariantStorage_upsert_col_visib_data `json:"argv"`	
}

//Exported model metadata
var VariantStorage_upsert_col_visib_data_md fields.FieldCollection

func VariantStorage_upsert_col_visib_data_Model_init() {	
	VariantStorage_upsert_col_visib_data_md = fields.GenModelMD(reflect.ValueOf(VariantStorage_upsert_col_visib_data{}))
}

//
type VariantStorage_upsert_col_visib_data struct {
	Storage_name fields.ValText `json:"storage_name" required:"true"`
	Variant_name fields.ValText `json:"variant_name" required:"true"`
	Col_visib fields.ValText `json:"col_visib" required:"true"`
	Default_variant fields.ValBool `json:"default_variant" required:"true"`
}
