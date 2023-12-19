package menu

//Controller method model
import (
	"osbe/fields"
)

//
type VariantStorage_get_col_order_data struct {
	Storage_name fields.ValText `json:"storage_name" required:"true"`
	Variant_name fields.ValText `json:"variant_name"`
}

type VariantStorage_get_col_order_data_argv struct {
	Argv *VariantStorage_get_col_order_data `json:"argv"`	
}


