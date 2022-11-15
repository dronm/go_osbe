package constants

/**
 * Andrey Mikhalevich 16/12/21
 * This file is part of the OSBE framework
 */

//Controller method model
import (
	"osbe/fields"
)

type Constant_get_values_argv struct {
	Argv *Constant_get_values `json:"argv"`	
}

//
type Constant_get_values struct {
	Id_list fields.ValText `json:"id_list" required:"true"`
	Field_sep fields.ValText `json:"field_sep"`
}

