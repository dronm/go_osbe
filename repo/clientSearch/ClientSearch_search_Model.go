package clientSearch

//Controller method model
import (
		
	"osbe/fields"
)

type ClientSearch_search struct {
	Query fields.ValText `json:"query" required:"true"`
}
type ClientSearch_search_argv struct {
	Argv *ClientSearch_search `json:"argv"`	
}

