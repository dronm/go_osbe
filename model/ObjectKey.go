package model

//Keys for delete/get object

import (
	"osbe/fields"
)

type Object_keys struct {
	Id fields.ValInt `json:"id"`
}
type Object_keys_argv struct {
	Argv *Object_keys `json:"argv"`	
}

